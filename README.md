This is a simple Kubernetes operator built using kubebuilder. We are attempting to create a simple platform as a service
offering where a user can provide URL to a git repo, and platform takes care of provisioning backend and database for it.

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
    -p 8443:443@loadbalancer -p 8080:80@loadbalancer
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

`127.0.0.1 php-box.k99s-pass.com  go-box.k99s-pass.com`

Open http://php-box.k99s-pass.com:8080 in browser.