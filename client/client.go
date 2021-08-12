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
  log.Println(err)

  if(action == "Ask") {
    query, err := ask.Run()
    log.Println(err)

    client.Ask(query)

    Game(client, actions, ask, guess)
  } else {
    answer, err := guess.Run()
    log.Println(err)

    client.Guess(answer)

    if(client.Won) {
      fmt.Println("you won!")
    } else {
      fmt.Println("you lost!")
    }
  }
}

type Client struct {
  id string
  Won bool
  client *http.Client
}

func (c *Client) Start() {
  response, err := c.client.Do(http.NewRequest("POST", "/start", nil);
  log.Println(err)

  id, err := strconv.ParseInt(string(response.Body), 10, 64)
  log.Println(err)

  c.id = id 
}

func (c *Client) Ask(query string) {
  response, err := c.client.Do(http.NewRequest("GET", "/ask?query=" + query, nil);
  log.Println(err)
}

func (c *Client) Guess(answer string) {
  response, err := c.client.Do(http.NewRequest("POST", "/guess", strings.NewReader("answer=" + answer));
  log.Println(err)
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
  log.Println(err)

  game(client, actions, ask, guess)
}
