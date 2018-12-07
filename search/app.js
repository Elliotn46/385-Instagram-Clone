const express = require('express');
const bodyParser = require('body-parser');
const routes = require('./routes/router');
const cassandra = require('./cassandra/query');
const app = express();

app.use(bodyParser.json())

app.use('/search', routes)

const port = process.env.PORT || 3001;

const sleep = ms => {
  return new Promise(resolve => setTimeout(resolve, ms))
}

(async () => {
  while (true) {
    try {
      await cassandra.client.connect();
      break;
    } catch (e) {
      console.log("Could not connect to cassandra...sleeping")
      await sleep(5000);
    }
  }
  console.log('Connected')
  app.listen(port, () => {
      console.log(`App is running on port ${port}`);
  })
})()
