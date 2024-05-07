# ExternalDNS - Anexia Webhook Provider

[![License](https://img.shields.io/github/license/probstenhias/external-dns-anexia-webhook?style=for-the-badge)](LICENSE.md)
[![Build](https://img.shields.io/github/actions/workflow/status/probstenhias/external-dns-anexia-webhook/pull_request.yml?style=for-the-badge)](https://github.com/probstenhias/external-dns-anexia-webhook/actions/workflows/pull_request.yml)
[![GoReport](https://goreportcard.com/badge/github.com/probstenhias/external-dns-anexia-webhook?style=for-the-badge)](https://goreportcard.com/report/github.com/probstenhias/external-dns-anexia-webhook)
[![Coverage](https://img.shields.io/coverallsCoverage/github/ProbstenHias/external-dns-anexia-webhook?style=for-the-badge)](https://coveralls.io/github/ProbstenHias/external-dns-anexia-webhook?branch=main))

The Anexia Webhook Provider for [ExternalDNS](https://github.com/kubernetes-sigs/external-dns) allows you to use Anexia's DNS API to manage DNS records for your domains.

The provider is heavily inspired by the [ExternalDNS - IONOS Webhook](https://github.com/ionos-cloud/external-dns-ionos-webhook) and some inspiration taken from the [External DNS - Adguard Home Provider](https://github.com/muhlba91/external-dns-provider-adguard/tree/main).

## Configuration

See [cmd/webhook/init/configuration/configuration.go](cmd/webhook/init/configuration/configuration.go) for all available configuration options for the webhook sidecar, and [internal/anexia/configuration.go](internal/anexia/configuration.go) for all available configuration options for the Anexia provider.

## Kubernetes Deployment

The Anexia Webhook Provider is provided as  an OCI image in [ghcr.io/probstenhias/external-dns-anexia-webhook](https://ghcr.io/probstenhias/external-dns-anexia-webhook).

The following is an example deployment for the Anexia Webhook Provider:

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami

# create the anexia configuration
kubectl create secret generic anexia-configuration \
    --from-literal=url='<ANEXIA_API_URL>' \
    --from-literal=token='<ANEXIA_API_TOKEN>'

# create the helm values file
cat <<EOF > external-dns-anexia-values.yaml
provider: webhook

extraArgs:
  webhook-provider-url: http://localhost:8888

sidecars:
  - name: anexia-webhook
    image: ghcr.io/probstenhias/external-dns-anexia-webhook:$RELEASE_VERSION
    ports:
      - containerPort: 8888
        name: http
    livenessProbe:
      httpGet:
        path: /healthz
        port: http
      initialDelaySeconds: 10
      timeoutSeconds: 5
    readinessProbe:
      httpGet:
        path: /healthz
        port: http
      initialDelaySeconds: 10
      timeoutSeconds: 5
    env:
      - name: LOG_LEVEL
        value: debug
      - name: ANEXIA_API_URL
        valueFrom:
          secretKeyRef:
            name: anexia-configuration
            key: url
      - name: ANEXIA_API_TOKEN
        valueFrom:
          secretKeyRef:
            name: aneixa-configuration
            key: token
      - name: SERVER_HOST
        value: "0.0.0.0"
      - name: DRY_RUN
        value: "false"
EOF

# install external-dns with helm
helm install external-dns-anexia bitnami/external-dns -f external-dns-anexia-values.yaml
```
