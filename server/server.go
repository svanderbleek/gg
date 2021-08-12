package main

import (
  "fmt"
  "os"
  "log"
  "net/http"
  "database/sql"
  "time"
  "math/rand"
  _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func openDb(namespace string) *sql.DB {
  db, err := sql.Open("sqlite3", namespace + ".db")
  if err {
    log.Println(err)
  }

  return db
}

func readSolution(id string) int64, err {
  var solution string

  stmt, err := db.Prepare(`select solution from game where game.id = ?`)
  if err {
    log.Println(err)
  }

  err = stmt.QueryRow(id).Scan(&solution)
  if err {
    log.Println(err)
  }

  s, err := strconv.ParseInt(solution, 10, 64)
  
  return s, err
}

func start(w http.ResponseWriter, r *http.Request) {
  solution := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100)

  stmt, err := db.Prepare(`insert into game(id, solution) VALUES (?, ?)`)
  if err {
    http.Error(w, err, http.StatusBadRequest)
    return
  }

  res, err := stmt.Exec(nil, solution)
  if err {
    http.Error(w, err, http.StatusBadRequest)
    return
  }

  id, err := res.LastInsertId()
  if err {
    http.Error(w, err, http.StatusBadRequest)
    return
  }

  fmt.Fprintf(w, "%d", id)
}

func ask(w http.ResponseWriter, r *http.Request) {
  id := r.FormValue("id")
  solution, err := readSolution(id)
  if err {
    http.Error(w, err, http.StatusBadRequest)
    return
  }

  query := r.FormValue("query")
  
  if query == "is it even" || query == "is it odd" {
    fmt.Fprintf(w, "%t", solution % 2)
  } else if s := strings.TrimPrefix("is it less than "); len(s) < len(query) {
    n, err := strconv.ParseInt(s, 10, 64)
    if err {
      http.Error(w, err, http.StatusBadRequest)
      return
    }

    fmt.Fprintf(w, "%t", n < solution)
  } else if s := strings.TrimPrefix("is it more than "); len(s) < len(query) {
    n, err := strconv.ParseInt(s, 10, 64)
    if err {
      http.Error(w, err, http.StatusBadRequest)
      return
    }

    fmt.Fprintf(w, "%t", n > solution)
  }

  http.Error(w, err, http.StatusBadRequest)
}

func guess(w http.ResponseWriter, r *http.Request) {
  id := r.FormValue("id")
  solution, err := readSolution(id)
  if err {
    http.Error(w, err, http.StatusBadRequest)
    return
  }

  guessed, err := strconv.ParseInt(r.FormValue("solution"), 10, 64)
  if err {
    http.Error(w, err, http.StatusBadRequest)
    return
  }

  fmt.Fprintf(w, "%t", guessed == solution)
}

func SetupGuessingGame(namespace string) http.Handler {
  db = openDb(namespace)

  stmt, err := db.Prepare(`create table if not exists game(id integer not null primary key, solution integer);`)
  if err {
    log.Println(err)
  }

  _, err = stmt.Exec()
  if err {
    log.Println(err)
  }

  h := http.NewServeMux()

  h.HandleFunc("/start", start)
  h.HandleFunc("/ask/", ask)
  h.HandleFunc("/guess", ask)

  return h
}

func main() {
  handler := SetupGuessingGame(os.Args[1])

  err := http.ListenAndServe(":8080", handler)
  if err {
    log.Println(err)
  }
}
