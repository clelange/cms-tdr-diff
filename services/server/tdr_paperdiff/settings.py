from starlette.config import Config
from starlette.datastructures import CommaSeparatedStrings, Secret, URL

config = Config(".env")

# Define frontend origin for CORS via environment with default fallback value
FRONTEND_ORIGIN = config('FRONTEND_ORIGIN', cast=CommaSeparatedStrings, default=['http://localhost:3000'])
GITLAB_TOKEN = config('GITLAB_TOKEN', cast=Secret)
TRIGGER_TOKEN = config('TRIGGER_TOKEN', cast=Secret)
GITLAB_URL = config('GITLAB_URL', cast=URL, default='https://gitlab.cern.ch')
GITLAB_PROJECT="56283"
PIPELINE_URL = f"https://gitlab.cern.ch/api/v4/projects/{GITLAB_PROJECT}/trigger/pipeline"

GROUP_IDS = {'papers': 16555,
            'notes': 16803,
            # 'reports': 19085
            }