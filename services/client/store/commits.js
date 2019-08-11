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
    console.log(commits.commits.length)
    console.log(commits.project_info.length)
    state.commitList = commits.commits
    state.projectInfo = commits.project_info
    // console.log(state.commits)
    // console.log(state.commitList)
    // console.log(state.commits.categoryName[0])
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
    // commit('setApiStatus', null)
    await this.$axios.$get(categoryName+'/'+analysisId).then(
      commits => {
        // commit('setApiStatus', 'good')
        // console.log(commits)
          // var e = document.createElement('div');
          // e.innerHTML = input;
          // return e.childNodes.length === 0 ? "" : e.childNodes[0].nodeValue;
        commit('updateCommits', { commits, categoryName })
      },
      (err) => {
        console.log('error')
        console.log(err)
        // commit('setApiStatus', 'bad')
      }
    )
  }
}
