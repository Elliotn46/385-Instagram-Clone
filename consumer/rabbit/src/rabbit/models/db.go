package models

import (
  "github.com/gocql/gocql"
  "log"
  "sync"
  "time"
  "strconv"
  "os"
)

var CassandraSession *gocql.Session

func Init_cassandra(keyspace string) {
  cluster := gocql.NewCluster(os.Getenv("CASSANDRA_DB"))
  cluster.Consistency = gocql.One
  if keyspace != "" {
    cluster.Keyspace = keyspace
  }
  for true {
    sess, err := cluster.CreateSession()
    if err != nil {
      log.Println("Attempting to connect to cassandra")
      time.Sleep(2 * time.Second)
    } else {
      log.Printf("Connected to Cassandra Successfully\n")
      CassandraSession = sess
      break
    }
  }
}

func Do_not_call_drop_keyspace() error {
  if err := CassandraSession.Query(`DROP KEYSPACE IF EXISTS instagram`).Exec(); err != nil {
    log.Printf("Error Dropping Keyspace, ", err)
    return err
  }
  return nil
}

//Aggregation query used without partition key
// ok for tests
func Get_user_timeline_length() int {
  var count string
  if err := CassandraSession.Query(`select COUNT(*) from user_post_timeline`).Consistency(gocql.One).Scan(&count); err != nil {
    log.Fatal(err)
  }

  int, error := strconv.Atoi(count)
  if error != nil {
    return 0
  }
  return int
}

func Update_user_timelines(data []byte) {
  post := New_post{}
  err := post.unmarshall(data)
  if err != nil || post.User_id == "" || post.Post_id == "" || post.Tag == "" || post.Caption == ""  {
    log.Println("Malformed JSON object, ", err)
    return
  }
  if post.User_id == "" || post.Post_id == "" || /*post.Tag == "" ||*/ post.Caption == "" {
    return
  }

  y, m, _ := time.Now().Date()
  monthyear := m.String() + "/" + strconv.Itoa(y)

  var wg sync.WaitGroup
  wg.Add(2)

  go post.update_followers_timelines(&wg, monthyear)
  go post.update_subscribers_timelines(&wg, monthyear)
  wg.Wait()
}
