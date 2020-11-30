import { Howl } from 'howler';

let sfxCache = {};

// NOTE: Be sure to add new entries here to internal/pkg/sfx/sfx.go as well.
export const INVENTORY_DRAG_START = 'INVENTORY_DRAG_START';
export const INVENTORY_DRAG_STOP = 'INVENTORY_DRAG_STOP';
export const PICKUP_ITEM = 'PICKUP_ITEM';
export const SELL_BUY_ITEM = 'SELL_BUY_ITEM';
export const CAT_MEOW = 'CAT_MEOW';

function preload(soundFile) {
    const sfx = new Howl({
        src: [soundFile],
        autoplay: false,
        loop: false,
        preload: true,
    });

    sfx.once('load', () => {
        console.log(`Sound ready: ${soundFile}`);
    });

    return sfx;
}

export default {
    install: function(Vue) {
        // Preload.
        sfxCache[INVENTORY_DRAG_START] = preload('sfx/mouse-click.wav');
        sfxCache[INVENTORY_DRAG_STOP] = preload('sfx/mouse-release.wav');
        sfxCache[PICKUP_ITEM] = preload('sfx/pickup.wav');
        sfxCache[SELL_BUY_ITEM] = preload('sfx/sell-buy-item.wav');
        sfxCache[CAT_MEOW] = preload('sfx/cat_meow.wav');

        Vue.prototype.$soundEvent = function(event) {
            sfxCache[event].play();
        };
    }
}