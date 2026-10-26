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
	gam, ent := newP1Test(vgeo.NewXY[float32](32, 32))
	gam.In().Dir = vgeo.NewXY[int8](1, 1)
	gam.Poll().DeltaMillis = 3000
	ent.Update(gam)
	if want := vgeo.NewXY[float32](40, 40); ent.XY != want {
		t.Fatalf("position after a long frame = %v, want %v", ent.XY, want)
	}
}

func TestP1EntHitsWallUsesDrawnOrigin(t *testing.T) {
	gam, ent := newP1Test(vgeo.NewXY[float32](32, 32))
	for y := range 8 {
		gam.Board().Tiles[y*8+4] = vboards.NewTile(1, true)
	}
	box := gam.Atlas().Meta[ent.Tag].Hurtbox.Cast[float32]()
	lo, hi := -box.Min.X, 64-box.Max.X
	for _, test := range []struct {
		name string
		x    float32
		want bool
	}{
		{"right before rounding", hi + .49, false},
		{"right rounded overlap", hi + .5, true},
		{"left half tie", lo - .5, false},
		{"left rounded overlap", lo - .51, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := ent.hitsWall(gam, vgeo.NewXY(test.x, float32(32)),
				tags.P1WalkRight); got != test.want {
				t.Fatalf("hitsWall(%v) = %v, want %v", test.x, got, test.want)
			}
		})
	}
}

func TestP1EntSynchronizesDiagonalAfterCardinalMovement(t *testing.T) {
	gam, ent := newP1Test(vgeo.NewXY[float32](32, 32))
	for i, move := range []struct {
		dx, dy           int8
		by, wantX, wantY float32
	}{
		{1, 0, .25, 32, 32},
		{1, -1, .25, 32, 32},
		{1, -1, .5, 32, 32},
		{1, -1, .25, 33, 31},
		{1, -1, 1, 34, 30},

		{-1, 0, .25, 34, 30},
		{-1, -1, .25, 34, 30},
		{-1, -1, .5, 34, 30},
		{-1, -1, .25, 33, 29},
		{-1, -1, 1, 32, 28},

		{-1, 0, .25, 32, 28},
		{-1, 1, .25, 32, 28},
		{-1, 1, .5, 32, 28},
		{-1, 1, .25, 31, 29},
		{-1, 1, 1, 30, 30},

		{1, 0, .25, 30, 30},
		{1, 1, .25, 30, 30},
		{1, 1, .5, 30, 30},
		{1, 1, .25, 31, 31},
		{1, 1, 1, 32, 32},
	} {
		stepP1Test(t, gam, &ent, vgeo.NewXY(move.dx, move.dy), move.by)
		want := vgeo.NewXY(move.wantX, move.wantY)
		if ent.XY != want || ent.Tag != tags.P1WalkRight ||
			ent.FlipX() != (move.dx < 0) {
			t.Fatalf("move %d: sprite %+v, want position %v and horizontal facing",
				i, ent.Spr, want)
		}
	}
}

func TestP1EntResetsChangedMovementAxes(t *testing.T) {
	gam, ent := newP1Test(vgeo.NewXY[float32](32, 32))
	for i, move := range []struct {
		dx, dy   int8
		by, x, y float32
	}{
		{1, 0, .75, 32, 32},
		{1, 1, .25, 32, 32},
		{-1, -1, .25, 32, 32},
		{-1, -1, .75, 31, 31},
		{0, 0, 1, 31, 31},
		{1, 0, .5, 31, 31},
		{0, 0, 0, 31, 31},
		{1, 0, .5, 31, 31},
		{1, 0, .5, 32, 31},
	} {
		stepP1Test(t, gam, &ent, vgeo.NewXY(move.dx, move.dy), move.by)
		if want := vgeo.NewXY(move.x, move.y); ent.XY != want {
			t.Fatalf("move %d: position %v, want %v", i, ent.XY, want)
		}
	}
}

