package game

import (
	"github.com/oidoid/void-template/src/boards"
	"github.com/oidoid/void/src/void/veng"
)

type lvl0 struct {
	p1          P1Ent
	titleButton TitleButtonEnt
}

func routeLvl0(gam *Game) *lvl0 {
	if len(boards.Lvl0P1Spawns) != 1 {
		panic("lvl0 must have exactly one P1 spawn")
	}
	gam.Texts().Clear()
	gam.SetBoard(&boards.Lvl0Board)
	this := &lvl0{
		p1:          NewP1Ent(boards.Lvl0P1Spawns[0]),
		titleButton: NewTitleButtonEnt(),
	}
	gam.Router().Update = this.update
	return this
}

func (this *lvl0) update(gam *Game) veng.Status {
	stat := this.p1.Update(gam)
	stat |= this.titleButton.Update(gam)
	return stat
}
