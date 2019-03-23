package main

import (
  "container/list"

  "fmt"
  "database/sql"

  _ "github.com/lib/pq"
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
  addUserStatement := `INSERT INTO users (id, name, email) VALUES ($1, $2, $3) ON CONFLICT NO ACTION`

  db.QueryRow(addUserStatement, user.Id, user.Name, user.Email)
}


func DbInsertBlab(blab Blab)  {
  addUserStatement := `INSERT INTO users (id, name, email) VALUES ($1, $2, $3) ON CONFLICT NO ACTION`
  addBlabStatement := `INSERT INTO blabs (id, postTime, author, message) VALUES ($1, $2, $3, $4)`

  db.QueryRow(addUserStatement, blab.Author.Id, blab.Author.Name, blab.Author.Email)
  db.QueryRow(addBlabStatement, blab.Id, blab.PostTime, blab.Author.Id, blab.Message)
}


func DbBlabs() *list.List  {
  blabs := list.New()

  queryStatement := "SELECT blabs.id, blabs.postTime, users.name, users.email, blabs.message FROM blabs LEFT JOIN users ON blabs.author = users.id"

  rows, err := db.Query(queryStatement)
  if err != nil {
    // TODO: handle this error better than this
    panic(err)
  }

  defer rows.Close()
  for rows.Next() {
    var blab Blab
    var author User
    err = rows.Scan(&blab.Id,
      &blab.PostTime,
      &author.Name,
      &author.Email,
      &blab.Message)

    if err != nil {
      // TODO: handle this error
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


func DbRemove(id string) int64 {
  removeStatement := `DELETE FROM blabs WHERE id=$1`
  res, err := db.Exec(removeStatement, id)
  if err != nil {
    panic(err)
  }

  rowsRemoved, err := res.RowsAffected()
  if err != nil {
    panic(err)
  }

  return rowsRemoved
}
