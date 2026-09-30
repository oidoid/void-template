package game

import (
	"testing"

	"github.com/oidoid/void-template/src/boards"
	"github.com/oidoid/void-template/src/gfx"
	"github.com/oidoid/void/src/void/veng"
	"github.com/oidoid/void/src/void/vgeo"
	"github.com/oidoid/void/src/void/vin"
)

func TestTouchInputRequestsLandscape(t *testing.T) {
	tests := []struct {
		name   string
		device vin.PtrDevice
		want   veng.FullscreenReq
	}{
		{name: "mouse", device: vin.PtrDevMouse, want: veng.FullscreenReqNone},
		{name: "touch", device: vin.PtrDevTouch, want: veng.FullscreenReqLandscape},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gam := NewGame()
			poll := &gam.Poll().InPoll
			poll.PtrsLen = 1
			poll.Ptrs[0] = vin.PtrPoll{
				Device: test.device, Primary: true, Clicks: vin.ClickPrimary,
			}
			gam.Update()
			if got := gam.FullscreenReq(); got != int32(test.want) {
				t.Fatalf("fullscreen request = %d, want %d", got, test.want)
			}
		})
	}
}

func TestNewPopulatesTiledCursor(t *testing.T) {
	gam := NewGame()
	spawn := boards.TitleCursorSpawns[0]
	cursor := gam.Cursor()
	if cursor == nil {
		t.Fatal("cursor is nil")
	}
	if cursor != &gam.cursor.CursorEnt {
		t.Fatal("engine cursor does not reference the persistent game cursor")
	}
	if got, want := cursor.Spr, spawn.Spawn.Spr(); got != want {
		t.Errorf("cursor sprite = %#v, want %#v", got, want)
	}
	if got, want := cursor.KbdVel, float32(spawn.KbdVel); got != want {
		t.Errorf("cursor keyboard velocity = %v, want %v", got, want)
	}
	poll := &gam.Poll().InPoll
	poll.PtrsLen = 1
	poll.Ptrs[0] = vin.PtrPoll{
		ID: 1, Device: vin.PtrDevMouse, Primary: true,
		Phy: vgeo.XYWH[float32](8, 16, 1, 1),
	}
	gam.Update()
	if got := len(gam.Layer(gfx.LayerCursor).Sprs); got != 1 {
		t.Fatalf("cursor sprite count = %d, want 1", got)
	}
}

func TestRouteLoadsLvl0(t *testing.T) {
	gam := NewGame()
	if gam.Board() != &boards.TitleBoard {
		t.Fatal("app did not initialize the title level")
	}
	if gam.BoardLevel() != uint16(boards.TitleLevel) {
		t.Fatalf("board level = %d, want title level %d", gam.BoardLevel(), boards.TitleLevel)
	}
	if gam.BoardTilesLen() != uint32(len(boards.TitleBoard.Tiles)) {
		t.Fatalf("exported title tile count = %d, want %d", gam.BoardTilesLen(), len(boards.TitleBoard.Tiles))
	}

	routeLvl0(gam)
	if gam.Board() != &boards.Lvl0Board {
		t.Fatal("route did not load lvl0")
	}
	if gam.BoardLevel() != uint16(boards.Lvl0Level) {
		t.Fatalf("board level = %d, want lvl0 level %d", gam.BoardLevel(), boards.Lvl0Level)
	}
	if gam.BoardTilesLen() != uint32(len(boards.Lvl0Board.Tiles)) {
		t.Fatalf("exported lvl0 tile count = %d, want %d", gam.BoardTilesLen(), len(boards.Lvl0Board.Tiles))
	}
	if boards.TitleLevel == boards.Lvl0Level {
		t.Fatal("generated boards have the same level ID")
	}
}
