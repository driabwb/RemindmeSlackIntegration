package main

import (
  "encoding/json"
  "net/http"
  "log"
  "os"
)

func handleNewRequest(w http.ResponseWriter, r *http.Request){
  log.Println("handle request")
}

func handleReminderTrigger(reminder *Reminder){
  log.Println("handle reminder trigger");
}

func httpserver(port int, done chan bool){
  log.Println("Server has started on " + port)
  http.HandleFunc("/api/reminder", handleNewRequest)
  log.Fatal(http.ListenAndServe(":"+port, nil))
  done <- true
}

func main(){
  done := make(chan bool, 1)
  port := os.Getenv("PORT")
  if port == "" {
    port = 8080
  }
  go httpserver(port, done)
  <-done
}

