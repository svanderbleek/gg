package main

import (
  "fmt"
  "log"
  "os"
  "net/http"
  "database/sql"
  "time"
  "math/rand"
  "strconv"
  "strings"
  _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func openDb(namespace string) *sql.DB {
  db, err := sql.Open("sqlite3", namespace + ".db")
  if err != nil {
    panic(err)
  }

  return db
}

func readSolution(id string) (int64, error) {
  var solution string

  stmt, err := db.Prepare(`select solution from game where game.id = ?`)
  if err != nil {
    return 0, err
  }

  err = stmt.QueryRow(id).Scan(&solution)
  if err != nil {
    return 0, err
  }

  s, err := strconv.ParseInt(solution, 10, 64)
  
  return s, err
}

func start(w http.ResponseWriter, r *http.Request) {
  solution := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100)

  stmt, err := db.Prepare(`insert into game(id, solution) VALUES (?, ?)`)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  res, err := stmt.Exec(nil, solution)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  id, err := res.LastInsertId()
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  log.Printf("started %d with solution %d", id, solution)

  fmt.Fprintf(w, "%d", id)
}

func ask(w http.ResponseWriter, r *http.Request) {
  id := r.FormValue("id")
  solution, err := readSolution(id)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  query := r.FormValue("query")

  log.Printf("asking %s for game %s with solution %d", query, id, solution)
  
  if query == "even" {
    fmt.Fprintf(w, "%t", solution % 2 == 0)
    return
  } else if query == "odd" {
    fmt.Fprintf(w, "%t", solution % 2 == 1)
    return
  } else if s := strings.TrimPrefix(query, "less"); len(s) < len(query) {
    n, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
    }

    fmt.Fprintf(w, "%t", solution < n)
    return
  } else if s := strings.TrimPrefix(query, "more"); len(s) < len(query) {
    n, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
    }

    fmt.Fprintf(w, "%t", solution > n)
    return
  }

  http.Error(w, err.Error(), http.StatusBadRequest)
}

func guess(w http.ResponseWriter, r *http.Request) {
  id := r.FormValue("id")

  solution, err := readSolution(id)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  guessed, err := strconv.ParseInt(r.FormValue("solution"), 10, 64)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  log.Printf("guessing %d for game %s with solution %d", guessed, id, solution)

  fmt.Fprintf(w, "%t", guessed == solution)
}

func SetupGuessingGame(namespace string) http.Handler {
  db = openDb(namespace)

  stmt, err := db.Prepare(`create table if not exists game(id integer not null primary key, solution integer);`)
  if err != nil {
    panic(err)
  }

  _, err = stmt.Exec()
  if err != nil {
    panic(err)
  }

  h := http.NewServeMux()

  h.HandleFunc("/start", start)
  h.HandleFunc("/ask/", ask)
  h.HandleFunc("/guess", guess)

  return h
}

func main() {
  handler := SetupGuessingGame(os.Args[1])

  log.Println("server listening")

  err := http.ListenAndServe(":8080", handler)
  if err != nil {
    panic(err)
  }
}
