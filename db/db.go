package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // This is required by sqlx to be able to connect to postgres

	"fmt"
	"io/ioutil"
	"os"
)

// QuoteSave is a struct used to insert a row into the 'quotes' table
type QuoteSave struct {
	GuildID   string `db:"guild_id"`
	ChannelID string `db:"channel_id"`
	MessageID string `db:"message_id"`
	Content   string `db:"content"`
	Author    string `db:"author"`
	AddedBy   string `db:"added_by"`
	AddedAt   uint64 `db:"added_at"`
}

// QuotePreview is a struct used to hold a preview of a row from the 'quotes' table
type QuotePreview struct {
	Content string `db:"content"`
	Author  string `db:"author"`
}

// Quote is a struct used to hold a row from the 'quotes' table
type Quote struct {
	ID        uint64 `db:"id"`
	GuildID   string `db:"guild_id"`
	ChannelID string `db:"channel_id"`
	MessageID string `db:"message_id"`
	Content   string `db:"content"`
	Author    string `db:"author"`
	AddedBy   string `db:"added_by"`
	AddedAt   uint64 `db:"added_at"`
}

// DB is a Reusable variable for db connection
var DB *sqlx.DB

// Generates connection string using environment variables
func connectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
	)
}

// DBInit initialises the database connection
func DBInit() {
	var err error
	// Gen connection string
	connStr := connectionString()
	// Connect to postgres db using connection string
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		panic(err)
	}
	// Load schema from file
	schema, err := ioutil.ReadFile("./schema.sql")
	if err != nil {
		panic(err)
	}
	// Execute schema (panics on failure)
	DB.MustExec(string(schema))
}
