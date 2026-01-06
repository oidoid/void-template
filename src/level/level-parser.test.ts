import {test} from 'node:test'
import {assert} from '@oidoid/void/test'
import {parsePizza} from './level-parser.ts'

test('parsePizza()', () => {
  assert(parsePizza({}), {cheese: false})
  assert(parsePizza({cheese: true}), {cheese: true})
})
