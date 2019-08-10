// export const state = () => ({})

// export const actions = {
//   async nuxtServerInit({ dispatch }) {
//     await dispatch('tdr_types/load')
//   }
// }

export const state = () => ({
  myProjects: []
})

export const mutations = {
  updateProjects: (state, { projects, categoryName }) => {
    console.log('called updateProjects() in projects.js', categoryName)
    console.log(projects.data.length)
    state.myProjects = projects.data
    // console.log(state.projects)
    console.log(state.myProjects)
    // console.log(state.projects.categoryName[0])
  }
}

export const actions = {
  async load({ state, commit }, categoryName) {
    console.log('called load() in projects.js for category', categoryName)
    console.log(state.myProjects)
    if ((state.myProjects) && (categoryName in Object.keys(state.myProjects))) {
      console.log(state.myProjects[categoryName].length)
      return
    }
    console.log('axios projects')
    // commit('setApiStatus', null)
    await this.$axios.$get('/projects/'+categoryName).then(
      projects => {
        // commit('setApiStatus', 'good')
        // console.log(projects)
        commit('updateProjects', { projects, categoryName })
      },
      (err) => {
        console.log('error')
        console.log(err)
        // commit('setApiStatus', 'bad')
      }
    )
  }
}
