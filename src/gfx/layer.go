package gfx

import "github.com/oidoid/void/src/void/vgfx"

const (
	LayerTiles vgfx.Layer = iota
	LayerP1
	LayerUI
	LayerCursor
)

var (
	ZP1       vgfx.Z = LayerP1.Z(0)
	ZUIWidget vgfx.Z = LayerUI.Z(0)
	ZUIText   vgfx.Z = LayerUI.Z(1)
	ZCursor   vgfx.Z = LayerCursor.Z(0)
)
