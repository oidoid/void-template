import type * as V from '@oidoid/void'

export type PizzaEnt = V.SysEnt<PizzaSys>

export class PizzaSys implements V.Sys {
  readonly query = 'pizza'

  update(ent: PizzaEnt, _v: V.Void): void {
    console.log(`cheese: ${ent.pizza.cheese}`)
  }
}
