package main

import (
  "net/http"
  "log"
  "os"
  "time"
  //"strings"
)

var TOKEN string

type Reminder struct{
  Context []Message
  ReminderTime time.Time
  User_name string
}

type User struct{
  Id string
  FirstName string
  LastName string
  User_name string
}

var check chan Reminder

func handleNewRequest(w http.ResponseWriter, r *http.Request){
  if r.Method != "POST" {
    log.Print("There was a non-POST message sent to the endpoint")
    w.WriteHeader(405)
    w.Write([]byte("Error: POST only"))
    return
  }
  r.ParseForm()
  token := r.Form["token"]
  channel_id := r.Form["channel_id"]
  user_name := r.Form["user_name"]
  log.Print("Token: " + token[0] + "| channel_id: " + channel_id[0] + "| user_name: " + user_name[0])
  //history, err := getHistory(strings.Join(channel_id, ","))
  w.Write([]byte("Your reminder has been set."))
}

func handleReminderTrigger(reminder *Reminder){
  log.Println("handle reminder trigger")
}

func httpserver(port string, done chan bool){
  log.Println("Server has started on " + port)
  http.HandleFunc("/api/reminder", handleNewRequest)
  log.Fatal(http.ListenAndServe(":"+port, nil))
  done <- true
}

func getNextReminder() *Reminder{
  var db CassandraDB
  db.Init()
  defer db.Close()

  rem, err := db.ReadNextRequest()
  if nil != err {
    log.Fatal(err)
    return nil
  }
  return rem
}

func init(){
  setupEnvironment()
  TOKEN = os.Getenv("SLACK_TOKEN")
}

func main(){
  httpServerDone := make(chan bool, 1)
  nextReminderDone := make(chan bool, 1)
  output := make(chan Reminder)
  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }
  go httpserver(port, httpServerDone)
  first := *getNextReminder()
  go nextReminder(nextReminderDone, output, check, first)
  <-httpServerDone
  nextReminderDone <- true
  <-nextReminderDone
}

