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
# install Weaver 0.22.0
go install github.com/ServiceWeaver/weaver/cmd/weaver@v0.22.0
```


Install python packages to use the `manager` script:
```zsh
pip install -r requirements.txt
```


## LOCAL Deployment

### Running Locally with Service Weaver in Multi Process

Deploy datastores:

``` zsh
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
go generate
go build
weaver multi deploy weaver.toml
```

Deploy us_deployment:

``` zsh
cd ../eu_deployment
go generate
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
go generate
go build
weaver multi deploy weaver.toml
```

Deploy us_deployment:

``` zsh
cd ../eu_deployment
go generate
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
go generate
go build
weaver multi deploy weaver.toml
```

Deploy us_deployment:

``` zsh
cd ../eu_deployment
go generate
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

### Additional

#### Manual Testing of HTTP Requests

**Publish Post**: {post}

``` zsh
curl "localhost:12345/post_notification?post=POST"

# e.g.
curl "localhost:12345/post_notification?post=my_first_post"
```