func TestP1EntSlidesAcrossTileSeams(t *testing.T) {
	for _, test := range []struct {
		name     string
		dx, dy   int8
		wallCol  bool
		wallCell int
		tag      vatlas.Tag
		flip     bool
	}{
		{"right wall up", 1, -1, true, 4, tags.P1WalkUp, false},
		{"left wall down", -1, 1, true, 3, tags.P1WalkDown, false},
		{"bottom wall left", -1, 1, false, 4, tags.P1WalkRight, true},
		{"top wall right", 1, -1, false, 3, tags.P1WalkRight, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			gam, ent := newP1Test(vgeo.NewXY[float32](48, 48))
			for i := range 8 {
				x, y := i, test.wallCell
				if test.wallCol {
					x, y = test.wallCell, i
				}
				gam.Board().Tiles[y*8+x] = vboards.NewTile(1, true)
			}
			box := gam.Atlas().Meta[ent.Tag].Hurtbox.Cast[float32]()
			dir := vgeo.NewXY(test.dx, test.dy)
			free := dir
			if test.wallCol {
				ent.X, free.X = 64-box.Max.X, 0
				if test.dx < 0 {
					ent.X = 64 - box.Min.X
				}
			} else {
				ent.Y, free.Y = 64-box.Max.Y, 0
				if test.dy < 0 {
					ent.Y = 64 - box.Min.Y
				}
			}
			start := ent.XY
			ent.XY.AddTo(vgeo.NewXY[float32](.25, .25))
			stepP1Test(t, gam, &ent, dir.Sub(free), 4)
			if ent.XY != start {
				t.Fatalf("entered the wall: %v, want %v", ent.XY, start)
			}
			for i, move := range []struct {
				dir        vgeo.XY[int8]
				by, wantBy float32
			}{
				{dir, .5, 0}, {dir, .5, 1}, {dir, 7, 8}, {dir, 8, 16},
				{free, .5, 16}, {free, .5, 17}, {free, 7, 24}, {free, 8, 32},
			} {
				stepP1Test(t, gam, &ent, move.dir, move.by)
				want := start.Add(free.Cast[float32]().Mul(move.wantBy))
				if ent.XY != want || move.wantBy != 0 &&
					(ent.Tag != test.tag || ent.FlipX() != test.flip) {
					t.Fatalf("move %d: position %v, tag %d, flip %v; want %v, %d, %v",
						i, ent.XY, ent.Tag, ent.FlipX(), want, test.tag, test.flip)
				}
			}
		})
	}
}

func TestP1EntResolvesDiagonalTileCorners(t *testing.T) {
	type step struct{ by, wantX float32 }
	for _, test := range []struct {
		name   string
		x, y   float32
		dx, dy int8
		steps  []step
	}{
		{"north east", 57, 76, 1, -1,
			[]step{{.25, 57}, {.5, 57}, {.25, 58}, {1, 59}}},
		{"north west", 79, 76, -1, -1,
			[]step{{.25, 79}, {.5, 79}, {.25, 78}, {1, 77}}},
		{"south west", 79, 51, -1, 1,
			[]step{{.25, 79}, {.5, 79}, {.25, 78}, {1, 77}}},
		{"south east", 57, 51, 1, 1,
			[]step{{.25, 57}, {.5, 57}, {.25, 58}, {1, 59}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			gam, ent := newP1Test(vgeo.NewXY(test.x, test.y))
			gam.Board().Tiles[4*8+4] = vboards.NewTile(1, true)
			for i, move := range test.steps {
				stepP1Test(t, gam, &ent, vgeo.NewXY(test.dx, test.dy), move.by)
				if want := vgeo.NewXY(move.wantX, test.y); ent.XY != want {
					t.Fatalf("move %d: corner slide %v, want %v", i, ent.XY, want)
				}
			}
		})
	}
}

