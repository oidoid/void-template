package game

import (
	"github.com/oidoid/void-template/src/tags"
	"github.com/oidoid/void/src/void/vatlas"
	"github.com/oidoid/void/src/void/vboards"
	"github.com/oidoid/void/src/void/veng"
	"github.com/oidoid/void/src/void/vgeo"
	"github.com/oidoid/void/src/void/vgfx"
)

const (
	p1Vel     = float32(32)
	p1MaxMove = float32(8)
)

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
	by := min(p1Vel*float32(gam.DeltaSecs()), p1MaxMove)
	if dir.X != 0 {
		next := this.XY
		next.X += float32(dir.X) * by
		if !this.hitsWall(gam, next, tag) {
			this.X = next.X
		}
	}
	if dir.Y != 0 {
		moveTag := tag
		// prefer vertical animation when a wall blocks horizontal movement.
		if this.X == xy.X {
			moveTag = tags.P1WalkDown
			if dir.Y < 0 {
				moveTag = tags.P1WalkUp
			}
		}
		next := this.XY
		next.Y += float32(dir.Y) * by
		if !this.hitsWall(gam, next, moveTag) {
			this.Y = next.Y
			if this.X == xy.X {
				tag = moveTag
				dir.X = 0
			}
		}
	}
	// preserve facing when no movement.
	if this.XY != xy {
		this.SetTag(tag)
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

func (this *P1Ent) hitsWall(
	gam *Game, xy vgeo.XY[float32], tag vatlas.Tag,
) bool {
	box := gam.Atlas().Meta[tag].Hurtbox.Cast[float32]()
	box.MoveTo(xy)
	return gam.Board().HitsBox(box)
}
