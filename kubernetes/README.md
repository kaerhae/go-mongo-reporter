# /kubernetes

This directory contains deployments for reporter tool and MongoDB. To make these deployments, environment variables must be set to secrets. Reporter tool and MongoDB are requiring `reporter-secrets` and `mongo-secrets`, respectively. For the sake of simplicity, it might be recommended to first set `mongo-secrets` and apply MongoDB deployment, and then set `reporter-secrets` and apply reporter tool deployment, since reporter deployments environment variables are depending on mongo deployment.

Note: Before reporter deployment, built kaerhae/go-mongo-reporter image must be pushed to registry, so it can be fetched. If using `minikube`, image can be loaded to minikube with: `minikube image load kaerhae/go-mongo-reporter`

To build images, run:
```bash
docker build -t kaerhae/go-mongo-reporter -f ./docker/reporter/Dockerfile .
docker build -t kaerhae/go-mongo-reporter-migration -f ./docker/migrate/Dockerfile .
```

## Deploy MongoDB
To set `mongo-secrets`, run following command with suitable parameters:
```bash
kubectl create secret generic mongo-secrets \
--from-literal=MONGO_INITDB_ROOT_USERNAME="<USERNAME>" \
--from-literal=MONGO_INITDB_ROOT_PASSWORD="<PASSWORD>" \
--from-literal=MONGO_INITDB_DATABASE="<REPORTER_DB>"
```

MongoDB deployment is by default set to persisted volume in /data/mongodb. To set persistentVolume and persistentVolumeClaim, run:
```bash
kubectl apply -f mongo-pv.yml
kubectl apply -f mongo-pvc.yml
```


Then run deploy:
```bash
kubectl apply -f mongo-deployment.yml
```

To query MongoDB deployment with reporter tool, port must be forwarded:
`kubectl port-forward <MONGO_POD_NAME> 27017:27017`

## Deploy reporter

Run following, with suitable parameters:
```bash
kubectl create secret generic reporter-secrets \
--from-literal=MONGO_USER="<USERNAME>" \
--from-literal=MONGO_PASS="<PASSWORD>" \
--from-literal=MONGO_IP='<IP_ADDR' \
--from-literal=MONGO_PORT="<PORT>" \
--from-literal=SECRET_KEY="<SECRET_KEY>" \
--from-literal=DATABASE="<REPORTER_DB>" \
--from-literal=REPORTER_ROOT_USER="<ADMIN_USER>" \
--from-literal=REPORTER_ROOT_PASSWORD="<ADMIN_PASSWD>" \
```

Note: MONGO_IP can be retrieved with `kubectl describe pod <MONGO_POD_NAME> | grep 'IP'`

Run reporter deployment:
```bash
kubectl apply -f reporter-deployment.yml
```

To request reporter deployment, port must be forwarded:
`kubectl port-forward <REPORTER_POD_NAME> 8080:8080`
