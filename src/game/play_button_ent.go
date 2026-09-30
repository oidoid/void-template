package game

import (
	"github.com/oidoid/void/src/void/veng"
	"github.com/oidoid/void/src/void/vgeo"
)

type PlayButtonEnt struct {
	veng.ButtonEnt
}

func NewPlayButtonEnt() PlayButtonEnt {
	this := PlayButtonEnt{ButtonEnt: newButtonEnt("play")}
	this.ClipAnchor = veng.HUDEnt{
		Anchor: vgeo.DirS,
		Margin: vgeo.Edge[int16]{S: 32},
	}
	this.AnchorMode = veng.ButtonAnchorHUD
	return this
}

func (this *PlayButtonEnt) Update(gam *Game) veng.Status {
	stat := this.ButtonEnt.Update(&gam.Eng)
	if this.IsOffStart() {
		routeLvl0(gam)
	}
	return stat
}
