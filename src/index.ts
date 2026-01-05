import * as V from '@oidoid/void'
import {description} from './assets/manifest.json' // non-standard import to treeshake.
import config from './assets/void.game.json' with {type: 'json'}
import {LoaderSys} from './ents/loader.ts'

console.debug(
  `void-template v${V.bundle.version}+${V.bundle.published}.${V.bundle.hash} ───oidoid>°──`
)

const v = new V.Void({
  atlas: document.querySelector<HTMLImageElement>('#atlas'),
  config: config as V.VoidConfig,
  description,
  loader: {loader: {level: undefined}},
  loaderSys: new LoaderSys()
})
await v.register('add')
