package main

import (
  "time"
  "log"
)

func nextReminder(exit chan bool, output chan<- Reminder, check <-chan Reminder, next Reminder) {
  for {
    sendTime := time.After(next.ReminderTime.Sub(time.Now()))
    select {
    case <-sendTime:
      log.Println("Send Reminder")
      output <- next
      log.Println("Get Next Reminder")
      next = *getNextReminder()
    case newReminder := <-check:
      if newReminder.ReminderTime.Before(next.ReminderTime) {
        next = newReminder
      }
    case <-exit:
      exit<-true
      return
    }
  }
}
