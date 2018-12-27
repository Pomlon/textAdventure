package main

import "errors"

type player struct {
	heroes   []*Hero
	position *Node
	inTown   bool
}

func NewPlayer() *player {
	p := player{

		inTown: true,
	}
	p.heroes = append(p.heroes, NewHero())
	return &p
}

func (p *player) GenerateHeroes() {
	var newher []*Hero
	for i := 0; i <= 3; i++ {
		newher = append(newher, &Hero{Pawn{
			Hp:     Randomize(300, 500),
			Mana:   Randomize(500, 1000),
			Dead:   false,
			Atk:    Randomize(20, 30),
			Def:    Randomize(2, 10),
			Action: true,
		}})
	}
	p.heroes = newher
}

func (p *player) GetFirstLivingHero() (*Hero, error) {
	for i := 0; i < len(p.heroes); i++ {
		if p.heroes[i].Dead == false {
			her := p.heroes[i]
			return her, nil
		}
	}
	return &Hero{}, errors.New("All ded")
}
