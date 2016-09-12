package main

import (
  "encoding/json"
  "net/http"
  "log"
  "os"
  "time"
)

type Message struct{
  message string
}

type Reminder struct{
  context []Message
  reminderTime time
  user_name string
}

func handleNewRequest(w http.ResponseWriter, r *http.Request){
  if r.Method != "POST" {
    log.Print("There was a non-POST message sent to the endpoint")
  }
 r.ParseForm()
 token := r.Form["token"]
 channel_id := r.Form["channel_id"]
 user_name := r.Form["user_name"]
 w.Write("Your reminder has been set.")
}

func handleReminderTrigger(reminder *Reminder){
  log.Println("handle reminder trigger")
}

func httpserver(port int){
  log.Println("Server has started on " + port)
  http.HandleFunc("/api/reminder", handleNewRequest)
  log.Fatal(http.ListenAndServe(":"+port, nil))
}

func main(){
  port := os.Getenv("PORT")
  if port == "" {
    port = 8080
  }
  go httpserver(port)
}
