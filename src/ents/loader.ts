import type * as V from '@oidoid/void'
import levelJSON from '../assets/init.level.jsonc' with {type: 'json'}
import {parseLevel} from '../level/level-parser.ts'
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
  v.zoo.addDefaultSystems()
  const level = parseLevel(levelJSON, v.pool, v.preload)
  v.zoo.add(...level.ents, {draw: {}}) // to-do: support ent ordering. can't put this in level file right now.
  v.zoo.addSystem({
    cam: new CamSys(),
    draw: new DrawSys(),
    pizza: new PizzaSys()
  })
  // to-do: validate all ents on a system add.
  ent.loader.level = 'Init'
}
