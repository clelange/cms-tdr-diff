<template>
  <div>
    <b-loading :active.sync="isLoading" :can-cancel="false"></b-loading>
    <section class="section">
      <h1 class="title is-3">{{ categoryName }} / {{ $store.state.commits.projectInfo.name }}</h1>
      <h2 class="subtitle is-6">description: {{ $store.state.commits.projectInfo.description }}<br/>
      repository: <a :href=$store.state.commits.projectInfo.web_url>{{ $store.state.commits.projectInfo.web_url }}</a></h2>
      <p>You can filter by commit title and author name using the search box below. This will also show matches from other pages. It might also be useful to sort by date of last activity. Only commits of the last 90 days are shown.</p>
      <p>Select two commits, then hit the submit button to trigger the PDF diff pipeline. You can find the status of your jobs on the  <nuxt-link to="/statusboard">Status Board</nuxt-link> page.</p>
    </section>
    <div>
      <nav class="panel">
        <div class="panel-block">
          <b-field label="Filter by name">
            <p class="control has-icons-left">
              <b-input v-model="search_query" type="text" icon="magnify" placeholder="search"></b-input>
              <span class="icon is-small is-left">
                <i class="fas fa-search" aria-hidden="true"></i>
              </span>
            </p>
          </b-field>
        </div>
      </nav>
    </div>

    <!-- <span>{{ checkedRows }}</span> -->
    <div class="notification">
      <b-button
        v-for="(item, key, index) in checkedRows"
        v-on:click="removeElement(index)"
        :key="index"
        type="is-info"
        icon-right="delete"
      >{{ item.short_id }}</b-button>
      <b-button
        type="is-primary"
        size="is-large"
        v-bind:disabled="checkedRows.length != 2"
        @click="submitJob()"
      >Submit</b-button>
    </div>
    <section>
      <b-field grouped group-multiline>
        <button
          class="button field is-danger"
          @click="checkedRows = []"
          :disabled="!checkedRows.length"
        >
          <b-icon icon="close"></b-icon>
          <span>Clear selected</span>
        </button>
      </b-field>
      <b-tabs>
        <b-table
          :data="filtered"
          :hoverable="true"
          :striped="true"
          sort-icon="chevron-up"
          default-sort-direction="asc"
          :default-sort="['created_at', 'desc']"
          :checked-rows.sync="checkedRows"
          checkable
          :header-checkable="false"
          checkbox-position="left"
          style="width:90vw;"
          @click="(row) => toggleSelected(row)"
        >
          <template slot-scope="props">
            <b-table-column
              field="short_id"
              label="ID"
              width="40"
              sortable
              numeric
            >{{ props.row.short_id }}</b-table-column>
            <b-table-column field="title" label="Commit title" sortable>{{ props.row.title }}</b-table-column>
            <b-table-column field="created_at" label="Commit date" centered sortable>
              <span
                :class="
            [
                'tag',
                {'is-danger': ($dateFns.differenceInDays(new Date(), new Date(props.row.created_at)) >= 7) },
                {'is-success': ($dateFns.differenceInDays(new Date(), new Date(props.row.created_at)) < 7) }
            ]"
              >{{ $dateFns.format(new Date(props.row.created_at), 'DD/MM/YYYY') }}</span>
            </b-table-column>
            <b-table-column field="author_name" label="Author name">{{ props.row.author_name }}</b-table-column>
            <b-table-column field="author_email" label="Author email">{{ props.row.author_email }}</b-table-column>
          </template>
          <template slot="empty">
            <section class="section">
              <div class="content has-text-grey has-text-centered">
                <p>
                  <b-icon icon="emoticon-sad" size="is-large"></b-icon>
                </p>
                <p>No commits found in the last 90 days.</p>
              </div>
            </section>
          </template>
        </b-table>
      </b-tabs>
    </section>
  </div>
</template>

<script>
export default {
  data() {
    return {
      isLoading: !(this.$store.state.apiStatus),
      categoryName: this.$route.params.pathMatch.split('/')[0],
      search_query: '',
      checkedRows: [],
      commitList: [],
      currentPipeline: null
    }
  },
  computed: {
    filtered() {
      var query = this.search_query
      while (query.endsWith('\\')) {
        query = query.slice(0, query.lastIndexOf('\\') - 1)
      }
      var name_re = new RegExp(query, 'i')
      const myCommitList = this.commitList
      var tableData = []
      for (var i in myCommitList) {
        if (
          myCommitList[i].short_id.match(name_re) ||
          myCommitList[i].title.match(name_re) ||
          myCommitList[i].author_name.match(name_re)
        ) {
          tableData.push(myCommitList[i])
        }
      }
      return tableData
    }
  },
  mounted() {
    this.commitList = this.$store.state.commits.commitList
  },
  methods: {
    removeElement(index) {
      this.checkedRows.splice(index, 1)
    },
    toggleSelected(row) {
      const index = this.checkedRows.findIndex(p => p.short_id == row.short_id)
      console.log(row, index)

      if (index >= 0) {
        this.checkedRows.splice(index, 1)
      } else {
        this.checkedRows.push(row)
      }
    },
    success(payload) {
      this.$buefy.toast.open({
        duration: 5000,
        message:
          payload.status + ' Pipeline ID: ' + payload.pipeline_id.toString(),
        type: 'is-success'
      })
    },
    compare(a, b) {
      const timeA = new Date(a.created_at)
      const timeB = new Date(b.created_at)
      let comparison = 0
      if (timeA > timeB) {
        return 1
      } else return -1
    },
    async submitJob() {
      const sorted = this.checkedRows.sort(this.compare)
      // older comes first
      console.log(sorted, sorted[0].short_id, sorted[1].short_id)
      const postDict = {
        sha1: sorted[0].id,
        sha2: sorted[1].id,
        group: this.categoryName,
        project: this.$store.state.commits.projectInfo.name
      }
      this.$axios.setToken(this.$env.REQUEST_TOKEN)
      await this.$axios
        .$post('/trigger', postDict)
        .then(response => {
          console.log(response.pipeline_id)
          this.currentPipeline = response.pipeline_id
          this.$store.dispatch('jobs/load', response.pipeline_id)
          this.success(response)
        })
        .catch(error => {
          console.log(error)
        })
    }
  }
}
</script>

<style lang="scss" scoped></style>
