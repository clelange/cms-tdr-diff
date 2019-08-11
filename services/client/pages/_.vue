<template>
  <div>
    <div>
      <page-header />
    </div>
    <div>
      <!-- <h1>Viewing route: {{ $route.params }}</h1> -->
      <!-- <nuxt-link to="/">Home</nuxt-link> -->
      <!-- <nuxt-link to="/statusboard">Status Board</nuxt-link> -->
      <!-- <span>categoryPage: {{ categoryPage }}</span> -->
      <list-notes v-if="categoryPage" />
      <list-commits v-if="!categoryPage" />
      <!-- When using the store with fetch:
      <span>...{{ $store.state.projects.myProjects }}...</span>-->
      <!-- When using asyncData:
      <span>...{{ myProjects }}...</span>-->
    </div>
  </div>
</template>

<script>
import { mapState } from 'vuex'
import ListNotes from '~/components/ListNotes.vue'
import ListCommits from '~/components/ListCommits.vue'
import PageHeader from '~/components/Header.vue'

export default {
  components: {
    ListNotes,
    ListCommits,
    PageHeader
  },
  data() {
    return {
      categoryPage: false,
      categoryName: null,
      analysisId: null
      // myProjects: []
    }
  },
  // this works as well
  // async asyncData({ $axios, params }) {
  //   const res = params.pathMatch.split('/')
  //   const categoryName = res[0]
  //   if (res.length === 1) {
  //     let { data } = await $axios.$get('/projects/'+categoryName)
  //     return { myProjects: data }
  //   }
  // },
  async fetch({ store, params }) {
    const res = params.pathMatch.split('/')
    const categoryName = res[0]
    if (res.length === 1) {
      return await store.dispatch('projects/load', categoryName)
    } else if (res.length >= 1) {
      const analysisId = res[1]
      return await store.dispatch('commits/load', { categoryName, analysisId })
    }
  },
  computed: {
    ...mapState({
      proj: 'projects/myProjects',
      comm: 'commits/commitList',
      tdrTypes: 'tdrTypes'
    })
  },
  mounted() {
    const res = this.$route.params.pathMatch.split('/')
    if (res.length <= 1) {
      this.categoryPage = true
      this.categoryName = res[0]
    } else {
      this.categoryName = res[0]
      this.analysisId = res[1]
    }
  },
  validate({ params, query, store }) {
    const res = params.pathMatch.split('/')
    const types = store.state.tdrTypes.names
    console.log('types:', types)
    if (!types.includes(res[0])) {
      return false
    }
    if (res.length <= 1) {
      return true
    }
    // Now check for CADI/AN style
    if (res[0] in ['papers', 'notes']) {
      const notePattern = /^\S{2,}-\d{2}-\d{3}/
      const matchesPattern = notePattern.test(res[1])
      if (matchesPattern) {
        return true
      } else {
        throw new Error(res[2], 'is not a valid', res[1], 'pattern!')
      }
    } else {
      return true
    }
  }
}
</script>

<style lang="scss" scoped></style>
