package main

type TypeRes struct {
	Data     any `json:"Data"`
	Metadata any `json:"Metadata"`
}

type TypeAseResPlayer struct {
	Name  string
	Ping  int
	Score int
}
type TypeAseResPlayers []TypeAseResPlayer

type TypeAseRes struct {
	Header  string
	Info    string
	Players TypeAseResPlayers
}
