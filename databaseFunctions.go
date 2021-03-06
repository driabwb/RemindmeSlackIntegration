package main

import (
  "fmt"
  "log"
  "sync"
  "encoding/json"
  "time"

  "github.com/gocql/gocql"
);

type ReminderDB interface {
  Init() error
  Close()
  InsertRequest(*Reminder) error
  ReadNextRequest() (*Reminder, error)
  DeleteRequest(time.Time) error
}

type UserDB interface {
  InsertUserData(User) error
  GetUser(string) (*User, error)
}

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
  date := info.ReminderTime.YearDay()
  Context, err := json.Marshal(info.Context)
  if nil != err {
    log.Println("Insert Request Json Marshal Error")
    log.Println(err)
    return err
  }
  err = cdb.session.Query(
    `INSERT INTO Messages (alertTime, messages, date, user_name) VALUES (?, ?, ?, ?)`,
    info.ReminderTime, Context, date, info.User_name,
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
                                &reminder.ReminderTime,
                                &messages,
                                &reminder.User_name,
                               )
  if nil != err {
    log.Println("Read Request Error")
    log.Println(err)
    return nil, err
  }
  err = json.Unmarshal(messages, &reminder.Context)
  if nil != err {
    log.Println("Read Request JSON Unmarshal Error")
    log.Println(err)
    return nil, err
  }
  return reminder, err
}

func (cdb *CassandraDB) DeleteRequest(timestamp time.Time) error{
  date := timestamp.YearDay()
  err := cdb.session.Query(
    `DELETE FROM Messages WHERE date = ? AND alertTime = ?`,
    date, timestamp,
  ).Exec()
  if nil != err {
    log.Println("Delete Request Error")
    log.Println(err)
    return err
  }
  return err
}

func (cdb *CassandraDB) InsertUserData(user User) error{
  err := cdb.session.Query(
    `INSERT INTO Users (id, firstname, lastname, name) VALUES (?, ?, ?, ?)`,
    user.Id, user.FirstName, user.LastName, user.User_name,
  ).Exec()
  if nil != err {
    log.Println("Insert User Data")
    log.Println(err)
    return err
  }
  return err
}

func (cdb *CassandraDB) GetUser(id string) (*User, error){
  user := new(User)
  err := cdb.session.Query(
    `SELECT id, firstname, lastname, name FROM Users WHERE id = ?`,
    id,
  ).Consistency(gocql.One).Scan(
                                &user.Id,
                                &user.FirstName,
                                &user.LastName,
                                &user.User_name,
                              )
  if nil != err {
    log.Println("Get User Error")
    log.Println(err)
    return nil, err
  }
  return user, err
}

