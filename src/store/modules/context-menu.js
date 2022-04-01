export default {
    namespaced: true,
    state() {
        return {
            visible: false,
            items: [],
            objectName: '',
            objectColor: '#fff',
            objectBrackets: true,
            position: { x: 0, y: 0 },
        }
    },
    mutations: {
        SET_CONTEXT_MENU_VISIBLE: (state, { visible }) => {
            state.visible = visible;
        },

        SET_CONTEXT_MENU_ITEMS: (state, { items }) => {
            state.items = items;
        },

        SET_CONTEXT_MENU_OBJECT: (state, { name, color, subjectBrackets }) => {
            state.objectName = name;
            state.objectColor = color;
            state.objectBrackets = (typeof subjectBrackets !== 'undefined') ? subjectBrackets : true;
        },

        SET_CONTEXT_MENU_POSITION: (state, { pos }) => {
            state.position = { x: pos.x, y: pos.y };
        }
    },
    actions: {
        show: ({ commit }, payload) => {
            commit('SET_CONTEXT_MENU_ITEMS', { items: payload.items });
            commit('SET_CONTEXT_MENU_OBJECT', { name: payload.object.name, color: payload.object.color, subjectBrackets: payload.object.subjectBrackets });
            commit('SET_CONTEXT_MENU_POSITION', { pos: payload.at });
            commit('SET_CONTEXT_MENU_VISIBLE', { visible: true });
        },

        hide: ({ commit }) => {
            commit('SET_CONTEXT_MENU_VISIBLE', { visible: false });
        },
    }
}