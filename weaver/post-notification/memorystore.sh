gcloud redis instances create memorystore-primary --size=1 --region=europe-west3 --tier=STANDARD --async

gcloud redis instances create memorystore-standby --size=1 --region=us-central1 --tier=STANDARD --async

#update envoy.yaml

#on db-eu
sudo docker run --rm -d -p 8001:8001 -p 6381:1999 -v $(pwd)/envoy.yaml:/envoy.yaml envoyproxy/envoy:v1.21.0 -c /envoy.yaml

#on db-us
sudo docker run --rm -d -p 8001:8001 -p 6381:2000 -v $(pwd)/envoy.yaml:/envoy.yaml envoyproxy/envoy:v1.21.0 -c /envoy.yaml


#delete
gcloud redis instances delete memorystore-primary --region=europe-west3
gcloud redis instances delete memorystore-standby --region=us-central1