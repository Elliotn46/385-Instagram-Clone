const express = require('express');

const controllers = require('../cassandra/query')


const Router = express.Router();

Router.route('/ping')
  .get(controllers.ping)

Router.route('/')
  .post(controllers.getFromTags)
  .get(controllers.ping)

Router.route('/userTimeline/:user_id')
  .post(controllers.querySimply)

module.exports = Router
