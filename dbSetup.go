package main

import (
  "log"
  "os"
  "bufio"

  "github.com/gocql/gocql"
)

func setupDB(filename string){
  log.Println(os.Args[1])
  f, err := os.Open(filename)
  if nil != err {
    panic(err)
  }
  reader := bufio.NewScanner(f)
  cluster := gocql.NewCluster("127.0.0.1")
  cluster.ProtoVersion = 3
  session, _ := cluster.CreateSession()
  defer session.Close()
  for reader.Scan() {
    log.Println(reader.Text())
    err := session.Query(reader.Text()).Exec()
    if nil != err {
      panic(err)
    }
  }
}
