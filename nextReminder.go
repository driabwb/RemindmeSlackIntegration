package main

import (
  "time"
)

func nextReminder(exit chan bool, output <-chan Reminder, check chan<- Reminder, next Reminder) {
  for {
    sendTime := time.After(next.ReminderTime - time.Now())
    select {
    case <-sendTime:
      log.Println("Send Reminder")
      check <- next
      log.Println("Get Next Reminder")
    case newTime <-check:
      if newTime.Before(next) {
        next = newTime
      }
    case <-exit:
      done<-true
      return
  }
}
