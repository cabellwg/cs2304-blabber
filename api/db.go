package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"io/ioutil"
	"path/filepath"

	"github.com/lib/pq"
)

// Types

type BlabDb struct {
	*sql.DB
}

// Globals

const (
	host     = "blab_db"
	port     = 5432
	user     = "blabclient"
	dbname   = "blabdb"
)

// Functions

// Connect connects to the PostgreSQL database
func (db *BlabDb) Connect() {
	path, err := filepath.Abs("/run/secrets/blabber-db-password")
	if (err != nil) {
		panic(err)
	}
	data, err := ioutil.ReadFile(path)
	if (err != nil) {
		panic(err)
	}

	password := string(data)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db.DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
}


// Connected returns true if there is an open database connection
func (db *BlabDb) Connected() bool {
	return db.Ping() == nil
}


// InsertUser inserts a new user into the database
func (db *BlabDb) InsertUser(user User) {
	addUserStatement := `INSERT INTO users (id, name, email) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`

	db.QueryRow(addUserStatement, user.ID, user.Name, user.Email)
}

// InsertBlab adds a blab to the database
func (db *BlabDb) InsertBlab(blab Blab) {
	db.InsertUser(blab.Author)

	addBlabStatement := `INSERT INTO blabs (id, postTime, author, message) VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(addBlabStatement,
		blab.ID,
		pq.FormatTimestamp(blab.PostTime),
		blab.Author.ID,
		blab.Message)
	if err != nil {
		panic(err)
	}
}

// Blabs returns all blabs in the database created at or after the given time
func (db *BlabDb) Blabs(createdSince time.Time) []Blab {
	blabs := make([]Blab, 0)

	goDateFormat := "2006-01-02T15:04:05Z"
	postgresDateFormat := "yyyy-mm-ddThh:mm:ssZ"
	queryStatement := `SELECT blabs.id, blabs.postTime, TO_CHAR(blabs.postTime :: timestamp, $1), users.name, users.email, blabs.message FROM blabs LEFT JOIN users ON blabs.author = users.id WHERE blabs.postTime >= $2`

	rows, err := db.Query(queryStatement, postgresDateFormat, pq.FormatTimestamp(createdSince))
	if err != nil {
		log.Println("Get query failed")
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
		var blab Blab
		var author User
		var pgTimeString string
		var timestamp string
		err = rows.Scan(&blab.ID,
			&pgTimeString,
			&timestamp,
			&author.Name,
			&author.Email,
			&blab.Message)

		if err != nil {
			log.Println("Could not parse row into blab")
			panic(err)
		}

		blab.PostTime, err = time.Parse(goDateFormat, timestamp)

		if err != nil {
			panic(err)
		}

		blab.Author = author

		blabs = append(blabs, blab)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return blabs
}

// Remove removes a blab from the database by its ID
func (db *BlabDb) Remove(id uint32) int {
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
