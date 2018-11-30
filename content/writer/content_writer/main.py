"""
content service endpoint (writer methods)

POST:   Multipart file upload, including comments and tags
PUT:    Replaces description and tags of content
DELETE: Remove a piece of content

This service forwards the upload to the object store and queues
events for processing by other components of the system.

Configuration is passed through several environment variables.
The following variables must be set:
	- STORAGE_BUCKET: the Google Cloud Storage bucket name
	- RABBITMQ_SERVICE_SERVICE_HOST: the RabbitMQ host (provided by Kubernetes)
	- RABBITMQ_SERVICE_SERVICE_PORT: the RabbitMQ port (provided by Kubernetes)

Google Cloud service credentials should work automatically when running on
Google Cloud infrastructure.

To run the service in a Docker container, you can pass the credentials through
the environment, for example:

```
docker network create igclone
docker run -d --name rabbitmq --network igclone rabbitmq
# Wait for rabbitmq to start up (this could be handled better)
docker run -d \
	--name content_writer \
	--network igclone \
	-e RABBITMQ_SERVICE_SERVICE_HOST=rabbitmq \
	-e RABBITMQ_SERVICE_SERVICE_PORT=5672 \
	-e STORAGE_BUCKET=385ig \
	-e GOOGLE_CLOUD_PROJECT=cs385fa18 \
	-e GOOGLE_APPLICATION_CREDENTIALS=/application_default_credentials.json \
	-v $HOME/.config/gcloud/application_default_credentials.json:/application_default_credentials.json \
	content_writer
```

TODO:
    - request validation
    - logging
    - authz/authn
"""

"""
ISC License

Copyright (c) 2018, Ryan Moeller

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted, provided that the above
copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
"""

import os, sys
import json, uuid
from typing import Tuple

from aiohttp import ClientSession, MultipartReader, MultipartWriter, hdrs, web
import aioamqp
from aioamqp.channel import Channel

from google.cloud import storage


STORAGE_BUCKET = os.environ['STORAGE_BUCKET']


async def mq_connect(app: web.Application):
    transport, protocol = await aioamqp.connect(os.environ['RABBITMQ_SERVICE_SERVICE_HOST'],
                                                os.environ['RABBITMQ_SERVICE_SERVICE_PORT'])
    channel = await protocol.channel()
    await channel.queue_declare(queue_name='uploads', durable=True)
    await channel.queue_declare(queue_name='edits',   durable=True)
    await channel.queue_declare(queue_name='deletes', durable=True)
    await channel.close()
    app['mq'] = (transport, protocol)

def mq_channel(request: web.Request) -> Channel:
    return request.app['mq'][1].channel()

async def mq_publish(request: web.Request, metadata: dict, key: str):
    channel = await mq_channel(request)
    await channel.basic_publish(
        payload=json.dumps(metadata),
        exchange_name='',
        routing_key=key,
        properties={'delivery_mode': 2}
    )
    await channel.close()

async def mq_close(app: web.Application):
    transport, protocol = app['mq']
    await protocol.close()
    transport.close()


async def read_multipart(reader: MultipartReader) -> Tuple[dict, bytearray]:
    part1 = await reader.next()
    if part1 is None or part1.headers[hdrs.CONTENT_TYPE] != 'application/json':
        raise web.HTTPNotAcceptable(reason='expected metadata')
    metadata = await part1.json()
    part2 = await reader.next()
    if part2 is None or part2.headers[hdrs.CONTENT_TYPE] != 'image/jpeg':
        raise web.HTTPNotAcceptable(reason='expected image data')
    jpeg = await part2.read(decode=False)
    return (metadata, jpeg)


def storage_blob_name(identifier: str) -> str:
    return f'{identifier}.jpeg'

def storage_upload(request: web.Request, identifier: str, jpeg: bytes):
    storage_client = request.app['storage_client']
    bucket = storage_client.bucket(STORAGE_BUCKET)
    blob = bucket.blob(storage_blob_name(identifier))
    # TODO: get user_id from request
    blob.metadata = {'user_id': 'TODO', 'content_id': identifier}
    # XXX: bytes() probably makes a copy, AND the upload is NOT async!
    blob.upload_from_string(bytes(jpeg), content_type='image/jpeg')

def storage_delete(request: web.Request, identifier: str):
    storage_client = request.app['storage_client']
    bucket = storage_client.bucket(STORAGE_BUCKET)
    blob = bucket.blob(storage_blob_name(identifier))
    # XXX: NOT async!
    blob.delete()


async def post(request: web.Request) -> web.Response:
    reader = await request.multipart()
    metadata, jpeg = await read_multipart(reader)
    identifier = str(uuid.uuid4())
    storage_upload(request, identifier, bytes(jpeg))
    metadata['content_id'] = identifier
    await mq_publish(request, metadata, 'uploads')
    raise web.HTTPCreated(headers={'Location': f'/content/{identifier}'})

async def put(request: web.Request) -> web.Response:
    metadata = await request.json()
    metadata['content_id'] = request.match_info['content_id']
    await mq_publish(request, metadata, 'edits')
    raise web.HTTPNoContent()

async def delete(request: web.Request) -> web.Response:
    identifier = request.match_info['content_id']
    storage_delete(request, identifier)
    metadata = {'content_id': identifier}
    await mq_publish(request, metadata, 'deletes')
    raise web.HTTPNoContent()


async def init_app(argv=None) -> web.Application:

    app = web.Application()

    app['storage_client'] = storage.Client()

    app.on_startup.append(mq_connect)
    app.on_cleanup.append(mq_close)

    app.router.add_post('/', post)
    app.router.add_put('/{content_id}', put)
    app.router.add_delete('/{content_id}', delete)

    return app


def main(argv):
    app = init_app(argv)
    web.run_app(app)


if __name__ == '__main__':
    main(sys.argv[1:])
