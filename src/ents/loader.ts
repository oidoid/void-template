import * as V from '@oidoid/void'
import levelJSON from '../assets/init.level.jsonc' with {type: 'json'}
import {parseComponent} from '../level/level-parser.ts'
import {CamSys} from './cam.ts'
import {DrawSys} from './draw.ts'
import {PizzaSys} from './pizza.ts'

export class LoaderSys implements V.Sys {
  readonly query = 'loader'

  update(ent: V.LoaderEnt, v: V.Void): void {
    switch (ent.loader.level) {
      case undefined: {
        init(ent, v)
        return
      }
      case 'Init':
        return
      default:
        ent.loader.level satisfies never
    }
  }
}

function init(ent: V.LoaderEnt, v: V.Void): void {
  v.zoo.addSystem({
    button: new V.ButtonSys(),
    cam: new CamSys(),
    draw: new DrawSys(),
    cursor: new V.CursorSys(),
    hud: new V.HUDSys(),
    ninePatch: new V.NinePatchSys(),
    override: new V.OverrideSys(),
    pizza: new PizzaSys(),
    sprite: new V.SpriteSys(),
    textWH: new V.TextWHSys(),
    textXY: new V.TextXYSys()
  })
  const zoo = v.loadLevel(levelJSON, 'default', parseComponent)
  v.zoo.add(...zoo.start, ...zoo.default, ...zoo.end)
  ent.loader.level = 'Init'
}
