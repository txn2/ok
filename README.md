# GO OK

A Simple go web server listening on port **8080** that returns a bit of
useful data for testing and getting started with Kubernetes.

## Docker Run
```bash
# run version 1
docker run --rm -p 8080:8080 GIN_MODE=release cjimti/go-ok:v1

# run version 2
docker run --rm -p 8080:8080 GIN_MODE=release cjimti/go-ok:v2

```

Browse to http://localhost:8080

```json
{
    "client_ip": "172.17.0.1",
    "count": 3,
    "message": "ok",
    "time": "2018-03-05T08:38:03.936996398Z",
    "uuid_call": "dddb3561-7273-45ee-5f80-7b022d2bf2e9",
    "uuid_instance": "79defbd7-690e-4fc7-5652-354e1662ff7c",
    "version": 2,
    "version_msg": "version 2"
}
```

## Run Source

```bash
go get github.com/gin-gonic/gin
go get github.com/nu7hatch/gouuid

GIN_MODE=release go run ./ok.go
```

## Build Docker [Container]

```bash
$ docker build -t go-ok .
```

## Kubernetes Scripted

[See Official Kubernetes tutorial.][Official Tutorial]

### Create a [Deployment]

**Running** an [image] will create a [Pod] with one container using the
[kubectl] command.

```bash
$ kubectl run go-ok --image=go-ok:v1 --port=8080
```

View the [Deployment]:

```bash
$ kubectl get deployments
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
go-ok     1         1         1            0           2m
```

Looks like we have a problem since **0** [Pod]s are available. There are
a number of common reason why [Deployment]s fail.

- Read [10 Most Common Reasons Kubernetes Deployments Fail (Part 1)][Fail Article]

The #1 reason for failure is specifying the wrong image. Which we did
earlier, so take a look at our [Pod]s:

View the [Pod]s:

```bash
$ kubectl get pods
NAME                     READY     STATUS             RESTARTS   AGE
go-ok-8645cf567d-8zc62   0/1       ImagePullBackOff   0          22m
```

The the Pod status reports ImagePullBackOff (see [Pod Lifecycle]) and is
due to the fact that Kubernetes can not pull the image. We did not give
it the full and correct image name. Ue used go-ok:v1 and should have used
cjimti/go-ok:v1. Kubernetes needs to pull the [go-ok image] from a registry,
in this case hub.docker.com, the default Docker repository.

If we used Replication controllers we would want to perform a
`kubectl rolling-update` to push out the new image to each [Pod]. However
Kubernetes official documentation recommends:

>Note that kubectl rolling-update only supports Replication Controllers.
However, if you deploy applications with Replication Controllers, consider
switching them to Deployments. A Deployment is a higher-level controller
that automates rolling updates of applications declaratively, and therefore
is recommended. --[Rolling Update / Replication Controller](https://kubernetes.io/docs/tasks/run-application/rolling-update-replication-controller/)

Since we used a [Deployment], let's update the [Deployment] with the
correct image.

```bash
$ kubectl set image deployment/go-ok ok-ok=cjimti/go-ok:v1
deployment "go-ok" image updated
```

View the [Deployment] status:

```bash
$ kubectl get deployments
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
go-ok     1         1         1            1           1h
```

View the status of our [Pod]s:

```bash
NAME                     READY     STATUS    RESTARTS   AGE
go-ok-54c5f4d58f-gmzn2   1/1       Running   0          1m
```

Get a list of the latest Kubernetes events:

```bash
# kubectl get events
... (List of events)
```

### Create a [Service]

Since [Pod]s can come and go, a [Service] is a persistent link to a [pod].
>"A Kubernetes Service is an abstraction which defines a logical set of Pods
and a policy by which to access them - sometimes called a micro-service."
--Official [Service] Documentation

```bash
$ kubectl expose deployment go-ok --type=LoadBalancer
service "go-ok" exposed
```

> The --type=[LoadBalancer] flag indicates that you want to expose your
Service outside of the cluster. On cloud providers that support load
balancers, an external IP address would be provisioned to access the Service.
On Minikube, the [LoadBalancer] type makes the [Service] accessible through
the minikube service command.
--[Official Tutorial]

```bash
$ kubectl get services
NAME         TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
go-ok        LoadBalancer   10.111.21.238   <pending>     8080:31825/TCP   2m
kubernetes   ClusterIP      10.96.0.1       <none>        443/TCP          3d

```

Use the [minikube command] to automatically open up a browser window using a local IP address:

```bash
minikube service go-ok
```

### Local `minikube` Dashboard

Before we open the `minikube` dashboard let's enable **heapster**, which
gives us CPU and Memory graphs.

```bash
# List the minikube addons:
$ minikube addons list
- addon-manager: enabled
- coredns: disabled
- dashboard: enabled
- default-storageclass: enabled
- efk: disabled
- freshpod: disabled
- heapster: disabled
- ingress: disabled
- kube-dns: enabled
- registry: disabled
- registry-creds: disabled
- storage-provisioner: enabled

# enable heapster
$ minikube addons enable heapster
heapster was successfully enabled
```

Open the [Kubernetes Dashboard] with `minikube`:

```bash
$ minikube dashboard
```

View the [Pod]s and [Services] via command line:

```bash
$ kubectl get po,svc -n kube-system
```

[Heapster] is running as a service. You can open the web
interface with:

```bash
minikube addons open heapster
```

You get a [Grafana] web interface with pre-configured dashboards opened
up in your web browser.

### Cleanup

```bash
$ kubectl delete service hello-node
$ kubectl delete deployment hello-node

# Optionally, force removal of the Docker images created:

$ docker rmi go-ok:v1 go-ok:v2 -f

# Optionally, stop the Minikube VM:

$ minikube stop
$ eval $(minikube docker-env -u)

# Optionally, delete the Minikube VM:
$ minikube delete

```

## Next Steps

[Kubernetes 101](https://kubernetes.io/docs/user-guide/walkthrough/)


---

[go-ok image]: https://hub.docker.com/r/cjimti/go-ok/
[kubectl]: https://kubernetes.io/docs/reference/kubectl/overview/
[minikube command]: https://kubernetes.io/docs/getting-started-guides/minikube/
[Official Tutorial]: https://kubernetes.io/docs/tutorials/stateless-application/hello-minikube/
[Image]: https://kubernetes.io/docs/concepts/containers/images/
[Container]: https://kubernetes.io/docs/concepts/overview/components/#container-runtime
[Pod]: https://kubernetes.io/docs/concepts/workloads/pods/pod/
[Pod Lifecycle]: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/
[Deployment]: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/
[Service]: https://kubernetes.io/docs/concepts/services-networking/service/
[LoadBalancer]: https://kubernetes.io/docs/tasks/access-application-cluster/create-external-load-balancer/
[Kubernetes Dashboard]: https://github.com/kubernetes/dashboard
[Heapster]: https://github.com/kubernetes/heapster
[Grafana]: https://grafana.com/
[Fail Article]: https://kukulinski.com/10-most-common-reasons-kubernetes-deployments-fail-part-1/