<template>
  <div>
    <!-- <label>Filter by Name:</label> -->
    <section class="section">
      <h1 class="title is-3">{{ categoryName }}</h1>
      <!-- <h2 class="subtitle is-4">Actions:</h2> -->
      <p>You can filter by ID/name using the search box below. This will also show matches from other pages. It might also be useful to sort by date of last activity.</p>
      <p>Once you have found the desired analysis, click on its name.</p>
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
      <!-- <label>Min ID:</label> -->
      <!-- <code>{{ filtered }}</code> -->
      <!-- <code>{{ idRange }}</code>
    <section style="width:20vw;">
      <b-field>
        <b-slider v-model="idRange" :min="50000" :max="80000" :step="5000" ticks></b-slider>
      </b-field>
      </section>-->
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
        <b-tabs>
          <b-table
            :data="filtered"
            :hoverable="true"
            :striped="true"
            default-sort-direction="desc"
            :default-sort="['name', 'asc']"
            sort-icon="chevron-up"
            paginated
            backend-pagination
            per-page="10"
            backend-sorting
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
      // checkedRows: []
    }
  },
  computed: {
    filtered() {
      var query = this.search_query
      while (query.endsWith("\\")) {
        query = query.slice(0, query.lastIndexOf('\\')-1)
      }
      var name_re = new RegExp(query, 'i')

      const myProject = this.$store.state.projects.myProjects
      var tableData = []
      for (var i in myProject) {
        if (
          myProject[i].name.match(name_re)
          // (myProject[i].name.includes(query) ||
          //   (myProject[i].description &&
          //     myProject[i].description.includes(query)) ||
          //   myProject[i].id.toString().includes(query)) &&
          // (myProject[i].id < this.idRange[1] &&
          //   myProject[i].id > this.idRange[0])
        ) {
          tableData.push(myProject[i])
        }
      }
      return tableData
    }
  }
}
</script>

<style lang="scss" scoped></style>
