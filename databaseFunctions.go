package main

import (
  "fmt"
  "log"
  "sync"
  "encoding/json"
  "time"

  "github.com/gocql/gocql"
);

type CassandraDB struct{
  once sync.Once
  session *gocql.Session
}

func (cdb *CassandraDB) Init() error{
  cdb.once.Do(func() {
    cluster := gocql.NewCluster("127.0.0.1")
    cluster.Keyspace = "remindmeslackintegration"
    cluster.Consistency = gocql.Quorum
    cluster.ProtoVersion = 3
    cdb.session, _ = cluster.CreateSession()
  })
  if nil == cdb.session {
    err := fmt.Errorf("The database session was not created.")
    return err
  }

  return nil
}

func (cdb *CassandraDB) Close() {
  cdb.session.Close()
}

func (cdb *CassandraDB) InsertRequest(info *Reminder) error{
  date := info.reminderTime.YearDay()
  context, err := json.Marshal(info.context)
  if nil != err {
    log.Println("Insert Request Json Marshal Error")
    log.Println(err)
    return err
  }
  err = cdb.session.Query(
    `INSERT INTO Messages (alertTime, messages, date, user_name) VALUES (?, ?, ?, ?)`,
    info.reminderTime, context, date, info.user_name,
  ).Exec()
  if nil != err {
    log.Println("Insert Request Error:")
    log.Println(err)
  }
  return err
}

func (cdb *CassandraDB) ReadNextRequest() (*Reminder, error){
  reminder := new(Reminder)
  var messages []byte
  err := cdb.session.Query(
    `SELECT alertTime, messages, user_name FROM Messages WHERE date = ?`,
    time.Now().YearDay(),
  ).Consistency(gocql.One).Scan(
                                &reminder.reminderTime,
                                &messages,
                                &reminder.user_name,
                               )
  if nil != err {
    log.Println("Read Request Error")
    log.Println(err)
    return nil, err
  }
  err = json.Unmarshal(messages, &reminder.context)
  if nil != err {
    log.Println("Read Request JSON Unmarshal Error")
    log.Println(err)
    return nil, err
  }
  return reminder, err
}

