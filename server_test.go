package main

import (
  "testing"
  "httptest"
  "http"
  "io"
  "strconv"
)

func TestStart(t *testing.T) {
  request, err := http.NewRequest("POST", "/start", nil);
  recorder := httptest.NewRecorder();

  handler := SetupGuessingGame();
  handler.ServeHTTP(recorder, request);

  response := recorder.Result();

  body, err := io.ReadAll(response.Body);
  id, err := strconv.ParseInt(body);

  if(err != nil) {
    t.Errorf(err);
  }
}
