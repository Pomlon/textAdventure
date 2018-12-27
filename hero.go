package main

type hero struct {
	pawn
	hp   int
	mana int
	atk  int
	def  int
}

func NewHero() *hero {
	return &hero{}
}
