This is a simple Kubernetes operator built using kubebuilder. We are attempting to create a simple platform as a service
offering where a user can provide URL to a git repo, and platform takes care of provisioning backend and database for it.

A simple resource will look like this:

```yaml
apiVersion: paas.example.com/v1
kind: Box
metadata:
  name: php-box
spec:
  runtime: php
  gitURL: https://github.com/haisum/k99s.git
  gitSubPath: docker/php/src
  backend: mysql
  bootstrapSQL: |
    CREATE TABLE php_user(
      id int not null PRIMARY KEY AUTO_INCREMENT
    )
```

As of now, this project has implementation for php and go runtimes, and mysql backend.

## Cluster setup

1. Install k3d.

2. Create a file `registries.yaml` with contents:
```yaml
mirrors:
  "localhost:5000":
    endpoint:
      - http://k3d-registry.localhost:5000
```

3. Create cluster and registry
```bash
k3d registry create registry.localhost --port=5000
k3d cluster create --agents=2  --registry-config registries.yaml  --registry-use k3d-registry.localhost:5000 \
    -p 8443:443@loadbalancer -p 8088:80@loadbalancer
```

4. Build and push runtime images for php and go:
   
```bash
cd docker/php
docker build . -t localhost:5000/php:1.0
docker push localhost:5000/php:1.0

cd ../go
docker build . -t localhost:5000/go:1.0
docker push localhost:5000/go:1.0
```

## Operator setup

Install crd by running `make install`. Verify by running `kubectl get box`.

Do one of these:

- Run controller locally: `make run`
- Run controller on cluster:
   1. build and push operator image: `make docker-build && make docker-push`
   2. deploy operator on cluster: `make deploy`

## Try it

- `kubectl apply -f config/samples/config/samples/paas_v1_box.yaml` 
- `kubectl get box`

Put this in `/etc/hosts`:

`127.0.0.1 php-box.k99s-paas.com  go-box.k99s-paas.com`

Open http://php-box.k99s-paas.com:8088 in browser.