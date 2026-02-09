# prometheus-dellhw-exporter

A Helm chart for the dellhw_exporter

![Version: 1.2.2](https://img.shields.io/badge/Version-1.2.2-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: v2.0.0-rc.5](https://img.shields.io/badge/AppVersion-v2.0.0--rc.5-informational?style=flat-square)

## Get Repo Info

```console
helm repo add dellhw_exporter https://galexrt.github.io/dellhw_exporter
helm repo update
```

_See [helm repo](https://helm.sh/docs/helm/helm_repo/) for command documentation._

## Install Chart

To install the chart with the release name `my-release`:

```console
helm install --namespace <your-cluster-namespace> my-release dellhw_exporter/prometheus-dellhw-exporter
```

The command deploys dellhw_exporter on the Kubernetes cluster in the default configuration.

_See [configuration](#configuration) below._

_See [helm install](https://helm.sh/docs/helm/helm_install/) for command documentation._

### Development Build

To deploy from a local build from your development environment:

```console
cd charts/prometheus-dellhw-exporter
helm install --namespace <your-cluster-namespace> my-release . -f values.yaml
```

## Uninstall Chart

To uninstall/delete the my-release deployment:

```console
helm delete --namespace <your-cluster-namespace> my-release
```

This removes all the Kubernetes components associated with the chart and deletes the release.

_See [helm uninstall](https://helm.sh/docs/helm/helm_uninstall/) for command documentation._

## Configuration

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| additionalEnv | list | `[{"name":"HOSTNAME","valueFrom":{"fieldRef":{"fieldPath":"spec.nodeName"}}}]` | Additional environments to be added to the dellhw_exporter container, use this to configure the exporter (see https://github.com/galexrt/dellhw_exporter/blob/main/docs/configuration.md#environment-variables) |
| additionalVolumeMounts | list | `[]` | Additional volume mounts for the dellhw_exporter container. |
| additionalVolumes | list | `[]` | Additional volumes to be mounted in the dellhw_exporter container. |
| affinity | object | `{}` | Affinity for the DaemonSet |
| fullnameOverride | string | `""` | Override fully-qualified app name |
| image.pullPolicy | string | `"IfNotPresent"` | Override the `imagePullPolicy` |
| image.repository | string | `"quay.io/galexrt/dellhw_exporter"` | Image repository |
| image.tag | string | `""` | Overrides the image tag whose default is the chart appVersion. |
| imagePullSecrets | list | `[]` | ImagePullSecrets to add to the DaemonSet |
| nameOverride | string | `""` | Override chart name |
| nodeSelector | object | `{}` | NodeSelector for the DaemonSet |
| podAnnotations | object | `{}` | Annotations to add to the Pods created by the DaemonSet |
| podLabels | object | `{}` | Additional labels to add to the Pods created by the DaemonSet |
| podSecurityContext | object | `{}` | Kubernetes PodSecurityContext for the Pods |
| prometheusRule.additionalLabels | object | `{}` | Additional Labels for the PrometheusRule object |
| prometheusRule.enabled | bool | `false` | Specifies whether a prometheus-operator PrometheusRule should be created |
| prometheusRule.rules | list | `[]` | Checkout the https://github.com/galexrt/dellhw_exporter/blob/main/contrib/monitoring/prometheus-alerts/prometheus-alerts.yml for example alerts |
| psp.create | bool | `false` | Specifies whether a PodSecurityPolicy (PSP) should be created |
| psp.spec | object | `{"allowedHostPaths":[],"privileged":true,"volumes":["secret"]}` | PodSecurityPolicy spec |
| resources | object | `{}` | Resources for the dellhw_exporter container |
| securityContext | object | `{"privileged":true}` | SecurityContext for the container |
| service.port | int | `9137` | Service port |
| service.type | string | `"ClusterIP"` | Service type |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created |
| serviceAccount.name | string | `""` | The name of the service account to use. |
| serviceMonitor.additionalLabels | object | `{}` | Additional Labels for the ServiceMonitor object |
| serviceMonitor.enabled | bool | `false` | Specifies whether a prometheus-operator ServiceMonitor should be created |
| serviceMonitor.namespaceSelector | string | `nil` | ServiceMonitor namespace selector (check the comments below for examples) |
| serviceMonitor.scrapeInterval | duration | `"30s"` | Interval at which metrics should be scraped |
| serviceMonitor.scrapeTimeout | duration | `"20s"` | Timeout for scraping |
| tolerations | list | `[]` | Tolerations for the DaemonSet |
