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
	xy vgeo.XY[float32] // layer px.
}

func NewP1Ent(spawn vboards.Spawn) P1Ent {
	return P1Ent{Spr: spawn.Spr()}
}

func (this *P1Ent) Update(gam *Game) veng.Status {
	in := gam.In()
	dir := in.Dir
	idle := dir == (vgeo.XY[int8]{})
	xy := vgeo.RoundXY(this.XY)
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
	step := min(p1Vel*float32(gam.DeltaSecs()), p1MaxMove)
	by := vgeo.NewXY(float32(dir.X)*step, float32(dir.Y)*step)
	if dir.X != in.PrevDir.X {
		this.xy.X = 0
	}
	if dir.Y != in.PrevDir.Y {
		this.xy.Y = 0
	}
	this.xy.AddTo(by)
	if dir.X != 0 {
		probe := vgfx.SnapMove(xy.Add(this.xy), xy, vgeo.NewXY(by.X, float32(0)))
		if probe.X == xy.X {
			probe.X += float32(dir.X)
		}
		if this.hitsWall(gam, probe, tag) {
			by.X, this.xy.X = 0, 0
		}
	}
	next := vgfx.SnapMove(xy.Add(this.xy), xy, by)
	if dir.Y != 0 {
		moveTag := tag
		// prefer vertical animation when a wall blocks horizontal movement.
		if by.X == 0 {
			moveTag = tags.P1WalkDown
			if dir.Y < 0 {
				moveTag = tags.P1WalkUp
			}
		}
		probe := vgfx.SnapMove(xy.Add(this.xy), xy, vgeo.NewXY(float32(0), by.Y))
		probe.X = next.X
		if probe.Y == xy.Y {
			probe.Y += float32(dir.Y)
		}
		if this.hitsWall(gam, probe, moveTag) {
			by.Y, this.xy.Y = 0, 0
			next = vgfx.SnapMove(xy.Add(this.xy), xy, by)
		} else if by.X == 0 {
			tag = moveTag
			dir.X = 0
		}
	}
	this.xy = this.xy.Sub(next.Sub(xy))
	this.XY = next
	// single-axis input chooses facing; blocked diagonals keep the last facing.
	if this.XY != xy || !idle && (dir.X == 0 || dir.Y == 0) {
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
	box.MoveTo(vgeo.RoundXY(xy))
	return gam.Board().HitsBox(box)
}
