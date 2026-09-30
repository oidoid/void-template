package game

import (
	"github.com/oidoid/void/src/void/veng"
	"github.com/oidoid/void/src/void/vgeo"
)

type TitleButtonEnt struct {
	veng.ButtonEnt
}

func NewTitleButtonEnt() TitleButtonEnt {
	this := TitleButtonEnt{ButtonEnt: newButtonEnt("title")}
	this.ClipAnchor = veng.HUDEnt{
		Anchor: vgeo.DirNE,
		Margin: vgeo.Edge[int16]{E: 4, N: 4},
	}
	this.AnchorMode = veng.ButtonAnchorHUD
	return this
}

func (this *TitleButtonEnt) Update(gam *Game) veng.Status {
	stat := this.ButtonEnt.Update(&gam.Eng)
	if this.IsOffStart() {
		routeTitle(gam)
	}
	return stat
}
