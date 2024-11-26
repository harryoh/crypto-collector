const state = {
  currencyRate: {
    fixExchangeRate: 1300,
    exchangeRate: 1300
  },
  bybitPrice: {
    BTC: 0,
    ETH: 0,
    EOS: 0,
    XRP: 0
  },
  upbitPrice: {
    BTC: 0,
    ETH: 0,
    EOS: 0,
    XRP: 0
  },
  bithumbPrice: {
    BTC: 0,
    ETH: 0,
    EOS: 0,
    XRP: 0
  }
}

const mutations = {
  SET_PRICE: (state, { key, value }) => {
    // eslint-disable-next-line no-prototype-builtins
    if (state.hasOwnProperty(key)) {
      if (value) {
        state[key] = value
      }
    }
  },
  SET_COIN: (state, { key, coin, value }) => {
    // eslint-disable-next-line no-prototype-builtins
    if (state.hasOwnProperty(key)) {
      state[key][coin] = value
    }
  },
  SET_CURRENCY: (state, value) => {
    // state[key] = Object.assign(state[key], value)
    value.fixExchangeRate = state.currencyRate.fixExchangeRate
    value.Timestamp = new Date().getTime() / 1000
    state.currencyRate = value
  }
}

const actions = {
  setPrice ({ commit }, data) {
    commit('SET_PRICE', data)
  },
  setCoin ({ commit }, data) {
    commit('SET_COIN', data)
  },
  setCurrency ({ commit }, data) {
    commit('SET_CURRENCY', data)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
