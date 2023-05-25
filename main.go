package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func getCommands() map[string]command {
	var commands = map[string]command{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "see 20 locations on the map",
			callback:    commandMap,
		},
	}

	return commands
}

func commandMap() error {

	type Response struct {
		Count    int    `json:"count"`
		Next     string `json:"next"`
		Previous any    `json:"previous"`
		Results  []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"results"`
	}

	locationAreaEndPoint := "https://pokeapi.co/api/v2/location-area"
	res, err := http.Get(locationAreaEndPoint)

	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	response := Response{}
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n")
	for _, locationArea := range response.Results {
		fmt.Println(locationArea.Name)
	}
	fmt.Printf("\n\n")

	return nil
}

func commandHelp() error {
	fmt.Printf("\n Welcome to the Pokedex \n\n\n")

	for name, command := range getCommands() {
		fmt.Printf(" %s: %s \n", name, command.description)
	}

	fmt.Printf("\n\n")

	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

type command struct {
	name        string
	description string
	callback    func() error
}

func main() {
	fmt.Printf("\n hello, from the pokedex! \n\n")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("pokedex > ")
		scanner.Scan()
		line := scanner.Text()

		if line == "exit" {
			break
		}

		if command, ok := getCommands()[line]; ok {
			command.callback()
		} else {
			fmt.Printf("\n Invalid command. Use \"help\" to see available commands \n\n")
		}
	}
}
