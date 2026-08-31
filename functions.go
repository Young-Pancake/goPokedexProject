package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"log"
	"errors"
    "math/rand"
)

func startRepl(cfg *config, args ...string) {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
        args := []string{}
        if len(words) > 1 {
            args = words[1:]
        }
		command, exists := cfg.commands[commandName]
		if exists {
			err := command.callback(cfg, args...)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Shows 20 locations in sets",
			callback:    commandMap,
		},
		"mapb": {
			name:        "map",
			description: "Shows 20 locations in sets backwards",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore {location_area}",
			description: "Lists the pokemon in a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch {pokemon_name}",
			description: "Catches a pokemon if you score high",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect {pokemon_name}",
			description: "Inspects a pokemon if you caught it",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Shows all caught pokemon",
			callback:    commandPokedex,
		},
	}
}

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range cfg.commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}

func commandMap(cfg *config, args ...string) error {
    resp, err := cfg.pokeApiClient.ListLocationAreas(cfg.nextLocationAreaURL)
    if err != nil {
        return err
    }
    fmt.Println("Location areas:")
    for _, area := range resp.Results {
        fmt.Printf(" - %s\n", area.Name)
    }
    cfg.nextLocationAreaURL = resp.Next
    cfg.prevLocationAreaURL = resp.Previous
    return nil
}

func commandMapb(cfg *config, args ...string) error {
    if cfg.prevLocationAreaURL == nil {
        return errors.New("previous location is nil")
    }
    resp, err := cfg.pokeApiClient.ListLocationAreas(cfg.prevLocationAreaURL)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Location areas:")
    for _, area := range resp.Results {
        fmt.Printf(" - %s\n", area.Name)
    }
    cfg.nextLocationAreaURL = resp.Next
    cfg.prevLocationAreaURL = resp.Previous
    return nil
}

func commandExplore(cfg *config, args ...string) error {
    if len(args) != 1 {
        return errors.New("no location area provided")
    }
    locationAreaName := args[0]
    locationArea, err := cfg.pokeApiClient.ListLocationArea(locationAreaName)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Pokemon in %s:\n", locationArea.Name)
    for _, pokemon := range locationArea.PokemonEncounters {
        fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
    }
    return nil
}

func commandCatch(cfg *config, args ...string) error {
    if len(args) != 1 {
        return errors.New("no pokemon name provided")
    }
    pokemonName := args[0]
    pokemon, err := cfg.pokeApiClient.GetPokemon(pokemonName)
    if err != nil {
        log.Fatal(err)
    }

    const threshold = 50
    randNum := rand.Intn(pokemon.BaseExperience)
    fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
    fmt.Println(randNum, pokemon.BaseExperience, threshold)
    if randNum > threshold {
        return fmt.Errorf("Failed to catch %s\n", pokemonName)
    }

    fmt.Printf("%s was caught!\n", pokemonName)
    cfg.caughtPokemon[pokemonName] = pokemon
    return nil
}

func commandInspect(cfg *config, args ...string) error {
    if len(args) != 1 {
        return errors.New("no pokemon name provided")
    }
    pokemonName := args[0]
    pokemon, ok := cfg.caughtPokemon[pokemonName]
    if !ok {
        return errors.New("you havent caught this pokemon yet...")
    }

    fmt.Printf("Name: %s\n", pokemon.Name)
    fmt.Printf("Height: %v\n", pokemon.Height)
    fmt.Printf("Weight: %v\n", pokemon.Weight)
    fmt.Println("Stats:")
    for _, stat := range pokemon.Stats {
        fmt.Printf(" - %s: %v\n", stat.Stat.Name, stat.BaseStat)
    }
    fmt.Println("Types:")
    for _, typ := range pokemon.Types {
        fmt.Printf(" - %s\n", typ.Type.Name)
    }
    return nil
}

func commandPokedex(cfg *config, args ...string) error {
    if len(cfg.caughtPokemon) == 0 {
        return errors.New("Pokedex is empty")
    }
    fmt.Println("Showing caught pokemon...")
    for _, pokemon := range cfg.caughtPokemon {
        fmt.Printf(" - %s\n", pokemon.Name)
    }
    return nil
}

