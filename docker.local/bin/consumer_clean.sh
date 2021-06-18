#!/bin/sh
  
for i in $(seq 1 6);
do
  rm -r docker.local/consumer$i/*
done

