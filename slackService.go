package main

import (
  "encoding/json"
  "net/http"
  "log"
  "io"
  //"io/ioutil"
)

const APIBASEURL string = "https://slack.com/api/"
var APIENDPOINTS = map[string]string {
  "history": "channels.history",
  "user": "users.info",
}

type HistoryAPIResponse struct {
  Messages []Message
}

type Message struct {
  Text string
  User string
  Ts string
  username string
}

type UserAPIResponse struct {
  Id string
  User string
  Profile UserAPIResponseProfile
}

type UserAPIResponseProfile struct {
  First_name string
  Last_name string
}

func getJson(json_to_decode io.Reader, target interface{}) error {
  return json.NewDecoder(json_to_decode).Decode(&target)
}

func makeAPICall(endpoint string, params map[string]string, dest interface{}) error {
  url := APIBASEURL + endpoint + "?token=" + TOKEN
  if len(params) != 0 {
    for param, val := range params {
      url += "&" + param + "=" + val
    }
  }

  res, err := http.Get(url)
  defer res.Body.Close()
  if err != nil {
    return err
  }

  return getJson(res.Body, dest)
}

func getHistory(channel string) (*HistoryAPIResponse, error) {
  response := new(HistoryAPIResponse)
  err := makeAPICall(APIENDPOINTS["history"], map[string]string{"channel": channel, "count": "5"}, response)

  if err != nil {
    log.Print("Error: HTTP GET failed")
    log.Print(err)
    return nil, err
  }

  return response, nil
}

func getUserData(user_id string) (*UserAPIResponse, error) {
  response := new(UserAPIResponse)
  err := makeAPICall(APIENDPOINTS["user"], map[string]string{"user": user_id}, response)

  if err != nil {
    log.Print("Error: HTTP GET failed")
    log.Print(err)
    return nil, err
  }

  return response, nil
}
