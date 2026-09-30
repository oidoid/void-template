// biome-ignore lint/correctness/noUndeclaredDependencies: installed via `link-void`.
import * as V from '@oidoid/void'
import wasm from '../../dist/index.wasm'

const eng = new V.Eng()
await eng.load(undefined, wasm, 0xe6e6e6ff)
eng.register()
