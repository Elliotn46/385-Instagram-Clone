package handlers

import (
  "log"
)

const (
  userKey     = ":user"
  timelineKey = ":timeline"
)


func get_timeline(uuid string) []string {
  redis_key := uuid + timelineKey
  els, err := Redis_client.LRange(redis_key, 0, -1).Result()

  if err != nil {
    log.Print("err is ", err)
  }

  return els
}
