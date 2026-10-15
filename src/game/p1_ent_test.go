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

func TestP1EntIdleAnimationStaysStillAcrossRedraws(t *testing.T) {
	tests := []struct {
		name string
		walk vin.Key
		idle vin.Key
	}{
		{name: "released", walk: vin.KeyLeft},
		{
			name: "left and right", walk: vin.KeyLeft,
			idle: vin.KeyLeft | vin.KeyRight,
		},
		{
			name: "up and down", walk: vin.KeyUp,
			idle: vin.KeyUp | vin.KeyDown,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gam := NewGame()
			lvl := routeLvl0(gam)
			poll := gam.Poll()
			poll.CanvasPhy = vgeo.NewWH[uint16](256, 144)
			poll.DeltaMillis = 16
			poll.NowMillis = 187.5
			poll.InPoll.Kbd.Keys = test.walk
			gam.Update()
			xy, tag, flip := lvl.p1.XY, lvl.p1.Tag(), lvl.p1.FlipX()

			poll.InPoll.Kbd.Keys = test.idle
			for _, now := range []float64{250, 312.5, 999, 1000, 1062.5} {
				poll.NowMillis = now
				if got := gam.Update(); got != veng.Pause {
					t.Fatalf("idle status = %v, want Pause", got)
				}
				spr := gam.Layer(gfx.LayerP1).Sprs[0]
				if spr.XY != xy || spr.Tag() != tag || spr.FlipX() != flip {
					t.Fatalf("idle P1 changed position or facing: %v", spr)
				}
				// resolve the rendered cel using the sprite shader's clock.
				cel := (uint64(spr.Cel()) + uint64(now/vatlas.CelMillis)) &
					uint64(vatlas.TagCelMask)
				if cel != 0 {
					t.Fatalf("idle cel at %g ms = %d, want 0", now, cel)
				}
			}

			poll.InPoll.Kbd.Keys = test.walk
			for i := uint64(1); i <= 2; i++ {
				poll.NowMillis += vatlas.CelMillis
				if got := gam.Update(); got != veng.Loop {
					t.Fatalf("walking status = %v, want Loop", got)
				}
				spr := gam.Layer(gfx.LayerP1).Sprs[0]
				if spr.XY == xy {
					t.Fatal("P1 did not resume moving")
				}
				cel := (uint64(spr.Cel()) +
					uint64(poll.NowMillis/vatlas.CelMillis)) &
					uint64(vatlas.TagCelMask)
				if cel != i {
					t.Fatalf("walking cel = %d, want %d", cel, i)
				}
			}
		})
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
