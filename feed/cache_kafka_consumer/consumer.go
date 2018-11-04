package main

import (
	"os/signal"
	"os"
	"time"
	"log"
  "github.com/go-redis/redis"
  "encoding/json"
  cluster "github.com/bsm/sarama-cluster"
  "strconv"
	// "sync"
)


// [2018-11-04 16:08:13,417] INFO [GroupCoordinator 1001]: Preparing to rebalance group test_group with old generation 1 (__consumer_offsets-12) (kafka.coordinator.group.GroupCoordinator)
// [2018-11-04 16:08:14,748] INFO [GroupCoordinator 1001]: Stabilized group test_group generation 2 (__consumer_offsets-12) (kafka.coordinator.group.GroupCoordinator)
// [2018-11-04 16:08:14,750] INFO [GroupCoordinator 1001]: Assignment received from leader for group test_group for generation 2 (kafka.coordinator.group.GroupCoordinator)

type instagram_post struct {
    User_id        int        `json:"id"`
    Img_url        string     `json:"url"`
    Comments_link  string     `json:"comment"`
}

var redis_client *redis.Client


// https://www.oreilly.com/library/view/kafka-the-definitive/9781491936153/ch04.html

func main() {

  log.Println("Connecting to redis")

  set_up_redis()
  defer redis_client.Close()

  log.Println("Connecting to kafka")
  config := cluster.NewConfig()
	config.Group.Mode = cluster.ConsumerModePartitions

	// init consumer
	brokers := []string{"kafka:9092"}
	topics := []string{"instagram_cache"}
	consumer, err := cluster.NewConsumer(brokers, "test_group", topics, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()
  log.Println("Connected to kafka")
	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume partitions
	for {
		select {
		case part, ok := <-consumer.Partitions():
			if !ok {
				return
			}

			// start a separate goroutine to consume messages
			go func(pc cluster.PartitionConsumer) {
        var ip instagram_post
				for msg := range pc.Messages() {

          log.Printf("%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)

          err := json.Unmarshal([]byte(msg.Value), &ip)

          if err != nil {
            log.Printf("Bad Kafka Value: %s\n", err)
          } else {
            insert_redis(ip.User_id, ip)
          }

          consumer.MarkOffset(msg, "")	// mark message as processed
				}
			}(part)
		case <-signals:
			return
		}
	}
}

func set_up_redis() {
  for true {
    log.Println("Connecting to redis")
    redis_client = redis.NewClient(&redis.Options{
  		Addr:     "redis:6379",
  		Password: "",
  		DB:       0,
  	})

  	_, err := redis_client.Ping().Result()

    if err != nil {
      log.Println("ERR")
      time.Sleep(5 * time.Second)
    } else {
      log.Println("Success")
      // rediscon_chan <- true
      return
    }
  }
}


func insert_redis(uuid int, post instagram_post) bool {

  redis_key := strconv.Itoa(uuid) + ":timeline"
  pst, err := json.Marshal(post)


  if err != nil {
    log.Printf("insert_redis(%s, %s) -> json.Marshal()", uuid, post)
    return false
  }

  size, err := redis_client.LPush(redis_key, string(pst)).Result()

  if err != nil {
    log.Printf("insert_redis(%s, %s) -> LPush -> %s\n", uuid, post, err)
    return false
  }

  if size >= 100 {
    redis_client.LTrim(redis_key, 0, 100)
  }

  return true
}
