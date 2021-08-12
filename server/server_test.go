package main

import (
  "net/http"
  "testing"
  "net/http/httptest"
  "io"
  "strconv"
  "os"
)

var namespace = "test";

func TestMain(m *testing.M) {
  code := m.Run()
  err := os.Remove(namespace + ".db")
  if err != nil {
    os.Exit(1)
  } else {
    os.Exit(code)
  } 
}

func TestStart(t *testing.T) {
  request, err := http.NewRequest("POST", "/start", nil)
  if(err != nil) {
    t.Error(err)
  }

  recorder := httptest.NewRecorder()

  handler := SetupGuessingGame(namespace)
  handler.ServeHTTP(recorder, request)

  response := recorder.Result()

  body, err := io.ReadAll(response.Body)
  if(err != nil) {
    t.Error(err)
  }

  id, err := strconv.ParseInt(string(body), 10, 64)
  if(err != nil) {
    t.Error(err)
  } 

  t.Log(id)
}
