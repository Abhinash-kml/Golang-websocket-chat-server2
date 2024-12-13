package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	DB, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost:5432?sslmode=disable")
	if err != nil {
		log.Fatal("Error opening postgress connection.\n")
		return
	}
	defer DB.Close()

	CreateTable(DB)
	InsertData(DB)
}

func CreateTable(DB *sql.DB) {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS channels(
	id SERIAL,
	name VARCHAR(64));`)

	if err != nil {
		DB.Close()
		log.Fatal("Error creating table channels.\n")
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS messages(
	id SERIAL,
	cid int,
	cname VARCHAR(128),
	sid int,
	sname VARCHAR(128),
	message VARCHAR(128),
	FOREIGN KEY(cid) REFERENCES channels(id),
	FOREIGN KEY(sid) REFERENCES users(id),
	FOREIGN KEY(cname) REFERENCES channels(name),
	FOREIGN KEY(sname) REFERENCES users(name));`)

	if err != nil {
		DB.Close()
		log.Fatal("Error creating table messages.\n")
	}
}

func InsertData(DB *sql.DB) {
	_, err := DB.Exec(`INSERT INTO channels(name) VALUES
	('general'),
	('english'),
	('hindi'),
	('bakchodi');`)

	if err != nil {
		fmt.Println("Error inserting rows into channels table")
		return
	}
}

func GetAllMessagesOfChannel(DB *sql.DB, channel string) []string {
	var messages []string

	// Check if the required channel exists in the table
	_, err := DB.Query(`SELECT COUNT(*) FROM channels WHERE name = $1;`, channel)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
	}

	// Get the rows
	rows, err := DB.Query(`SELECT sname, message FROM messages WHERE cname = $1;`, channel)
	if err != nil {
		log.Fatal("Error getting messages of channel ", channel)
		return nil
	}

	// Scan and append them to the returning slice
	for rows.Next() {
		var message, sname, smessage string
		if err = rows.Scan(&sname, &smessage); err != nil {
			fmt.Printf("Error scanning current row from rows.\n")
			continue
		}

		message = sname + ":" + smessage
		messages = append(messages, message)
	}

	return messages
}

func InsertMessageIntoChannel(DB *sql.DB, channel, sendername, message string) bool {
	_, err := DB.Exec("INSERT INTO messages(cname, sname, message) VALUES($1, $2, $3);", channel, sendername, message)
	if err != nil {
		fmt.Println("Error interting new message into channel ", channel)
		return false
	}

	return true
}
