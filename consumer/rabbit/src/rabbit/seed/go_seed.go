package seed

import (
  "log"
  "strconv"
  "time"
  "rabbit/models"
  "encoding/json"
  "github.com/gocql/gocql"
  "golang.org/x/crypto/bcrypt"
)




func Init_c_tables() error {
  tables := []string{
    // "DROP KEYSPACE IF EXISTS instagram",
    "CREATE KEYSPACE IF NOT EXISTS instagram WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor' : 3}",
    "CREATE KEYSPACE IF NOT EXISTS test WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor' : 3}",
    "CREATE TABLE IF NOT EXISTS user(user_id timeuuid, username text, email text, password text, PRIMARY KEY (user_id));",
    "CREATE TABLE IF NOT EXISTS subscription(user_id timeuuid, tag text, subscribe_date timestamp, PRIMARY KEY (tag, subscribe_date)) WITH CLUSTERING ORDER BY (subscribe_date DESC);",
    "CREATE TABLE IF NOT EXISTS follows(user_id timeuuid, follows timeuuid, follow_date timestamp, PRIMARY KEY (user_id, follow_date)) WITH CLUSTERING ORDER BY (follow_date DESC);",
    "CREATE TABLE IF NOT EXISTS followed_by(user_id timeuuid, followed_by timeuuid, follow_date timestamp, PRIMARY KEY (user_id, follow_date)) WITH CLUSTERING ORDER BY (follow_date DESC);",
    "CREATE TABLE IF NOT EXISTS user_post(user_id timeuuid, post_id timeuuid, tag text, caption text, PRIMARY KEY (user_id, post_id)) WITH CLUSTERING ORDER BY (post_id DESC);",
    "CREATE TABLE IF NOT EXISTS user_post_tag(user_id timeuuid, post_id timeuuid, tag text, caption text, PRIMARY KEY(tag, post_id)) WITH CLUSTERING ORDER BY (post_id DESC);",
    "CREATE TABLE IF NOT EXISTS user_post_comment(user_id timeuuid, post_id timeuuid, comment text, PRIMARY KEY (user_id, post_id)) WITH CLUSTERING ORDER BY (post_id DESC);",
    "CREATE TABLE IF NOT EXISTS user_post_timeline(user_id timeuuid, post_id timeuuid, monthyear text, caption text, PRIMARY KEY ((user_id, monthyear), post_id)) WITH CLUSTERING ORDER BY (post_id DESC);"}

  for _, query := range tables {
    if err := models.CassandraSession.Query(query).Exec(); err != nil {
      log.Printf("err with query: %s\n", query)
      return err
    }
  }
  return nil
}

func instagram_follow(user1 string, user2 string, query_value string, timestamp time.Time) {
  query := "insert into " + query_value + "(user_id, " + query_value + ", follow_date) VALUES (?, ?, ?)"
  if err := models.CassandraSession.Query(query, user1, user2, timestamp).DefaultTimestamp(true).Exec(); err != nil {
    log.Fatal(err)
  }
}

func instagram_subscribe(userid string, tag string) {
  query := "insert into subscription(user_id, tag, subscribe_date) VALUES (?, ?, ?)"
  if err := models.CassandraSession.Query(query, userid, tag, time.Now()).Exec(); err != nil {
    log.Fatal(err)
  }
}


func follow(user string, to_follow string) {
  stamp := time.Now()
  instagram_follow(user, to_follow, "follows", stamp)
  instagram_follow(to_follow, user, "followed_by", stamp)
}

func Users_create_fol_sub(users []models.User) {
  user_main := users[0]
  for i, value := range users {
    instagram_subscribe(value.User_id, "test_subscribe")
    if i == 0 {
      continue
    }
    follow(user_main.User_id, value.User_id)
  }
}

func Add_users_to_cassandra(users []models.User) error {
  query := `insert into user(user_id, username, email, password) VALUES (?, ?, ?, ?)`
  for _, value := range users {
    if err := models.CassandraSession.Query(query, value.User_id, value.Username, value.Email, value.Password).Exec(); err != nil {
      log.Fatal(err)
      return err
    }
  }
  return nil

}


func get_user(username string, password string, email string) (models.User, error) {
  uuid := gocql.TimeUUID()

  hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), 10)
  if err != nil {
    return models.User{}, err
  }

  return models.User{uuid.String(), username, email, string(hashed_password)}, nil
}

func Create_many_users() []models.User {
  user_struct := []models.User{}
  for i := 0; i < 10; i++ {
    iter := strconv.Itoa(i)
    user, _ := get_user("test" + iter, "test_" + iter, "test__" + iter + "@sonoma.edu")
    user_struct = append(user_struct, user)
  }
  return user_struct
}

func Get_mock_post(user models.User) []byte {
  post := models.New_post{user.User_id, gocql.TimeUUID().String(), "test_subscribe", "test post"}
  bt, err := json.Marshal(post)
  if err != nil {
    log.Println(err)
  }
  return bt
}
