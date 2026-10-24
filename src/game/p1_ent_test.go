package game

import (
	"testing"

	"github.com/oidoid/void-template/src/gfx"
	"github.com/oidoid/void-template/src/tags"
	"github.com/oidoid/void/src/void/vatlas"
	"github.com/oidoid/void/src/void/vboards"
	"github.com/oidoid/void/src/void/veng"
	"github.com/oidoid/void/src/void/vgeo"
	"github.com/oidoid/void/src/void/vin"
)

func TestP1EntPausesAndResumesAuthoredAnimation(t *testing.T) {
	for _, test := range []struct {
		name string
		idle vin.Key
	}{
		{name: "released"},
		{name: "opposing horizontal", idle: vin.KeyLeft | vin.KeyRight},
		{name: "opposing vertical", idle: vin.KeyUp | vin.KeyDown},
	} {
		t.Run(test.name, func(t *testing.T) {
			gam := NewGame()
			lvl := routeLvl0(gam)
			poll := gam.Poll()
			poll.CanvasPhy = vgeo.NewWH[uint16](256, 144)
			poll.DeltaMillis = 150
			poll.InPoll.Kbd.Keys = vin.KeyLeft
			gam.Update()
			if lvl.p1.Cel != 1 {
				t.Fatalf("150 ms cel = %d, want 1", lvl.p1.Cel)
			}
			xy, tag, cel := lvl.p1.XY, lvl.p1.Tag, lvl.p1.Cel
			poll.InPoll.Kbd.Keys = test.idle
			for _, delta := range []float64{250, 62.5, 686.5, 1, 62.5} {
				poll.NowMillis += delta
				poll.DeltaMillis = delta
				if got := gam.Update(); got != veng.Pause {
					t.Fatalf("idle status = %v", got)
				}
				spr := gam.Layer(gfx.LayerP1).Sprs[0]
				if spr.XY != xy || spr.Tag != tag || spr.Cel != cel || !spr.FlipX() {
					t.Fatalf("idle sprite changed: %+v", spr)
				}
			}
			poll.InPoll.Kbd.Keys = vin.KeyLeft
			poll.DeltaMillis = 49
			if got := gam.Update(); got != veng.Loop {
				t.Fatalf("resumed status = %v, want Loop", got)
			}
			if lvl.p1.Cel != cel || lvl.p1.XY == xy {
				t.Fatalf("resume did not preserve the partial cel and move: %+v", lvl.p1.Spr)
			}
			poll.DeltaMillis = 1
			gam.Update()
			if lvl.p1.Cel != 2 {
				t.Fatalf("resumed cel = %d, want 2", lvl.p1.Cel)
			}
		})
	}
}

func TestP1EntPrefersUnblockedMovementAnimation(t *testing.T) {
	gam := NewGame()
	lvl := routeLvl0(gam)
	spawn := lvl.p1
	maxX := float32(gam.Board().W) - float32(spawn.W)
	maxY := float32(gam.Board().H) - float32(spawn.H)
	gam.Poll().DeltaMillis = 125

	for _, test := range []struct {
		name         string
		x, y         float32
		dx, dy       int8
		wantX, wantY float32
		tag          vatlas.Tag
		flip         bool
	}{
		{"left edge up", 0, 32, -1, -1, 0, 28, tags.P1WalkUp, false},
		{"left edge down", 0, 32, -1, 1, 0, 36, tags.P1WalkDown, false},
		{"right edge up", maxX, 32, 1, -1, maxX, 28, tags.P1WalkUp, false},
		{"right edge down", maxX, 32, 1, 1, maxX, 36, tags.P1WalkDown, false},
		{"top edge left", 32, 0, -1, -1, 28, 0, tags.P1WalkRight, true},
		{"top edge right", 32, 0, 1, -1, 36, 0, tags.P1WalkRight, false},
		{"bottom edge left", 32, maxY, -1, 1, 28, maxY, tags.P1WalkRight, true},
		{"bottom edge right", 32, maxY, 1, 1, 36, maxY, tags.P1WalkRight, false},
		{"free up left", 32, 32, -1, -1, 28, 28, tags.P1WalkRight, true},
		{"free up right", 32, 32, 1, -1, 36, 28, tags.P1WalkRight, false},
		{"free down left", 32, 32, -1, 1, 28, 36, tags.P1WalkRight, true},
		{"free down right", 32, 32, 1, 1, 36, 36, tags.P1WalkRight, false},
		{"blocked corner", 0, 0, -1, -1, 0, 0, tags.P1WalkRight, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			ent := spawn
			ent.XY = vgeo.NewXY(test.x, test.y)
			ent.SetFlipX(true)
			dir := vgeo.NewXY(test.dx, test.dy)
			gam.In().Dir = dir
			if got := ent.Update(gam); got != veng.Loop {
				t.Fatalf("movement status = %v, want Loop", got)
			}
			if want := vgeo.NewXY(test.wantX, test.wantY); ent.XY != want {
				t.Errorf("position = %v, want %v", ent.XY, want)
			}
			if ent.Tag != test.tag || ent.FlipX() != test.flip {
				t.Errorf("animation = %d, flip = %v; want %d, %v",
					ent.Tag, ent.FlipX(), test.tag, test.flip)
			}
			if gam.In().Dir != dir {
				t.Errorf("input direction changed: %v, want %v", gam.In().Dir, dir)
			}
		})
	}
}

func TestP1EntDrawsOnlyInsideLayerClip(t *testing.T) {
	gam := NewGame()
	ent := NewP1Ent(vboards.Spawn{
		WH: vgeo.NewWH[float32](8, 13), Z: gfx.ZP1, Tag: tags.P1WalkRight,
	})
	layer := gam.Layer(gfx.LayerP1)
	layer.Clip = vgeo.XYWH[float32](100, 100, 10, 10)
	ent.Update(gam)
	if got := len(layer.Sprs); got != 0 {
		t.Fatalf("offscreen sprite count = %d, want 0", got)
	}

	layer.Clip = vgeo.XYWH[float32](0, 0, 10, 10)
	ent.Update(gam)
	if got := len(layer.Sprs); got != 1 {
		t.Fatalf("visible sprite count = %d, want 1", got)
	}
}

func TestP1EntKeepsMovingAfterAnimationEnds(t *testing.T) {
	gam := NewGame()
	lvl := routeLvl0(gam)
	atlas := gam.Atlas()
	meta := &atlas.Meta[tags.P1WalkRight]
	meta.Cels = 1
	meta.Repeat = 1
	atlas.Cels[meta.CelIndex].Millis = 1
	gam.Poll().DeltaMillis = 125
	gam.In().Dir = vgeo.NewXY[int8](1, 0)
	for range 2 {
		xy := lvl.p1.XY
		if got := lvl.p1.Update(gam); got != veng.Loop {
			t.Fatalf("movement status = %v, want Loop", got)
		}
		if !lvl.p1.Done() || lvl.p1.X <= xy.X {
			t.Fatalf("finished animation stopped movement: %+v", lvl.p1.Spr)
		}
	}
}
