package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	db, err := sql.Open("postgres", "postgresql://postgres:Abx305@localhost:5432?sslmode=disable")
	if err != nil {
		log.Fatal("Error opening postgress connection.\n")
		return
	}
	DB = db
	//defer DB.Close()

	CreateTable()
	InsertData()
}

func CreateTable() {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS channels(
	id SERIAL,
	name VARCHAR(64));`)

	if err != nil {
		DB.Close()
		fmt.Println(err)
		log.Fatal("Error creating table channels.\n")
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS messages(
	id SERIAL, 
	cname VARCHAR(128),
	sname VARCHAR(128),
	message VARCHAR(128));`)

	if err != nil {
		DB.Close()
		fmt.Println(err)
		log.Fatal("Error creating table messages.\n")
	}
}

func InsertData() {
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

func GetAllMessagesOfChannel(channel string) []string {
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

func InsertMessageIntoChannel(channel, sendername, message string) bool {
	result, err := DB.Exec("INSERT INTO messages(cname, sname, message) VALUES($1, $2, $3);", channel, sendername, message)
	if err != nil {
		fmt.Println("Error interting new message into channel ", channel)
		fmt.Println(err)
		return false
	}

	rowsEffected, err := result.RowsAffected()
	fmt.Println("Message added. Rows effected: ", rowsEffected)
	return true
}
