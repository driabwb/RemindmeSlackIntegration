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

  err = db.InsertRequest(reminder)
  if err != nil {
    return err
  }

  if err == nil {
    log.Print("here?")
  }

  /*
   *for _, message := range history.Messages {  // Check DB for users
   *  user_data, err := db.GetUser(message.User)
   *  if err != nil { // User not in DB?
   *    //TODO: Check error is id not found error
   *    slack_user, err := getUserData(message.User)
   *    if err != nil {
   *      return err
   *    }
   *    user_data = &User{
   *                  Id: slack_user.Id,
   *                  User_name: slack_user.User,
   *                  FirstName: slack_user.Profile.First_name,
   *                  LastName: slack_user.Profile.Last_name,
   *    }
   *    db.InsertUserData(*user_data)
   *  }
   *}
   */
  return nil
}
