package game

import (
	"github.com/oidoid/void/src/void/veng"
	"github.com/oidoid/void/src/void/vin"
)

func UpdateFullscreen(gam *Game) veng.Status {
	in := gam.In()
	if in.Ptr.Device() == vin.PtrDevTouch && in.On&^in.PrevOn != 0 {
		gam.ReqFullscreen(veng.FullscreenReqLandscape)
	}
	return veng.Pause
}
