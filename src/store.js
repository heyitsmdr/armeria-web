import Vue from 'vue'
import Vuex from 'vuex'
import { Room } from './models';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    isProduction: process.env.NODE_ENV === "production",
    isConnected: false,
    gameText: [],
    allowGlobalHotkeys: true,
    forceInputFocus: { forced: false, text: '' },
    minimapData: { name: '', rooms: [] },
    characterLocation: { x: 0, y: 0, z: 0 },
    roomObjects: [],
    roomTitle: 'Unknown',
    objectTargetUUID: '',
    objectEditorOpen: false,
    objectEditorData: {},
    autoLoginToken: window.localStorage.getItem('auto_login_token') || '',
    inventory: [],
    itemBeingDragged: false,
    permissions: [],
    playerInfo: { uuid: '', name: '' },
    commandHistory: [],
    itemTooltipUUID: '',
    itemTooltipVisible: false,
    itemTooltipCache: [],
    itemTooltipMouseCoords: { x: 0, y: 0 },
    money: '0',
    commandDictionary: [],
    sentKeepAlive: 0,
    pingTime: 0,
    settings: {},
    contextMenuVisible: false,
    contextMenuItems: [],
    contextMenuObjectName: '',
    contextMenuObjectColor: '#fff',
    contextMenuObjectBrackets: true,
    contextMenuPosition: { x: 0, y: 0 },
  },
  getters: {
    itemTooltipCache: (state) => (uuid) => {
      for(let i = 0; i < state.itemTooltipCache.length; i++) {
        const cacheItem = state.itemTooltipCache[i];
        if (cacheItem.uuid === uuid) {
          return cacheItem;
        }
      }
      return null;
    },

    hasPermission: (state) => (permission) => {
      return state.permissions.indexOf(permission) >= 0;
    }
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
        state.gameText.push({ id: state.gameText.length, html: '<br>Connection to the game server has been closed.' });
      } else {
        state.gameText.push({ id: state.gameText.length, html: 'A connection to the game server could not be established.' });
      }
    },

    SOCKET_ONERROR: () => {
      //console.log('an error occurred in the socket connection');
    },

    ADD_GAME_TEXT: (state, text) => {
      state.gameText.push({
        id: state.gameText.length,
        html: text
            .replace(/\n/g, "<br>")
            .replace(/\[b\]/g, "<span style='font-weight:600'>")
            .replace(/\[\/b\]/g, "</span>")
            .replace(/\[cmd=([^\]]*)\]/g, "<a href='#' class='inline-command' onclick='window.Armeria.$store.dispatch(\"sendSlashCommand\", {command:\"$1\"})'>")
            .replace(/\[\/cmd\]/g, "</a>")
      });
    },

    SET_ALLOW_GLOBAL_HOTKEYS: (state, allow) => {
      state.allowGlobalHotkeys = allow;
    },

    SET_MINIMAP_DATA: (state, minimapData) => {
      state.minimapData = {
        name: minimapData.name,
        rooms: [],
      };
      minimapData.rooms.forEach(r => {
        state.minimapData.rooms.push(new Room(r));
      });
    },

    SET_CHARACTER_LOCATION: (state, loc) => {
      state.characterLocation = loc;
    },

    SET_ROOM_OBJECTS: (state, objects) => {
      state.roomObjects = objects;
    },

    SET_ROOM_TITLE: (state, title) => {
      state.roomTitle = title;
    },

    SET_OBJECT_TARGET: (state, targetUUID) => {
      state.objectTargetUUID = targetUUID;
    },

    SET_OBJECT_EDITOR_OPEN: (state, open) => {
      state.objectEditorOpen = open;
    },

    SET_OBJECT_EDITOR_DATA: (state, data) => {
      state.objectEditorData = data;
    },

    SET_FORCE_INPUT_FOCUS: (state, data) => {
      state.forceInputFocus = data;
    },

    SET_AUTOLOGIN_TOKEN: (state, token) => {
      state.autoLoginToken = token;
      window.localStorage.setItem('auto_login_token', token);
      if (token.length > 0) {
        state.gameText.push({ id: state.gameText.length, html: '<br>You will now be automatically logged in to this character.' });
      } else {
        state.gameText.push({ id: state.gameText.length, html: '<br>You will no longer be automatically logged in to this character.' });
      }
    },

    SET_INVENTORY: (state, inventory) => {
      state.inventory = inventory;
    },

    SET_ITEM_BEING_DRAGGED: (state, being_dragged) => {
      state.itemBeingDragged = being_dragged;
    },

    SET_PERMISSIONS: (state, permissions) => {
      state.permissions = permissions.split(' ')
    },

    SET_PLAYER_INFO: (state, playerInfo) => {
      state.playerInfo = playerInfo;
    },

    APPEND_COMMAND_HISTORY: (state, command) => {
      state.commandHistory.push(command);
    },

    SET_ITEM_TOOLTIP_UUID: (state, uuid) => {
      state.itemTooltipUUID = uuid;
      state.itemTooltipVisible = true;
    },

    HIDE_ITEM_TOOLTIP: (state) => {
      state.itemTooltipVisible = false;
      state.itemTooltipUUID = '';
    },

    SET_ITEM_TOOLTIP_HTML: (state, data) => {
      for(let i = 0; i < state.itemTooltipCache.length; i++) {
        let cacheItem = state.itemTooltipCache[i];
        if (cacheItem.uuid === data.uuid) {
          cacheItem = data;
          return;
        }
      }

      state.itemTooltipCache.push(data);
    },

    SET_ITEM_TOOLTIP_MOUSE_COORDS: (state, mouseCoords) => {
      state.itemTooltipMouseCoords = mouseCoords;
    },

    CLEAR_ITEM_TOOLTIP_CACHE: (state) => {
      state.itemTooltipCache = [];
    },

    SET_MONEY: (state, money) => {
      state.money = money;
    },

    SET_COMMAND_DICTIONARY: (state, dictionary) => {
      state.commandDictionary = dictionary;
    },

    KEEP_ALIVE_RESPONSE: (state) => {
      state.pingTime = Date.now() - state.sentKeepAlive;
    },

    SET_SETTINGS: (state, settings) => {
      state.settings = settings;
    },

    SET_CONTEXT_MENU_VISIBLE: (state, visible) => {
      state.contextMenuVisible = visible;
    },

    SET_CONTEXT_MENU_ITEMS: (state, items) => {
        state.contextMenuItems = items;
    },

    SET_CONTEXT_MENU_OBJECT: (state, obj) => {
      state.contextMenuObjectName = obj.name;
      state.contextMenuObjectColor = obj.color;
      state.contextMenuObjectBrackets = (typeof obj.subjectBrackets !== 'undefined') ? obj.subjectBrackets : true;
    },

    SET_CONTEXT_MENU_POSITION: (state, pos) => {
      state.contextMenuPosition = { x: pos.x, y: pos.y };
    }
  },
  actions: {
    sendSlashCommand: ({ state , commit }, payload) => {
      if (!state.isConnected) {
        return;
      }

      commit('APPEND_COMMAND_HISTORY', payload.command);

      let echoCmd = payload.command;
      if (payload.command.indexOf('logintoken') === 1) {
        // Hide the actual token from the command echo'd to the main text area.
        echoCmd = `${payload.command.split(':')[0]}:&lt;redacted&gt;`
      }
      commit('ADD_GAME_TEXT', `<div class="inline-loopback">${echoCmd}</div>`);

      Vue.prototype.$socket.sendObj({
        type: "command",
        payload: payload.command
      });
    },

    sendKeepAlive: ({ state }) => {
      state.sentKeepAlive = Date.now();
      Vue.prototype.$socket.sendObj({
        type: "ping",
      });
    },

    setAllowGlobalHotkeys: ({ commit }, payload) => {
      commit('SET_ALLOW_GLOBAL_HOTKEYS', payload);
    },

    setObjectTarget: ({ commit }, payload) => {
      commit('SET_OBJECT_TARGET', payload);
    },

    debugAlterState: ({ commit }, payload) => {
      commit('DEBUG_ALTER_STATE', payload.key, payload.value);
    },

    setObjectEditorOpen: ({ commit }, payload) => {
      if (payload) {
        commit('SET_OBJECT_EDITOR_OPEN', true);
      } else {
        commit('SET_OBJECT_EDITOR_OPEN', false);
        commit('SET_OBJECT_EDITOR_DATA', {});
      }
    },

    setForceInputFocus: ({ commit }, payload) => {
      commit('SET_FORCE_INPUT_FOCUS', payload);
    },

    setItemBeingDragged: ({ commit }, payload) => {
      commit('SET_ITEM_BEING_DRAGGED', payload);
    },

    showItemTooltip: ({ commit }, payload) => {
      commit('SET_ITEM_TOOLTIP_UUID', payload);
    },

    hideItemTooltip: ({ commit }) => {
      commit('HIDE_ITEM_TOOLTIP');
    },

    moveItemTooltip: ({ commit }, payload) => {
      commit('SET_ITEM_TOOLTIP_MOUSE_COORDS', payload);
    },

    clearItemTooltipCache: ({ commit }) => {
      commit('CLEAR_ITEM_TOOLTIP_CACHE');
    },

    showContextMenu: ({ commit }, payload) => {
      commit('SET_CONTEXT_MENU_ITEMS', payload.items);
      commit('SET_CONTEXT_MENU_OBJECT', payload.object);
      commit('SET_CONTEXT_MENU_POSITION', payload.at);
      commit('SET_CONTEXT_MENU_VISIBLE', true);
    },

    hideContextMenu: ({ commit }) => {
      commit('SET_CONTEXT_MENU_VISIBLE', false);
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

    setRoomTitle: ({ commit }, payload) => {
      commit('SET_ROOM_TITLE', payload.data);
    },

    setObjectEditorData: ({ commit }, payload) => {
      commit('SET_OBJECT_EDITOR_DATA', JSON.parse(payload.data));
      commit('SET_OBJECT_EDITOR_OPEN', true);
    },

    disconnect: () => {
      Vue.prototype.$socket.close();
    },

    pong: ({ commit }) => {
      commit('KEEP_ALIVE_RESPONSE');
    },

    toggleAutoLogin: ({ state, commit }, payload) => {
      if (state.autoLoginToken !== '') {
        commit('SET_AUTOLOGIN_TOKEN', '');
      } else {
        commit('SET_AUTOLOGIN_TOKEN', payload.data);
      }
    },

    setInventory: ({ commit }, payload) => {
      commit('SET_INVENTORY', JSON.parse(payload.data) || []);
    },

    setPermissions: ({ commit }, payload) => {
      commit('SET_PERMISSIONS', payload.data);
    },

    setPlayerInfo: ({ commit }, payload) => {
      commit('SET_PLAYER_INFO', JSON.parse(payload.data));
    },

    setItemTooltipHTML: ({ commit }, payload) => {
      commit('SET_ITEM_TOOLTIP_HTML', JSON.parse(payload.data));
    },

    setCommandDictionary: ({ commit }, payload) => {
      commit('SET_COMMAND_DICTIONARY', JSON.parse(payload.data));
    },

    setMoney: ({ commit }, payload) => {
      commit('SET_MONEY', payload.data);
    },

    playSFX: (_, payload) => {
      const sfx = JSON.parse(payload.data);
      Vue.prototype.$soundEvent(sfx.id, sfx.volume);
    },

    setSettings: ({ commit }, payload) => {
      commit('SET_SETTINGS', JSON.parse(payload.data));
    },
  }
})
