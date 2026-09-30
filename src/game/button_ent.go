package game

import (
	"github.com/oidoid/void-template/src/gfx"
	"github.com/oidoid/void-template/src/tags"
	"github.com/oidoid/void/src/void/veng"
	"github.com/oidoid/void/src/void/vgeo"
	"github.com/oidoid/void/src/void/vgfx"
)

func newButtonEnt(label string) veng.ButtonEnt {
	edge := vgfx.Spr{TagCel: tags.WidgetEdgeLight.Cel(0)}
	fill := vgfx.Spr{TagCel: tags.WidgetFill.Cel(0)}
	this := veng.ButtonEnt{
		PatchByDir: [9]vgfx.Spr{
			vgeo.DirE:      edge,
			vgeo.DirN:      edge,
			vgeo.DirW:      edge,
			vgeo.DirS:      edge,
			vgeo.DirCenter: fill,
		},
		CornerWH: vgeo.NewWH[uint16](1, 1),
		Pals: veng.ButtonPals{
			Base:      tags.PalWidget,
			Focused:   tags.PalWidgetFocused,
			On:        tags.PalWidgetOn,
			FocusedOn: tags.PalWidgetFocusedOn,
		},
		TextPals: veng.ButtonPals{
			Base:      tags.PalText,
			Focused:   tags.PalText,
			On:        tags.PalTextLight,
			FocusedOn: tags.PalTextLight,
		},
		MinW: 48,
		Type: veng.ButtonTypeButton,
	}
	this.Text.SetText(label)
	this.Text.Z = gfx.ZUIText
	this.NinePatchEnt.SetZ(gfx.ZUIWidget)
	return this
}
