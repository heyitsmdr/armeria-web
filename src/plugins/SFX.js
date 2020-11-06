import { Howl } from 'howler';

let sfxCache = {};

// NOTE: Be sure to add new entries here to internal/pkg/sfx/sfx.go as well.
export const INVENTORY_DRAG_START = 'INVENTORY_DRAG_START';
export const INVENTORY_DRAG_STOP = 'INVENTORY_DRAG_STOP';
export const PICKUP_ITEM = 'PICKUP_ITEM';
export const SELL_BUY_ITEM = 'SELL_BUY_ITEM';

export default {
    install: function(Vue) {
        Vue.prototype.$playSound = function(soundFile, volume = 1.0) {
            const soundPath = `sfx/${soundFile}`;
            if (!sfxCache[soundFile]) {
                sfxCache[soundFile] = new Howl({
                    src: [soundPath],
                    autoplay: false,
                    loop: false,
                    volume: volume,
                });
            }

            sfxCache[soundFile].play();
        };

        Vue.prototype.$soundEvent = function(event, volume = 1.0) {
            switch(event) {
                case INVENTORY_DRAG_START:
                    this.$playSound('mouse-click.wav', volume);
                    break;
                case INVENTORY_DRAG_STOP:
                    this.$playSound('mouse-release.wav', volume);
                    break;
                case PICKUP_ITEM:
                    this.$playSound('pickup.wav', volume);
                    break;
                case SELL_BUY_ITEM:
                    this.$playSound('sell-buy-item.wav', volume);
                    break;
            }
        };
    }
}