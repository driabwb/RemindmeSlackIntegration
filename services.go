package main

import (
  "os"
  "bufio"
  "strings"
)

func Readln(r *bufio.Reader) (string, error) {
  var (isPrefix bool = true
       err error = nil
       line, ln []byte
      )
  for isPrefix && err == nil {
      line, isPrefix, err = r.ReadLine()
      ln = append(ln, line...)
  }
  return string(ln),err
}

func setupEnvironment() {
  f, err := os.Open(".env")
  if err != nil {
      panic(err)
  }
  reader := bufio.NewReader(f)
  s, e := Readln(reader)
  for e == nil {
    pair := strings.Split(s, "=")
    os.Setenv(pair[0], pair[1])

    s,e = Readln(reader)
  }
}

func init() {
  setupEnvironment()
}
