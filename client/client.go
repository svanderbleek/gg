package main

import (
  "fmt"
  "log"
  "net/http"
  "strconv"
  "github.com/manifoldco/promptui"
)

func Game(client Client, actions, ask, guess) {
  client.Start()

  _, action, err := actions.Run()
  if err {
    log.Println(err)
  }

  if(action == "Ask") {
    query, err := ask.Run()
    if err {
      log.Println(err)
    }

    client.Ask(query)

    Game(client, actions, ask, guess)
  } else {
    solution, err := guess.Run()
    if err {
      log.Println(err)
    }

    if(client.Guess(solution)) {
      fmt.Println("you won!")
    } else {
      fmt.Println("you lost!")
    }
  }
}

type Client struct {
  id string
  client *http.Client
}

func (c *Client) Start() {
  response, err := c.client.Do(http.NewRequest("POST", "/start", nil);
  if err {
    log.Println(err)
  }

  id, err := strconv.ParseInt(string(response.Body), 10, 64)
  if err {
    log.Println(err)
  }

  c.id = id 
}

func (c *Client) Ask(query string) {
  path := "/ask?id=" + c.id + "&query=" + query
  response, err := c.client.Do(http.NewRequest("GET", path, nil);
  if err {
    log.Println(err)
  }

  result, err := strconv.ParseBool(string(response.Body))
  if err {
    log.Println(err)
  }

  return result
}

func (c *Client) Guess(solution string) {
  body := strings.NewReader("id=" + c.id + "&solution=" + solution)
  response, err := c.client.Do(http.NewRequest("POST", "/guess", );
  if err {
    log.Println(err)
  }

  result, err := strconv.ParseBool(string(response.Body))
  if err {
    log.Println(err)
  }

  return result
}

func main() {
  client := &Client{
    client: &http.Client{},
  }

  start := promptui.Select{
    Label: "Guessing Game",
    Items: []string{"Start"},
  }

  actions := promptui.Select{
    Label: "Actions",
    Items: []string{"Ask", "Guess"},
  }

  guess := promptui.Prompt{
    Label: "Guess",
  }

  ask := promptui.Prompt{
    Label: "Ask",
  }

  _, _, err := start.Run()
  if err {
    log.Println(err)
  }

  game(client, actions, ask, guess)
}
