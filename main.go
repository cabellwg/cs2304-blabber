package main

import (
  "fmt"
  "html"
  "log"
  "net/http"
  "database/sql"

  _ "github.com/lib/pq"
)

const (
  host = "blab_db"
  port = 5432
  user = "blabclient"
  password = "r$J89ka&36"
  dbname = "blabs"
)

func main() {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

  db, err := sql.Open("postgres", psqlInfo)

  if err != nil {
    panic(err)
  }
}
