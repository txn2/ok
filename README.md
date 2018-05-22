# OK

A Simple web server listening on port **8080** that returns a bit of
useful data for testing and diagnosing ingress with Kubernetes.

## Docker Run
```bash
# run
docker run --rm -p 8080:8080 -e txn2/ok

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
go run ./ok.go
```

## Build Docker [Container]

```bash
$ docker build -t ok .
```

## Kubernetes Scripted

[See Official Kubernetes tutorial.][Official Tutorial]

### Create a [Deployment]

**Running** an [image] will create a [Pod] with one container using the
[kubectl] command. See the [kubectl cheatsheet] for a list of common commands.

```bash
$ kubectl run ok --image=ok --port=8080
```

View the [Deployment]:

```bash
$ kubectl get deployments
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
ok        1         1         1            0           2m
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
ok-8645cf567d-8zc62      0/1       ImagePullBackOff   0          22m
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
$ kubectl set image deployment/ok go-ok=txn2/ok
deployment "ok" image updated
```

View the [Deployment] status:

```bash
$ kubectl get deployments
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
ok        1         1         1            1           1h
```

View the status of our [Pod]s:

```bash
NAME                     READY     STATUS    RESTARTS   AGE
ok-54c5f4d58f-gmzn2      1/1       Running   0          1m
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
$ kubectl expose deployment ok --type=NodePort
service "ok" exposed
```

```bash
$ kubectl get services
NAME         TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
ok           NodePort       10.111.21.238   <none>        8080:31825/TCP   2m
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

View the [Pod]s and [Services] for the system namespace:

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

### Export Configuration

View the [Pod]s and [Services] for the default:

```bash
$ kubectl get po,svc
```

Get the YAML configuration for go-ok.

### Configuration

Using a declarative model of setting up out service and pods we can
leverage version control and have a more stable way of describing the
system to others and Kubernetes.

Rather than type a configuration file from scratch we can export the
current [Deployment] and take advantage of the work we have already done.


```bash
# get the deployment and service in one file
$ kubectl get deployment,service ok -o yaml --export >k8s-dev-local.yml
```

We can remove the key `nodePort: 30712` from the [Service] definition. The
port number will be different for you since this port was randomly generated
when we created the service. Allowing Kubernetes to assign a random port  is
useful for others who may be running of other service and may be using that
that port in their Minikube.

#### Testing Configuration

Since the configuration file matches our current [Deployment] we can use it
to delete the [Deployment]

```bash
$ kubectl delete -f k8s-dev-local.yml
```

Re-create the deployment using the new configuration file:

```bash
$ kubectl create -f k8s-dev-local.yml
```


### Manual Cleanup

```bash
# Delete service and deployment
$ kubectl delete service ok
$ kubectl delete deployment ok
```

## Next Steps

[Kubernetes 101](https://kubernetes.io/docs/user-guide/walkthrough/)


---

[ok image]: https://hub.docker.com/r/txn2/ok/
[kubectl]: https://kubernetes.io/docs/reference/kubectl/overview/
[kubectl cheatsheet]: https://kubernetes.io/docs/reference/kubectl/cheatsheet/
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