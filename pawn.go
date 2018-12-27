package main

type pawn struct {
	hp   int
	mana int
	atk  int
	def  int
}

type IFpawn interface {
	getHit()
	die()
}

func (h *pawn) getHit() {

}

func (h *pawn) die() {

}
