package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jingen11/pokedexcli/command"
	"github.com/jingen11/pokedexcli/network"
	"github.com/jingen11/pokedexcli/repl"
)

func main() {
	pokedex := map[string]network.PokemonCatch{}
	mapper := command.NewCommand(pokedex)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		arr := repl.CleanInput(scanner.Text())
		if len(arr) == 0 {
			continue
		}
		cliComm, ok := mapper.Commands[arr[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		cliComm.Callback(command.Params{
			Arguments: arr,
			Params:    cliComm.Params,
		})
	}
}
