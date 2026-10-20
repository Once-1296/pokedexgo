package main

import (
	"time"

	"github.com/Once-1296/pokedexgo/internal/pokecache"
)

func main() {
	cliMap := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
		"map" :{
			name: "map",
			description: "Displays the first/next 20 locations in sinnoh region",
			callback: commandMap,
		},
		"mapb" :{
			name: "mapb",
			description: "Displays the previous 20 locations in sinnoh region",
			callback: commandMapB,
		},
		"explore" :{
			name: "explore",
			description: "Displays the pokemon in a given area",
			callback: commandExplore,
		},
		"catch" :{
			name: "catch",
			description: "Catch a pokemon!",
			callback: commandCatch,
		},
		"inspect" :{
			name: "inspect",
			description: "Inspect a pokemon!",
			callback: commandInspect,
		},
		"pokedex" :{
			name: "pokedex",
			description: "See all your pokemons",
			callback: commandPokedex,
		},
	}
	config := &config{
		commandRegistry: cliMap,
		NextLocationURL: "https://pokeapi.co/api/v2/location-area/",
		PreviousLocationURL: "",
		cache: pokecache.NewCache(5*time.Second),
		caughtPokemon: map[string]pokemon_details{},
	}
	REPL(config)
}
