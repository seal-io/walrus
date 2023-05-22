package deployer

import (
	"text/template"
)

var (
	tmplOpencost            = template.Must(template.New("opencost").Parse(tmplOpencostContent))
	tmplPrometheusScrapeJob = template.Must(
		template.New("prometheusScrapeJob").Parse(tmplPrometheusScrapeJobContent))
)

// source: https://github.com/opencost/opencost/blob/v1.100.2/kubernetes/opencost.yaml
var tmplOpencostContent = `apiVersion: v1
kind: Namespace
metadata:
    name: {{.Namespace}}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: {{.Name}}
rules:
  - apiGroups:
      - ''
    resources:
      - configmaps
      - deployments
      - nodes
      - pods
      - services
      - resourcequotas
      - replicationcontrollers
      - limitranges
      - persistentvolumeclaims
      - persistentvolumes
      - namespaces
      - endpoints
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - extensions
    resources:
      - daemonsets
      - deployments
      - replicasets
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - apps
    resources:
      - statefulsets
      - deployments
      - daemonsets
      - replicasets
    verbs:
      - list
      - watch
  - apiGroups:
      - batch
    resources:
      - cronjobs
      - jobs
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - autoscaling
    resources:
      - horizontalpodautoscalers
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - policy
    resources:
      - poddisruptionbudgets
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - storage.k8s.io
    resources:
      - storageclasses
    verbs:
      - get
      - list
      - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: {{.Name}}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: {{.Name}}
subjects:
  - kind: ServiceAccount
    name: {{.Name}}
    namespace: {{.Namespace}}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
  labels:
    app: {{.Name}}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{.Name}}
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: {{.Name}}
    spec:
      restartPolicy: Always
      serviceAccountName: {{.Name}}
      containers:
        - image: quay.io/kubecost1/kubecost-cost-model:latest
          name: opencost
          resources:
            requests:
              cpu: "10m"
              memory: "55M"
            limits:
              cpu: "999m"
              memory: "1G"
          env:
            - name: PROMETHEUS_SERVER_ENDPOINT
              value: {{.PrometheusEndpoint}}
            - name: CLOUD_PROVIDER_API_KEY
              # The GCP Pricing API requires a key. This is supplied just for evaluation.
              value: "AIzaSyD29bGxmHAVEOBYtgd8sYM2gM2ekfxQX4U" 
            - name: CLUSTER_ID
              value: {{.ClusterID}} # Default cluster ID to use if cluster_id is not set in Prometheus metrics.
            - name: PRICING_CONFIGMAP_NAME
              value: {{.Name}}
            - name: KUBECOST_NAMESPACE
              value: {{.Namespace}}
          imagePullPolicy: Always
---
kind: Service
apiVersion: v1
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
spec:
  selector:
    app: {{.Name}}
  type: ClusterIP
  ports:
    - name: {{.Name}}
      port: 9003
      targetPort: 9003`

// source: https://github.com/opencost/opencost/blob/v1.100.2/kubernetes/prometheus/extraScrapeConfigs.yaml
var tmplPrometheusScrapeJobContent = `- job_name: {{.Name}}
  honor_labels: true
  scrape_interval: 1m
  scrape_timeout: 10s
  metrics_path: /metrics
  scheme: http
  dns_sd_configs:
  - names:
    - {{.Name}}.{{.Namespace}}
    type: 'A'
    port: 9003`
