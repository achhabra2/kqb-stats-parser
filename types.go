package main

type Player struct {
	Name    string `json:"Name"`
	Kills   int    `json:"Kills"`
	Deaths  int    `json:"Deaths"`
	Berries int    `json:"Berries"`
	Snail   int    `json:"Snail"`
	Queen   int    `json:"Queen"`
}

type Team struct {
	Color   string   `json:"Color"`
	MapsWon int      `json:MapsWon`
	Players []Player `json:"Players"`
}

type Set struct {
	Teams []Team `json:"Teams"`
}
