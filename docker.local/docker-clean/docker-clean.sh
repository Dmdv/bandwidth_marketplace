#!/bin/sh

echo "cleaning 6 consumers..."
for i in $(seq 1 6)
do
  echo "deleting consumer$i logs"
  rm -rf ./consumer$i/log/*
done

echo "cleaned up"
