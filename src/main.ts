import Vue from 'vue'
// @ts-ignore
import VueNativeSock from 'vue-native-websocket'
import App from './App.vue'
// @ts-ignore
import store from './store'
// @ts-ignore
import SFX from './plugins/SFX'

Vue.config.productionTip = false

let connectionString;
if (process.env.NODE_ENV === "production") {
  connectionString = `wss://${window.location.hostname}/ws`;
} else {
  connectionString = "ws://localhost:8081/ws";
}

Vue.use(VueNativeSock, connectionString, { store: store, format: 'json' })
Vue.use(SFX)

// @ts-ignore
window['Armeria'] = new Vue({
  store,
  render: h => h(App)
}).$mount('#app')
