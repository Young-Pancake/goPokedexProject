package main

import (
    "github.com/Young-Pancake/goPokedexProject/internal/pokeApi"
    "time"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

type config struct {
    commands            map[string]cliCommand
    pokeApiClient       pokeApi.Client
    nextLocationAreaURL *string
    prevLocationAreaURL *string
    caughtPokemon       map[string]pokeApi.Pokemon
}


func main() {
    cfg := &config{
        commands: getCommands(),
        pokeApiClient: pokeApi.NewClient(time.Hour),
        caughtPokemon: make(map[string]pokeApi.Pokemon),
    }
    startRepl(cfg)
}

