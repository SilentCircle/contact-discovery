package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/docopt/docopt-go"
	"github.com/gorilla/mux"
)

type HashRequest struct {
	Prefixes []string
}

type JSONResponse map[string]interface{}

func (r JSONResponse) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

type Hashes struct {
	HashList map[string][]string
}

func getContactsView(w http.ResponseWriter, r *http.Request) {
	var request HashRequest

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Fprint(w, JSONResponse{"result": "error", "error": err.Error()})
		return
	}

	hashes := make([]string, 0)
	for _, prefix := range request.Prefixes {
		for _, hash := range getHashesForPrefix(prefix) {
			hashes = append(hashes, hash)
		}
	}

	fmt.Fprint(w, JSONResponse{"result": "success", "hashes": hashes})
	return
}

func addHashView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	hash := vars["hash"]

	insertHash(hash)

	fmt.Fprint(w, JSONResponse{"result": "success"})
	return
}

func deleteHashView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	hash := vars["hash"]

	deleteHash(hash)

	fmt.Fprint(w, JSONResponse{"result": "success"})
	return
}

func main() {
	usage := `contact-discovery

Usage:
  contact-discovery [options]

Options:
  -d --database=<filename>  The filename of the database file [default: contacts.sqlite3]
  -m --prefix-length=<num>  The minimum prefix length to accept [default: 4]
  -s --hash-length=<num>    The length of the hash to return [default: 20]
  -p --port=<port>          The port to listen to [default: 8080]
  -h --help                 Show this screen
  --version                 Show version`

	arguments, _ := docopt.Parse(usage, nil, true, "Contact Discovery Server", false)

	initDatabase()
	port, err := strconv.Atoi(arguments["--port"].(string))
	if err != nil {
		log.Fatal("Invalid port.")
	}

	DB_FILENAME = arguments["--database"].(string)

	prefixLength, err := strconv.Atoi(arguments["--prefix-length"].(string))
	if err != nil {
		log.Fatal("Invalid prefix length.")
	}
	MIN_PREFIX_LENGTH = prefixLength

	hashLength, err := strconv.Atoi(arguments["--hash-length"].(string))
	if err != nil {
		log.Fatal("Invalid hash length.")
	}
	MAX_HASH_LENGTH = hashLength

	fmt.Printf("Server starting up on port [0;32m%v[0m...\n", port)

	r := mux.NewRouter()
	r.HandleFunc("/contacts/", getContactsView).Methods("POST")
	r.HandleFunc("/hashes/{hash:[0-9a-f]+}/", addHashView).Methods("POST")
	r.HandleFunc("/hashes/{hash:[0-9a-f]+}/", deleteHashView).Methods("DELETE")
	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
