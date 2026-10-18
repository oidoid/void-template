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
	in := gam.In()
	idle := in.Dir == (vgeo.XY[int8]{})
	by := p1Vel * float32(gam.DeltaSecs())
	this.X += float32(in.Dir.X) * by
	this.Y += float32(in.Dir.Y) * by
	this.X = vmath.Clamp(0, max(float32(gam.Board().W)-float32(this.W), 0), this.X)
	this.Y = vmath.Clamp(0, max(float32(gam.Board().H)-float32(this.H), 0), this.Y)
	tag := this.Tag
	switch {
	case in.Dir.X != 0:
		tag = tags.P1WalkRight
	case in.Dir.Y < 0:
		tag = tags.P1WalkUp
	case in.Dir.Y > 0:
		tag = tags.P1WalkDown
		// to-do: idle tag.
	}
	this.SetTag(tag)
	if !idle {
		this.SetFlipX(in.Dir.X < 0)
	}
	this.Anim.Update(gam.Atlas(), gam.DeltaMs())
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
