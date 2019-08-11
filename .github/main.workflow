workflow "New workflow" {
  on = "push"
  resolves = ["Docker Push"]
}

action "Docker Registry" {
  uses = "actions/docker/login@86ab5e854a74b50b7ed798a94d9b8ce175d8ba19"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "GitHub Action for Docker" {
  uses = "actions/docker/cli@86ab5e854a74b50b7ed798a94d9b8ce175d8ba19"
  needs = ["Docker Registry"]
  args = "build -t tdr-diff-backend"
  runs = "sh -c \"cd services/server && docker $*\""
}

action "Docker Tag" {
  uses = "actions/docker/tag@86ab5e854a74b50b7ed798a94d9b8ce175d8ba19"
  needs = ["GitHub Action for Docker"]
  args = "tdr-diff-backend ${DOCKER_USERNAME}/tdr-diff-backend:${GITHUB_SHA}"
}

action "Docker Push" {
  uses = "actions/docker/cli@86ab5e854a74b50b7ed798a94d9b8ce175d8ba19"
  needs = ["Docker Tag"]
  args = "push ${DOCKER_USERNAME}/tdr-diff-backend:${GITHUB_SHA}"
  secrets = ["DOCKER_USERNAME"]
}
