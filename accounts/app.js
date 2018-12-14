const express = require('express');
const bodyParser = require('body-parser');
const routes = require('./routes/router');
const cassandra = require('./cassandra/query');
require('dotenv').config()

const app = express();

app.use(bodyParser.json())

app.use('/account', routes)

// XXX: THIS SHOULD REQUIRE AUTHENTICATION!
app.get('/cpu', (req, res) => {
   //Compute last digit of Pi
   var Pi=0;
   var n=1;
   for (i=0;i<=10000000000;i++) {
      Pi=Pi+(4/n)
      n=n+2
      Pi=Pi-(4/n)
      n=n+2
   }
   //function fibonacci(n) { return n < 1 ? 0 : n <= 2 ? 1 : fibonacci(n - 1) + fibonacci(n - 2); }
   //let result = fibonacci(Math.floor(Math.random() * 300));
   res.send(`Pi is ${Pi}`);
})

if (!process.env.JWT_SECRET)
  throw new Error("JWT_SECRET missing from process environment")

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
