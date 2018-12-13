package models

import (
  "log"
  "time"
  "os"
  "github.com/streadway/amqp"
)

var MQConnection  *amqp.Connection
var MQChannel     *amqp.Channel

func Init_mq() {
  for true {
	  conn, err := amqp.Dial("amqp://guest:guest@" + os.Getenv("RABBITMQ_SERVICE_SERVICE_HOST") + ":" + os.Getenv("RABBITMQ_SERVICE_SERVICE_PORT"))
    if err != nil {
      log.Println("Attempting to connect to rabbit mq")
      log.Println(err)
      time.Sleep(2 * time.Second)
    } else {
      MQConnection = conn
      log.Println("Connected to rabbit mq successfully")
      return
    }
  }
}

func Init_mq_chan() {
  ch, err := MQConnection.Channel()
  if err != nil {
    log.Fatalf("%s: %s", err, "Could Not Open Channel")
  }
  MQChannel = ch
}
