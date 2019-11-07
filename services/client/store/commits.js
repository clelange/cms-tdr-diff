// export const state = () => ({})

// export const actions = {
//   async nuxtServerInit({ dispatch }) {
//     await dispatch('tdr_types/load')
//   }
// }

export const state = () => ({
  commitList: [],
  projectInfo: []
})

export const mutations = {
  updateCommits: (state, { commits, categoryName, analysisId }) => {
    console.log('called updateCommits() in commits.js', categoryName, analysisId)
    state.commitList = commits.commits
    state.projectInfo = commits.project_info
  }
}

export const actions = {
  async load({ state, commit }, { categoryName, analysisId }) {
    console.log('called load() in commits.js for category', categoryName, 'and analysis', analysisId)
    console.log(state.commitList)
    // FIXME
    if ((state.commitList) && (categoryName in Object.keys(state.commitList))) {
      console.log(state.commitList[categoryName].length)
      return
    }
    console.log('axios commits')
    commit('setApiStatus', null, { root: true })
    this.$axios.setToken(process.env.requestToken)
    await this.$axios.$get('commits/'+categoryName+'/'+analysisId).then(
      commits => {
        commit('setApiStatus', 'good', { root: true })
        commit('updateCommits', { commits, categoryName, analysisId })
      },
      (err) => {
        console.log('error')
        console.log(err)
        commit('setApiStatus', 'bad', { root: true })
      }
    )
  }
}
