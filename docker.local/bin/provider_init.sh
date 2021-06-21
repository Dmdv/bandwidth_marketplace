#!/bin/sh

for i in $(seq 1 6)
do
  mkdir -p "docker.local/provider$i/log"

  file="docker.local/provider$i/prometheus.yml"
  cp config/prov-prometheus.yml "$file"

  match="- targets:"
  node_addr="        - 198.18.0.1$i:507$i"
  sed -i "s/$match/$match\n${node_addr}/"  "$file"
done
