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

func TestP1EntCapsLongFrames(t *testing.T) {
	gam := NewGame()
	board := vboards.Board{
		WH:    vgeo.NewWH[int32](128, 128),
		Tile:  vgeo.NewWH[uint8](16, 16),
		Tiles: make([]vboards.Tile, 64),
	}
	gam.SetBoard(&board)
	gam.In().Dir = vgeo.NewXY[int8](1, 1)
	gam.Poll().DeltaMillis = 3000
	ent := NewP1Ent(vboards.Spawn{
		XY: vgeo.NewXY[float32](32, 32), WH: vgeo.NewWH[float32](8, 13),
		Z: gfx.ZP1, Tag: tags.P1WalkRight,
	})
	ent.Update(gam)
	if want := vgeo.NewXY[float32](40, 40); ent.XY != want {
		t.Fatalf("P1 position after a long frame = %v, want %v", ent.XY, want)
	}
}

func TestP1EntSlidesAcrossTileSeams(t *testing.T) {
	for _, test := range []struct {
		name     string
		xy       vgeo.XY[float32]
		dir      vgeo.XY[int8]
		wallCol  bool
		wallCell int
	}{
		{"right wall down", vgeo.NewXY[float32](57, 32),
			vgeo.NewXY[int8](0, 1), true, 4},
		{"right wall up", vgeo.NewXY[float32](57, 80),
			vgeo.NewXY[int8](0, -1), true, 4},
		{"left wall down", vgeo.NewXY[float32](63, 32),
			vgeo.NewXY[int8](0, 1), true, 3},
		{"left wall up", vgeo.NewXY[float32](63, 80),
			vgeo.NewXY[int8](0, -1), true, 3},
		{"bottom wall right", vgeo.NewXY[float32](32, 51),
			vgeo.NewXY[int8](1, 0), false, 4},
		{"bottom wall left", vgeo.NewXY[float32](80, 51),
			vgeo.NewXY[int8](-1, 0), false, 4},
		{"top wall right", vgeo.NewXY[float32](32, 60),
			vgeo.NewXY[int8](1, 0), false, 3},
		{"top wall left", vgeo.NewXY[float32](80, 60),
			vgeo.NewXY[int8](-1, 0), false, 3},
		{"hold into right wall", vgeo.NewXY[float32](57, 32),
			vgeo.NewXY[int8](1, 1), true, 4},
		{"hold into left wall", vgeo.NewXY[float32](63, 32),
			vgeo.NewXY[int8](-1, 1), true, 3},
		{"hold into bottom wall", vgeo.NewXY[float32](32, 51),
			vgeo.NewXY[int8](1, 1), false, 4},
		{"hold into top wall", vgeo.NewXY[float32](32, 60),
			vgeo.NewXY[int8](1, -1), false, 3},
	} {
		t.Run(test.name, func(t *testing.T) {
			gam := NewGame()
			board := vboards.Board{
				WH:    vgeo.NewWH[int32](160, 160),
				Tile:  vgeo.NewWH[uint8](16, 16),
				Tiles: make([]vboards.Tile, 100),
			}
			for i := range 10 {
				x, y := i, test.wallCell
				if test.wallCol {
					x, y = test.wallCell, i
				}
				board.Tiles[y*10+x] = vboards.NewTile(1, true)
			}
			gam.SetBoard(&board)
			gam.Poll().DeltaMillis = 125
			ent := NewP1Ent(vboards.Spawn{
				XY: test.xy, WH: vgeo.NewWH[float32](8, 13),
				Z: gfx.ZP1, Tag: tags.P1WalkRight,
			})
			// walk into the wall before checking movement along it.
			into := int8(1)
			if test.wallCell == 3 {
				into = -1
			}
			if test.wallCol {
				gam.In().Dir = vgeo.NewXY(into, int8(0))
			} else {
				gam.In().Dir = vgeo.NewXY(int8(0), into)
			}
			ent.Update(gam)
			if ent.XY != test.xy {
				t.Fatalf("P1 entered the wall: %v, want %v", ent.XY, test.xy)
			}
			gam.In().Dir = test.dir
			want := test.xy
			for frame := range 8 {
				if test.wallCol {
					want.Y += float32(test.dir.Y) * 4
				} else {
					want.X += float32(test.dir.X) * 4
				}
				ent.Update(gam)
				if ent.XY != want {
					t.Fatalf("frame %d: P1 stopped at %v, want %v",
						frame, ent.XY, want)
				}
			}
		})
	}
}

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
			gam.SetBoard(&vboards.Board{
				WH:    vgeo.NewWH[int32](256, 144),
				Tile:  vgeo.NewWH[uint8](16, 16),
				Tiles: make([]vboards.Tile, 16*9),
			})
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

