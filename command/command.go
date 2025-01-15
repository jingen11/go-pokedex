package command

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/jingen11/pokedexcli/network"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func(Params) error
	Params      interface{}
}

type Command struct {
	Commands map[string]CliCommand
}

type Params struct {
	Params    interface{}
	Arguments []string
}

type HelpParams struct {
	commands map[string]struct {
		name        string
		description string
	}
}

type MapParams struct {
	client *network.NetworkClient
	offset int
}

type ExploreParams struct {
	client *network.NetworkClient
}

type CatchParams struct {
	client  *network.NetworkClient
	pokedex map[string]network.PokemonCatch
}

type InspectParams struct {
	pokedex map[string]network.PokemonCatch
}
type PokedexParams struct {
	pokedex map[string]network.PokemonCatch
}

func NewCommand(pokedex map[string]network.PokemonCatch) Command {
	mapParams := MapParams{
		client: network.NewNetworkClient(),
		offset: 0,
	}
	exploreParams := ExploreParams{
		client: network.NewNetworkClient(),
	}
	catchParams := CatchParams{
		client:  network.NewNetworkClient(),
		pokedex: pokedex,
	}
	inspectParams := InspectParams{
		pokedex: pokedex,
	}
	pokedexParams := PokedexParams{
		pokedex: pokedex,
	}
	helpParams := HelpParams{
		commands: map[string]struct {
			name        string
			description string
		}{},
	}
	command := Command{
		Commands: map[string]CliCommand{
			"exit": {
				Name:        "exit",
				Description: "Exit the Pokedex",
				Callback:    commandExit,
			},
			"help": {
				Name:        "help",
				Description: "Displays a help message",
				Callback:    commandHelp,
				Params:      &helpParams,
			},
			"map": {
				Name:        "map",
				Description: "Get location areas from the map",
				Callback:    commandMap,
				Params:      &mapParams,
			},
			"mapb": {
				Name:        "mapb",
				Description: "Get last page location areas from the map",
				Callback:    commandMapPrev,
				Params:      &mapParams,
			},
			"explore": {
				Name:        "explore",
				Description: "Explore the provided location",
				Callback:    commandExplore,
				Params:      &exploreParams,
			},
			"catch": {
				Name:        "catch",
				Description: "Catch a pokemon",
				Callback:    commandCatch,
				Params:      &catchParams,
			},
			"inspect": {
				Name:        "inspect",
				Description: "Inspect a caught pokemon",
				Callback:    commandInspect,
				Params:      &inspectParams,
			},
			"pokedex": {
				Name:        "pokedex",
				Description: "Get pokedex",
				Callback:    commandPokedex,
				Params:      &pokedexParams,
			},
		},
	}

	for k, v := range command.Commands {
		helpParams.commands[k] = struct {
			name        string
			description string
		}{
			name:        v.Name,
			description: v.Description,
		}
	}
	return command
}

func commandExit(args Params) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(args Params) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for _, v := range args.Params.(*HelpParams).commands {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}

	return nil
}

func commandMap(args Params) error {
	res, err := args.Params.(*MapParams).client.GetLocations(20, args.Params.(*MapParams).offset)
	if err != nil {
		return err
	}
	for _, location := range res {
		fmt.Println(location.Name)
	}

	if len(res) == 20 {
		args.Params.(*MapParams).offset += 20
	} else {
		fmt.Println("~~~~~~~End of records~~~~~~~~~~")
	}

	fmt.Println("")
	return nil
}

func commandMapPrev(args Params) error {
	isFirstPage := (args.Params.(*MapParams).offset-40 <= 0)
	if !isFirstPage {
		args.Params.(*MapParams).offset -= 40
	} else {
		args.Params.(*MapParams).offset = 0
	}

	res, err := args.Params.(*MapParams).client.GetLocations(20, args.Params.(*MapParams).offset)
	if err != nil {
		return err
	}
	for _, location := range res {
		fmt.Println(location.Name)
	}

	fmt.Println("")
	if len(res) == 20 {
		if isFirstPage {
			args.Params.(*MapParams).offset += 20
		}
	} else {
		fmt.Println("~~~~~~~End of records~~~~~~~~~~")
	}

	return nil
}

func commandExplore(args Params) error {
	if len(args.Arguments) < 2 {
		return errors.New("no area found")
	}
	area := args.Arguments[1]
	res, err := args.Params.(*ExploreParams).client.GetPokemons(area)
	if err != nil {
		return err
	}
	for _, v := range res.PokemonEncounters {
		fmt.Println(v.Pokemon.Name)
	}
	return nil
}

func commandCatch(args Params) error {
	if len(args.Arguments) < 2 {
		return errors.New("no pokemon found")
	}
	pokemon := args.Arguments[1]
	res, err := args.Params.(*CatchParams).client.CatchPokemon(pokemon)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)

	benchmark := 200 // increase to catch easier

	pow := float64(res.BaseExperience) / float64(benchmark)
	chance := math.Pow(0.5, pow)
	round := rand.Float64()

	if round > (1.0 - chance) {
		args.Params.(*CatchParams).pokedex[res.Name] = res
		fmt.Printf("%s was caught!\n", pokemon)
	} else {
		fmt.Printf("%s escaped!\n", pokemon)
	}

	return nil
}

func commandInspect(args Params) error {
	if len(args.Arguments) < 2 {
		return errors.New("no pokemon found")
	}
	pokemon := args.Arguments[1]

	p, ok := args.Params.(*InspectParams).pokedex[pokemon]
	if !ok {
		return errors.New("no caught pokemon found")
	}
	stats := network.PokemonStatMap{}

	for _, v := range p.Stats {
		switch v.Stat.Name {
		case "hp":
			stats.Hp = v.BaseStat
		case "attack":
			stats.Attack = v.BaseStat
		case "defense":
			stats.Defence = v.BaseStat
		case "special-attack":
			stats.SpecialAttack = v.BaseStat
		case "special-defense":
			stats.SpecialDefence = v.BaseStat
		case "speed":
			stats.Speed = v.BaseStat
		}
	}
	fmt.Println("Name: " + p.Name)
	fmt.Println("Height: " + strconv.Itoa(p.Height))
	fmt.Println("Weight: " + strconv.Itoa(p.Weight))
	fmt.Println("Stats: ")
	fmt.Println("  -hp: " + strconv.Itoa(stats.Hp))
	fmt.Println("  -attack: " + strconv.Itoa(stats.Attack))
	fmt.Println("  -defence: " + strconv.Itoa(stats.Defence))
	fmt.Println("  -special-attack: " + strconv.Itoa(stats.SpecialAttack))
	fmt.Println("  -special-defence: " + strconv.Itoa(stats.SpecialDefence))
	fmt.Println("  -speed: " + strconv.Itoa(stats.Speed))
	fmt.Println("Types: ")
	for _, v := range p.Types {
		fmt.Println("  - " + v.Type.Name)
	}
	return nil
}

func commandPokedex(args Params) error {
	fmt.Println("Your Pokedex:")
	for k := range args.Params.(*PokedexParams).pokedex {
		fmt.Println(" - " + k)
	}
	return nil
}
