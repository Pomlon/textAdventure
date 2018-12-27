package main

type player struct {
	heroes   []*hero
	position *Node
	inTown   bool
}

func NewPlayer() *player {
	p := player{
		heroes: make([]*hero, 0),
		inTown: true,
	}
	p.heroes = append(p.heroes, NewHero())
	return &p
}
