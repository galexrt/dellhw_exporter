# dellhw_exporter

![Version: 0.2.13](https://img.shields.io/badge/Version-0.2.13-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: v1.13.12](https://img.shields.io/badge/AppVersion-v1.13.12-informational?style=flat-square)

A Helm chart for the dellhw_exporter

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Affinity for the DaemonSet |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"IfNotPresent"` | Override the `imagePullPolicy` |
| image.repository | string | `"quay.io/galexrt/dellhw_exporter"` | Image repository |
| image.tag | string | `""` | Overrides the image tag whose default is the chart appVersion. |
| imagePullSecrets | list | `[]` | ImagePullSecrets to add to the DaemonSet |
| nameOverride | string | `""` |  |
| nodeSelector | object | `{}` | NodeSelector for the DaemonSet |
| podAnnotations | object | `{}` | Annotations to add to the Pods created by the DaemonSet |
| podLabels | object | `{}` | Additional labels to add to the Pods created by the DaemonSet |
| podSecurityContext | object | `{}` | Kubernetes PodSecurityContext for the Pods |
| prometheusRule.additionalLabels | object | `{}` | Additional Labels for the PrometheusRule object |
| prometheusRule.enabled | bool | `false` | Specifies whether a prometheus-operator PrometheusRule should be created |
| prometheusRule.rules | list | `[]` | Checkout the `/contrib/prometheus-alerts/prometheus-alerts.yml` file for example alerts |
| psp.create | bool | `true` | Specifies whether a PodSecurityPolicy (PSP) should be created |
| psp.spec | object | `{"allowedHostPaths":[],"privileged":true,"volumes":["secret"]}` | PodSecurityPolicy spec |
| resources | object | `{}` |  |
| securityContext | object | `{"privileged":true}` | SecurityContext for the container |
| service.port | int | `9137` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created |
| serviceAccount.name | string | `""` | The name of the service account to use. |
| serviceMonitor.additionalLabels | object | `{}` | Additional Labels for the ServiceMonitor object |
| serviceMonitor.enabled | bool | `false` | Specifies whether a prometheus-operator ServiceMonitor should be created |
| serviceMonitor.namespaceSelector | string | `nil` |  |
| serviceMonitor.scrapeInterval | string | `"30s"` |  |
| tolerations | list | `[]` | Tolerations for the DaemonSet |

