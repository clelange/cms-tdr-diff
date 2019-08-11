<template>
  <div>
    <div>
      <page-header />
    </div>
    <section class="section">
      <h1 class="title is-3">Job status board</h1>
      <h2 class="subtitle is-4">Actions:</h2>
      <div class="button">
        <b-button @click="updatePipelines()" rounded>Update Pipeline Status</b-button>
        <b-button @click="injectPipeline()" rounded>Inject Pipeline</b-button>
      </div>
    </section>
    <section class="section">
      <b-tabs>
        <b-table
          :data="filtered"
          :loading="!loaded"
          :hoverable="true"
          :striped="true"
          sort-icon="chevron-up"
          default-sort-direction="asc"
          :default-sort="['pipeline_id', 'asc']"
          :header-checkable="false"
          checkbox-position="left"
          style="width:90vw;"
        >
          <template slot-scope="props">
            <b-table-column
              field="pipeline_id"
              label="PipelineID"
              width="100"
              sortable
              numeric
            >{{ props.row.pipelineId }}</b-table-column>
            <b-table-column
              field="job_id"
              label="JobID"
              width="40"
              sortable
              numeric
            >{{ props.row.jobId }}</b-table-column>
            <b-table-column field="status" label="Status" width="40" sortable>
              <span :class="props.row.status_style">{{ props.row.status }}</span>
            </b-table-column>
            <b-table-column
              field="created_at"
              label="Created"
              width="200"
              centered
              sortable
            >{{ props.row.created_at }} ago</b-table-column>
            <b-table-column
              field="duration"
              label="Duration"
              width="150"
              centered
            >{{ props.row.duration }}</b-table-column>
            <b-table-column
              field="artifacts_expire_at"
              label="Expires"
              width="150"
              centered
              sortable
            >{{ props.row.artifacts_expire_at }}</b-table-column>
            <b-table-column field="web_url" label="Job URL">
              <a :href="props.row.web_url">Link</a>
            </b-table-column>
            <b-table-column field="artifact" label="Artifacts">
              <a :href="props.row.artifacts_link">{{ props.row.artifacts }}</a>
            </b-table-column>
          </template>
          <template slot="empty">
            <section class="section">
              <div class="content has-text-grey has-text-centered">
                <p>
                  <b-icon icon="emoticon-sad" size="is-large"></b-icon>
                </p>
                <p>No pipelines found.</p>
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
      activeStep: 0
    }
  }
}
</script>
</template>

<script>
import { mapState } from 'vuex'
import PageHeader from '~/components/Header.vue'

export default {
  components: {
    PageHeader
  },
  data() {
    return {
      posts: [
        { pipelineId: 1, job_status: { id: 12 } },
        { pipelineId: 2, job_status: { id: 1122 } }
      ]
    }
  },
  computed: {
    ...mapState({
      tdrTypes: 'tdrTypes',
      apiStatus: 'apiStatus',
      myJobs: 'jobs/pipelineStatus'
    }),
    loaded() {
      return this.$store.state.jobs.status
    },
    filtered() {
      console.log('filtered status')
      let massagedPipelines = []
      for (var i in this.$store.state.jobs.pipelineStatus) {
        const currentPipeline = this.$store.state.jobs.pipelineStatus[i]
        var jobDict = {}
        jobDict['pipelineId'] = currentPipeline.pipelineId
        jobDict['jobId'] = currentPipeline.jobStatus.id
        jobDict['status'] = currentPipeline.jobStatus.status
        switch (jobDict['status']) {
          case 'pending':
            jobDict['status_style'] = 'tag is-warning'
            jobDict['duration'] = '-'
            jobDict['artifacts_expire_at'] = '-'
            jobDict['artifacts_link'] = ''
            jobDict['artifacts'] = ''
            break
          case 'running':
            jobDict['status_style'] = 'tag is-warning'
            jobDict['duration'] =
              Number.parseInt(currentPipeline.jobStatus.duration) + ' s'
            jobDict['artifacts_expire_at'] = '-'
            jobDict['artifacts_link'] = ''
            jobDict['artifacts'] = ''
            break
          case 'success':
            jobDict['status_style'] = 'tag is-success'
            jobDict['duration'] =
              Number.parseInt(currentPipeline.jobStatus.duration) + ' s'
            if (
              this.$dateFns.differenceInSeconds(
                new Date(),
                new Date(currentPipeline.jobStatus.artifacts_expire_at)
              ) > 0
            ) {
              jobDict['artifacts_expire_at'] = 'expired'
              jobDict['artifacts_link'] = ''
              jobDict['artifacts'] = ''
            } else {
              jobDict['artifacts_expire_at'] =
                'in ' +
                this.$dateFns.distanceInWordsToNow(
                  new Date(currentPipeline.jobStatus.artifacts_expire_at)
                )
              jobDict['artifacts_link'] =
                currentPipeline.jobStatus.web_url + '/artifacts/download'
              jobDict['artifacts'] =
                currentPipeline.jobStatus.artifacts[0].filename
            }
            break
          default:
            jobDict['status_style'] = 'tag is-danger'
            jobDict['duration'] =
              Number.parseInt(currentPipeline.jobStatus.duration) + ' s'
            jobDict['artifacts_expire_at'] = '-'
            jobDict['artifacts_link'] = ''
            jobDict['artifacts'] = ''
        }
        jobDict['created_at'] =
          this.$dateFns.distanceInWordsToNow(
            new Date(currentPipeline.jobStatus.created_at)
          ) + ' ago'
        jobDict['web_url'] = currentPipeline.jobStatus.web_url
        massagedPipelines.push(jobDict)
      }
      return massagedPipelines
    }
  },
  methods: {
    async updatePipelines() {
      console.log('updatePipelines')
      await this.$store.dispatch('jobs/update')
    },
    async injectPipeline() {
      console.log('injectPipeline')
      const lines = [1031242, 1031241, 1031218, 1031195, 1031160, 672230]
      const pipelineId = lines[Math.floor(Math.random() * lines.length)]
      await this.$store.dispatch('jobs/load', pipelineId)
    }
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
