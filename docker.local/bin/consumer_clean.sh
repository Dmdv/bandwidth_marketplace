#!/bin/sh
  
for i in $(seq 1 6)
do
  rm docker.local/consumer$i/log/*
done

