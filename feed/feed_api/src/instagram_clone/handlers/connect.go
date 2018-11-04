package handlers

import (
  "log"
  "github.com/go-redis/redis"
  "time"
)

var Redis_client *redis.Client

type instagram_post struct {
    User_id        int        `json:"id"`
    Img_url        string     `json:"url"`
    Comments_link  string     `json:"comment"`
}

func Set_up_client(/*rediscon_chan chan<- bool*/) {
  for true {
    Redis_client = redis.NewClient(&redis.Options{
  		Addr:     "redis:6379",
  		Password: "",
  		DB:       0,
  	})

  	_, err := Redis_client.Ping().Result()

    if err != nil {
      log.Println(err)
      time.Sleep(5 * time.Second)
    } else {
      // rediscon_chan <- true
      return
    }
  }

}
