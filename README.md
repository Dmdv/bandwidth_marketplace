
# Bandwidth marketplace setup

## Table of Contents

- [Initial setup](#initial-setup)

- [Building and starting nodes](#building-and-starting-nodes)

- [Connect to other network](#connect-to-other-network)

- [Logs](#logs)

- [Miscellaneous](#miscellaneous) - [Cleanup](#cleanup)

## Initial setup

### Directory setup

Go to the `git/bandwidth-marketplace` directory and run the following commands:

- For setting up Consumer:

```
./docker.local/bin/consumer_init.sh
```

- For setting up Provider:

```
./docker.local/bin/provider_init.sh
```

You also can clean logs by running the following commands:

- For Consumer:

```
./docker.local/bin/consumer_clean.sh
```

- For Provider:

```
./docker.local/bin/provider_clean.sh
```

## Building and starting nodes


1. Set up a network called `testnet0` for each of these node containers to talk to each other.

 ```
docker network create --driver=bridge --subnet=198.18.0.0/15 --gateway=198.18.0.255 testnet0
```

2. Go to the `git/bandwidth-marketplace` directory to build containers using:


- For building Consumer:

```
./docker.local/bin/consumer_build.sh
```

- For building Provider:

```
./docker.local/bin/provider_build.sh
```

**_Note: you can skip building step if you want to use existing docker images._**

3. To start nodes working in the `git/bandwidth-marketplace` run the following commands:

- For starting consumer:

```
cd docker.local/consumer1
```
```
../bin/consumer_start_bls.sh
```

- For starting provider:

```
cd docker.local/provider1
```
```
../bin/provider_start_bls.sh
```

## Connect to other network

- Your network connection depends on the `block_worker` url you give in the `git/bandwith-marketplace/config/prov-config.yaml` and `cons-config.yaml` config files.

`block_worker: http://198.18.0.98:9091`

This works as a dns service, You need to know the above url for any network you want to connect, just replace it in the above mentioned file.

For example: If you want to connect to one network

```
block_worker: http://one.devnet-0chain.net/dns
```

See [0dns](https://github.com/0chain/0dns) if you want to deploy your own dns server.

See [0chain](https://github.com/0chain/0chain) if you want to deploy your own `sharder` and `miner` nodes.

## Logs

Logs are stored in the `git/bandwidth-marketplace/docker.local/${node}/log` directory.

## Miscellaneous

### Cleanup

1. Get rid of old unused docker resources:

```
docker system prune
```

2. To get rid of all the docker resources and start afresh:

```
docker system prune -a
```

3. Stop All Containers

```
docker stop $(docker ps -a -q)
```

4. Remove All Containers

```
docker rm $(docker ps -a -q)
```
