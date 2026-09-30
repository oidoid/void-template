package game

import (
	"github.com/oidoid/void-template/src/assets"
	"github.com/oidoid/void-template/src/boards"
	"github.com/oidoid/void-template/src/gfx"
	"github.com/oidoid/void-template/src/tags"
	"github.com/oidoid/void/src/void/vatlas"
	"github.com/oidoid/void/src/void/veng"
	"github.com/oidoid/void/src/void/vgeo"
	"github.com/oidoid/void/src/void/vgfx"
	"github.com/oidoid/void/src/void/vtext"
)

type Game struct {
	veng.Eng
	veng.GameHooks[*Game]
	cursor CursorEnt
}

var Version string

func NewGame() *Game {
	font := vtext.MemProp5x6
	font.FirstTag = tags.MemProp5x600
	this := &Game{Eng: *veng.NewEng(&veng.EngOpts{
		Atlas: vatlas.DecodeAtlas(assets.AtlasBin),
		Font:  font, RenderMode: vgfx.RenderModePixel,
	})}
	for _, layer := range []vgfx.Layer{
		gfx.LayerTiles, gfx.LayerP1, gfx.LayerUI, gfx.LayerCursor,
	} {
		config := this.Layer(layer)
		config.CamMode = vgfx.LayerCamModeFixed
		config.ScaleMode = vgfx.LayerScaleModeAutoInt
		config.AutoscaleMinClip = vgeo.NewWH[uint16](256, 144)
	}
	this.Layer(gfx.LayerTiles).Shader = vgfx.ShaderTiles
	this.Layer(gfx.LayerUI).Depth = true
	this.In().MapDefaults()
	this.RegisterPreupdate(UpdateFullscreen)
	routeTitle(this)
	this.cursor = NewCursorEnt(boards.TitleCursorSpawns[0], this.Atlas())
	this.SetCursor(&this.cursor.CursorEnt)
	this.RegisterUpdate(this.cursor.Update)
	return this
}

func (this *Game) Update() veng.Status {
	stat := this.Eng.BeginTick()
	stat |= this.GameHooks.Preupdate(&this.Eng, this)
	stat |= this.GameHooks.Update(this)
	stat |= this.Router().Update(this)
	return this.Eng.EndTick(stat)
}
