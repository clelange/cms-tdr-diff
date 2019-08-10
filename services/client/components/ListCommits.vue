<template>
  <div>
    <span>
      <ul v-for="(item, key) in $store.state.commits.projectInfo" :key="item.id">
        <li>{{ key }}: {{ item }}</li>
      </ul>
    </span>
    <!-- <div>
        <label class="typo__label">Simple select / dropdown</label>
        <multiselect v-model="value" :options="commitList" :multiple="true" :close-on-select="false" :clear-on-select="false" :preserve-search="true" placeholder="Pick some" label="short_id" track-by="short_id" :preselect-first="false">
        <template slot="selection" slot-scope="{ values, search, isOpen }"><span class="multiselect__single" v-if="values.length &amp;&amp; !isOpen">{{ values.length }} options selected</span></template>
        </multiselect>
        <pre class="language-json"><code>{{ value  }}</code></pre>
    </div>-->
    <!-- {{ $store.state.commits.commitList[0] }} -->
    <label>Filter by Name:</label>
    <b-field label="Search by name">
      <b-input v-model="search_query"></b-input>
    </b-field>
    <!-- <label>Min ID:</label> -->
    <!-- <code>{{ filtered }}</code> -->
    <!-- <code>{{ idRange }}</code> -->
    <!-- <section style="width:20vw;">
      <b-field>
        <b-slider v-model="idRange" :min="50000" :max="80000" :step="5000" ticks></b-slider>
      </b-field>
    </section>-->

    <!-- <span>{{ checkedRows }}</span> -->
    <div class="notification">
      <b-button
        v-for="(item, key, index) in checkedRows"
        v-on:click="removeElement(index)"
        :key="index"
        type="is-danger"
        icon-right="delete"
      >{{ item.short_id }}</b-button>
      <b-button type="is-primary" size="is-large" v-bind:disabled="checkedRows.length != 2" @click="submitJob()">Submit</b-button>
    </div>
    <section>
      <b-field grouped group-multiline>
        <button
          class="button field is-danger"
          @click="checkedRows = []"
          :disabled="!checkedRows.length"
        >
          <b-icon icon="close"></b-icon>
          <span>Clear all</span>
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
            <b-table-column field="title" label="Title" sortable>{{ props.row.title }}</b-table-column>
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
        </b-table>
      </b-tabs>
    </section>
    <!-- <th>Name</th>
        <th>last activity</th>
    </thead>
    <tbody slot="body" slot-scope="{displayData}">
        <tr v-for="row in displayData" :key="row.id">
          <td>{{ row.id }}</td>
          <td>{{ row.name }}</td>
          <td>{{ row.last_activity_at }}</td>
        </tr>
    </tbody>-->

    <!-- <h1>Listing notes for {{ categoryName }}</h1>
    <div v-for="project in $store.state.projects.myProjects" :key="project.id" class="project">
      <h3>
        <a :href="project.web_url">{{ project.name }}</a>
      </h3>
    </div>-->
  </div>
</template>

<script>
export default {
  data() {
    console.log('called data() in listNotes.vue')
    return {
      categoryName: this.$route.params.pathMatch.split('/')[0],
      search_query: '',
      // idRange: [50000, 80000],
      // filters: {
      //   name: { value: '', keys: ['name'] },
      //   id: { value: [70000, 75000], custom: this.ageFilter }
      // },
      checkedRows: [],
      commitList: []
    }
  },
  computed: {
    filtered() {
      // var name_re = new RegExp(this.search_query, 'i')
      const query = this.search_query
      const myCommitList = this.commitList
      var tableData = []
      for (var i in myCommitList) {
        if (
          myCommitList[i].short_id.includes(query) ||
          myCommitList[i].title.includes(query) ||
          myCommitList[i].author_name.includes(query)
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
    removeElement: function(index) {
      this.checkedRows.splice(index, 1)
    },
    toggleSelected(row) {
      console.log('toggleSelected')
      console.log(this)
      const index = this.checkedRows.findIndex(p => p.short_id == row.short_id)
      console.log(row, index)

      if (index >= 0) {
        this.checkedRows.splice(index, 1)
      } else {
        this.checkedRows.push(row)
      }
    },
    compare(a, b) {
      const timeA = new Date(a.created_at);
      const timeB = new Date(b.created_at);
      let comparison = 0;
      if (timeA > timeB) {
        return 1;
      } else
        return -1;
    },
    submitJob() {
      const sorted = this.checkedRows.sort(this.compare)
      // older comes first
      console.log(sorted, sorted[0].short_id, sorted[1].short_id)
    }
  }
}
</script>

<style lang="scss" scoped></style>
