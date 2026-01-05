import * as V from '@oidoid/void'
import levelJSON from '../assets/init.level.jsonc' with {type: 'json'}
import {CamSys} from '../ents/cam.ts'
import {DrawSys} from '../ents/draw.ts'
import {PizzaSys} from '../ents/pizza.ts'
import {parseComponent} from './level-parser.ts'

export class Loader implements V.Loader {
  cursor: V.CursorEnt | undefined
  #lvl: 'Init' | undefined
  readonly #systems: V.SysMap = {
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
  }
  #zoo: V.Zoo = {start: new Set(), default: new Set(), end: new Set()}

  update(v: V.Void): void {
    switch (this.#lvl) {
      case undefined: {
        this.#init(v)
        break
      }
      case 'Init':
        break
      default:
        this.#lvl satisfies never
    }

    for (const zoo of Object.values(this.#zoo))
      V.zooUpdate(zoo, this.#systems, v)
  }

  #init(v: V.Void): void {
    this.#zoo = v.loadLevel(levelJSON, 'default', parseComponent)
    this.cursor = V.zooFindByID(this.#zoo.default, 'Cursor')
    this.#lvl = 'Init'
  }
}
