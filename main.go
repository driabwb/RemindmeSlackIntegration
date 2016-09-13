package main

import (
 "encoding/json"
  "net/http"
  "log"
  "os"
  "time"
)

var TOKEN string = os.Getenv("SLACK_TOKEN")
const APIBASEURL string = "https://slack.com/api/"
var APIENDPOINTS = map[string]map[string]string{
  "history": {"endpoint": "channels.history", "param": "channel"},
  "user": {"endpoint": "users.info", "param": "user"},
}

type JsonObject struct{
  Obj map[string]string
}

type Message struct{
  message string
}

type Reminder struct{
  context []Message
  reminderTime time.Time
  user_name string
}

func getJson(url string, target interface{}) error {
  r, err := http.Get(url)
  if err != nil {
      return err
  }
  defer r.Body.Close()

  return json.NewDecoder(r.Body).Decode(target)
}

func makeAPICall(endpoint map[string]string, param map[string]string) *JsonObject {
  result := new(JsonObject)
  url := APIBASEURL + endpoint["endpoint"] + "?token=" + TOKEN
  if endpoint_param, ok := endpoint["param"]; ok {
    url += "&" + endpoint_param + "=" + param[endpoint_param]
  }
  getJson(url, result)
  return result
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
 log.Print("Token: " + token[0] + "| channel_id: " + channel_id[0] + "| user_name: " + user_name[0])
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

func main(){
  done := make(chan bool, 1)
  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }
  go httpserver(port, done)
  <-done
}

