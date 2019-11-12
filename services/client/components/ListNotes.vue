<template>
  <div>
    <b-loading :active.sync="isLoading" :can-cancel="false"></b-loading>
    <!-- <label>Filter by Name:</label> -->
    <section class="section">
      <h1 class="title is-3">{{ categoryName }}</h1>
      <!-- <h2 class="subtitle is-4">Actions:</h2> -->
      <p>You can filter by ID/name using the search box below. This will show matches for the full category, i.e. not only the ones shown on the respective page. If you would like to see all analyses/matches, you can also disable pagination. It might also be useful to sort by name instead of date of last activity.</p>
      <p>Once you have found the desired analysis, click on its name.</p>
    </section>

    <nav class="panel">
      <div class="panel-block">
        <b-field label="Filter by name" label-position="on-border" grouped>
          <b-input v-model="search_query" type="text" icon="magnify" placeholder="search"></b-input>
          <p class="control">
            <button
              class="button is-primary"
              size="is-medium"
              v-bind:disabled="search_query == ''"
              @click="clearSearchQuery()"
            >Clear filter</button>
          </p>
        </b-field>
      </div>
    </nav>
    <div>
      <section>
        <b-tabs>
          <b-field grouped group-multiline>
            <b-select v-model="perPage" :disabled="!isPaginated">
              <option value="10">10 per page</option>
              <option value="20">20 per page</option>
              <option value="50">50 per page</option>
            </b-select>
            <div class="control is-flex">
              <b-switch v-model="isPaginated">Paginated</b-switch>
            </div>
          </b-field>
          <b-table
            :data="filtered"
            :paginated="isPaginated"
            :per-page="perPage"
            pagination-position="top"
            :hoverable="true"
            :striped="true"
            default-sort-direction="asc"
            :default-sort="['last_activity_at', 'desc']"
            sort-icon="chevron-up"
          >
            <template slot-scope="props">
              <b-table-column field="id" label="ID" width="40" sortable numeric>{{ props.row.id }}</b-table-column>
              <b-table-column field="name" label="Name" sortable>
                <nuxt-link :to="props.row.name" append>{{ props.row.name }}</nuxt-link>
              </b-table-column>
              <b-table-column field="last_activity_at" label="Last activity" centered sortable>
                <span
                  :class="
            [
                'tag',
                {'is-danger': ($dateFns.differenceInDays(new Date(), new Date(props.row.last_activity_at)) >= 7) },
                {'is-success': ($dateFns.differenceInDays(new Date(), new Date(props.row.last_activity_at)) < 7) }
            ]"
                >{{ $dateFns.distanceInWordsToNow(new Date(props.row.last_activity_at)) }}</span>
              </b-table-column>
              <b-table-column field="description" label="Description">{{ props.row.description }}</b-table-column>
              <b-table-column field="web_url" label="GitLab repository">
                <a :href="props.row.web_url">{{ props.row.web_url }}</a>
              </b-table-column>
            </template>
          </b-table>
        </b-tabs>
      </section>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    console.log('called data() in listNotes.vue')
    return {
      isLoading: !this.$store.state.apiStatus,
      categoryName: this.$route.params.pathMatch.split('/')[0],
      // search_query: this.$store.state.preferences.search_query,
      perPage: 10,
      isPaginated: true
    }
  },
  computed: {
    search_query: {
      get() {
        return this.$store.state.preferences.search_query
      },
      set(value) {
        this.$store.commit('preferences/updateSearchQuery', value)
      }
    },
    filtered() {
      var query = this.$store.state.preferences.search_query
      while (query.endsWith('\\')) {
        query = query.slice(0, query.lastIndexOf('\\') - 1)
      }
      var name_re = new RegExp(query, 'i')

      const myProject = this.$store.state.projects.myProjects
      var tableData = []
      for (var i in myProject) {
        if (myProject[i].name.match(name_re)) {
          tableData.push(myProject[i])
        }
      }
      return tableData
    },
    loaded() {
      return this.$store.state.jobs.status
    }
  },
  methods: {
    clearSearchQuery() {
      this.search_query = ''
    }
  }
}
</script>

<style lang="scss" scoped></style>
