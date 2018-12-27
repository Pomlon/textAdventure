package main

type Pawn struct {
	Hp     int
	Mana   int
	Atk    int
	Def    int
	Dead   bool
	Action bool
}

type IFPawn interface {
	getHit(amnt int)
	die()
}

func (p *Pawn) getHit(amnt int) int {
	var hitAfterCalcs int
	hitAfterCalcs = amnt - p.Def
	if hitAfterCalcs <= 0 {
		hitAfterCalcs = 1
	}
	p.Hp -= hitAfterCalcs
	if p.Hp <= 0 {
		p.die()
	}
	return hitAfterCalcs
}

func (p *Pawn) die() {
	p.Dead = true
	//todo drop items prolly
}
