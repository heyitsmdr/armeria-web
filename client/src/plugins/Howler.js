import { Howl, Howler } from 'howler';

let sounds = {}

export default {
    install: function(Vue) {
        Vue.prototype.$playSound = function(soundFile, volume = 1.0) {
            const soundPath = `sfx/${soundFile}`;
            if (!sounds[soundFile]) {
                sounds[soundFile] = new Howl({
                    src: [soundPath],
                    autoplay: false,
                    loop: false,
                    volume: volume,
                });
            }

            sounds[soundFile].play();
        }
    }
}