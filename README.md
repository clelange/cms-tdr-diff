# cms-tdr-diff

A service to highlight changes between CMS paper/note versions.

## Setup

### Local development

```shell
cd services/server
python example.py
```

Now you can point your browser to <http://0.0.0.0:8000> to see the JSON response.

```shell
cd services/client
npm run dev
```

The application should then be running at <http://localhost:3000>.

### Running with docker-compose

```shell
docker-compose up -d --build
```

The frontend will be available at [http://localhost:8080](http://localhost:8080) and the backend at [http://localhost:8000](http://localhost:8000). To stop the service, run `docker-compose down`.

## Running on minikube

```shell
brew install kubectl
brew install docker-machine-driver-hyperkit
sudo chown root:wheel /usr/local/opt/docker-machine-driver-hyperkit/bin/docker-machine-driver-hyperkit
sudo chmod u+s /usr/local/opt/docker-machine-driver-hyperkit/bin/docker-machine-driver-hyperkit
brew cask install minikube
minikube start --vm-driver=hyperkit
minikube dashboard
```

```shell
docker build -t clelange/starlette-nuxt-kubernetes-client ./services/client
docker build -t clelange/starlette-nuxt-kubernetes-server ./services/server
```

```shell
minikube addons enable ingress
```

https://alligator.io/vuejs/working-with-environment-variables/

## Running on kubernetes

[Connect a Front End to a Back End Using a Service](https://kubernetes.io/docs/tasks/access-application-cluster/connecting-frontend-backend/)

## Backend

For now, need a token created at <https://gitlab.cern.ch/clange/tdr-diff/-/settings/ci_cd#js-pipeline-triggers>.
