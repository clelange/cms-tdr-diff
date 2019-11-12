export const state = () => ({
  search_query: ''
})

export const mutations = {
  updateSearchQuery(state, search_query) {
    // console.log('updateSearchQuery', search_query)
    state.search_query = search_query
  }
}
