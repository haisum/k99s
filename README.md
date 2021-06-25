Registry file `registries.yaml`:
```yaml
mirrors:
  "localhost:5000":
    endpoint:
      - http://k3d-registry.localhost:5000
```
Create cluster and registry

```bash
k3d registry create registry.localhost --port=5000
k3d cluster create --agents=2  --registry-config registries.yaml  --registry-use k3d-registry.localhost:5000 \
    -p 8443:443@loadbalancer -p 8080:80@loadbalancer
```

To install CRDs: `make install`

To build image: `make docker-build && make docker-push`

To run controller locally: `make run`

To deploy operator on cluster: `make deploy`

To install a sample resource `kubectl apply -f config/samples/paas_v1_box.yaml`

Build docker images:

```bash
cd docker/php
docker build . -t localhost:5000/php:1.0
docker push localhost:5000/php:1.0

cd ../go
docker build . -t localhost:5000/go:1.0
docker push localhost:5000/go:1.0
```