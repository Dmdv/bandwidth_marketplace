#!/bin/sh

for i in $(seq 1 6)
do
  mkdir -p docker.local/provider$i/log
done
