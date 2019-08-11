<template>
  <b-navbar>
    <template slot="brand">
      <b-navbar-item class="marker" to="/" tag="nuxt-link">
        <b-icon icon="vector-difference" size="is-small" />CMS Paper Diff
      </b-navbar-item>
    </template>

    <template slot="start">
      <b-navbar-item to="/" tag="nuxt-link">Home</b-navbar-item>
      <b-navbar-item
        v-for="item in tdrTypes.names"
        :key="item"
        :to="'/'+item"
        tag="nuxt-link"
      >{{ item }}</b-navbar-item>
    </template>

    <template slot="end">
      <b-navbar-item href="#">
        <h1 v-if="apiStatus === 'good'" style="color: green;" class="button">API status OK</h1>
        <h1
          v-else-if="apiStatus === 'bad'"
          style="color: red;"
          class="button"
        >Cannot connect to backend API</h1>
        <h1 v-else class="button">Loading...</h1>
      </b-navbar-item>
      <b-navbar-item to="/statusboard" tag="nuxt-link">
        <div class="buttons">
          <strong>Status Board</strong>
        </div>
      </b-navbar-item>
    </template>
  </b-navbar>
</template>

<script>
import { mapState } from 'vuex'

export default {
  computed: {
    ...mapState(['apiStatus', 'tdrTypes'])
  }
}
</script>

<style lang="scss" scoped>
.marker {
  font-family: 'Permanent Marker', -apple-system, BlinkMacSystemFont, 'Segoe UI',
    Roboto, 'Helvetica Neue', Arial, sans-serif;
  font-weight: normal;
}
</style>