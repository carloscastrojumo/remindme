# remindme
Reminde how to use commands by showing snippets

```sh
$ rmm add -t xargs -c "k get deploy -L helm.sh/chart | grep microservice | grep 3.0.25 | xargs kubectl get pods"

$ rmm show -t xargs
k get deploy -L helm.sh/chart | grep microservice | grep 3.0.25 | xargs kubectl get pods

$ rmm rm -t xargs -c "k get deploy -L helm.sh/chart | grep microservice | grep 3.0.25 | xargs kubectl get pods"
```
