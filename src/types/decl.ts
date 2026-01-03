// biome-ignore lint/correctness/useJsonImportAttributes:;
import type gameJSON from '../assets/void.game.json'
import type {Pizza} from '../ents/ent.ts'
import type {PizzaSchema} from '../level/level-schema.ts'

declare module '@oidoid/void' {
  interface Ent {
    pizza?: Pizza
  }

  interface EntSchema {
    pizza?: PizzaSchema
  }

  interface Loader {
    level: undefined | 'Init'
  }

  interface ReturnTag {
    // biome-ignore lint/style/useShorthandFunctionType:;
    (): keyof typeof gameJSON.atlas.anim
  }

  interface Void {}
}
