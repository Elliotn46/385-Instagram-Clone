package main

import (
  "testing"
  "rabbit/models"
  "rabbit/seed"
)

func TestLogic(t *testing.T) {
  models.Init_cassandra("test")
  defer models.CassandraSession.Close()

  seed.Init_c_tables()
  users := seed.Create_many_users()

  err := seed.Add_users_to_cassandra(users)

  if err != nil {
    t.Error("Could not add users to cassandra: ", err)
  }

  seed.Users_create_fol_sub(users)

  mock_post := seed.Get_mock_post(users[0])

  models.Update_user_timelines(mock_post)

  count := models.Get_user_timeline_length()

  if count != 10 {
    t.Error("Incorrect number of timelines - error in logic")
  }
}
