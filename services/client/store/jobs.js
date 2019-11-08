export const state = () => ({
  pipelineStatus: []
})

export const mutations = {
  addPipeline: (state, { jobStatus, pipelineId }) => {
    state.pipelineStatus.push({
      pipelineId: pipelineId,
      jobStatus: jobStatus.job_status
    })
  },
  updatePipeline: (state, { jobStatus, currentPipelineId }) => {
    const myJobStatus = jobStatus.job_status
    const index = state.pipelineStatus.findIndex(
      p => p.pipelineId == Number.parseInt(currentPipelineId)
    )
    if (index < 0) {
      return
    }
    state.pipelineStatus[index]['jobStatus'] = myJobStatus
  }
}

export const actions = {
  async load({ state, commit }, pipelineId) {
    const index = state.pipelineStatus.findIndex(
      p => p.pipelineId == pipelineId
    )
    if (index >= 0) {
      console.log('Pipeline already in store:', pipelineId)
      return
    }
    commit('setApiStatus', null, { root: true })
    this.$axios.setToken(this.$env.REQUEST_TOKEN)
    await this.$axios.$get('/status/pipeline/' + pipelineId).then(
      jobStatus => {
        commit('setApiStatus', 'good', { root: true })
        commit('addPipeline', { jobStatus, pipelineId })
      },
      err => {
        console.log('error')
        console.log(err)
        commit('setApiStatus', 'bad', { root: true })
      }
    )
  },
  async update({ state, commit }) {
    if (state.pipelineStatus.length < 1) {
      console.log('No pipelines found')
      return
    }
    for (var i in state.pipelineStatus) {
      const currentPipelineId = Number.parseInt(
        state.pipelineStatus[i].pipelineId
      )
      commit('setApiStatus', null, { root: true })
      this.$axios.setToken(this.$env.REQUEST_TOKEN)
      await this.$axios.$get('/status/pipeline/' + currentPipelineId).then(
        jobStatus => {
          commit('setApiStatus', 'good', { root: true })
          commit('updatePipeline', { jobStatus, currentPipelineId })
        },
        err => {
          console.log('error')
          console.log(err)
          commit('setApiStatus', 'bad', { root: true })
        }
      )
    }
  }
}
