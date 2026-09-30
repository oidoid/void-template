// codegen by packboards.
package boards

import (
	"github.com/oidoid/void/src/void/vboards"
)

type P1Spawn = vboards.Spawn

type CursorSpawn struct {
	vboards.Spawn
	KbdVel int32
}

type TextSpawn = vboards.TextSpawn
