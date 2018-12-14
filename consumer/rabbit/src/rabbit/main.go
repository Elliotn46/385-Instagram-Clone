package main

import (
  "log"
  "sync"
  _ "rabbit/seed"
  "rabbit/models"
  "rabbit/graphics"
  "github.com/streadway/amqp"
)



func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}

func consume(q *amqp.Queue) {
  msgs, err := models.MQChannel.Consume(
    q.Name, // queue
    "",     // consumer
    false,   // auto-ack
    false,  // exclusive
    false,  // no-local
    false,  // no-wait
    nil,    // args
  )

  failOnError(err, "Failed to register a consumer")

  forever := make(chan bool)

  go func() {
    for d := range msgs {
      post := New_post{}
      err := post.unmarshall(d.Body)
      if err != nil
        || post.User_id == ""
        || post.Post_id == ""
        || post.Tag == ""
        || post.Caption == ""  {
        if err == nil {
          err = "missing required field"
        }
        log.Println("Malformed JSON object: ", err)
      }
      if graphics.Resize(post.Post_id) != nil {
        models.Update_user_timelines(post)
        d.Ack(false)
      }
    }
  }()

  log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
  <-forever
}

func main() {

  models.Init_mq()
  models.Init_cassandra("instagram")
  defer models.MQConnection.Close()
  defer models.CassandraSession.Close()

  models.Init_mq_chan()
  defer models.MQChannel.Close()

  q, err := models.MQChannel.QueueDeclare(
    "instagram", // name
    false,       // durable
    false,       // delete when unused
    false,       // exclusive
    false,       // no-wait
    nil,         // arguments
  )

  failOnError(err, "Failed to declare a queue")


  err = models.MQChannel.Qos(
    1,     // prefetch count
    0,     // prefetch size
    false, // global
  )
  failOnError(err, "Failed to change settings")

  var wg sync.WaitGroup
  wg.Add(1)

  go consume(&q);
  go consume(&q);

  wg.Wait()


}
