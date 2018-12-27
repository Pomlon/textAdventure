package main

type hero struct {
	pawn
}

func NewHero() *hero {
	return &hero{}
}
