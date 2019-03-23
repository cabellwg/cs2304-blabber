package main

import (
  "fmt"
  "time"
  "log"

  "crypto/sha256"

  "encoding/json"
  "net/http"

  "github.com/julienschmidt/httprouter"
)

// Types

type User struct {
  Id [32]byte
  Name string
  Email string
}

type Blab struct {
  Id [32]byte
  PostTime int64
  Author User
  Message string
}

// Entrypoint

func main() {
  DbConnect()

  router := httprouter.New()
  router.DELETE("/blabs/:id", RemoveBlab)
  router.GET("/blabs", GetBlabs)
  router.POST("/blabs", AddBlab)

  log.Fatal(http.ListenAndServe(":5000", nil))
}

// Functions

func RemoveBlab(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  rowsRemoved := DbRemove(ps.ByName("id"))

  if rowsRemoved == 0 {
    fmt.Fprintf(w, http.StatusText(404))
    return
  }

  fmt.Fprintf(w, http.StatusText(200))
}

func GetBlabs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  blabs := DbBlabs()

  b, err := json.Marshal(blabs)
  if err != nil {
    panic(err)
  }

  fmt.Fprintf(w, string(b))
}

func AddBlab(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  decoder := json.NewDecoder(r.Body)
  var blab Blab
  err := decoder.Decode(&blab)
  if err != nil {
      panic(err)
  }

  postTime := time.Now().Unix()
  blab.PostTime = postTime

  blab.Id = sha256.Sum256([]byte(fmt.Sprintf("%v", blab)))
  blab.Author.Id = sha256.Sum256([]byte(fmt.Sprintf("%v", blab.Author)))

  DbInsertBlab(blab)

  fmt.Fprintf(w, http.StatusText(201))
}
