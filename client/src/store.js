import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    isConnected: false,
    gameText: [],
  },
  mutations: {
    SOCKET_ONOPEN: (state, event) => {
      console.log('socket connection is open')
      state.isConnected = true;
    },
    SOCKET_ONCLOSE: (state, event) => {
      console.log('socket connection is now closed')
      state.isConnected = false;
    },
    SOCKET_ONERROR: (state, event) => {
      console.log('socket connection is now closed due to an error')
      state.isConnected = false;
    },
    ADD_GAME_TEXT: (state, text) => {
      state.gameText.push(text);
    }
  },
  actions: {
    sendSlashCommand: ({ state }, payload) => {
      if (!state.isConnected) {
        return;
      }

      Vue.prototype.$socket.sendObj({
        type: "command",
        payload: payload.command
      });
    },

    showText: ({ commit }, payload) => {
      commit('ADD_GAME_TEXT', payload.data);
    }
  }
})
