package main

import (
  "fmt"
  "time"
  "log"
  "strconv"

  "crypto/sha256"
  "encoding/binary"

  "encoding/json"
  "net/http"

  "github.com/julienschmidt/httprouter"
)

// Types

type User struct {
  Id uint32
  Name string
  Email string
}

type Blab struct {
  Id uint32
  PostTime time.Time
  Author User
  Message string
}

// Entrypoint

func main() {
  time.Sleep(200000000) // 2s, for db startup

  DbConnect()

  router := httprouter.New()
  router.GET("/", HelloWorld)
  router.DELETE("/blabs/:id", RemoveBlab)
  router.GET("/blabs", GetBlabs)
  router.POST("/blabs", AddBlab)

  log.Fatal(http.ListenAndServe(":80", router))
}

// Functions

func HelloWorld(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
  fmt.Fprintf(w, "Hello, world!\n")
}

func RemoveBlab(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  rowsRemoved := DbRemove(ps.ByName("id"))

  if rowsRemoved == 0 {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "Blab not found")
    return
  }

  w.WriteHeader(http.StatusOK)
  fmt.Fprintf(w, "Blab deleted successfully")
}

func GetBlabs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  keys, ok := r.URL.Query()["createdSince"]

  if !ok {
    keys = []string{"0"}
  }

  sinceInt, err := strconv.ParseInt(keys[0], 10, 64)
  if err != nil {
    panic(err)
  }
  since := time.Unix(sinceInt, 0)
  blabs := DbBlabs(since)

  b, err := json.Marshal(blabs)
  if err != nil {
    log.Println("Could not parse blab into json\n")
    panic(err)
  }

  w.WriteHeader(http.StatusOK)
  fmt.Fprintf(w, string(b))
}

func AddBlab(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  decoder := json.NewDecoder(r.Body)
  var blab Blab
  err := decoder.Decode(&blab)
  if err != nil {
      panic(err)
  }

  blab.PostTime = time.Now()

  blabHash := sha256.Sum256([]byte(fmt.Sprintf("%v", blab)))
  blab.Id = binary.BigEndian.Uint32(blabHash[:]) >> 1

  authorHash := sha256.Sum256([]byte(fmt.Sprintf("%v", blab.Author)))
  blab.Author.Id = binary.BigEndian.Uint32(authorHash[:]) >> 1

  DbInsertBlab(blab)

  b, err := json.Marshal(blab)
  if err != nil {
    log.Println("Could not parse blab into json\n")
    panic(err)
  }

  w.WriteHeader(http.StatusCreated)
  fmt.Fprintf(w, string(b))
}
