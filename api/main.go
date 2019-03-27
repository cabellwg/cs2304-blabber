package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"crypto/sha256"
	"encoding/binary"

	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Entrypoint

func main() {
	time.Sleep(200000000) // 2s, for db startup

	DbConnect()

	router := httprouter.New()
	router.GET("/", helloWorld)
	router.DELETE("/blabs/:id", removeBlab)
	router.GET("/blabs", getBlab)
	router.POST("/blabs", addBlab)

	log.Fatal(http.ListenAndServe(":80", router))
}

// Functions

func helloWorld(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Hello, world!\n")
}

func removeBlab(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	rowsRemoved := DbRemove(ps.ByName("id"))

	if rowsRemoved == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Blab not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Blab deleted successfully")
}

func getBlab(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		log.Println("Could not parse blab into json")
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(b))
}

func addBlab(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var blab Blab
	err := decoder.Decode(&blab)
	if err != nil {
		panic(err)
	}

	blab.PostTime = time.Now()

	blabHash := sha256.Sum256([]byte(fmt.Sprintf("%v", blab)))
	blab.ID = binary.BigEndian.Uint32(blabHash[:]) >> 1

	authorHash := sha256.Sum256([]byte(fmt.Sprintf("%v", blab.Author)))
	blab.Author.ID = binary.BigEndian.Uint32(authorHash[:]) >> 1

	DbInsertBlab(blab)

	b, err := json.Marshal(blab)
	if err != nil {
		log.Println("Could not parse blab into json")
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(b))
}
