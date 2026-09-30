package game

import (
	"testing"

	"github.com/oidoid/void-template/src/gfx"
	"github.com/oidoid/void-template/src/tags"
	"github.com/oidoid/void/src/void/vboards"
	"github.com/oidoid/void/src/void/vgeo"
)

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
