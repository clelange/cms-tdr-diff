workflow "New workflow" {
  on = "push"
  resolves = ["Docker Push Backend", "Docker Push Frontend"]
}

action "Docker Registry" {
  uses = "actions/docker/login@86ab5e854a74b50b7ed798a94d9b8ce175d8ba19"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "Build Backend" {
  uses = "actions/docker/cli@86ab5e854a74b50b7ed798a94d9b8ce175d8ba19"
  needs = ["Docker Registry"]
  args = "build -t tdr-diff-backend ./services/server"
}

action "Docker Tag Backend" {
  uses = "actions/docker/tag@86ab5e854a74b50b7ed798a94d9b8ce175d8ba19"
  needs = ["Build Backend"]
  args = "tdr-diff-backend ${DOCKER_USERNAME}/tdr-diff-backend"
  secrets = ["DOCKER_USERNAME"]
}

action "Docker Push Backend" {
  uses = "actions/docker/cli@86ab5e854a74b50b7ed798a94d9b8ce175d8ba19"
  needs = ["Docker Tag Backend"]
  args = "push ${DOCKER_USERNAME}/tdr-diff-backend:${IMAGE_SHA}"
  secrets = ["DOCKER_USERNAME"]
}

action "Build Frontend" {
  uses = "actions/docker/cli@86ab5e854a74b50b7ed798a94d9b8ce175d8ba19"
  needs = ["Docker Registry"]
  args = "build -t tdr-diff-client ./services/client"
}

action "Docker Tag Frontend" {
  uses = "actions/docker/tag@86ab5e854a74b50b7ed798a94d9b8ce175d8ba19"
  needs = ["Build Frontend"]
  args = "tdr-diff-client ${DOCKER_USERNAME}/tdr-diff-client"
  secrets = ["DOCKER_USERNAME"]
}

action "Docker Push Frontend" {
  uses = "actions/docker/cli@86ab5e854a74b50b7ed798a94d9b8ce175d8ba19"
  needs = ["Docker Tag Frontend"]
  args = "push ${DOCKER_USERNAME}/tdr-diff-client:${IMAGE_SHA}"
  secrets = ["DOCKER_USERNAME"]
}