package game

import (
	"github.com/oidoid/void-template/src/boards"
	"github.com/oidoid/void/src/void/veng"
)

type titleLvl struct {
	playButton PlayButtonEnt
}

func routeTitle(gam *Game) *titleLvl {
	gam.SetBoard(&boards.TitleBoard)
	gam.Texts().Clear()
	for _, spawn := range boards.TitleTextSpawns {
		if spawn.Hidden {
			continue
		}
		gam.Texts().Add(veng.NewTextEnt(spawn))
	}
	this := &titleLvl{playButton: NewPlayButtonEnt()}
	gam.Router().Update = this.update
	return this
}

func (this *titleLvl) update(gam *Game) veng.Status {
	return this.playButton.Update(gam)
}
