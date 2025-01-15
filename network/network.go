package network

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jingen11/pokedexcli/pokecache"
)

type PokemonListResponse struct {
	Count   int        `json:"count"`
	Results []Location `json:"results"`
}

type LocationArea struct {
	id                     int
	name                   string
	game_index             string
	encounter_method_rates []EncounterMethodRate
	Location               Location `json:"location"`
	names                  []Name
	PokemonEncounters      []PokemonEncounter `json:"pokemon_encounters"`
}

type EncounterMethodRate struct {
	name string
	url  string
}

type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Name struct {
	name string
}

type PokemonEncounter struct {
	Pokemon General `json:"pokemon"`
}

type General struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokemonCatch struct {
	Id             int           `json:"id"`
	Name           string        `json:"name"`
	BaseExperience int           `json:"base_experience"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Stats          []PokemonStat `json:"stats"`
	Types          []PokemonType `json:"types"`
}

type PokemonStatMap struct {
	Hp             int
	Attack         int
	Defence        int
	SpecialAttack  int
	SpecialDefence int
	Speed          int
}

type PokemonStat struct {
	BaseStat int     `json:"base_stat"`
	Effort   int     `json:"effort"`
	Stat     General `json:"stat"`
}

type PokemonType struct {
	Slot int     `json:"slot"`
	Type General `json:"type"`
}

type NetworkClient struct {
	client *http.Client
	cache  pokecache.Cache
}

func NewNetworkClient() *NetworkClient {
	client := NetworkClient{
		client: &http.Client{},
		cache:  pokecache.NewCache(time.Duration(5) * time.Second),
	}

	return &client
}

func (n *NetworkClient) GetLocations(limit, offset int) ([]Location, error) {
	url := "https://pokeapi.co/api/v2/location-area" + "?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset)
	var result PokemonListResponse
	if b, ok := n.cache.Get(url); ok {
		err := json.Unmarshal(b, &result)
		if err != nil {
			return nil, err
		}
		return result.Results, nil
	}
	res, err := n.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

func (n *NetworkClient) GetPokemons(area string) (LocationArea, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + area
	var result LocationArea
	if b, ok := n.cache.Get(url); ok {
		err := json.Unmarshal(b, &result)
		if err != nil {
			return result, err
		}
		return result, nil
	}
	res, err := n.client.Get(url)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (n *NetworkClient) CatchPokemon(pokemon string) (PokemonCatch, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemon
	var result PokemonCatch
	if b, ok := n.cache.Get(url); ok {
		err := json.Unmarshal(b, &result)
		if err != nil {
			return result, err
		}
		return result, nil
	}
	res, err := n.client.Get(url)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}
