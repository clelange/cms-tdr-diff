// export const state = () => ({})
export const state = () => ({
  tdrTypes: [],
  apiStatus: null
})

export const actions = {
  async loadTdr({ state, commit }) {
    console.log('computed')
    if (state.tdrTypes.length) return
    commit('setApiStatus', null)
    await this.$axios.$get('/types').then(
      tdrTypes => {
        commit('setApiStatus', 'good')
        commit('updateTdr', tdrTypes)
      },
      (err) => {
        console.log(err)
        commit('setApiStatus', 'bad')
      }
    )
  }
}

export const mutations = {
  updateTdr(state, payload) {
    state.tdrTypes = payload
  },
  setApiStatus(state, status) {
    state.apiStatus = status
  }
}
