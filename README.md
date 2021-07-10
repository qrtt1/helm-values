# helm-values

[![Go](https://github.com/qrtt1/friendly-yaml/actions/workflows/go.yml/badge.svg)](https://github.com/qrtt1/friendly-yaml/actions/workflows/go.yml)

## Extract configurations

List all configuration names in the values file

```
helm-values -f values.yaml
```

List all configurations with the filter

```
helm-values -f values.yaml -e "jupyter.*resources$"
```
```
jupyterhub.hub.resources
jupyterhub.prePuller.hook.resources
jupyterhub.prePuller.resources
jupyterhub.proxy.chp.resources
jupyterhub.scheduling.userScheduler.resources
```

Dump to yaml `-y` with the filter

```
helm-values -f values.yaml -e "jupyter.*resources$" -y
```
```yaml
jupyterhub:
  hub:
    resources:
      limits:
        cpu: 1000m
        memory: 512Mi
      requests:
        cpu: 100m
        memory: 512Mi
  prePuller:
    hook:
      resources:
        limits:
          cpu: 50m
          memory: 256Mi
        requests:
          cpu: 50m
          memory: 256Mi
    resources:
      limits:
        cpu: 50m
        memory: 256Mi
      requests:
        cpu: 50m
        memory: 256Mi
  proxy:
    chp:
      resources:
        limits:
          cpu: 200m
          memory: 512Mi
        requests:
          cpu: 200m
          memory: 512Mi
  scheduling:
    userScheduler:
      resources:
        limits:
          cpu: 50m
          memory: 256Mi
        requests:
          cpu: 50m
          memory: 256Mi
```
