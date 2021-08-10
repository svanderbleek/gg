package main

import (
  "fmt"
  "log"
  "net/http"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

func OpenDb() *sql.DB {
  db, err := sql.Open("sqlite3", "./gg.db")
  log.Println(err);

  return db;
}


func start(w http.ResponseWriter, r *http.Request) {
  db := OpenDb();

  guess := 42;

  stmt, err := db.Prepare(`insert into game(id, guess) VALUES (?, ?)`)
  log.Println(err);

  res, err := stmt.Exec(nil, guess)
  log.Println(err); 

  id, err := res.LastInsertId();
  log.Println(err);

  fmt.Fprintf(w, "%d", id);
}

func ask(w http.ResponseWriter, r *http.Request) {
  query := r.FormValue("query");
  log.Println(query); 
}

func guess(w http.ResponseWriter, r *http.Request) {
  guess := r.FormValue("guess");
  log.Println(guess);
}

func SetupGuessingGame() http.Handler {
  db := OpenDb();

  stmt, err := db.Prepare(`create table if not exists game(id integer not null primary key, guess integer);`);;
  log.Println(err);

  _, err = stmt.Exec();
  log.Println(err);

  h := http.NewServeMux();

  h.HandleFunc("/start", start)
  h.HandleFunc("/ask/", ask)
  h.HandleFunc("/guess", ask)

  return h;
}

func main() {
  handler := SetupGuessingGame();

  err := http.ListenAndServe(":8080", handler);
  log.Println(err)
}
