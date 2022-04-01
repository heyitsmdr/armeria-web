import { createApp } from 'vue'
// @ts-ignore
import VueNativeSock from 'vue-native-websocket-vue3'
// @ts-ignore
import VueAnimXYZ from '@animxyz/vue'
import App from './App.vue'
// @ts-ignore
import { store } from './store'
// @ts-ignore
import SFX from './plugins/SFX'

let connectionString;
if (process.env.NODE_ENV === "production") {
  connectionString = `wss://${window.location.hostname}/ws`;
} else {
  connectionString = `ws://${window.location.hostname}:8081/ws`;
}

export const app = createApp(App);

app.use(VueNativeSock, connectionString, { store: store, format: 'json' })
app.use(SFX);
app.use(VueAnimXYZ);
app.use(store);

app.mount('#app');
