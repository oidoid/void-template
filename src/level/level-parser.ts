import type * as V from '@oidoid/void'
import type {Pizza} from '../ents/ent.ts'
import type {PizzaSchema} from './level-schema.ts'

export const parseComponent: V.ComponentHook = (_ent, json, k) => {
  if (json[k] == null) throw Error('no component val')
  switch (k) {
    case 'pizza':
      return parsePizza(json[k]) satisfies V.Ent[typeof k]
  }
}

/** @internal */
export function parsePizza(json: Readonly<PizzaSchema>): Pizza {
  return {cheese: json.cheese ?? false}
}
