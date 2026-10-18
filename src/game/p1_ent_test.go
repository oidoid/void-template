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
			poll.DeltaMillis = 100
			poll.InPoll.Kbd.Keys = vin.KeyLeft
			gam.Update()
			if lvl.p1.Cel != 1 {
				t.Fatalf("100 ms cel = %d, want 1", lvl.p1.Cel)
			}
			xy, tag, cel := lvl.p1.XY, lvl.p1.Tag, lvl.p1.Cel
			poll.InPoll.Kbd.Keys = test.idle
			for _, now := range []float64{250, 312.5, 999, 1000, 1062.5} {
				poll.NowMillis = now
				if got := gam.Update(); got != veng.Pause {
					t.Fatalf("idle status = %v", got)
				}
				spr := gam.Layer(gfx.LayerP1).Sprs[0]
				if spr.XY != xy || spr.Tag != tag || spr.Cel != cel || !spr.FlipX() {
					t.Fatalf("idle sprite changed: %+v", spr)
				}
			}
			poll.InPoll.Kbd.Keys = vin.KeyLeft
			gam.Update()
			if lvl.p1.Cel != 2 {
				t.Fatalf("resumed cel = %d, want 2", lvl.p1.Cel)
			}
		})
	}
}

func TestP1EntUsesAtlasPlaybackSettings(t *testing.T) {
	gam := NewGame()
	lvl := routeLvl0(gam)
	atlas := gam.Atlas()
	atlas.Meta[tags.P1WalkRight].Dir = vatlas.AnimDirReverse
	atlas.Meta[tags.P1WalkRight].Repeat = 1
	atlas.Cels[atlas.Meta[tags.P1WalkRight].CelIndex+7].Millis = 237
	poll := gam.Poll()
	poll.CanvasPhy = vgeo.NewWH[uint16](256, 144)
	poll.InPoll.Kbd.Keys = vin.KeyRight
	poll.DeltaMillis = 236
	gam.Update()
	if lvl.p1.Cel != 7 {
		t.Fatalf("reverse start = %d", lvl.p1.Cel)
	}
	poll.DeltaMillis = 1
	gam.Update()
	if lvl.p1.Cel != 6 {
		t.Fatalf("authored duration = %d", lvl.p1.Cel)
	}
	poll.DeltaMillis = 10000
	if gam.Update() != veng.Loop {
		t.Fatal("movement must keep updating after animation finishes")
	}
	if lvl.p1.Cel != 0 || !lvl.p1.Done() {
		t.Fatal("one-shot did not stop")
	}
}

func TestP1EntKeepsLeftFacingWhenInputStops(t *testing.T) {
	gam := NewGame()
	ent := NewP1Ent(vboards.Spawn{
		WH: vgeo.NewWH[float32](8, 13), Z: gfx.ZP1, Tag: tags.P1WalkRight,
	})

	gam.In().Dir.X = -1
	ent.Update(gam)
	if !ent.FlipX() {
		t.Fatal("P1 is not facing left while Left is pressed")
	}

	gam.In().Dir = vgeo.XY[int8]{}
	ent.Update(gam)
	if !ent.FlipX() {
		t.Fatal("P1 stopped facing left when Left was released")
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
