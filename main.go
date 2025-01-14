package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		for scanner.Scan() {
			fmt.Println("replied" + scanner.Text())
		}
	}
}
