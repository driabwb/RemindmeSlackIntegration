package main

import (
  "time"
  "log"
)

func createReminder(channel_id string, user_id string, time_delta int) error {
  var db CassandraDB
  db.Init()
  defer db.Close()

  trigger_time := time.Now().Add(time.Duration(time_delta) * time.Minute)

  reminder := new(Reminder)
  reminder.User_id = user_id
  reminder.ReminderTime = trigger_time

  history, err := getHistory(channel_id)
  if err != nil {
    log.Print("Failed to get channel history")
    return err
  }

  reminder.Context = history.Messages

  for _, message := range history.Messages {  // Check DB for users
    user_data, err := db.GetUser(message.User)
    if err == nil { // User in DB?
      message.username = user_data.User_name
    } else {
      //TODO: Check error is id not found error
      slack_user, err := getUserData(message.User)
      if err != nil {
        return err
      }
      user_data = &User{
                    Id: slack_user.User.Id,
                    User_name: slack_user.User.Name,
                    FirstName: slack_user.User.Profile.First_name,
                    LastName: slack_user.User.Profile.Last_name,
      }
      db.InsertUserData(*user_data)
      message.username = user_data.User_name
      log.Print(message)
    }
    reminder.Context = append(reminder.Context, message)

  }

  err = db.InsertRequest(reminder)
  if err != nil {
    return err
  }

  return nil
}
