const express = require('express');

const controllers = require('../cassandra/query')


const Router = express.Router();

Router.route('/ping')
  .get(controllers.pong)

Router.route('/register')
  .post(controllers.register)

Router.route('/token')
  .post(controllers.login)

module.exports = Router
