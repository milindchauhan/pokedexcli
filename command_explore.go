package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/milindchauhan/pokedexcli/internal/pokecache"
)

type exploreResp struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func commandExplore(config *Config, cache *pokecache.PokeCache) error {
	args := *config.args

	if len(args) < 2 {
		return errors.New(" missing argument. which area to explore?")
	}

	name := args[1]
	fmt.Printf(" Exploring %s", name)

	locationAreaEndPoint := "https://pokeapi.co/api/v2/location-area"
	url := locationAreaEndPoint + "/" + name

	var body []byte
	body, ok := cache.Get(url)
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		body, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return err
		}

	}

	response := exploreResp{}
	err := json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	for _, pe := range response.PokemonEncounters {
		fmt.Println(pe.Pokemon.Name)
	}

	return nil
}
