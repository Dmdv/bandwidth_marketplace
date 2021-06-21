#!/bin/sh

for i in $(seq 1 6);
do
  mkdir -p "docker.local/consumer$i/log"

  config_file="docker.local/consumer$i/prometheus.yml"
  cp config/cons-prometheus.yml "$config_file"

  match="- targets:"
  node_addr="        - 198.18.0.91:505$i"
  sed -i "s/$match/$match\n${node_addr}/"  "$config_file"
done
