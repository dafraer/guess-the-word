package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/fatih/color"
)

// Checks if current letter from player's guess is found in the word
func contains(s string, c string) bool {
	for i := 0; i < len(s); i++ {
		if string(s[i]) == c {
			return true
		}
	}
	return false
}

func getRandomWord() (word string, hiddenWord string) {
	res, err := http.Get("https://random-word-api.herokuapp.com/word")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		panic("API is not currently available:/")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	word = string(body[2 : len(body)-2])
	for i := 0; i < len(word); i++ {
		hiddenWord += "*"
	}
	return
}

// Outputs the guess, highlighting correctly guessed letters
func output(word string, guess string, checker string) {
	for i := 0; i < len(guess); i++ {
		switch {
		case word[i] == guess[i]:
			color.Set(color.FgHiGreen)
			fmt.Printf("%s", string(guess[i]))
			color.Unset()
		case contains(checker, string(guess[i])):
			color.Set(color.FgHiYellow)
			fmt.Printf("%s", string(guess[i]))
			color.Unset()
		default:
			color.Set(color.FgHiRed)
			fmt.Printf("%s", string(guess[i]))
			color.Unset()
		}
	}
}

func main() {
	word, hiddenWord := getRandomWord()
	checker := word
	var guess string
	var win, lose bool
	tries := 10
	//Game process:
	fmt.Printf("Guess the Word!\nTries remaining: %v\n", tries)
	fmt.Println(hiddenWord)
out:
	for win == false && lose == false {
		fmt.Scan(&guess)
		tries--
		switch {
		case guess == word:
			win = true
			break out
		case len(guess) > len(word) || len(guess) < len(word):
			fmt.Printf("The word must contain %v letters. Try again\n", len(word))
			tries++
			continue
		case tries <= 0:
			lose = true
			break out
		}
		for i := 0; i < len(word); i++ {
			if word[i] == guess[i] && i < len(word)-1 {
				checker = checker[:i] + "*" + checker[i+1:]
			} else if word[i] == guess[i] {
				checker = checker[:i] + "*"
			}
		}
		output(word, guess, checker)
		fmt.Printf("  Tries remaining: %v\n", tries)
	}
	if win == true {
		color.HiGreen("YOU WIN!")
	} else {
		color.HiRed("YOU LOSE :/")
		fmt.Printf("The word was: %s", word)
	}
}
