package game

import (
	"testing"

	"github.com/oidoid/void-template/src/boards"
	"github.com/oidoid/void-template/src/gfx"
	"github.com/oidoid/void/src/void/veng"
	"github.com/oidoid/void/src/void/vgeo"
	"github.com/oidoid/void/src/void/vin"
)

func TestTextLifecycle(t *testing.T) {
	gam := NewGame()
	if got, want := gam.Texts().Len(), 1; got != want {
		t.Fatalf("title text count = %d, want %d", got, want)
	}
	routeLvl0(gam)
	if got := gam.Texts().Len(); got != 0 {
		t.Fatalf("lvl0 text count = %d, want 0", got)
	}
}

func TestPlayButtonRoutesToLvl0(t *testing.T) {
	gam := NewGame()
	gam.Poll().CanvasPhy = vgeo.NewWH[uint16](256, 144)
	titleLvl := routeTitle(gam)
	gam.Update()
	clickButton(gam, &titleLvl.playButton.ButtonEnt)
	if gam.Board() != &boards.Lvl0Board {
		t.Fatal("Play button did not route to lvl0")
	}
	gam.Update()
	p1Sprs := gam.Layer(gfx.LayerP1).Sprs
	if got, want := len(p1Sprs), 1; got != want {
		t.Fatalf("P1 sprite count = %d, want %d", got, want)
	}
	spawn := boards.Lvl0P1Spawns[0]
	if p1Sprs[0].Tag() != spawn.Tag || p1Sprs[0].WH != spawn.WH.Cast[uint16]() {
		t.Fatalf("P1 sprite does not match its Tiled representation: %v", p1Sprs[0])
	}
	if p1Sprs[0].XY != boards.Lvl0P1Spawns[0].XY {
		t.Fatalf(
			"P1 position = %v, want %v",
			p1Sprs[0].XY,
			boards.Lvl0P1Spawns[0].XY,
		)
	}
}

func TestTitleButtonRoutesToTitle(t *testing.T) {
	gam := NewGame()
	gam.Poll().CanvasPhy = vgeo.NewWH[uint16](256, 144)
	lvl0 := routeLvl0(gam)
	gam.Update()

	ui := gam.Layer(gfx.LayerUI)
	bounds := lvl0.titleButton.AnchorBox()
	if got, want := bounds.Max.X, ui.Clip.Max.X-4; got != want {
		t.Fatalf("Title button right edge = %v, want %v", got, want)
	}
	if got, want := bounds.Min.Y, ui.Clip.Min.Y+4; got != want {
		t.Fatalf("Title button top edge = %v, want %v", got, want)
	}
	clickButton(gam, &lvl0.titleButton.ButtonEnt)
	if gam.Board() != &boards.TitleBoard {
		t.Fatal("Title button did not route to title")
	}
	if got, want := gam.Texts().Len(), 1; got != want {
		t.Fatalf("title text count = %d, want %d", got, want)
	}
}

func clickButton(gam *Game, button *veng.ButtonEnt) {
	ui := gam.Layer(gfx.LayerUI)
	bounds := button.AnchorBox()
	center := vgeo.NewXY(
		(bounds.Min.X+bounds.Max.X)/2,
		(bounds.Min.Y+bounds.Max.Y)/2,
	)
	phy := ui.LayerToPhy(center)
	poll := &gam.Poll().InPoll
	poll.PtrsLen = 1
	poll.Ptrs[0] = vin.PtrPoll{
		ID: 1, Device: vin.PtrDevMouse, Primary: true,
		Phy:    vgeo.XYWH(phy.X, phy.Y, float32(1), float32(1)),
		Clicks: vin.ClickPrimary,
	}
	gam.Update()
	poll.Ptrs[0].Clicks = 0
	gam.Update()
}
