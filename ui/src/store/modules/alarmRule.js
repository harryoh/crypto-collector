const state = {
  alarmRules: []
}

const mutations = {
  SET_RULES: (state, rules) => {
    state.alarmRules = rules
  }
}

const actions = {
  setRules ({ commit }, data) {
    commit('SET_RULES', data)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
