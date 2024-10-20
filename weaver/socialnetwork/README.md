# DeathStarBench SocialNetwork @ Service Weaver

Implementation of [DeathStarBench](https://github.com/delimitrou/DeathStarBench) SocialNetwork benchmark using Service Weaver framework, drawing inspiration from the [Blueprint](https://gitlab.mpi-sws.org/cld/blueprint)'s repository.

# Table of Contents
- [rainTicket @ Service Weaver](#deathstarbench-socialnetwork--service-weaver)
- [Table of Contents](#table-of-contents)
- [1. Requirements](#1-requirements)
  - [1.1. Install Python Dependencies](#11-install-python-dependencies)
- [2. Configuration](#2-configuration)
  - [2.1. GCP Configuration](#21-gcp-configuration)
  - [2.1. Workload Configuration](#21-workload-configuration)
- [3. Application Deployment](#3-application-deployment)
  - [3.1. GCP Deployment](#31-gcp-deployment)
- [4. Complementary Information](#4-complementary-information)
  - [4.1. Manual Testing of HTTP Workload Generator](#41-manual-testing-of-http-workload-generator)
  - [4.1. Manual Testing of HTTP Requests](#41-manual-testing-of-http-requests)

# 1. Requirements

- [Terraform >= v1.6.6](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)
- [Ansible >= v2.15.2](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html)
- [GCloud Cli](https://cloud.google.com/sdk/docs/install)
- [Golang >= 1.21.5](https://go.dev/doc/install)
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


## 1.1. Install Python Dependencies

Install python packages to use the `manager.py` script:
```zsh
pip install -r requirements.txt
```

**OPTIONAL**: if HTTP Workload Generator (wrk2) is going to run locally
- `Lua >= 5.1.5`
- `LuaRocks >= 3.8.0` (with `lua-json`, `luasocket`, `penlight` packages)
- `OpenSSL >= 3.0.2`

```zsh
sudo apt-get install -y libssl-dev luarocks
sudo luarocks install lua-json
sudo luarocks install luasocket
sudo luarocks install penlight
```

# 2. Configuration

## 2.1. GCP Configuration

1. Ensure that you have a GCP project created and setup as your default project in `gcloud cli`:
``` zsh
# list all projects
gcloud projects list
# select the desired project id
gcloud config set project YOUR_PROJECT_ID
# verify it is now set as default
gcloud config get-value project
```
1. Ensure that [Compute Engine API](https://console.cloud.google.com/marketplace/product/google/compute.googleapis.com) is enabled in GCP
2. Go to `weaver-dsb/socialnetwork/gcp/config.yml` and place you GCP `project_id` and any desired `username` for accessing GCP machines using that hostname
3. Create new firewall rules
``` zsh
./manager.py configure --gcp
```
1. Setup a new Service Account key for authenticating (for more information: https://developers.google.com/workspace/guides/create-credentials)
    - Go to [IAM & Admin -> Service Accounts](https://console.cloud.google.com/iam-admin/serviceaccounts) of your project
    - Select your compute engine default service account
    - Go to the keys tab and select `ADD KEY` to create a new key in JSON
    - Place your JSON file as `credentials.json` in `weaver-dsb/socialnetwork/gcp/credentials.json`

## 2.2. Workload Configuration

Generate workload binary:

```zsh
cd wrk2
make
```

# 3. Application Deployment

## 3.1. GCP Deployment

Deploy, and start your application (datastores + services). You can also display some info for docker swarm and hosts of gcp machines
``` zsh
./manager.py deploy --gcp
./manager.py start --gcp
./manager.py info --gcp
```

Run workload and automatically gather metrics to `evaluation` directory:
``` zsh
# default params: 2 threads, 2 clients, 30 duration (in seconds), 50 rate
./manager.py wrk2 --local
# you can also specify other parameters
./manager.py wrk2 --local -t THREADS -c CLIENTS -d DURATION -r RATE
```

Restart datastores and application:
``` zsh
./manager.py restart --gcp
```

Otherwise, to clean all gcp resources at the end, do:

``` zsh
./manager.py clean --gcp
```

# 4. Complementary Information

## 4.1. Manual Testing of HTTP Workload Generator

Compose Posts

```zsh
cd wrk2
./wrk -D exp -t <num-threads> -c <num-conns> -d <duration> -L -s ./scripts/social-network/compose-post.lua http://localhost:12345/wrk2-api/post/compose -R <reqs-per-sec>

# e.g.
./wrk -D exp -t 1 -c 1 -d 1 -L -s ./scripts/social-network/compose-post.lua http://localhost:12345/wrk2-api/post/compose -R 1
```

Read Home Timelines

```zsh
cd wrk2
./wrk -D exp -t <num-threads> -c <num-conns> -d <duration> -L -s ./scripts/social-network/read-home-timeline.lua http://localhost:12345/wrk2-api/home-timeline/read -R <reqs-per-sec>
```

Read User Timelines

```zsh
cd wrk2
./wrk -D exp -t <num-threads> -c <num-conns> -d <duration> -L -s ./scripts/social-network/read-user-timeline.lua http://localhost:12345/wrk2-api/user-timeline/read -R <reqs-per-sec>
```

## 4.2. Manual Testing of HTTP Requests

**Register User**: {username, first_name, last_name, password} [user_id]

``` zsh
curl -X POST "localhost:12345/wrk2-api/user/register" -d "username=USERNAME&user_id=USER_ID&first_name=FIRST_NAME&last_name=LAST_NAME&password=PASSWORD"

# e.g.
curl -X POST "localhost:12345/wrk2-api/user/register" -d "username=ana&user_id=0&first_name=ana1&last_name=ana2&password=123"
curl -X POST "localhost:12345/wrk2-api/user/register" -d "username=bob&user_id=1&first_name=bob1&last_name=bob2&password=123"
```

**Follow User**: [{user_id, followee_id}, {user_name, followee_name}]

``` zsh
curl -X POST "localhost:12345/wrk2-api/user/follow" -d "user_id=USER_ID&followee_id=FOLLOWEE_ID"
OR
curl -X POST "localhost:12345/wrk2-api/user/follow" -d "user_name=USER_NAME&followee_name=FOLLOWEE_AME"

# e.g.
curl -X POST "localhost:12345/wrk2-api/user/follow" -d "user_name=ana&followee_name=bob"
curl -X POST "localhost:12345/wrk2-api/user/follow" -d "user_id=1&followee_id=0"
```

**Unfollow User**: [{user_id, followee_id}, {username, followee_name}]

``` zsh
curl -X POST "localhost:12345/wrk2-api/user/unfollow" -d "user_id=USER_ID&followee_id=FOLLOWEE_ID"

# e.g.
curl -X POST "localhost:12345/wrk2-api/user/unfollow" -d "user_id=1&followee_id=0"
```

**Compose Post**: {user_id, text, username, post_type} [media_types, media_ids]

``` zsh
curl -X POST "localhost:12345/wrk2-api/post/compose" -d "user_id=USER_ID&text=TEXT&username=USER_ID&post_type=POST_TYPE"

# e.g.
curl -X POST "localhost:12345/wrk2-api/post/compose" -d "user_id=0&text=helloworld_0&username=ana&post_type=0&media_types=["png"]&media_ids=[0]"
curl -X POST "localhost:12345/wrk2-api/post/compose" -d "user_id=1&text=helloworld_0&username=username_1&post_type=0&media_types=["png"]&media_ids=[0]"
```

**Read User Timeline**: {user_id} [start, stop]

``` zsh
curl "localhost:12345/wrk2-api/user-timeline/read" -d "user_id=USER_ID"

# e.g.
curl "localhost:12345/wrk2-api/user-timeline/read" -d "user_id=0"
curl "localhost:12345/wrk2-api/user-timeline/read" -d "user_id=1"
```

**Read Home Timeline**: {user_id} [start, stop]

``` zsh
curl "localhost:12345/wrk2-api/home-timeline/read" -d "user_id=USER_ID"

# e.g.
curl "localhost:12345/wrk2-api/home-timeline/read" -d "user_id=1"
curl "localhost:12345/wrk2-api/home-timeline/read" -d "user_id=88"
```

