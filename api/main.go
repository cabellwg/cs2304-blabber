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
	router.GET("/blabs", getBlabs)
	router.POST("/blabs", addBlab)

	log.Fatal(http.ListenAndServe(":80", router))
}

// Functions

func helloWorld(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte("Hello, world!"))
}

func removeBlab(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid ID"))
		return
	}

	rowsRemoved := DbRemove(uint32(id))

	if rowsRemoved == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Blab not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Blab deleted successfully"))
}

func getBlabs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	for k, v := range r.URL.Query() {
		fmt.Println("k:", k, "v:", v)
	}
	keys, ok := r.URL.Query()["createdSince"]

	if !ok {
		keys = []string{"0"}
	}

	log.Println(keys)

	sinceInt, err := strconv.ParseInt(keys[0], 10, 64)
	if err != nil {
		log.Printf("Unrecognized timestamp: %s\n", keys[0])
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}
	since := time.Unix(sinceInt, 0)
	blabs := DbBlabs(since)

	b, err := json.Marshal(blabs)
	if err != nil {
		log.Println("Could not parse blab into json")
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func addBlab(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println(ps)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}
