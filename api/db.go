package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

// Globals

const (
	host     = "blab_db"
	port     = 5432
	user     = "blabclient"
	password = "r$J89ka&36"
	dbname   = "blabdb"
)

// TODO: implement with method handlers instead
var db *sql.DB

// Functions

// DbConnect connects to the PostgreSQL database
func DbConnect() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
}

// DbInsertUser inserts a new user into the database
func DbInsertUser(user User) {
	addUserStatement := `INSERT INTO users (id, name, email) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`

	db.QueryRow(addUserStatement, user.ID, user.Name, user.Email)
}

// DbInsertBlab adds a blab to the database
func DbInsertBlab(blab Blab) {
	DbInsertUser(blab.Author)

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

// DbBlabs returns all blabs in the database created at or after the given time
func DbBlabs(createdSince time.Time) []Blab {
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

// DbRemove removes a blab from the database by its ID
func DbRemove(id uint32) int {
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
