import { Howl } from 'howler';

let sfxCache = {};

// NOTE: Be sure to add new entries here to internal/pkg/sfx/sfx.go as well.
export const INVENTORY_DRAG_START = 'INVENTORY_DRAG_START';
export const INVENTORY_DRAG_STOP = 'INVENTORY_DRAG_STOP';
export const PICKUP_ITEM = 'PICKUP_ITEM';
export const SELL_BUY_ITEM = 'SELL_BUY_ITEM';
export const CAT_MEOW = 'CAT_MEOW';
export const TELEPORT = 'TELEPORT';

function preload(soundFile, volume) {
    const sfx = new Howl({
        src: [soundFile],
        autoplay: false,
        loop: false,
        preload: true,
        volume: volume,
    });

    sfx.once('load', () => {
        console.log(`Sound preloaded: ${soundFile}`);
    });

    return sfx;
}

export default {
    install: function(app) {
        // Preload.
        const preloadMap = {
            INVENTORY_DRAG_START: ['sfx/mouse-click.wav', 1],
            INVENTORY_DRAG_STOP: ['sfx/mouse-release.wav', 1],
            PICKUP_ITEM: ['sfx/pickup.wav', 1],
            SELL_BUY_ITEM: ['sfx/sell-buy-item.wav', 1],
            CAT_MEOW: ['sfx/cat_meow.wav', 1],
            TELEPORT: ['sfx/teleport.wav', 0.1],
        };

        Object.keys(preloadMap).forEach(k => {
            sfxCache[k] = preload(preloadMap[k][0], preloadMap[k][1]);
        });

        app.config.globalProperties.$soundEvent = (event) => {
            sfxCache[event].play();
        }
    }
}
