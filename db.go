package main

import (
  "container/list"

  "fmt"
  "time"
  "database/sql"

  "github.com/lib/pq"
)


// Globals

const (
  host = "blab_db"
  port = 5432
  user = "blabclient"
  password = "r$J89ka&36"
  dbname = "blabs"
)

// TODO: implement with method handlers instead
var db *sql.DB


// Functions

func DbConnect()  {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

  var err error
  db, err = sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
}


func DbInsertUser(user User)  {
  addUserStatement := `INSERT INTO users (id, name, email) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`

  db.QueryRow(addUserStatement, user.Id, user.Name, user.Email)
}


func DbInsertBlab(blab Blab)  {
  addUserStatement := `INSERT INTO users (id, name, email) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
  addBlabStatement := `INSERT INTO blabs (id, postTime, author, message) VALUES ($1, $2, $3, $4)`

  _, err := db.Exec(addUserStatement,
    blab.Author.Id,
    blab.Author.Name,
    blab.Author.Email)
  if err != nil {
    panic(err)
  }

  _, err = db.Exec(addBlabStatement,
    blab.Id,
    pq.FormatTimestamp(blab.PostTime),
    blab.Author.Id,
    blab.Message)
  if err != nil {
    panic(err)
  }
}


func DbBlabs(createdSince time.Time) *list.List  {
  blabs := list.New()

  queryStatement := "SELECT blabs.id, blabs.postTime, users.name, users.email, blabs.message FROM blabs LEFT JOIN users ON blabs.author = users.id WHERE blabs.postTime > $1"

  rows, err := db.Query(queryStatement, pq.FormatTimestamp(createdSince))
  if err != nil {
    fmt.Printf("Get query failed")
    panic(err)
  }

  defer rows.Close()
  for rows.Next() {
    var blab Blab
    var author User
    var pgTimeString string
    err = rows.Scan(&blab.Id,
      &pgTimeString,
      &author.Name,
      &author.Email,
      &blab.Message)
    if err != nil {
      fmt.Printf("Could not parse row into blab")
      panic(err)
    }

    fmt.Printf(pgTimeString)
    blab.PostTime, err = pq.ParseTimestamp(time.FixedZone("UTC-8", 0), pgTimeString)
    if err != nil {
      fmt.Printf("Could not parse timestamp")
      panic(err)
    }

    blab.Author = author

    blabs.PushBack(blab)
  }
  // get any error encountered during iteration
  err = rows.Err()
  if err != nil {
    panic(err)
  }

  return blabs
}


func DbRemove(id string) int {
  removeStatement := `DELETE FROM blabs WHERE id=$1`
  res, err := db.Exec(removeStatement, id)
  if err != nil {
    panic(err)
  }

  rowsRemoved, err := res.RowsAffected()
  if err != nil {
    panic(err)
  }

  return int(rowsRemoved)
}
