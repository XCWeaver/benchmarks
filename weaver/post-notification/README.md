# Post-Notification / Service Weaver

Implementation of a geo-replicated Post-Notification application that shows the occurence of cross-service inconsistencies.

## Requirements

- [Terraform >= v1.6.6](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)
- [Ansible >= v2.15.2](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html)
- [GCloud Cli](https://cloud.google.com/sdk/docs/install)
- [Golang >= 1.21](https://go.dev/doc/install)
```zsh
# install Golang v1.21.5
sudo wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin:$HOME/.go/bin' >> ~/.bashrc
source ~/.bashrc
```
- [Service Weaver >= v0.22.0](https://serviceweaver.dev/docs.html#installation)
```zsh
# install Weaver 0.24.3
go install github.com/ServiceWeaver/weaver/cmd/weaver@v0.24.3
```


Install python packages to use the `manager` script:
```zsh
pip install -r requirements.txt
```


## LOCAL Deployment

### Running Locally with Service Weaver in Multi Process

Deploy datastores:

``` zsh
./manager.py storage-build --local
./manager.py storage-run --local
```

> **_NOTE:_**  This command will deploy 4 different clusters of 4 different datastore types ([RabbitMq](https://www.rabbitmq.com/), [Redis](https://redis.io/), [MongoDB](https://www.mongodb.com/) and [MySQL](https://www.mysql.com/)).

You can deploy the Post_notification application using tree different storage types

#### Deploy Post-Notification with a Redis datastore:

``` zsh
cd redis-rabbitmq
```

Deploy eu_deployment:
``` zsh
cd eu_deployment
go build
weaver multi deploy weaver.toml
```

On another terminal deploy us_deployment:

``` zsh
cd ../us_deployment
go build
weaver multi deploy weaver.toml
```

Run benchmark:

``` zsh
cd ..
./manager.py wrk2 --local
```

Gather metrics:
``` zsh
./manager.py metrics --local
```

#### Deploy Post-Notification with a MongoDB datastore:

``` zsh
cd mongo-rabbitmq
```

Deploy eu_deployment:
``` zsh
cd eu_deployment
go build
weaver multi deploy weaver.toml
```

Deploy us_deployment:

``` zsh
cd ../us_deployment
go build
weaver multi deploy weaver.toml
```

Run benchmark:

``` zsh
cd ..
./manager.py wrk2 --local
```

Gather metrics:
``` zsh
./manager.py metrics --local
```

#### Deploy Post-Notification with a MySQL datastore:

``` zsh
cd mysql-rabbitmq
```

Deploy eu_deployment:
``` zsh
cd eu_deployment
go build
weaver multi deploy weaver.toml
```

Deploy us_deployment:

``` zsh
cd ../us_deployment
go build
weaver multi deploy weaver.toml
```

Run benchmark:

``` zsh
cd ..
./manager.py wrk2 --local
```

Gather metrics:
``` zsh
./manager.py metrics --local
```

### Running on GCP with Service Weaver kube deployer

Deploy the vm to run the wrk2 and the datastores:
``` zsh
./manager.py deploy --gcp
./manager.py start-datastores --gcp
```

Get redis hosts:
``` zsh
./manager.py redis-hosts --gcp
```

Go to deploy/memorystorage and on the file envoy.yaml update the line ... with the memorystore primary host and the line ... with the memorystore secundary host

Update envoy.yaml on wrk2 vm and replicate datastores:
``` zsh
./manager.py update-envoy-file --gcp
./manager.py replicate-datastores --gcp
```

#### Deploy Post-Notification with a Redis datastore:

#### EU region deployment:

Create and connect to a Kubernetes cluster:
``` zsh
./manager.py cluster-eu --gcp
```

Get redis hosts:
``` zsh
./manager.py redis-hosts --gcp
```

``` zsh
cd redis-rabbitmq
cd eu_deployment
```
Update weaver.toml file with the memorystore primary host on the rabbitmq_address and redis_address fields.

Deploy eu_deployment:
``` zsh
go build
weaver kube deploy config_redis.yaml
kubectl apply -f <file>
``` 

Get loadbalencer Eu address:
``` zsh
kubectl get all
```

#### US region deployment:

Go back to post-notification directory:
``` zsh
cd ..
```

Create and connect to a Kubernetes cluster:
``` zsh
./manager.py cluster-us --gcp
```

Get redis hosts:
``` zsh
./manager.py redis-hosts --gcp
```

``` zsh
cd redis-rabbitmq
cd us_deployment
```
Update weaver.toml file with the memorystore secondary host on the rabbitmq_address and redis_address fields.

Deploy us_deployment:
``` zsh
go build
weaver kube deploy config_redis.yaml
kubectl apply -f <file>
``` 

Get loadbalencer Us address:
``` zsh
kubectl get all
```

Run benchmark:

``` zsh
cd ..
./manager.py wrk2 --gcp -hteu <host-eu> -htus <host-us>
```
> **_NOTE:_**  host-eu is the host of the Eu loadbalancer and host-us is the host of the Us loadbalancer.

#### Deploy Post-Notification with a MongoDB datastore:

#### EU region deployment:

Create and connect to a Kubernetes cluster:
``` zsh
./manager.py cluster-eu --gcp
```

Get mongoDB hosts:
``` zsh
./manager.py info --gcp
```

``` zsh
cd mongo-rabbitmq
cd eu_deployment
```
Update weaver.toml file with the storage in europe-west6-a host on the rabbitmq_address and mongo_address fields.

Deploy eu_deployment:
``` zsh
go build
weaver kube deploy config_mongo.yaml
kubectl apply -f <file>
``` 

Get loadbalencer Eu address:
``` zsh
kubectl get all
```

#### US region deployment:

Go back to post-notification directory:
``` zsh
cd ..
```

Create and connect to a Kubernetes cluster:
``` zsh
./manager.py cluster-us --gcp
```

Get mongoDB hosts:
``` zsh
./manager.py info --gcp
```

``` zsh
cd mongo-rabbitmq
cd us_deployment
```
Update weaver.toml file with the storage in us-central1-a host on the rabbitmq_address and mongo_address fields.

Deploy us_deployment:
``` zsh
go build
weaver kube deploy config_mongo.yaml
kubectl apply -f <file>
``` 

Get loadbalencer Us address:
``` zsh
kubectl get all
```

Run benchmark:

``` zsh
cd ..
./manager.py wrk2 --gcp -hteu <host-eu> -htus <host-us>
```
> **_NOTE:_**  host-eu is the host of the Eu loadbalancer and host-us is the host of the Us loadbalancer.

#### Deploy Post-Notification with a MySQL datastore:

#### EU region deployment:

Create and connect to a Kubernetes cluster:
``` zsh
./manager.py cluster-eu --gcp
```

Get MySQL hosts:
``` zsh
./manager.py info --gcp
```

``` zsh
cd mysql-rabbitmq
cd eu_deployment
```
Update weaver.toml file with the storage in europe-west6-a host on the rabbitmq_address and mysql_address fields.

Deploy eu_deployment:
``` zsh
go build
weaver kube deploy config_mysql.yaml
kubectl apply -f <file>
``` 

Get loadbalencer Eu address:
``` zsh
kubectl get all
```

#### US region deployment:

Go back to post-notification directory:
``` zsh
cd ..
```

Create and connect to a Kubernetes cluster:
``` zsh
./manager.py cluster-us --gcp
```

Get mysql hosts:
``` zsh
./manager.py info --gcp
```

``` zsh
cd mysql-rabbitmq
cd us_deployment
```
Update weaver.toml file with the storage in us-central1-a host on the rabbitmq_address and mysql_address fields.

Deploy us_deployment:
``` zsh
go build
weaver kube deploy config_mysql.yaml
kubectl apply -f <file>
``` 

Get loadbalencer Us address:
``` zsh
kubectl get all
```

Run benchmark:

``` zsh
cd ..
./manager.py wrk2 --gcp -hteu <host-eu> -htus <host-us>
```
> **_NOTE:_**  host-eu is the host of the Eu loadbalancer and host-us is the host of the Us loadbalancer.

### Additional

#### Manual Testing of HTTP Requests locally

**Publish Post**: {post}

``` zsh
curl "localhost:12345/post_notification?post=POST"

# e.g.
curl "localhost:12345/post_notification?post=my_first_post"
```

#### Manual Testing of HTTP Requests on Kubernetes

**Publish Post**: {post}

``` zsh
curl "<host>/post_notification?post=POST"

# e.g.
curl "232.345.56.12/post_notification?post=my_first_post"
```
