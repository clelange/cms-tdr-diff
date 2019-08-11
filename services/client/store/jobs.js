export const state = () => ({
  pipelineStatus: []
})

export const mutations = {
  addPipeline: (state, { jobStatus, pipelineId }) => {
    console.log('called addPipeline() in jobs.js', pipelineId)
    console.log(jobStatus)
    state.pipelineStatus.push({pipelineId: pipelineId, jobStatus: jobStatus.job_status })
    console.log(state.pipelineStatus)
  },
  updatePipeline: (state, { jobStatus, currentPipelineId }) => {
    console.log('called updatePipeline() in jobs.js', currentPipelineId)
    const myJobStatus = jobStatus.job_status
    const index = state.pipelineStatus.findIndex(p => p.pipelineId == Number.parseInt(currentPipelineId))
    console.log("updating index:", index)
    if (index < 0) {
      return
    }
    state.pipelineStatus[index]['jobStatus'] = myJobStatus
  }
}

export const actions = {
  async load({ state, commit }, pipelineId) {
    console.log('called load() in jobs.js for pipeline', pipelineId)
    console.log(state.pipelineStatus)
    const index = state.pipelineStatus.findIndex(p => p.pipelineId == pipelineId)
    if (index >= 0) {
      console.log("Pipeline already in store:", pipelineId)
      return
    }
    console.log('pipelineStatus axios jobs')
    // commit('setApiStatus', null)
    await this.$axios.$get('/status/pipeline/'+pipelineId).then(
        jobStatus => {
        // commit('setApiStatus', 'good')
        // console.log(jobs)
          // var e = document.createElement('div');
          // e.innerHTML = input;
          // return e.childNodes.length === 0 ? "" : e.childNodes[0].nodeValue;
        commit('addPipeline', { jobStatus, pipelineId })
      },
      (err) => {
        console.log('error')
        console.log(err)
        // commit('setApiStatus', 'bad')
      }
    )
  },
  async update({ state, commit }) {
    console.log('called update() in jobs.js')
    console.log(state.pipelineStatus)
    if (state.pipelineStatus.length < 1) {
      console.log("No pipelines found")
      return
    }
    for (var i in state.pipelineStatus) {
      console.log(state.pipelineStatus)
      const currentPipelineId = Number.parseInt(state.pipelineStatus[i].pipelineId)
      console.log("this pipe:", currentPipelineId)
      // commit('setApiStatus', null)
      await this.$axios.$get('/status/pipeline/'+currentPipelineId).then(
      jobStatus => {
        // commit('setApiStatus', 'good')
        // console.log(jobs)
          // var e = document.createElement('div');
          // e.innerHTML = input;
          // return e.childNodes.length === 0 ? "" : e.childNodes[0].nodeValue;
          console.log('about to commit', currentPipelineId)
        commit('updatePipeline', { jobStatus, currentPipelineId })
      },
      (err) => {
        console.log('error')
        console.log(err)
        // commit('setApiStatus', 'bad')
        }
      )
    }
  }
}
