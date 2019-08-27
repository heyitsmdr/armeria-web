import { Howl } from 'howler';

let sfxCache = {};

export const INVENTORY_DRAG_START = 'INVENTORY_DRAG_START';
export const INVENTORY_DRAG_STOP = 'INVENTORY_DRAG_STOP';
export const PICKUP_ITEM = 'PICKUP_ITEM';

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

        Vue.prototype.$soundEvent = function(event) {
            switch(event) {
                case INVENTORY_DRAG_START:
                    this.$playSound('mouse-click.wav', 1.0);
                    break;
                case INVENTORY_DRAG_STOP:
                    this.$playSound('mouse-release.wav', 1.0);
                    break;
                case PICKUP_ITEM:
                    this.$playSound('pickup.wav', 1.0);
                    break;
            }
        };
    }
}