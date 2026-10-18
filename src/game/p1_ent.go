package game

import (
	"github.com/oidoid/void-template/src/tags"
	"github.com/oidoid/void/src/void/vboards"
	"github.com/oidoid/void/src/void/veng"
	"github.com/oidoid/void/src/void/vgeo"
	"github.com/oidoid/void/src/void/vgfx"
	"github.com/oidoid/void/src/void/vmath"
)

const p1Vel = float32(32)

type P1Ent struct {
	vgfx.Anim
}

func NewP1Ent(spawn vboards.Spawn) P1Ent {
	return P1Ent{Anim: vgfx.Anim{Spr: spawn.Spr()}}
}

func (this *P1Ent) Update(gam *Game) veng.Status {
	dir := gam.In().Dir
	idle := dir == (vgeo.XY[int8]{})
	xy := this.XY
	by := p1Vel * float32(gam.DeltaSecs())
	this.X += float32(dir.X) * by
	this.Y += float32(dir.Y) * by
	this.X = vmath.Clamp(0, max(float32(gam.Board().W)-float32(this.W), 0), this.X)
	this.Y = vmath.Clamp(0, max(float32(gam.Board().H)-float32(this.H), 0), this.Y)
	// prefer vertical animation when a board edge blocks horizontal movement.
	if this.X == xy.X && this.Y != xy.Y {
		dir.X = 0
	}
	tag := this.Tag
	switch {
	case dir.X != 0:
		tag = tags.P1WalkRight
	case dir.Y < 0:
		tag = tags.P1WalkUp
	case dir.Y > 0:
		tag = tags.P1WalkDown
		// to-do: idle tag.
	}
	this.SetTag(tag)
	if !idle {
		this.SetFlipX(dir.X < 0)
	}
	delta := gam.DeltaMs()
	if idle {
		delta = 0
	}
	this.Anim.Update(gam.Atlas(), delta)
	layer := gam.Layer(this.Z.Layer())
	if layer.Clip.HitsBox(vgeo.XYWH(
		this.X, this.Y, float32(this.W), float32(this.H),
	)) {
		layer.Sprs = append(layer.Sprs, this.Spr)
	}
	if !idle {
		return veng.Loop
	}
	return veng.Pause
}
