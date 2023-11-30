# notepad-demo

### Docker

To build your very own docker image type:
```shell
docker build -t notepad-demo:latest .
```
Don't forget to change the image in the [deployment file](./kubernetes/notepad/notepad-deployment.yaml).

### Kubernetes

To install notepad deployment, svc and mongo pvc, statefulset, service type the following:
```shell
kubernetes apply -f kubernetes/
```

For monitoring, it is required to install the following prometheus chart:
```shell
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus prometheus-community/prometheus
```

To access the web app type:
```shell
kubectl get svc -n default notepad-service -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
```

Navigate to the `<output-ip>/note/all` page and there you go.