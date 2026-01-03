import * as V from '@oidoid/void'
import {description} from './assets/manifest.json' // non-standard import to treeshake.
import config from './assets/void.game.json' with {type: 'json'}
import {LoaderSys} from './ents/loader.ts'

console.debug(
  `void-template v${V.bundle.version}+${V.bundle.published}.${V.bundle.hash} ───oidoid>°──`
)

const v = new V.Void({
  config: config as V.GameConfig,
  description,
  preloadAtlas: document.querySelector<HTMLImageElement>('#atlas'),
  loader: {loader: {level: undefined}},
  loaderSys: new LoaderSys()
})
await v.register('add')
