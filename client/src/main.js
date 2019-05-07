import Vue from 'vue'
import VueNativeSock from 'vue-native-websocket'
import App from './App.vue'
import store from './store'
import Howler from './plugins/Howler'

Vue.config.productionTip = false

Vue.use(VueNativeSock, 'ws://localhost:8081/ws', { store: store, format: 'json' })
Vue.use(Howler)

new Vue({
  store,
  render: h => h(App)
}).$mount('#app')
