# Kubernetes Commands Cheatsheet

## Cluster Setup (kind)

| Command | Description |
|---------|-------------|
| `brew install kind` | Install kind |
| `kind create cluster --name <name>` | Create a local cluster |
| `kind delete cluster --name <name>` | Delete a cluster |
| `kind get clusters` | List all kind clusters |
| `kind load docker-image <image> --name <cluster>` | Load a local Docker image into kind |
| `kubectl cluster-info --context kind-<name>` | Verify connection to cluster |

---

## Resource Types

| Resource | Description |
|----------|-------------|
| `Pod` | Smallest deployable unit — one or more containers sharing network and storage |
| `Deployment` | Manages a set of identical Pods; handles rolling updates and restarts |
| `Service` | Stable network address for a set of Pods; survives Pod restarts |
| `ConfigMap` | Non-sensitive key-value configuration injected into Pods |
| `Secret` | Sensitive key-value data (base64-encoded) injected into Pods |
| `Namespace` | Virtual cluster within a cluster — isolates resources by team or environment |
| `StatefulSet` | Like a Deployment but for stateful apps (databases) — stable Pod identity |
| `DaemonSet` | Runs one Pod per node — used for logging agents, monitoring |
| `Job` | Runs a Pod to completion — used for one-off tasks like migrations |
| `CronJob` | Runs a Job on a schedule |
| `Ingress` | Routes external HTTP/HTTPS traffic to Services |
| `HorizontalPodAutoscaler` | Automatically scales replicas based on CPU/memory |

---

## Manifest Structure

Every Kubernetes resource requires these four top-level fields:

```yaml
apiVersion: apps/v1     # API group and version
kind: Deployment        # resource type
metadata:
  name: my-app          # resource name (unique within namespace)
spec:                   # desired state (varies by resource type)
  ...
```

### ConfigMap
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  SERVER_PORT: "8081"
  REDIS_ADDR: "redis:6379"
```

### Secret
```yaml
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: app-secret
data:
  DATABASE_URL: <base64-encoded-value>   # echo -n "value" | base64
  JWT_SECRET: <base64-encoded-value>
```

### Service
```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-app
spec:
  type: ClusterIP        # ClusterIP | NodePort | LoadBalancer
  selector:
    app: my-app          # routes to Pods with this label
  ports:
    - port: 8081         # Service port
      targetPort: 8081   # container port
```

### Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
        - name: my-app
          image: my-app:latest
          imagePullPolicy: Never   # use local image, don't pull from registry
          ports:
            - containerPort: 8081
          env:
            - name: SERVER_PORT
              valueFrom:
                configMapKeyRef:
                  name: app-config
                  key: SERVER_PORT
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: app-secret
                  key: DATABASE_URL
```

---

## kubectl — Viewing Resources

| Command | Description |
|---------|-------------|
| `kubectl get pods` | List all Pods in current namespace |
| `kubectl get pods -A` | List Pods across all namespaces |
| `kubectl get pods -w` | Watch Pods in real time |
| `kubectl get deployment` | List Deployments |
| `kubectl get service` | List Services |
| `kubectl get configmap` | List ConfigMaps |
| `kubectl get secret` | List Secrets |
| `kubectl get all` | List all common resources |
| `kubectl get nodes` | List cluster nodes |
| `kubectl get namespaces` | List all namespaces |

---

## kubectl — Inspecting Resources

| Command | Description |
|---------|-------------|
| `kubectl describe pod <name>` | Full details of a Pod including events |
| `kubectl describe deployment <name>` | Full details of a Deployment |
| `kubectl describe service <name>` | Full details of a Service |
| `kubectl logs <pod-name>` | View logs from a Pod |
| `kubectl logs -l app=<label>` | View logs from all Pods matching a label |
| `kubectl logs -f <pod-name>` | Follow live logs |
| `kubectl logs --previous <pod-name>` | Logs from a crashed/restarted Pod |
| `kubectl exec -it <pod-name> -- sh` | Open a shell inside a running Pod |
| `kubectl top pods` | Show CPU and memory usage per Pod (requires metrics-server) |
| `kubectl top nodes` | Show CPU and memory usage per node |

---

## kubectl — Applying and Deleting

| Command | Description |
|---------|-------------|
| `kubectl apply -f <file.yaml>` | Create or update a resource |
| `kubectl apply -f <directory/>` | Apply all manifests in a directory |
| `kubectl delete -f <file.yaml>` | Delete resource defined in a file |
| `kubectl delete pod <name>` | Delete a specific Pod (Deployment recreates it) |
| `kubectl delete pods --all` | Delete all Pods in current namespace |
| `kubectl delete deployment <name>` | Delete a Deployment and its Pods |

---

## kubectl — Deployments and Rollouts

| Command | Description |
|---------|-------------|
| `kubectl scale deployment <name> --replicas=<n>` | Scale to n replicas |
| `kubectl rollout status deployment <name>` | Watch rollout progress |
| `kubectl rollout history deployment <name>` | View revision history |
| `kubectl rollout undo deployment <name>` | Roll back to previous revision |
| `kubectl rollout undo deployment <name> --to-revision=<n>` | Roll back to specific revision |
| `kubectl set image deployment/<name> <container>=<image>` | Update container image |

---

## kubectl — Namespaces

| Command | Description |
|---------|-------------|
| `kubectl get pods -n <namespace>` | List Pods in a specific namespace |
| `kubectl apply -f file.yaml -n <namespace>` | Apply manifest to a namespace |
| `kubectl config set-context --current --namespace=<ns>` | Set default namespace |

---

## kubectl — Debugging

| Command | Description |
|---------|-------------|
| `kubectl describe pod <name>` | Check Events section for error reasons |
| `kubectl get events --sort-by=.lastTimestamp` | View all cluster events sorted by time |
| `kubectl run debug --image=busybox -it --rm -- sh` | Spin up a temporary debug Pod |
| `kubectl port-forward pod/<name> 8081:8081` | Forward local port to a Pod |
| `kubectl port-forward service/<name> 8081:8081` | Forward local port to a Service |

---

## Service Types

| Type | Description |
|------|-------------|
| `ClusterIP` | Internal only — accessible within the cluster (default) |
| `NodePort` | Exposes service on each node's IP at a static port (30000-32767) |
| `LoadBalancer` | Provisions a cloud load balancer (AWS ELB, GCP LB) — not available in kind |

---

## Encoding Secrets

```bash
# encode a value
echo -n "my-secret-value" | base64

# decode a value
echo -n "bXktc2VjcmV0LXZhbHVl" | base64 -d
```

The `-n` flag suppresses the trailing newline — always use it, otherwise the newline gets encoded into the value and causes auth failures.
