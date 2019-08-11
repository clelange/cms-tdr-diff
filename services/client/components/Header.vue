<template>
  <b-navbar>
    <template slot="brand">
      <b-navbar-item>
        <nuxt-link to="/" class="marker">
          <b-icon icon="vector-difference" size="is-small" />
          CMS Paper Diff
        </nuxt-link>
      </b-navbar-item>
    </template>

    <template slot="start">
      <b-navbar-item>
        <nuxt-link to="/">Home</nuxt-link>
      </b-navbar-item>
      <b-navbar-item v-for="item in tdrTypes.names" :key="item">
        <nuxt-link :to="item">{{ item }}</nuxt-link>
      </b-navbar-item>
      <!-- <b-navbar-item href="#">
                About
      </b-navbar-item>-->
    </template>

    <template slot="end">
      <b-navbar-item tag="div">
        <h1 v-if="apiStatus === 'good'" style="color: green;" class="button">API status OK</h1>
      <h1 v-else-if="apiStatus === 'bad'" style="color: red;" class="button">Cannot connect to backend API</h1>
    <h1 v-else class="button">Loading...</h1>
        <div class="buttons">
          <nuxt-link to="/statusboard" class="button is-primary">
            <strong>Status Board</strong>
          </nuxt-link>
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