package main

import (
	"bufio"
	"encoding/json"
	"errors"
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
		"mapb": {
			name:        "mapb",
			description: "see previous 20 location areas on the map",
			callback:    commandMapB,
		},
	}

	return commands
}

func commandMapB(config *Config) error {

	type Response struct {
		Count    int     `json:"count"`
		Next     *string `json:"next"`
		Previous *string `json:"previous"`
		Results  []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"results"`
	}

	if config.prev == nil {
		return errors.New("you are already on the first page")
	}
	locationAreaEndPoint := *config.prev

	resp, err := http.Get(locationAreaEndPoint)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	response := Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n")
	for _, loc := range response.Results {
		fmt.Printf(" %s\n", loc.Name)
	}
	fmt.Printf("\n")

	config.next = response.Next
	config.prev = response.Previous

	return nil
}

func commandMap(config *Config) error {

	type Response struct {
		Count    int     `json:"count"`
		Next     *string `json:"next"`
		Previous *string `json:"previous"`
		Results  []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"results"`
	}

	var locationAreaEndPoint string

	if config.next != nil {
		locationAreaEndPoint = *(config.next)
	} else {
		locationAreaEndPoint = "https://pokeapi.co/api/v2/location-area"
	}

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

	config.next = response.Next
	config.prev = response.Previous
	return nil
}

func commandHelp(config *Config) error {
	fmt.Printf("\n Welcome to the Pokedex \n\n\n")

	for name, command := range getCommands() {
		fmt.Printf(" %s: %s \n", name, command.description)
	}

	fmt.Printf("\n\n")

	return nil
}

func commandExit(config *Config) error {
	os.Exit(0)
	return nil
}

type command struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	next *string
	prev *string
}

func main() {
	fmt.Printf("\n hello, from the pokedex! \n\n")

	scanner := bufio.NewScanner(os.Stdin)
	config := Config{}

	for {
		fmt.Printf("pokedex > ")
		scanner.Scan()
		line := scanner.Text()

		if line == "exit" {
			break
		}

		if command, ok := getCommands()[line]; ok {
			err := command.callback(&config)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("\n Invalid command. Use \"help\" to see available commands \n\n")
		}
	}
}
