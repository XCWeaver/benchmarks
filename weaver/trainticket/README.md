# TrainTicket @ Service Weaver

Implementation of the [TrainTicket](https://github.com/FudanSELab/train-ticket) benchmark using [Google's Service Weaver](https://serviceweaver.dev/). Drawing inspiration from the [Blueprint](https://gitlab.mpi-sws.org/cld/blueprint)'s repository.

# Table of Contents
- [rainTicket @ Service Weaver](#trainticket--service-weaver)
- [Table of Contents](#table-of-contents)
- [1. Requirements](#1-requirements)
  - [1.1. Install Python Dependencies](#11-install-python-dependencies)
- [2. Configuration](#2-configuration)
  - [2.1. GCP Configuration](#21-gcp-configuration)
- [3. Application Deployment](#3-application-deployment)
  - [3.1. GCP Deployment](#31-gcp-deployment)
  - [3.2. Local Deployment using Weaver Multi Process](#32-local-deployment-using-weaver-multi-process)
- [4. Complementary Information](#4-complementary-information)
  - [4.1. Manually Testing HTTP Requests](#41-manually-testing-http-requests)


# 1. Requirements

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


## 1.1. Install Python Dependencies

Install python packages to use the `manager.py` script:
```zsh
pip install -r requirements.txt
```


# 2. Configuration

## 2.1. GCP Configuration

1. Ensure that you have a GCP project created and setup as your default project in `gcloud cli`:
``` zsh
# initialize gcloud (set default region e.g. as europe-west3-a)
gcloud init
# list all projects
gcloud projects list
# select the desired project id
gcloud config set project YOUR_PROJECT_ID
# verify it is now set as default
gcloud config get-value project
```
1. Ensure that [Compute Engine API](https://console.cloud.google.com/marketplace/product/google/compute.googleapis.com) is enabled in GCP
2. Go to `tese/weaver/trainticket/gcp/config.yml` and place you GCP `project_id` and any desired `username` for accessing GCP machines using that hostname
3. Configure GCP firewalls and SSH keys for Compute Engine
``` zsh
./manager.py configure --gcp
```
4. Setup a new Service Account key for authenticating (for more information: https://developers.google.com/workspace/guides/create-credentials)
    - Go to [IAM & Admin -> Service Accounts](https://console.cloud.google.com/iam-admin/serviceaccounts) of your project
    - Select your compute engine default service account
    - Go to the keys tab and select `ADD KEY` to create a new key in JSON
    - Place your JSON file as `credentials.json` in `tese/weaver/trainticket/gcp/credentials.json`

# 3. Application Deployment

## 3.1. GCP Deployment

Use the following commands to deploy and start the application in GCP machines and display info for hosts of GCP machines:
``` zsh
./manager.py deploy --gcp
./manager.py start --gcp
./manager.py info --gcp
```

Run workload and automatically gather metrics to `evaluation` directory. If not specified, the default parameters are 2 threads, 2 clients, 30 duration (in seconds), 50 rate
``` zsh
./manager.py --gcp wrk2 -t THREADS -c CLIENTS -d DURATION -r RATE
./manager.py --gcp wrk2
```

> **_NOTE:_**  It will be generated a file inside evaluation/gcp that contains the metrics

Restart datastores and application:
``` zsh
./manager.py --gcp restart
```

Otherwise, to clean all gcp resources at the end, do:

``` zsh
./manager.py --gcp clean
```


## 3.2. Running Locally with Service Weaver in Multi Process

Deploy datastores:

``` zsh
./manager.py storage-run --local
```

Start uptrace:

``` zsh
git clone https://github.com/uptrace/uptrace.git
cd uptrace/examples/docker
sudo docker-compose pull
sudo docker-compose up -d
```


Deploy and run application:
``` zsh
go generate
go build
weaver multi deploy weaver-local.toml
```

Run workload and gather metrics:

``` zsh
./manager.py wrk2 --local
```

> **_NOTE:_**  It will be generated a file inside evaluation/gcp that contains the metrics

Gather metrics:
``` zsh
./manager.py metrics --local
```


# 4. Complementary Information

## 4.1. Manual Testing HTTP Requests

**Register User**: {token, username, password, gender, documentType, documentNum, email} [accountId]

``` zsh
curl -X POST "localhost:9000/wrk2-api/user/registerUser" -d "token=TOKEN&username=USERNAME&password=PASSWORD&gender=GENDER&documentType=DOCUMENT_TYPE&documentNum=DOCUMENT_NUM&email=EMAIL"

# e.g.
curl -X POST "localhost:9000/wrk2-api/user/registerUser" -d "token=&username=bob&password=mypassword&gender=0&documentType=1&documentNum=1234567&email=xcweaver@gmail.com"
```

**Login**: {username, password, verificationCode} [token]

``` zsh
curl -X POST "localhost:9000/wrk2-api/user/login" -d "username=USERNAME&password=PASSWORD&verificationCode=VERIFICATION_CODE"

# e.g.
curl -X POST "localhost:9000/wrk2-api/user/login" -d "username=bob&password=mypassword&verificationCode="
```

**Add Order**: {token, boughtDate, travelDate, accountId, contactsName, documdocumentTypentNum, contactsDocumentNumber, trainNumber, coachNumber, seatClass, seatNumber, from, to, status, price} [OrderID]

``` zsh
curl -X POST "localhost:9000/wrk2-api/admin/adminAddOrder" -d "token=TOKEN&boughtDate=BOUGHT_DATE&travelDate=TRAVEL_DATE&accountId=ACCOUNT_ID&contactsName=CONTACTS_NAME&documentType=DOCUMENT_TYPE&contactsDocumentNumber=CONTACTS_DOCUMENT_NUMBER&trainNumber=TRAIN_NUMBER&coachNumber=COACH_NUMBER&seatClass=SEAT_CLASS&seatNumber=SEAT_NUMBER&from=FROM&to=TO&status=STATUS&price=PRICE"

# e.g.
curl -X POST "localhost:9000/wrk2-api/admin/adminAddOrder" -d "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOiJlN2ZlOTkwNy02MjViLTRiNGItOGQ4OC1mNmUxY2JiMWYwMjUiLCJVc2VybmFtZSI6ImJvYiIsIlRpbWVzdGFtcCI6MTcxNzU5NDkxNTIxOSwiVHRsIjozNjAwLCJSb2xlIjoicm9sZTEiLCJleHAiOjE3MTc1OTg1MTV9.zy5erWHnQqNoYZkg3PcUxH_iS_jMTYoPcaQ131DCuE4&boughtDate=17/05/2024&travelDate=17/08/2024&accountId=e7fe9907-625b-4b4b-8d88-f6e1cbb1f025&contactsName=bob&documentType=1&contactsDocumentNumber=1234567&trainNumber=12345&coachNumber2=&seatClass=1&seatNumber=34&from=Lisboa&to=Guarda&status=0&price=10.20"
```

**Cancel Ticket**: {token, orderId, loginId}

``` zsh
curl -X POST "localhost:9000/wrk2-api/user/cancelTicket" -d "token=TOKEN&orderId=ORDER_ID&loginId=LOGIN_ID"

# e.g.
curl -X POST "localhost:9000/wrk2-api/user/cancelTicket" -d "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOiJlN2ZlOTkwNy02MjViLTRiNGItOGQ4OC1mNmUxY2JiMWYwMjUiLCJVc2VybmFtZSI6ImJvYiIsIlRpbWVzdGFtcCI6MTcxNzU5NDkxNTIxOSwiVHRsIjozNjAwLCJSb2xlIjoicm9sZTEiLCJleHAiOjE3MTc1OTg1MTV9.zy5erWHnQqNoYZkg3PcUxH_iS_jMTYoPcaQ131DCuE4&orderId=71e2d5da-76ac-4404-a97e-9a26659511fd&loginId=e7fe9907-625b-4b4b-8d88-f6e1cbb1f025"
```
