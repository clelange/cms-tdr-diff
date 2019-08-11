# import time
# import json
from datetime import datetime, timedelta
from starlette.config import environ
from starlette.applications import Starlette
from starlette.responses import JSONResponse
from starlette.middleware.cors import CORSMiddleware
from starlette.exceptions import HTTPException
from starlette.background import BackgroundTask
import requests
# import requests_cache
import gitlab
import uvicorn

from tdr_paperdiff import settings

print("FRONTEND_ORIGIN defined as:", settings.FRONTEND_ORIGIN)

# Cache requests to GitLab API
# requests_cache.install_cache() # use 'redis' eventually

print('Initialising GitLab')
gl = gitlab.Gitlab(settings.GITLAB_URL, private_token=str(settings.GITLAB_TOKEN), timeout=120)

app = Starlette(debug=True)

# A list of origins that should be permitted to make cross-origin requests
app.add_middleware(CORSMiddleware,
    allow_origins=settings.FRONTEND_ORIGIN,
)

@app.route('/')
async def homepage(request):
    # time.sleep(10)
    return JSONResponse({'hello': 'world'})

@app.route('/types')
async def get_types(request):
    return JSONResponse({'names': list(settings.GROUP_IDS.keys())})

@app.route('/projects/{group_id}')
async def get_group_projects(request):
    group_id = request.path_params['group_id']
    if group_id not in settings.GROUP_IDS:
        return HTTPException(404, detail=f"{group_id} not found in groups.")
    group = gl.groups.get(settings.GROUP_IDS[group_id])
    group_projects = group.projects.list(order_by='name', sort='asc', simple=True, all=True)
    # print(group_projects, dir(group_projects[0]))
    # print(group_projects[0].attributes)
    projects = {'data': [{'id': group_project.id,
                 'name': group_project.name,
                 'web_url': group_project.web_url,
                 'last_activity_at': group_project.last_activity_at,
                 'description': group_project.description}
                for group_project in group_projects]}
    return JSONResponse(projects)


@app.route('/{group}/{project}')
async def get_project_and_commits(request):
    group = request.path_params['group']
    project = request.path_params['project'].upper()
    since_days = 90
    use_date = datetime.today()
    use_date = use_date + timedelta(days=-since_days)
    try:
        project = gl.projects.get(f'tdr/{group}/{project}')
    except gitlab.exceptions.GitlabGetError:
        print(gitlab.exceptions.GitlabGetError)
        return HTTPException(404, detail=f"tdr/{group}/{project} not found.")
    project_info = {
        'id': project.id,
        'name': project.name,
        'description': project.description,
        'web_url': project.web_url,
        'tag_list': project.tag_list,
        'last_activity_at': project.last_activity_at,
    }
    gl_commits = project.commits.list(since=use_date.isoformat())
    commits = [{
                'id': glc.id,
                'short_id': glc.short_id,
                'created_at': glc.created_at,
                'title': glc.title,
                'author_name': glc.author_name,
                'author_email': glc.author_email,
               }
               for glc in gl_commits
              ]
    response = {'project_info': project_info, 'commits': commits}
    return JSONResponse(response)

@app.route('/trigger', methods=['POST'])
async def forward_ci_trigger(request):
    data = await request.json()
    print(data)
    sha1 = data['sha1']
    sha2 = data['sha2']
    group = data['group']
    project = data['project']
    # task = BackgroundTask(trigger_ci, sha1=sha1, sha2=sha2, group=group, project=project)
    # task = trigger_ci(sha1=sha1, sha2=sha2, group=group, project=project)
    payload = {
        'token': str(settings.TRIGGER_TOKEN),
        'ref': 'master',
        'variables[REPO_PROJECT]': project,
        'variables[REPO_GROUP]': group,
        'variables[GIT_SHA1]': sha1,
        'variables[GIT_SHA2]': sha2
    }
    r = requests.post(settings.PIPELINE_URL, data=payload)
    # print(r.json(), r.json().keys())
    pipeline_id = {'pipeline_id': r.json()['id']}
    print(pipeline_id)
    # print('task:', task)
    message = {'status': 'Pipeline triggered successfully!',
               'pipeline_id': r.json()['id']}
    return JSONResponse(message)


async def trigger_ci(sha1, sha2, group, project):
    payload = {
        'token': str(settings.TRIGGER_TOKEN),
        'ref': 'master',
        'variables[REPO_PROJECT]': project,
        'variables[REPO_GROUP]': group,
        'variables[GIT_SHA1]': sha1,
        'variables[GIT_SHA2]': sha2
    }
    r = requests.post(settings.PIPELINE_URL, data=payload)
    print(r.json(), r.json().keys())
    pipeline_id = r.json()['id']
    print(pipeline_id)
    project = gl.projects.get('clange/tdr-diff')
    pipeline = project.pipelines.get(pipeline_id)
    print(pipeline.variables.list())
    print(pipeline.attributes)
    return pipeline_id


@app.route('/status/pipeline/{pipeline_id}')
async def get_job_status(request):
    pipeline_id = request.path_params['pipeline_id']
    print("Received pipeline_id", pipeline_id)
    project = gl.projects.get('clange/tdr-diff')
    pipeline = project.pipelines.get(pipeline_id)
    print('Pipeline status:', pipeline.status)
    job = pipeline.jobs.list()[0]
    # print(job.attributes)
    job_status = {
        'id': job.id,
        'status': job.status,
        'created_at': job.created_at,
        'duration': job.duration,
        'web_url': job.web_url,
        'artifacts': job.artifacts,
        'artifacts_expire_at': job.artifacts_expire_at
    }
    return JSONResponse({'job_status': job_status})

if __name__ == '__main__':
    uvicorn.run(app, host='0.0.0.0', port=8000)