func TestP1EntFacingAtBlockedCorners(t *testing.T) {
	for _, test := range []struct {
		name     string
		dx, dy   int8
		vertical bool
	}{
		{"north east", 1, -1, true}, {"north west", -1, -1, false},
		{"south west", -1, 1, true}, {"south east", 1, 1, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			gam, ent := newP1Test(vgeo.NewXY[float32](32, 32))
			box := gam.Atlas().Meta[ent.Tag].Hurtbox.Cast[float32]()
			corner := vgeo.NewXY(128-box.Max.X, 128-box.Max.Y)
			verticalTag := tags.P1WalkDown
			if test.dx < 0 {
				corner.X = -box.Min.X
			}
			if test.dy < 0 {
				corner.Y, verticalTag = -box.Min.Y, tags.P1WalkUp
			}
			ent.XY = corner
			tag, flip := tags.P1WalkRight, test.dx < 0
			if test.vertical {
				ent.Y -= float32(test.dy) * 4
				tag, flip = verticalTag, false
			} else {
				ent.X -= float32(test.dx) * 4
			}
			for _, move := range []struct {
				name   string
				dx, dy int8
				tag    vatlas.Tag
				flip   bool
			}{
				{"slide to corner", test.dx, test.dy, tag, flip},
				{"hold diagonal", test.dx, test.dy, tag, flip},
				{"release vertical", test.dx, 0, tags.P1WalkRight, test.dx < 0},
				{"release horizontal", 0, 0, tags.P1WalkRight, test.dx < 0},
				{"press horizontal again", test.dx, 0, tags.P1WalkRight, test.dx < 0},
				{"press vertical", 0, test.dy, verticalTag, false},
				{"release vertical", 0, 0, verticalTag, false},
				{"press horizontal from idle", test.dx, 0, tags.P1WalkRight, test.dx < 0},
			} {
				cel := ent.Cel
				stepP1Test(t, gam, &ent, vgeo.NewXY(move.dx, move.dy), 4)
				if ent.XY != corner || ent.Tag != move.tag || ent.FlipX() != move.flip {
					t.Fatalf("%s: position %v, tag %d, flip %v; want %v, %d, %v",
						move.name, ent.XY, ent.Tag, ent.FlipX(),
						corner, move.tag, move.flip)
				}
				if move.dx == 0 && move.dy == 0 && ent.Cel != cel {
					t.Fatalf("%s: idle animation advanced from %d to %d",
						move.name, cel, ent.Cel)
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
	gam, ent := newP1Test(vgeo.NewXY[float32](32, 32))
	atlas := gam.Atlas()
	meta := &atlas.Meta[tags.P1WalkRight]
	meta.Cels = 1
	meta.Repeat = 1
	atlas.Cels[meta.CelIndex].Millis = 1
	gam.Poll().DeltaMillis = 125
	gam.In().Dir = vgeo.NewXY[int8](1, 0)
	for range 2 {
		xy := ent.XY
		if got := ent.Update(gam); got != veng.Loop {
			t.Fatalf("movement status = %v, want Loop", got)
		}
		if !ent.Done() || ent.X <= xy.X {
			t.Fatalf("finished animation stopped movement: %+v", ent.Spr)
		}
	}
}

func TestP1EntKeepsAnimatingWhileSliding(t *testing.T) {
	gam, ent := newP1Test(vgeo.NewXY[float32](32, 32))
	box := gam.Atlas().Meta[ent.Tag].Hurtbox
	ent.X = 128 - float32(box.Max.X)
	for _, test := range []struct {
		delta float64
		cel   uint8
	}{{150, 1}, {49, 1}, {1, 2}} {
		gam.In().PrevDir = gam.In().Dir
		gam.In().Dir = vgeo.NewXY[int8](1, 1)
		gam.Poll().DeltaMillis = test.delta
		ent.Update(gam)
		if ent.Tag != tags.P1WalkDown || ent.Cel != test.cel {
			t.Fatalf("sliding animation = %d, cel = %d; want %d, %d",
				ent.Tag, ent.Cel, tags.P1WalkDown, test.cel)
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

func newP1Test(xy vgeo.XY[float32]) (*Game, P1Ent) {
	gam := NewGame()
	gam.SetBoard(&vboards.Board{
		WH: vgeo.NewWH[int32](128, 128), Tile: vgeo.NewWH[uint8](16, 16),
		Tiles: make([]vboards.Tile, 64),
	})
	gam.Layer(gfx.ZP1.Layer()).Clip = vgeo.XYWH[float32](-16, -16, 176, 176)
	return gam, NewP1Ent(vboards.Spawn{
		XY: xy, WH: vgeo.NewWH[float32](8, 13),
		Z: gfx.ZP1, Tag: tags.P1WalkRight,
	})
}

func stepP1Test(
	t *testing.T, gam *Game, ent *P1Ent, dir vgeo.XY[int8], by float32,
) {
	t.Helper()
	layer := gam.Layer(gfx.ZP1.Layer())
	layer.Sprs = layer.Sprs[:0]
	in := gam.In()
	in.PrevDir, in.Dir = in.Dir, dir
	gam.Poll().DeltaMillis = float64(by/p1Vel) * 1000
	ent.Update(gam)
	if in.Dir != dir {
		t.Fatalf("input direction changed: %v, want %v", in.Dir, dir)
	}
	if len(layer.Sprs) != 1 || layer.Sprs[0].XY != ent.XY {
		t.Fatalf("drawn sprite differs from the resolved position: %v", layer.Sprs)
	}
	box := gam.Atlas().Meta[ent.Tag].Hurtbox.Cast[float32]()
	box.MoveTo(layer.Sprs[0].XY)
	if gam.Board().HitsBox(box) {
		t.Fatalf("drawn sprite overlaps a wall: %v", box)
	}
}
