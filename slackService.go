package main

import (
  "encoding/json"
  "net/http"
  "net/url"
  "log"
  "io"
  //"io/ioutil"
)

const APIBASEURL string = "https://slack.com/api/"
var APIENDPOINTS = map[string]string {
  "history": "channels.history",
  "user": "users.info",
  "message": "chat.postMessage",
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
  User UserAPIResponseUser
}

type UserAPIResponseUser struct {
  Id string
  Name string
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
  params["token"] = TOKEN
  var Url *url.URL
  Url, err := url.Parse(APIBASEURL)
  if err != nil {
    log.Print("Base API Parse failure")
    return err
  }
  Url.Path += endpoint
  parameters := url.Values{}
  for param, val := range params {
    parameters.Add(param, val)
  }

  Url.RawQuery = parameters.Encode()

  log.Print(Url.String())
  res, err := http.Get(Url.String())
  defer res.Body.Close()
  if err != nil {
    return err
  }

  if dest == nil {
    return nil
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

func sendMessage(channel string, text string) error {
  err := makeAPICall(APIENDPOINTS["message"], map[string]string{"username": "RemindMe Bot", "channel": channel, "text": text, "as_user": "False", "icon_emoji": ":sums:"}, nil)

  if err != nil {
    log.Print("Error: Sending messaage failed")
    log.Print(err)
    return err
  }
  return nil
}
