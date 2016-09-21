package main

import (
  "net/http"
  "log"
  "os"
  "time"
  "strconv"

  "gopkg.in/tylerb/graceful.v1"
)

var TOKEN string
var SECRET_TOKEN string

type Reminder struct{
  Context []Message
  ReminderTime time.Time
  User_name string
  User_id string
}

type User struct{
  Id string
  FirstName string
  LastName string
  User_name string
}

func reminderRequest(check chan Reminder) func (http.ResponseWriter, *http.Request) {

  return func (w http.ResponseWriter, r *http.Request){
    if r.Method != "POST" {
      log.Print("There was a non-POST message sent to the endpoint")
      w.WriteHeader(405)
      w.Write([]byte("Error: POST only"))
      return
    }
    err := r.ParseForm()
    if err != nil {
      log.Print(err)
      w.WriteHeader(400)
      w.Write([]byte("Error: Invalid form encoding"))
      return
    }

    token := r.Form["token"]
    channel_id := r.Form["channel_id"]
    user_id:= r.Form["user_id"]
    time_delta := r.Form["text"]

    if len(token) == 0 {
      log.Print("No token provided")
      w.WriteHeader(405)
      w.Write([]byte("Error: Token required"))
      return
    }
    if token[0] != SECRET_TOKEN {
      log.Print("Unauthorized POST attempt")
      w.WriteHeader(405)
      w.Write([]byte("Error: Invalid Token"))
      return
    }

    td, err := strconv.Atoi(time_delta[0])
    if err != nil {
      log.Print("Error: Time delta integer conversion failed")
      w.WriteHeader(400)
      w.Write([]byte("Invalid time"))
      return
    }

    w.Write([]byte("Your reminder has been set."))
    createReminder(channel_id[0], user_id[0], td, check)
    log.Print("Token: " + token[0] + "| channel_id: " + channel_id[0] + "| user_id: " + user_id[0])
  }
}

func handleReminderTrigger(reminder *Reminder){
  log.Println("handle reminder trigger")
}

func httpserver(port string, check chan Reminder, done chan bool){
  log.Println("Server has started on " + port)
  mux := http.NewServeMux()
  mux.HandleFunc("/api/reminder", reminderRequest(check))
  graceful.Run(":"+port, time.Second, mux)
  done <- true
}

func getNextReminder() *Reminder{
  var db CassandraDB
  db.Init()
  defer db.Close()

  rem, err := db.ReadNextRequest()
  if rem == nil {
    rem = new(Reminder)
    rem.ReminderTime = time.Now().Add(1000 * time.Hour)
  }
  if nil != err {
    log.Print(err)
    return nil
  }
  return rem
}

func init(){
  setupEnvironment()
  TOKEN = os.Getenv("SLACK_TOKEN")
  SECRET_TOKEN = os.Getenv("SECRET_TOKEN")
}

func main(){
  httpServerDone := make(chan bool, 1)
  nextReminderDone := make(chan bool, 1)
  //sendReminderDone := make(chan bool, 1)
  output := make(chan Reminder)
  check := make(chan Reminder)

  if 2 <= len(os.Args) && "-s" == os.Args[1] {
    setupDB("cassandraInitScript.cql")
  }
  if 2 <= len(os.Args) && "-d" == os.Args[1] {
    setupDB("cassandraDestroyScript.cql")
    return
  }

  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }
  go httpserver(port, check, httpServerDone)
  first := *getNextReminder()
  go nextReminder(nextReminderDone, output, check, first)
  //go sendReminder(output, sendReminderDone)
  <-httpServerDone
  nextReminderDone <- true
  //sendReminderDone <- true
  <-nextReminderDone
}

