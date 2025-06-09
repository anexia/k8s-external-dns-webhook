# external-dns for Anexia CloudDNS

[![License](https://img.shields.io/github/license/anexia/k8s-external-dns-webhook?style=for-the-badge)](LICENSE.md)
[![Build](https://img.shields.io/github/actions/workflow/status/anexia/k8s-external-dns-webhook/pull_request.yml?style=for-the-badge)](https://github.com/anexia/k8s-external-dns-webhook/actions/workflows/pull_request.yml)
[![GoReport](https://goreportcard.com/badge/github.com/anexia/k8s-external-dns-webhook?style=for-the-badge)](https://goreportcard.com/report/github.com/anexia/k8s-external-dns-webhook)
[![Coverage](https://img.shields.io/coverallsCoverage/github/anexia/k8s-external-dns-webhook?style=for-the-badge)](https://coveralls.io/github/anexia/k8s-external-dns-webhook?branch=main)

The Anexia Webhook Provider for [ExternalDNS](https://github.com/kubernetes-sigs/external-dns) allows you to use Anexia's DNS API to manage DNS records for your domains.

The provider is heavily inspired by the [ExternalDNS - IONOS Webhook](https://github.com/ionos-cloud/external-dns-ionos-webhook) and some inspiration taken from the [External DNS - Adguard Home Provider](https://github.com/muhlba91/external-dns-provider-adguard/tree/main).

The initial work was by done by @ProbstenHias, who still serves as the primary maintainer. Thanks a lot! :purple_heart:

## Configuration

See [cmd/webhook/init/configuration/configuration.go](cmd/webhook/init/configuration/configuration.go) for all available configuration options for the webhook sidecar, and [internal/anexia/configuration.go](internal/anexia/configuration.go) for all available configuration options for the Anexia provider.

## Kubernetes Deployment

The Anexia Webhook Provider is provided as  an OCI image in [ghcr.io/anexia/k8s-external-dns-webhook](https://ghcr.io/anexia/k8s-external-dns-webhook).

The following is an example deployment for the Anexia Webhook Provider:

```bash
helm repo add external-dns https://kubernetes-sigs.github.io/external-dns/

# create the anexia configuration
kubectl create secret generic anexia-credentials \
    --from-literal=token='<ANEXIA_API_TOKEN>'

# create the helm values file
cat <<EOF > external-dns-anexia-values.yaml

# -- Currently this image is only published to a private repository.
global:
  # -- Add a global image pull secret.
  imagePullSecrets:
    - name: ghcr-creds

# -- ExternalDNS Log level.
logLevel: debug # reduce in production

# -- if true, _ExternalDNS_ will run in a namespaced scope (Role and Rolebinding will be namespaced too).
namespaced: false

# -- _Kubernetes_ resources to monitor for DNS entries.
sources:
  - ingress
  - service
  - crd

extraArgs:
  ## must override the default value with port 8888 with port 8080 because this is hard-coded in the helm chart
  - --webhook-provider-url=http://localhost:8080
  ## You should filter the domains that can be requested to limit the amount of requests done to the anxia engine.
  ## This will help you avoid running into rate limiting
  - --domain-filter=your.domain.com

provider:
  name: webhook
  webhook:
    image: ghcr.io/anexia/k8s-external-dns-webhook
    tag: v0.2.0
    env:
      - name: LOG_LEVEL
        value: debug # reduce in production
      - name: ANEXIA_API_URL
        value: <ANEXIA_API_URL>
      - name: ANEXIA_API_TOKEN
        valueFrom:
          secretKeyRef:
            name: anexia-credentials
            key: token
      - name: SERVER_HOST
        value: "0.0.0.0"
      - name: SERVER_PORT
        value: "8080"
      - name: DRY_RUN
        value: "false"
EOF

# install external-dns with helm
helm upgrade -i external-dns-anexia external-dns/external-dns -f external-dns-anexia-values.yaml
```
