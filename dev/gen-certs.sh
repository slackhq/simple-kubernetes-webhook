#!/bin/bash

set -euo pipefail
umask 077

dev_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
generated_dir="$dev_dir/generated"
tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT

cat > "$tmp_dir/ca.conf" <<'EOF'
[req]
distinguished_name = dn
x509_extensions = ca
[dn]
[ca]
basicConstraints = critical, CA:TRUE
keyUsage = critical, keyCertSign, cRLSign
EOF

openssl genrsa -out "$tmp_dir/ca.key" 2048

openssl req -new -x509 -days 365 -sha256 -key "$tmp_dir/ca.key" \
  -config "$tmp_dir/ca.conf" \
  -subj "/C=AU/CN=simple-kubernetes-webhook" \
  -out "$tmp_dir/ca.crt"

openssl req -newkey rsa:2048 -nodes -keyout "$tmp_dir/server.key" \
  -subj "/C=AU/CN=simple-kubernetes-webhook" \
  -out "$tmp_dir/server.csr"

openssl x509 -req -sha256 \
  -extfile <(printf '%s\n' \
    "subjectAltName=DNS:simple-kubernetes-webhook.default.svc,DNS:localhost,IP:127.0.0.1" \
    "basicConstraints=critical,CA:FALSE" \
    "keyUsage=critical,digitalSignature,keyEncipherment" \
    "extendedKeyUsage=serverAuth") \
  -days 365 \
  -in "$tmp_dir/server.csr" \
  -CA "$tmp_dir/ca.crt" -CAkey "$tmp_dir/ca.key" -CAcreateserial \
  -out "$tmp_dir/server.crt"

ca_bundle=$(base64 < "$tmp_dir/ca.crt" | tr -d '\n')
tls_cert=$(base64 < "$tmp_dir/server.crt" | tr -d '\n')
tls_key=$(base64 < "$tmp_dir/server.key" | tr -d '\n')

cat > "$tmp_dir/webhook.tls.secret.yaml" <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: simple-kubernetes-webhook-tls
  namespace: default
type: kubernetes.io/tls
data:
  tls.crt: $tls_cert
  tls.key: $tls_key
EOF

for config in mutating validating; do
  sed "s|caBundle:.*|caBundle: $ca_bundle|" \
    "$dev_dir/manifests/cluster-config/$config.config.yaml" \
    > "$tmp_dir/$config.config.yaml"
done

mkdir -p "$generated_dir"
cp "$tmp_dir/ca.crt" "$tmp_dir/"*.yaml "$generated_dir/"
echo "TLS certificates and manifests written to $generated_dir"
