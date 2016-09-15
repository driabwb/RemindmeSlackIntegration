package main

import (
  "net/http"
  "log"
  "os"
  "time"
  //"strings"
)

var TOKEN string
var SECRET_TOKEN string

type Reminder struct{
  Context []Message
  ReminderTime time.Time
  User_name string
}

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

  if len(token) == 0 {
    log.Print("No token provided")
    w.WriteHeader(405)
    w.Write([]byte("Error: Token required"))
  }
  if token[0] != SECRET_TOKEN {
    log.Print("Unauthorized POST attempt")
    w.WriteHeader(405)
    w.Write([]byte("Error: Invalid Token"))
    return
  }

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

func init(){
  setupEnvironment()
  TOKEN = os.Getenv("SLACK_TOKEN")
  SECRET_TOKEN = os.Getenv("SECRET_TOKEN")
}

func main(){
  done := make(chan bool, 1)
  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }
  go httpserver(port, done)
  <-done
}

