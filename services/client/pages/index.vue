<template>
  <div>
    <div>
      <page-header />
    </div>
    <div class="container">
      <div>
        <cms-logo />
        <h1 class="title">CMS paper and notes diff</h1>
        <h1>Choose a category to browse:</h1>
        <div class="links">
          <span v-for="item in tdrTypes.names" :key="item">
            <nuxt-link :to="item" class="button--grey">{{ item }}</nuxt-link>
          </span>
        </div>
      </div>
    </div>
    <div>
      <page-footer />
    </div>
  </div>
</template>

<script>
import { mapState } from 'vuex'
import CmsLogo from '~/components/CmsLogo.vue'
import PageHeader from '~/components/Header.vue'
import PageFooter from '~/components/Footer.vue'

export default {
  components: {
    CmsLogo,
    PageHeader,
    PageFooter,
  },
  computed: {
    ...mapState(['tdrTypes', 'apiStatus'])
  }
}
</script>

<style>
.container {
  margin: 0 auto;
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  text-align: center;
}

.title {
  font-family: 'Quicksand', 'Source Sans Pro', -apple-system, BlinkMacSystemFont,
    'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  display: block;
  font-weight: 300;
  font-size: 100px;
  color: #35495e;
  letter-spacing: 1px;
}

.subtitle {
  font-weight: 300;
  font-size: 42px;
  color: #526488;
  word-spacing: 5px;
  padding-bottom: 15px;
}

.links {
  padding-top: 15px;
}
</style>
