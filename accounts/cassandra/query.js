const cassandra = require('cassandra-driver')
const bcrypt = require('bcrypt')
const client = new cassandra.Client({ contactPoints: ['cassandra'], keyspace: 'instagram' })
const jwt = require('jsonwebtoken')
require('dotenv').config()

const saltRounds = 10

const signToken = (username, email, user_id) => {
  return new Promise((resolve, reject) => jwt.sign(
    { username, email, user_id },
    process.env.JWT_SECRET,
    (err, token) => {
      if (err) reject(err)
      else resolve(token)
    }
  ))
}

exports.pong = async (req, res, next) => {
  res.send("PONG")
}

exports.register = async (req, res, next) => {
  // TODO: rate limit to prevent abuse
  const query = `select * from instagram.user where username = ? allow filtering`
  try {
    //This is bad - Should use zookeeper
    const { rows } = await client.execute(query, [ req.body.username ])
    if (rows.length === 0) {
      const query = `insert into instagram.user(user_id, username, email, password) VALUES (?, ?, ?, ?)`
      const timeuuid = cassandra.types.TimeUuid.now()
      const params = [
        timeuuid,
        req.body.username,
        req.body.email,
        await bcrypt.hash(req.body.password, saltRounds)
      ]
      await client.execute(query, params, { prepare: true })
      const token = signToken(req.body.username, req.body.email, timeuuid)
      res.status(200).json({ "status": "success", token })
    } else {
      res.status(400).json({ "status": "bad request", "error": "invalid credentials or account already exists" })
    }
  } catch (err) {
    res.status(500).json({ "status": "internal server error", "error": err })
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
        const token = await signToken(username, email, user_id)
        res.status(200).json({ "status": "success", token })
      } else {
        res.set("WWW-Authenticate", "login")
        res.status(401).json({ "status": "unauthorized", "error": "bad credentials" })
      }
    }
  } catch (err) {
    res.status(500).json({ "status": "internal server error", "error": err })
  }

}

exports.client = client;
