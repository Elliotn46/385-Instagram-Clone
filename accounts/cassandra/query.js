const cassandra = require('cassandra-driver')
const bcrypt = require('bcrypt')
const client = new cassandra.Client({ contactPoints: ['localhost'], keyspace: 'instagram' })
require('dotenv').config()

const signToken = (username, email, user_id) => {
  return jwt.sign({
    username,
    email,
    user_id
  }, process.env.jwt_token)
}

exports.pong = async (req, res, next) => {
  res.send("PONG")
}

exports.register = async (req, res, next) => {

  const query =  `select * from instagram.user where username = ? allow filtering`

  try {
    //This is bad - Should use zookeeper
    const { rows } = await client.execute(query, [ req.body.username || "foo" ])
    if (rows.length === 0) {
      const query = `insert into instagram.user(user_id, username, email, password) VALUES (?, ?, ?, ?)`
      const timeuuid = cassandra.types.TimeUuid.now()
      const params = [
        timeuuid,
        req.body.username || "testtt",
        req.body.email || "testtt",
        bcrypt.hashSync(req.body.password || "testtt", 10)
      ]
      await client.execute(query, params, { prepare: true  })
      res.cookie('jwt_token', signToken(req.body.username, req.body.email, timeuuid))

      res.status(200).json({"Success": true})

    } else {
      throw "Username Taken"
    }
  } catch (err) {
    res.status(500).json({"error": err})
  }

}

exports.login = async (req, res, next) => {
  const query = `select * from instagram.user where username = ? allow filtering`
  try {
    const { rows } = await client.execute(query, req.body.username)
    if (rows.length === 1) {
      const match = await bcrypt.compare(req.body.password, rows[0].password);
      if (match) {
        let {
          username,
          email,
          user_id
        } = rows[0]

        res.cookie('jwt_token', signToken(username, email, user_id))
        res.status(200).json({"Success": true})
        return
      } else {
        throw "Could not find account"
      }
    }
  } catch (err) {
    res.status(500).json({"error": err})
  }

}

exports.client = client;
