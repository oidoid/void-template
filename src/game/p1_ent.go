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
	vgfx.Spr
}

func NewP1Ent(spawn vboards.Spawn) P1Ent {
	return P1Ent{Spr: spawn.Spr()}
}

func (this *P1Ent) Update(gam *Game) veng.Status {
	in := gam.In()
	by := p1Vel * float32(gam.DeltaSecs())
	this.X += float32(in.Dir.X) * by
	this.Y += float32(in.Dir.Y) * by
	this.X = vmath.Clamp(0, max(float32(gam.Board().W)-float32(this.W), 0), this.X)
	this.Y = vmath.Clamp(0, max(float32(gam.Board().H)-float32(this.H), 0), this.Y)
	switch {
	case in.Dir.X < 0:
		this.SetTag(tags.P1WalkRight)
		this.SetFlipX(true)
	case in.Dir.X > 0:
		this.SetTag(tags.P1WalkRight)
		this.SetFlipX(false)
	case in.Dir.Y < 0:
		this.SetTag(tags.P1WalkUp)
		this.SetFlipX(false)
	case in.Dir.Y > 0:
		this.SetTag(tags.P1WalkDown)
		this.SetFlipX(false)
	}
	layer := gam.Layer(this.Z.Layer())
	if layer.Clip.HitsBox(vgeo.XYWH(
		this.X, this.Y, float32(this.W), float32(this.H),
	)) {
		layer.Sprs = append(layer.Sprs, this.Spr)
	}
	if in.Dir != (vgeo.XY[int8]{}) {
		return veng.Loop
	}
	return veng.Pause
}