func TestP1EntPrefersUnblockedMovementAnimation(t *testing.T) {
	gam := NewGame()
	board := vboards.Board{
		WH: vgeo.NewWH[int32](128, 128), Tile: vgeo.NewWH[uint8](16, 16),
		Tiles: make([]vboards.Tile, 64),
	}
	gam.SetBoard(&board)
	hurtbox := gam.Atlas().Meta[tags.P1WalkRight].Hurtbox.Cast[float32]()
	minX, maxX := -hurtbox.Min.X, 128-hurtbox.Max.X
	minY, maxY := -hurtbox.Min.Y, 128-hurtbox.Max.Y
	gam.Poll().DeltaMillis = 125

	for _, test := range []struct {
		name         string
		x, y         float32
		dx, dy       int8
		wantX, wantY float32
		tag          vatlas.Tag
		flip, solid  bool
	}{
		{"left edge up", minX, 32, -1, -1, minX, 28, tags.P1WalkUp, false, false},
		{"left edge down", minX, 32, -1, 1, minX, 36, tags.P1WalkDown, false, false},
		{"right edge up", maxX, 32, 1, -1, maxX, 28, tags.P1WalkUp, false, false},
		{"right edge down", maxX, 32, 1, 1, maxX, 36, tags.P1WalkDown, false, false},
		{"top edge left", 32, minY, -1, -1, 28, minY, tags.P1WalkRight, true, false},
		{"top edge right", 32, minY, 1, -1, 36, minY, tags.P1WalkRight, false, false},
		{"bottom edge left", 32, maxY, -1, 1, 28, maxY, tags.P1WalkRight, true, false},
		{"bottom edge right", 32, maxY, 1, 1, 36, maxY, tags.P1WalkRight, false, false},
		{"free up left", 32, 32, -1, -1, 28, 28, tags.P1WalkRight, true, false},
		{"free up right", 32, 32, 1, -1, 36, 28, tags.P1WalkRight, false, false},
		{"free down left", 32, 32, -1, 1, 28, 36, tags.P1WalkRight, true, false},
		{"free down right", 32, 32, 1, 1, 36, 36, tags.P1WalkRight, false, false},
		{"blocked corner", minX, minY, -1, -1, minX, minY, tags.P1WalkRight, true, false},
		{"solid right up", 48 - hurtbox.Max.X, 49, 1, -1, 48 - hurtbox.Max.X, 45, tags.P1WalkUp, false, true},
		{"solid right down", 48 - hurtbox.Max.X, 49, 1, 1, 48 - hurtbox.Max.X, 53, tags.P1WalkDown, false, true},
		{"solid left down", 64 - hurtbox.Min.X, 49, -1, 1, 64 - hurtbox.Min.X, 53, tags.P1WalkDown, false, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			board.Tiles[3*8+3] = 0
			if test.solid {
				board.Tiles[3*8+3] = vboards.NewTile(1, true)
			}
			ent := NewP1Ent(vboards.Spawn{
				XY: vgeo.NewXY(test.x, test.y), WH: vgeo.NewWH[float32](8, 13),
				Z: gfx.ZP1, Tag: tags.P1WalkRight,
			})
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

func TestP1EntKeepsAnimatingWhileSliding(t *testing.T) {
	gam := NewGame()
	gam.SetBoard(&vboards.Board{
		WH:    vgeo.NewWH[int32](128, 128),
		Tile:  vgeo.NewWH[uint8](16, 16),
		Tiles: make([]vboards.Tile, 64),
	})
	hurtbox := gam.Atlas().Meta[tags.P1WalkRight].Hurtbox
	ent := NewP1Ent(vboards.Spawn{
		XY:  vgeo.NewXY(128-float32(hurtbox.Max.X), 32),
		WH:  vgeo.NewWH[float32](8, 13),
		Z:   gfx.ZP1,
		Tag: tags.P1WalkRight,
	})
	gam.In().Dir = vgeo.NewXY[int8](1, 1)
	for _, test := range []struct {
		delta float64
		cel   uint8
	}{{150, 1}, {49, 1}, {1, 2}} {
		gam.Poll().DeltaMillis = test.delta
		ent.Update(gam)
		if ent.Tag != tags.P1WalkDown || ent.Cel != test.cel {
			t.Fatalf("sliding animation = %d, cel = %d; want %d, %d",
				ent.Tag, ent.Cel, tags.P1WalkDown, test.cel)
		}
	}
}

func TestP1EntKeepsFacingWhenCornerBlocksMovement(t *testing.T) {
	for _, solid := range []bool{false, true} {
		for _, dir := range []vgeo.XY[int8]{
			vgeo.NewXY[int8](1, -1), vgeo.NewXY[int8](-1, -1),
			vgeo.NewXY[int8](1, 1), vgeo.NewXY[int8](-1, 1),
		} {
			for _, vertical := range []bool{false, true} {
				gam := NewGame()
				board := vboards.Board{
					WH:    vgeo.NewWH[int32](128, 128),
					Tile:  vgeo.NewWH[uint8](16, 16),
					Tiles: make([]vboards.Tile, 64),
				}
				lo, hi := float32(0), float32(128)
				if solid {
					lo, hi = 16, 112
					for i := range 8 {
						for _, cell := range []int{i, 56 + i, i * 8, i*8 + 7} {
							board.Tiles[cell] = vboards.NewTile(1, true)
						}
					}
				}
				gam.SetBoard(&board)
				box := gam.Atlas().Meta[tags.P1WalkRight].Hurtbox.Cast[float32]()
				corner := vgeo.NewXY(lo-box.Min.X, lo-box.Min.Y)
				if dir.X > 0 {
					corner.X = hi - box.Max.X
				}
				if dir.Y > 0 {
					corner.Y = hi - box.Max.Y
				}
				start := corner
				tag, flip := tags.P1WalkRight, dir.X < 0
				if vertical {
					start.Y -= float32(dir.Y) * 4
					tag, flip = tags.P1WalkDown, false
					if dir.Y < 0 {
						tag = tags.P1WalkUp
					}
				} else {
					start.X -= float32(dir.X) * 4
				}
				ent := NewP1Ent(vboards.Spawn{
					XY: start, WH: vgeo.NewWH[float32](8, 13),
					Z: gfx.ZP1, Tag: tags.P1WalkRight,
				})
				gam.In().Dir = dir
				gam.Poll().DeltaMillis = 125
				for frame := range 4 {
					ent.Update(gam)
					if ent.XY != corner || ent.Tag != tag || ent.FlipX() != flip {
						t.Fatalf("solid %v, dir %v, vertical %v, frame %d: position %v, tag %d, flip %v; want %v, %d, %v",
							solid, dir, vertical, frame, ent.XY, ent.Tag,
							ent.FlipX(), corner, tag, flip)
					}
				}
			}
		}
	}
}
