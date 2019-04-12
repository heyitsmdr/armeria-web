import Vue from 'vue'
import VueNativeSock from 'vue-native-websocket'
import App from './App.vue'
import store from './store'

Vue.config.productionTip = false

Vue.use(VueNativeSock, 'ws://localhost:8081/ws')

new Vue({
  store,
  render: h => h(App)
}).$mount('#app')
