package main

import (
  "testing"
  "time"
)

func TestNextReminderExit(t *testing.T){
  exit := make(chan bool)
  output := make(chan Reminder)
  check := make(chan Reminder)
  var r Reminder
  r.ReminderTime = time.Now().Add(time.Second * 10)
  go nextReminder(exit, output, check, r)
  exit <- true
  select {
    case <-exit:
      return
    case <- time.After(time.Second * 5):
      t.Error("The Go routine did not exit")
  }
}

func TestNextReminderSwap(t *testing.T) {
  exit := make(chan bool)
  output := make(chan Reminder)
  check := make(chan Reminder)
  var r Reminder
  var r2 Reminder
  r.ReminderTime = time.Now().Add(time.Second * 10)
  r2.ReminderTime = time.Now().Add(time.Second * 2)
  r.User_name = "Test Reminder 1"
  r2.User_name = "Test Reminder 2"
  go nextReminder(exit, output, check, r)
  check <- r2
  select {
    case got := <-output:
      if got.User_name != r2.User_name {
        t.Error("The user name returned did not match the expected")
      }
    case <- time.After(time.Second * 15):
      t.Error("The Go Routine exceeded expected completion time")
  }
  exit<-true
  <-exit
}

