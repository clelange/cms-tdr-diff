<template>
  <div>
    <h1>Viewing route: {{ $route.params }}</h1>
    <nuxt-link to="/">Home</nuxt-link>
    <span>--{{ myProjects }}__</span>
    <!-- <div v-for="project in projects" :key="project.id" class="project">
      <h3>
        <a :href="project.web_url">{{ project.name }}</a>
      </h3>
    </div> -->
  </div>
</template>

<script>
// import { mapState } from 'vuex'
// import ListNotes from '~/components/ListNotes.vue'

export default {
  data() {
    return {
      myProjects: ['none']
    }
  },
  async asyncData({ $axios }) {
    // await this.$axios.$get('/projects/'+categoryName)
    let { data } = await $axios.$get('/projects/reports')
    console.log(data)
    return { myProjects: data }
  },
}
</script>

<style lang="scss" scoped></style>
