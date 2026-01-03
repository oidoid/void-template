import * as V from '@oidoid/void'
import type {Pizza} from '../ents/ent.ts'
import type {PizzaSchema} from './level-schema.ts'

export function parseLevel(
  json: Readonly<V.LevelSchema>,
  pools: Readonly<V.PoolMap>,
  atlas: Readonly<V.Atlas>
): V.Level {
  return V.parseLevel(json, pools, parseComponent, atlas)
}

/** @internal */
export const parseComponent: V.ComponentHook = (_ent, json, k) => {
  if (json[k] == null) throw Error('no component val')
  switch (k) {
    case 'pizza':
      return parsePizza(json[k]) satisfies V.Ent[typeof k]
  }
}

// to-do: unit test.
/** @internal */
export function parsePizza(json: Readonly<PizzaSchema>): Pizza {
  return {cheese: json.cheese ?? false}
}
