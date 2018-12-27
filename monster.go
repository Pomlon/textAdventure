package main

import "errors"

type Monster struct {
	Pawn
}

type MonsterHandler struct {
	Monsters []*Monster
}

func (m *MonsterHandler) DeathStatus() bool {
	for _, mon := range m.Monsters {
		if mon.Hp > 0 {
			return false
		}
	}

	return true
}

func (m *MonsterHandler) GenerateMobs() {
	if len(m.Monsters) == 0 {
		for i := 0; i < 5; i++ {
			m.Monsters = append(m.Monsters, &Monster{Pawn{
				Hp:     Randomize(90, 130),
				Atk:    Randomize(1, 5),
				Def:    Randomize(1, 6),
				Dead:   false,
				Action: true,
			}})
		}
	}

}

func (m *MonsterHandler) GetFirstLivingThing() (*Monster, error) {
	for i := 0; i < len(m.Monsters); i++ {
		if m.Monsters[i].Dead == false {
			return m.Monsters[i], nil
		}
	}
	return &Monster{}, errors.New("Fak, no mobs")
}
