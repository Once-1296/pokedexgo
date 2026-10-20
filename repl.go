package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"math/rand"

	"github.com/Once-1296/pokedexgo/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

type config struct {
	commandRegistry     map[string]cliCommand
	NextLocationURL     string
	PreviousLocationURL string
	cache               pokecache.Cache
	caughtPokemon map[string]pokemon_details
}

type location struct {
	Name string
	Url  string
}

type locationAPIResponse struct {
	Count    int
	Next     string
	Previous string
	Results  []location
}

type version struct {
	Name string
	Url string
}

type encounter_method struct {
	Name string
	Url string
}

type verson_detail struct {
	Rate int
	Version version
}


type encounter_method_rate struct {
	Encounter_method encounter_method
	Version_details []verson_detail
}

type language struct {
	Name string
	Url string
}

type city_name struct {
	Name string
	Language language
}

type pokemon struct {
	Name string
	Url string
}

type catch_method struct {
	Name string
	Url string
}

type catch_condition struct {
	Name string
	Url string
}

type encounter_detail struct {
	Min_level int
	Max_level int
	Chance int
	Method catch_method
	Condition_values []catch_condition
	Pokemon_details string
}

type encounter_version_detail struct {
	Version version
	Max_chance int
	Encounter_details []encounter_detail
}

type pokemon_encounter struct {
	Pokemon pokemon
	Version_details []encounter_version_detail
}

type locationDetails struct {
	Id int
	Name string
	Game_index int
	Encounter_method_rates []encounter_method_rate
	Location location
	Names []city_name
	Pokemon_encounters []pokemon_encounter
}

type stat struct {
	Name string
	Url string
}

type stat_detail struct {
	Base_stat int
	Effort int
	Stat stat
}

type pokemon_type struct {
	Name string
	Url string
}

type type_detail struct {
	slot int
	Type pokemon_type
}

type pokemon_details struct {
	Name string
	Height int
	Weight int
	Base_experience int
	Stats []stat_detail
	Types []type_detail
}

func cleanInput(test string) []string {
	var fields []string
	cutset := " \n\t\v\r\f"
	test = strings.TrimLeft(test, cutset)
	for len(test) > 0 {
		len_of_portion := len(test)
		word := ""
		j := len_of_portion
		for i := 0; i < len_of_portion; i++ {
			if strings.Contains(cutset, string(test[i])) {
				j = i
				break
			}
			word += string(test[i])
		}
		fields = append(fields, strings.ToLower(word))
		test = test[j:]
		test = strings.TrimLeft(test, cutset)
	}
	return fields
}

func commandExit(config *config,args ...string) error {
	if len(args) != 0 {
		fmt.Printf("Usage: exit\n")
		return fmt.Errorf("Usage: exit\n")
	}
	
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(config *config, args ...string) error {
	if len(args) != 0 {
		fmt.Printf("Usage: help\n")
		return fmt.Errorf("Usage: help\n")
	}
	
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, comm := range config.commandRegistry {
		fmt.Printf("%s: %s\n", comm.name, comm.description)
	}
	return nil
}

func commandMap(config *config, args ... string) error {
	if len(args) != 0 {
		fmt.Printf("Usage: map\n")
		return fmt.Errorf("Usage: map\n")
	}
	
	if config.NextLocationURL == "" {
		fmt.Print("you're on the last page\n")
		return fmt.Errorf("you're on the last page\n")
	}
	var data locationAPIResponse
	var bytes []byte
	var err error
	var res *http.Response
	var ok bool
	bytes, ok = config.cache.Get(config.NextLocationURL)
	if !ok {
		res, err = http.Get(config.NextLocationURL)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return err
		}
		defer res.Body.Close()
		bytes, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return err
		}
		config.cache.Add(config.NextLocationURL, bytes)
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	for _, loc := range data.Results {
		fmt.Printf("%s\n", loc.Name)
	}
	config.NextLocationURL = data.Next
	config.PreviousLocationURL = data.Previous
	return nil
}

