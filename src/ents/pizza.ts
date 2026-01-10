import type * as V from '@oidoid/void'

export type PizzaEnt = V.HookEnt<PizzaHook>

export class PizzaHook implements V.Hook {
  readonly query = 'pizza'

  update(ent: PizzaEnt, _v: V.Void): void {
    console.log(`cheese: ${ent.pizza.cheese}`)
  }
}
