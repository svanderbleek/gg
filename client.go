package main

import (
  "fmt"
  "log"
  "github.com/manifoldco/promptui"
)

func Play(id string) {
  actions := promptui.Select{
    Label: "Game #" + id,
    Items: []string{"Ask", "Guess"},
  }

  _, action, err := actions.Run()
  log.Println(err)

  if(action == "Ask") {
    ask := promptui.Prompt{
      Label: "Ask",
    }

    query, err := ask.Run()
    log.Println(err)

    fmt.Println(query)

    Play(id)
  } else {
    guess := promptui.Prompt{
      Label: "Guess",
    }

    answer, err := guess.Run()
    log.Println(err)

    fmt.Println(answer)

    if(answer == "42") {
      fmt.Println("you won!")
    } else {
      fmt.Println("you lost!")
    }
  }
}

func main() {
  start := promptui.Select{
    Label: "Guessing Game",
    Items: []string{"Start"},
  }

  _, _, err := start.Run()
  log.Println(err)

  Play("1")
}
