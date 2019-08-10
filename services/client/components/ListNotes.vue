<template>
  <div>
    <label>Filter by Name:</label>
    <b-field label="Search by name">
      <b-input v-model="search_query"></b-input>
    </b-field>
    <label>Min ID:</label>
    <code>{{ filtered }}</code>
    <code>{{ idRange }}</code>
    <section style="width:20vw;">
      <b-field>
        <b-slider v-model="idRange" :min="50000" :max="80000" :step="5000" ticks></b-slider>
      </b-field>
    </section>
    <!-- <!-- <InputSpinner
        v-model="filters.id.value.min"
        :min="68000"
        :max="filters.id.value.max"
    />-->

    <!-- <label>Max ID:</label>
      <InputSpinner
        v-model="filters.id.value.max"
        :min="filters.id.value.min"
        :max="72000"
    />-->
    <section>
      <b-field grouped group-multiline>
        <button
          class="button field is-danger"
          @click="checkedRows = []"
          :disabled="!checkedRows.length"
        >
          <b-icon icon="close"></b-icon>
          <span>Clear checked</span>
        </button>
      </b-field>
      <b-tabs>
        <b-table
          :data="filtered"
          :hoverable="true"
          :striped="true"
          :default-sort="name"
          :sort-icon="chevron-up"
          :default-sort-direction="desc"
          :checked-rows.sync="checkedRows"
          checkable
          :checkbox-position="left"
          style="width:50vw;"
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
                {'is-danger': props.row.id >= 70000 },
                {'is-success': props.row.id < 70000 }
            ]"
              >{{ new Date(props.row.last_activity_at).toLocaleDateString() }}</span>
            </b-table-column>
            <b-table-column field="web_url" label="GitLab repository">
              <a :href="props.row.web_url">{{ props.row.web_url }}</a>
            </b-table-column>
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
      idRange: [50000, 80000],
      // filters: {
      //   name: { value: '', keys: ['name'] },
      //   id: { value: [70000, 75000], custom: this.ageFilter }
      // },
      checkedRows: []
    }
  },
  // methods: {
  //   ageFilter(filterValue, row) {
  //     return row.id >= filterValue[0] && row.id <= filterValue[1]
  //   }
  // },
  computed: {
    filtered() {
      var name_re = new RegExp(this.search_query, 'i')
      var tableData = []
      for (var i in this.$store.state.projects.myProjects) {
        if (
          (this.$store.state.projects.myProjects[i].name.match(name_re) ||
            this.$store.state.projects.myProjects[i].id
              .toString()
              .match(name_re)) &&
          (this.$store.state.projects.myProjects[i].id < this.idRange[1] &&
            this.$store.state.projects.myProjects[i].id > this.idRange[0])
        ) {
          tableData.push(this.$store.state.projects.myProjects[i])
        }
      }
      return tableData
    }
  }
}
</script>

<style lang="scss" scoped></style>
