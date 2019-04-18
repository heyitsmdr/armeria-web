import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    isConnected: false,
    gameText: [],
  },
  mutations: {
    SOCKET_ONOPEN: (state) => {
      console.log('socket connection is open');
      state.isConnected = true;
    },
    SOCKET_ONCLOSE: (state) => {
      console.log('socket connection is now closed');
      if (state.isConnected) {
        state.isConnected = false;
        state.gameText.push('<br>Connection to the game server has been closed.');
      } else {
        state.gameText.push('A connection to the game server could not be established.');
      }
    },
    SOCKET_ONERROR: (state) => {
      console.log('an error occurred in the socket connection');
    },
    ADD_GAME_TEXT: (state, text) => {
      state.gameText.push(
        text
          .replace(/\n/g, "<br>")
          .replace(/\[b\]/g, "<strong>")
          .replace(/\[\/b\]/g, "</strong>")
      );
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
