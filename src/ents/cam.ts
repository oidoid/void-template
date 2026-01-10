import type * as V from '@oidoid/void'

export class CamHook implements V.Hook {
  readonly query = 'cam'

  update(_ent: V.CamEnt, v: V.Void): void {
    v.cam.update(v.canvas)
  }
}
