const cassandra = require('cassandra-driver');
const client = new cassandra.Client({ contactPoints: ['localhost'], keyspace: 'instagram' });

exports.ping = async (req, res, next) => {
  res.send("PONG")
}

exports.getFromTags = async (req, res, next) => {

  try {
    var query
    var name
    if (req.body.tags) {
      name = req.body.name;
      query = `SELECT post_id, tag, caption from instagram.user_post_tag where tag = ? LIMIT 200`
    } else {
      name = req.body.user;
      query = `SELECT post_id, tag, caption from instagram.user_post where user_id = ? LIMIT 200`
    }
    const { rows } = await client.execute(query, [ name ])
    res.status(200).json(rows)

  } catch (err) {
    res.status(500).json({"error": err})
  }
}

exports.querySimply = async (req, res, next) => {

  const query = `SELECT post_id, caption from instagram.user_post_timeline where user_id = ? LIMIT 200 ALLOW FILTERING`

  try {
    const { rows } = await client.execute(query, [ req.params.user_id ])
    res.status(200).json(rows)

  } catch (err) {
    res.status(500).json({"error": err})
  }
}

exports.client = client;
