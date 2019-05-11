import Vue from 'vue'
import VueNativeSock from 'vue-native-websocket'
import App from './App.vue'
import store from './store'
import Howler from './plugins/Howler'

Vue.config.productionTip = false

let websocketHostname = window.location.hostname;
if (websocketHostname === "localhost") {
  websocketHostname += ":8081";
}

Vue.use(VueNativeSock, `ws://${websocketHostname}/ws`, { store: store, format: 'json' })
Vue.use(Howler)

new Vue({
  store,
  render: h => h(App)
}).$mount('#app')
