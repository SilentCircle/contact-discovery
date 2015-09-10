package main

import (
	"database/sql"
	"log"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Ignore the specified prefix if it's less than this.
var MIN_PREFIX_LENGTH int = 4

// Truncated the returned hashes to this length.
var MAX_HASH_LENGTH int = 20

var DB_FILENAME string = "./contacts.sqlite3"

func initDatabase() {
	db, err := sql.Open("sqlite3", DB_FILENAME)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "main"."hashes" ("hash" CHAR PRIMARY KEY  NOT NULL  UNIQUE);`)
	if err != nil {
		log.Fatal("Could not create table.")
	}

	_, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS "main"."hashes_prefix" ON "hashes" ("hash" ASC);`)
	if err != nil {
		log.Fatal("Could not create index.")
	}
}

func insertHash(hash string) error {
	db, err := sql.Open("sqlite3", DB_FILENAME)
	if err != nil {
		return err
	}
	defer db.Close()

	hash = strings.TrimSpace(strings.ToLower(hash))

	_, err = db.Exec("INSERT INTO hashes (hash) VALUES (?);", hash)
	if err != nil {
		return err
	}
	return nil
}

func deleteAllHashes(hash string) error {
	db, err := sql.Open("sqlite3", DB_FILENAME)
	if err != nil {
		return err
	}
	defer db.Close()

	hash = strings.TrimSpace(strings.ToLower(hash))

	_, err = db.Exec("DELETE FROM hashes;")
	if err != nil {
		return err
	}
	return nil
}

func deleteHash(hash string) error {
	db, err := sql.Open("sqlite3", DB_FILENAME)
	if err != nil {
		return err
	}
	defer db.Close()

	hash = strings.TrimSpace(strings.ToLower(hash))

	_, err = db.Exec("DELETE FROM hashes WHERE hash = ?;", hash)
	if err != nil {
		return err
	}
	return nil
}

func getHashesForPrefix(prefix string) []string {
	// If the prefix is shorted than the minimum length, return.
	if len(prefix) < MIN_PREFIX_LENGTH {
		return []string{}
	}

	prefix = strings.TrimSpace(strings.ToLower(prefix))

	// If there are any characters other than hex digits, return.
	if matched, _ := regexp.MatchString("^[a-f0-9]+$", prefix); !matched {
		return []string{}
	}

	db, err := sql.Open("sqlite3", DB_FILENAME)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT hash FROM hashes WHERE hash LIKE ?", prefix+"%")
	if err != nil {
		log.Fatal(err)
	}

	results := make([]string, 0)

	defer rows.Close()
	for rows.Next() {
		var hash string
		rows.Scan(&hash)
		results = append(results, hash[:MAX_HASH_LENGTH])
	}

	return results
}
