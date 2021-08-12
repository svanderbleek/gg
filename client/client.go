package main

import (
  "fmt"
  "net/http"
  "strconv"
  "io/ioutil"
  "github.com/manifoldco/promptui"
)

var base_url = "http://localhost:8080/"

func Game(client *Client, actions *promptui.Select, ask *promptui.Prompt, guess *promptui.Prompt) {
  _, action, err := actions.Run()
  if err != nil {
    panic(err)
  }

  if(action == "Ask") {
    query, err := ask.Run()
    if err != nil {
      panic(err)
    }

    if client.Ask(query) {
      fmt.Println("yes")
    } else {
      fmt.Println("no")
    }

    Game(client, actions, ask, guess)
  } else {
    solution, err := guess.Run()
    if err != nil {
      panic(err)
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
  request, err := http.NewRequest("POST", base_url + "start", nil)
  if err != nil {
    panic(err)
  }

  response, err := c.client.Do(request)
  if err != nil {
    panic(err)
  }

  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    panic(err)
  }

  c.id = string(body)
}

func (c *Client) Ask(query string) bool {
  path := "ask?id=" + c.id + "&query=" + query
  request, err := http.NewRequest("GET", base_url + path, nil)
  if err != nil {
    panic(err)
  }

  response, err := c.client.Do(request)
  if err != nil {
    panic(err)
  }

  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    panic(err)
  }

  result, err := strconv.ParseBool(string(body))
  if err != nil {
    panic(err)
  }

  return result
}

func (c *Client) Guess(solution string) bool {
  path := "guess?id=" + c.id + "&solution=" + solution
  request, err := http.NewRequest("POST", base_url + path, nil)
  if err != nil {
    panic(err)
  }

  response, err := c.client.Do(request)
  if err != nil {
    panic(err)
  }

  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    panic(err)
  }

  result, err := strconv.ParseBool(string(body))
  if err != nil {
    panic(err)
  }

  return result
}

func main() {
  client := &Client{
    client: &http.Client{},
  }

  start := &promptui.Select{
    Label: "Guessing Game",
    Items: []string{"Start"},
  }

  actions := &promptui.Select{
    Label: "Actions",
    Items: []string{"Ask", "Guess"},
  }

  guess := &promptui.Prompt{
    Label: "Guess",
  }

  ask := &promptui.Prompt{
    Label: "Ask",
  }

  _, _, err := start.Run()
  if err != nil {
    panic(err)
  }
  
  client.Start()

  Game(client, actions, ask, guess)
}