func commandMapB(config *config,args ...string) error {
	if len(args) != 0 {
		fmt.Printf("Usage: mapb\n")
		return fmt.Errorf("Usage: mapb\n")
	}
	if config.PreviousLocationURL == "" {
		fmt.Print("you're on the first page\n")
		return fmt.Errorf("you're on the first page\n")
	}
	var data locationAPIResponse
	var bytes []byte
	var err error
	var res *http.Response
	var ok bool
	bytes, ok = config.cache.Get(config.PreviousLocationURL)
	if !ok {
		res, err = http.Get(config.PreviousLocationURL)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return err
		}
		defer res.Body.Close()
		bytes, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return err
		}
		config.cache.Add(config.PreviousLocationURL, bytes)
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	for _, loc := range data.Results {
		fmt.Printf("%s\n", loc.Name)
	}
	config.NextLocationURL = data.Next
	config.PreviousLocationURL = data.Previous
	return nil
}

func commandExplore(config *config, args ...string) error {
	if len(args) != 1 {
		fmt.Printf("Usage: explore <area-name>\n")
		return fmt.Errorf("Usage: explore <area-name>\n")
	}
	area := args[0]
	baseAPI := "https://pokeapi.co/api/v2/location-area/"
	baseAPI += area
	var data locationDetails
	var bytes []byte
	var err error
	var res *http.Response
	var ok bool
	bytes, ok = config.cache.Get(baseAPI)
	if !ok {
		res, err = http.Get(baseAPI)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return err
		}
		defer res.Body.Close()
		bytes, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return err
		}
		config.cache.Add(baseAPI, bytes)
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	fmt.Printf("Exploring %s...\n", data.Name)
	for _, encounter := range data.Pokemon_encounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *config, args ...string) error{
	if len(args) != 1 {
		fmt.Printf("Usage: catch <pokemon>\n")
		return fmt.Errorf("Usage: catch <pokemon>\n")
	}
	pokemon_name := args[0]
	baseAPI := "https://pokeapi.co/api/v2/pokemon/"
	baseAPI += pokemon_name
	var data pokemon_details
	var bytes []byte
	var err error
	var res *http.Response
	var ok bool
	bytes, ok = config.cache.Get(baseAPI)
	if !ok {
		res, err = http.Get(baseAPI)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return err
		}
		defer res.Body.Close()
		bytes, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return err
		}
		config.cache.Add(baseAPI, bytes)
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", data.Name)
	num := 36 + rand.Intn(750 - 36 + 1)
	if num >= data.Base_experience {
		fmt.Printf("%s was caught!\n", data.Name)
		config.caughtPokemon[data.Name] = data
	} else {
		fmt.Printf("%s escaped!\n",data.Name)
	}
	return nil
}

func commandInspect(config *config, args ...string) error {
	if len(args) != 1 {
		fmt.Printf("Usage: inspect <pokemon>\n")
		return fmt.Errorf("Usage: inspect <pokemon>\n")
	}
	pokemon_name := args[0]
	v, ok := config.caughtPokemon[pokemon_name]
	if !ok {
		fmt.Printf("you have not caught that pokemon\n")
		return fmt.Errorf("you have not caught that pokemon\n")
	}
	fmt.Printf("Name: %s\n",v.Name)
	fmt.Printf("Height: %d\n",v.Height)
	fmt.Printf("Weight: %d\n", v.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range v.Stats {
		fmt.Printf(" -%s: %d\n",stat.Stat.Name, stat.Base_stat)
	}
	fmt.Printf("Types:\n")
	for _, stat := range v.Types {
		fmt.Printf(" - %s\n",stat.Type.Name)
	}
	return nil
}


func commandPokedex(config *config, args ...string) error {
	if len(args) != 0 {
		fmt.Printf("Usage: pokedex\n")
		return fmt.Errorf("Usage: pokedex\n")
	}
	fmt.Printf("Your Pokedex: \n")
	for k,_ := range config.caughtPokemon {
		fmt.Printf(" - %s\n",k)
	}
	return nil
}

func REPL(config *config) error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		txt := scanner.Text()
		cleantxt := cleanInput(txt)
		// fmt.Printf("Your command was: %s\n",cleantxt[0])
		val, ok := config.commandRegistry[cleantxt[0]]
		if !ok {
			fmt.Print("Unknown command\n")
		} else {
			val.callback(config, cleantxt[1:]...)
		}
		scanner.Err()
	}
}

