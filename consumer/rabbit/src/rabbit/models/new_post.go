package models

import (
  "log"
  "encoding/json"
  "sync"
)


type New_post struct {
  User_id  string `json:"user_id"`
  Post_id  string `json:"post_id"`
  Tag      string `json:"tag"`
  Caption  string `json:"caption"`
}

func (post* New_post) Unmarshall(data []byte) error {
  return json.Unmarshal(data, &post)
}

func (post* New_post) write_timeline_to_cassandra(uuid_to_add string, monthyear string) {
  query := `INSERT INTO user_post_timeline(user_id, post_id, monthyear, caption) VALUES (?, ?, ?, ?)`
  if err := CassandraSession.Query(query, uuid_to_add, post.Post_id, monthyear, post.Caption).Exec(); err != nil {
    log.Fatal(err)
  }
}

func (post* New_post) update_subscribers_timelines(wg *sync.WaitGroup, monthyear string) {
  defer wg.Done()

  if post.Tag == "" {
    return
  }

  var user_id string
  iter := CassandraSession.Query(`SELECT user_id from subscription where tag = ? LIMIT 10000`, post.Tag).Iter()
  for iter.Scan(&user_id) {
    post.write_timeline_to_cassandra(user_id, monthyear)
  }
  if err := iter.Close(); err != nil {
    log.Println(err)
  }
}

func (post* New_post) update_followers_timelines(wg *sync.WaitGroup, monthyear string) {
  defer wg.Done()
  var user_id string
  iter := CassandraSession.Query(`SELECT followed_by from followed_by where user_id = ? LIMIT 10000`, post.User_id).Iter()
  for iter.Scan(&user_id) {
    post.write_timeline_to_cassandra(user_id, monthyear)
  }
  if err := iter.Close(); err != nil {
    log.Println(err)
  }
}
