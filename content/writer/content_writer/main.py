"""
content service endpoint (writer methods)

POST:   Multipart file upload, including comments and tags
PUT:    Replaces description and tags of content
DELETE: Remove a piece of content

This service forwards the upload to the object store and queues
events for processing by other components of the system.

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

from aiohttp import ClientSession, MultipartWriter, hdrs, web
import aioamqp
from aioamqp.channel import Channel
import uuid


UPLOAD_BUCKET_URL = 'https://www.googleapis.com/upload/storage/v1/b/uploads'


async def mq_connect(app: web.Application):
    transport, protocol = await aioamqp.connect(os.environ['RABBITMQ_SERVICE_SERVICE_HOST'],
                                                os.environ['RABBITMQ_SERVICE_SERVICE_PORT'])
    app['mq'] = (transport, protocol)
    async with protocol.channel() as chan:
        await chan.declare_queue(queue_name='uploads')
        await chan.declare_queue(queue_name='edits')
        await chan.declare_queue(queue_name='deletes')

async def mq_channel(request: web.Request) -> Channel:
    return request.app['mq'][1].channel()

async def mq_close(app: web.Application):
    transport, protocol = app['mq']
    await protocol.close()
    transport.close()


async def read_multipart(request: web.Request):
    reader = await request.multipart()
    part1 = await reader.next()
    if part1 is None or part1.headers[hdrs.CONTENT_TYPE] != 'application/json':
        raise web.HTTPNotAcceptable()
    metadata = await part1.json()
    part2 = await reader.next()
    if part2 is None or part2.headers[hdrs.CONTENT_TYPE] != 'image/jpeg':
        raise web.HTTPNotAcceptable()
    jpeg = part1
    return (metadata, jpeg)

async def write_multipart(identifier, metadata, jpeg):
    with MultipartWriter('mixed/related', boundary='-- --') as mpwriter:
        mpwriter.append_json(metadata)
        mpwriter.append(jpeg, {'Content-Type': 'image/jpeg'})
        async with ClientSession() as session:
            async with session.post(f'{UPLOAD_BUCKET_URL}/{identifier}_original.jpeg',
                    params={'uploadType': 'multipart'}, data=mpwriter) as resp:
                if not resp.ok:
                    raise web.HTTPServiceUnavailable()


async def post(request: web.Request):
    metadata, jpeg = await read_multipart(request)
    identifier = uuid.uuid4()
    await write_multipart(identifier, {'original filename': jpeg.filename}, jpeg)
    metadata['content_id'] = identifier
    async with mq_channel(request) as channel:
        await channel.basic_publish(payload=metadata, exchange_name='', routing_key='uploads')
    raise web.HTTPAccepted()

async def put(request: web.Request) -> web.Response:
    metadata = await request.json()
    metadata['content_id'] = request.match_info['content_id']
    async with mq_channel(request) as chan:
        await chan.basic_publish(payload=metadata, exchange_name='', routing_key='edits')
    raise web.HTTPAccepted()

async def delete(request: web.Request) -> web.Response:
    identifier = request.match_info['content_id']
    async with mq_channel(request) as chan:
        await chan.basic_publish(payload=identifier, exchange_name='', routing_key='deletes')
    raise web.HTTPNoContent()


async def init_app(argv=None) -> web.Application:
    
    app = web.Application()
    
    app.on_startup.append(mq_connect)
    app.on_cleanup.append(mq_close)

    app.router.add_post('/', post)
    app.router.add_put('/', put)
    app.router.add_delete('/', delete)

    return app


def main(argv):
    app = init_app(argv)
    web.run_app(app)


if __name__ == '__main__':
    main(sys.argv[1:])
