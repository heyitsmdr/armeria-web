import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    isConnected: false,
    gameText: [],
    allowGlobalHotkeys: true,
    minimapData: {},
    characterLocation: { x: 0, y: 0, z: 0 },
    roomObjects: [],
    objectTarget: '',
    objectEditorOpen: false,
    objectEditorData: {},
  },
  mutations: {
    DEBUG_ALTER_STATE: (state, key, val) => {
      state[key] = val;
    },

    SOCKET_ONOPEN: (state) => {
      state.isConnected = true;
    },

    SOCKET_ONCLOSE: (state) => {
      if (state.isConnected) {
        state.isConnected = false;
        state.gameText.push('<br>Connection to the game server has been closed.');
      } else {
        state.gameText.push('A connection to the game server could not be established.');
      }
    },

    SOCKET_ONERROR: () => {
      console.log('an error occurred in the socket connection');
    },

    ADD_GAME_TEXT: (state, text) => {
      state.gameText.push(
        text
          .replace(/\n/g, "<br>")
          .replace(/\[b\]/g, "<span style='font-weight:600'>")
          .replace(/\[\/b\]/g, "</span>")
      );
    },

    SET_ALLOW_GLOBAL_HOTKEYS: (state, allow) => {
      state.allowGlobalHotkeys = allow;
    },

    SET_MINIMAP_DATA: (state, minimapData) => {
      state.minimapData = minimapData;
    },

    SET_CHARACTER_LOCATION: (state, loc) => {
      state.characterLocation = loc;
    },

    SET_ROOM_OBJECTS: (state, objects) => {
      state.roomObjects = objects;
    },

    SET_OBJECT_TARGET: (state, target) => {
      state.objectTarget = target;
    },

    SET_OBJECT_EDITOR_OPEN: (state, open) => {
      state.objectEditorOpen = open;
    },
    SET_OBJECT_EDITOR_DATA: (state, data) => {
      state.objectEditorData = data;
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

    setAllowGlobalHotkeys: ({ commit }, payload) => {
      commit('SET_ALLOW_GLOBAL_HOTKEYS', payload);
    },

    setObjectTarget: ({ commit }, payload) => {
      commit('SET_OBJECT_TARGET', payload.target);
    },

    debugAlterState: ({ commit }, payload) => {
      commit('DEBUG_ALTER_STATE', payload.key, payload.value);
    },

    setObjectEditorOpen: ({ commit }, payload) => {
      commit('SET_OBJECT_EDITOR_OPEN', payload);
    },

    //
    // Server-triggered actions below
    //

    showText: ({ commit }, payload) => {
      commit('ADD_GAME_TEXT', payload.data);
    },

    setMapData: ({ commit }, payload) => {
      commit('SET_MINIMAP_DATA', JSON.parse(payload.data));
    },

    setCharacterLocation: ({ commit }, payload) => {
      commit('SET_CHARACTER_LOCATION', JSON.parse(payload.data));
    },

    setRoomObjects: ({ commit }, payload) => {
      commit('SET_ROOM_OBJECTS', JSON.parse(payload.data));
    },

    setObjectEditorData: ({ commit }, payload) => {
      commit('SET_OBJECT_EDITOR_DATA', JSON.parse(payload.data));
      commit('SET_OBJECT_EDITOR_OPEN', true);
    },

    disconnect: () => {
      Vue.prototype.$socket.close();
    }
  }
})
