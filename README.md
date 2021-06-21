
# Bandwidth marketplace setup

## Table of Contents

- [Initial setup](#initial-setup)

- [Building and starting nodes](#building-and-starting-nodes)

- [Connect to other network](#connect-to-other-network)

- [Logs](#logs)

- [Miscellaneous](#miscellaneous) - [Cleanup](#cleanup)

- [Metrics](#metrics)

## Initial setup

### Directory setup

Go to the `git/bandwidth-marketplace` directory and run the following commands:

- For setting up Consumer:

```
make cons-clean-init
```

- For setting up Provider:

```
make prov-clean-init
```

- For setting up Magma

```
make magma-clean-init
```

## Building and starting nodes


1. Set up a network called `testnet0` for each of these node containers to talk to each other.

 ```
docker network create --driver=bridge --subnet=198.18.0.0/15 --gateway=198.18.0.255 testnet0
```

2. Go to the `git/bandwidth-marketplace` directory to build containers using:


- For building Consumer:

```
make cons-build
```

- For building Provider:

```
make prov-build
```

- For building Magma:

```
make magma-build
```

**_Note: you can skip building step if you want to use existing docker images._**

3. To start nodes working in the `git/bandwidth-marketplace` run the following commands:

- For starting consumer:

```
CONSUMER_ID=N make cons-start
```

- For starting provider:

```
PROVIDER_ID=N make prov-start
```

- For starting magma:

```
make magma-start
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

## Metrics

For metrics collecting Prometheus is used.
See [docs](https://prometheus.io/download/) for more details.

For collecting metrics about GRPC processes used default:
https://github.com/grpc-ecosystem/go-grpc-prometheus

### Magma

Address of the Prometheus strongly configured in `git/bandwidth-marketplace/magma-b0docker-compose.yml` file. 
The Prometheus server is running at:
`http://localhost:9079`.

Supported metrics:

`magma_session_data_uploaded` - used for collecting info about all data, uploaded by user, for each session.
Labels: `session_id`.

`magma_session_data_downloaded` - used for collecting info about all data, downloaded by user, for each session.
Labels: `session_id`.

`magma_session_started` - used for counts all started sessions.

`magma_session_stopped` - used for counts all stopped sessions. 

### Consumer

Address of the Prometheus strongly configured in `git/bandwidth-marketplace/cons-b0docker-compose.yml` file. 
The Prometheus server is running at:
`http://localhost:907${CONSUMER_INDEX}`.

Supported metrics:

`proxy_terms_accepted` - used for counts all accepted terms.
Labels: `user_id`.

`proxy_terms_declined` - used for counts all declined terms.
Labels: `user_id`.

### Provider

Address of Prometheus strongly configured in `git/bandwidth-marketplace/prov-b0docker-compose.yml` file. 
The Prometheus server is running at:
`http://localhost:908${PROVIDER_INDEX}`.

Supported metrics:

`proxy_acknowledgment_verified` - used for counts all verified acknowledgments.

`proxy_acknowledgment_unverified` - used for counts all unverified acknowledgments.

`proxy_session_data_uploaded` - used for collecting info about all data, uploaded by users.
Labels: `session_id`.

`proxy_session_data_downloaded` - used for collecting info about all data, downloaded by users.
Labels: `session_id`.