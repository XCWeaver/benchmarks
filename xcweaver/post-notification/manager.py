#!/usr/bin/env python3
import argparse
import datetime
import json
import os
import re
import shutil
import socket
import statistics
import sys
import threading
import time
from plumbum import FG
from tqdm import tqdm
import requests
import random
import string
import googleapiclient.discovery
from google.oauth2 import service_account
from time import sleep
import yaml
import redis
import hdrh.dump
import hdrh.histogram
import hdrh
import io


APP_PORT                    = 12345
NUM_DOCKER_SWARM_SERVICES   = 8
NUM_DOCKER_SWARM_NODES      = 3
BASE_DIR                    = os.path.dirname(os.path.realpath(__file__))
APP_FOLDER_NAME             = "post-notification"

# -----------
# GCP profile
# -----------
# TBD
GCP_PROJECT_ID                  = None
GCP_USERNAME                    = None
GCP_CREDENTIALS                 = None
GCP_COMPUTE                     = None
GCP_REDIS                       = None

# ---------------------
# GCP app configuration
# ---------------------
# same as in terraform
APP_FOLDER_NAME           = "post-notification"
GCP_INSTANCE_APP_WRK2     = "xcweaver-pn-app-wrk2"
GCP_INSTANCE_DB_MANAGER   = "xcweaver-pn-db-manager"
GCP_INSTANCE_DB_EU        = "xcweaver-pn-db-eu"
GCP_INSTANCE_DB_US        = "xcweaver-pn-db-us"
GCP_MEMORY_STORAGE_EU = "memorystore-primary"
GCP_MEMORY_STORAGE_US = "memorystore-standby"
GCP_ZONE_MANAGER          = "europe-west6-a"
GCP_ZONE_EU               = "europe-west6-a"
GCP_ZONE_US               = "us-central1-a"
GCP_REGION_EU               = "europe-west6"
GCP_REGION_US               = "us-central1"

# --------------------
# Helpers
# --------------------

def load_gcp_profile():
  import yaml
  global GCP_PROJECT_ID, GCP_USERNAME, GCP_COMPUTE, GCP_REDIS
  try:
    with open('gcp/config.yml', 'r') as file:
      config = yaml.safe_load(file)
      GCP_PROJECT_ID  = str(config['project_id'])
      GCP_USERNAME    = str(config['username'])
    GCP_CREDENTIALS   = service_account.Credentials.from_service_account_file("gcp/credentials.json")
    GCP_COMPUTE = googleapiclient.discovery.build('compute', 'v1', credentials=GCP_CREDENTIALS)
    GCP_REDIS = googleapiclient.discovery.build('redis', 'v1', credentials=GCP_CREDENTIALS)
  except Exception as e:
      print(f"[ERROR] error loading gcp profile: {e}")
      exit(-1)

def display_progress_bar(duration, info_message):
  print(f"[INFO] {info_message} for {duration} seconds...")
  for _ in tqdm(range(int(duration))):
    sleep(1)

def get_instance_host(instance_name, zone):
  instance = GCP_COMPUTE.instances().get(project=GCP_PROJECT_ID, zone=zone, instance=instance_name).execute()
  network_interface = instance['networkInterfaces'][0]
  # public, private
  return network_interface['accessConfigs'][0]['natIP']

def get_redis_instance_host(instance_name, region):
  instance_name = f'projects/{GCP_PROJECT_ID}/locations/{region}/instances/{instance_name}'
  instance = GCP_REDIS.projects().locations().instances().get(name=instance_name).execute()
  
  return instance.get('host')

def is_port_open(address, port):
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    result = sock.connect_ex((address, port))
    sock.close()
    return result == 0

def get_consistency_window(port):
  url = f'http://localhost:{port}/consistency_window'
  response = requests.get(url)

  if response.status_code == 200:
      data = json.loads(response.text)
      response = requests.get(url)
      if response.status_code == 200:
        data2 = json.loads(response.text)
        while data2 == data:
          response = requests.get(url)
          if response.status_code == 200:
            data2 = json.loads(response.text)
      #sometimes only one réplica read notifications from the queue
      if data == None:
        values = data2
      elif data2 == None:
        values = data
      else:
        values = data + data2

      if values == None:
        return None
      # Calculate the median
      median = statistics.median(values)
      return median
  else:
      print(f"Failed to get data: {response.status_code}")
      return None

def metrics(deployment, timestamp):
  from plumbum.cmd import xcweaver
  import re

  if timestamp== None:
    timestamp = datetime.datetime.now().strftime("%Y-%m-%d_%H:%M:%S")

  pattern = re.compile(r'^.*│.*│.*│.*│\s*(\d+\.?\d*)\s*│.*$', re.MULTILINE)

  def get_filter_metrics(metric_name):
    return xcweaver['multi', 'metrics', metric_name]()

  # wkr2 api
  inconsitencies_metrics = get_filter_metrics('sn_inconsistencies')
  inconsistencies_count = sum(int(value) for value in pattern.findall(inconsitencies_metrics))
  requests_metrics = get_filter_metrics('requests')
  requests = sum(int(value) for value in pattern.findall(requests_metrics))
  pc_inconsistencies = "{:.2f}".format((inconsistencies_count / requests) * 100)
  write_post_duration_metrics = get_filter_metrics('sn_write_post_duration_ms')
  write_post_duration_metrics_values = pattern.findall(write_post_duration_metrics)
  print(write_post_duration_metrics_values)
  write_post_duration_avg_ms = sum(float(value) for value in write_post_duration_metrics_values if value != 0)/2 if write_post_duration_metrics_values else 0
  notifications_sent_metrics = get_filter_metrics('sn_notificationsSent')
  notifications_sent = sum(int(value) for value in pattern.findall(notifications_sent_metrics))
  notifications_received_metrics = get_filter_metrics('notificationsReceived')
  notifications_received = sum(int(value) for value in pattern.findall(notifications_received_metrics))
  percentage_notifications_received = "{:.2f}".format((notifications_received / notifications_sent) * 100)
  read_post_duration_metrics = get_filter_metrics('sn_read_post_duration_ms')
  read_post_duration_metrics_values = pattern.findall(read_post_duration_metrics)
  print(read_post_duration_metrics_values)
  read_post_duration_avg_ms = sum(float(value) for value in read_post_duration_metrics_values if value != 0)/2 if read_post_duration_metrics_values else 0
  queue_duration_metrics = get_filter_metrics('sn_queue_duration_ms')
  queue_duration_metrics_values = pattern.findall(queue_duration_metrics)
  print(queue_duration_metrics_values)
  queue_duration_avg_ms = sum(float(value) for value in queue_duration_metrics_values if value != 0)/2 if queue_duration_metrics_values else 0
  consistency_window_metrics = get_filter_metrics('sn_consistency_window_ms')
  consistency_window_metrics_values = pattern.findall(consistency_window_metrics)
  print(consistency_window_metrics_values)
  consistency_window_avg_ms = sum(float(value) for value in consistency_window_metrics_values if value != 0)/2 if consistency_window_metrics_values else 0

  consistency_window_median_ms = get_consistency_window(12346)
  
  consistency_window_median_ms = "{:.2f}".format(consistency_window_median_ms)
  read_post_duration_avg_ms = "{:.2f}".format(read_post_duration_avg_ms)
  queue_duration_avg_ms = "{:.2f}".format(queue_duration_avg_ms)
  consistency_window_avg_ms = "{:.2f}".format(consistency_window_avg_ms)
  write_post_duration_avg_ms = "{:.2f}".format(write_post_duration_avg_ms)

  results = {
    'num_requests': int(requests),
    'num_notifications_received': int(notifications_received),
    'per_notifications_received': float(percentage_notifications_received),
    'num_inconsistencies': int(inconsistencies_count),
    'per_inconsistencies': float(pc_inconsistencies),
    'avg_write_post_duration_ms': float(write_post_duration_avg_ms),
    'avg_read_post_duration_ms': float(read_post_duration_avg_ms),
    'avg_queue_duration_ms': float(queue_duration_avg_ms),
    'avg_consistency_window_ms': float(consistency_window_avg_ms),
    'consistency_window_median_ms': float(consistency_window_median_ms),
  }

  print(results)

  filepath = f"evaluation/{deployment}/{timestamp}/metrics.yml"
  os.makedirs(os.path.dirname(filepath), exist_ok=True)
  with open(filepath, 'w') as outfile:
    yaml.dump(results, outfile, default_flow_style=False)
  print(yaml.dump(results, default_flow_style=False))
  print(f"[INFO] evaluation results saved at {filepath}")


def run_workload(duration, url, rate, throughputs, latencies_results, postIds, requests_executed_list, index):

  start_time = time.time()
  interval = 1.0 / int(rate)

  def execute_request(session, url, params):
    request_time_ms = int(time.time() * 1000)
    try:
      with session.get(url, params=params) as response:
        if response.status_code != 200:
          print(f"Request failed with status code {response.status_code}.")
          return
        latencies.append(int(time.time() * 1000) - request_time_ms)
        postIds[index].append(response.text)
    except requests.RequestException as e:
      print(f"Request failed: {e}")
    except json.JSONDecodeError as e:
      print(f"JSON decoding failed: {e}")
    
  latencies = []
  start_time = time.time()
  interval = 1.0 / int(rate)
  requests_executed = 0
  with requests.Session() as session:
    first_request_time = int(time.time() * 1000)
    while True:
      post = ''.join(random.choice(string.ascii_letters + string.digits) for _ in range(15))
      params = {
        "post": post
      }

      execute_request(session, url, params)
      requests_executed += 1
      
      elapsed_time = time.time() - start_time
      if elapsed_time >= int(duration):
        break

      time.sleep(interval - ((time.time() - start_time) % interval))
    
  test_time = int(time.time() * 1000) - first_request_time
  throughput = requests_executed / (test_time / 1000)

  throughputs[index] = throughput
  latencies_results[index] = latencies
  requests_executed_list[index] = requests_executed

def gen_ansible_vars(workload_timestamp=None, deployment_type=None, deployment_folder=None, duration=None, threads=None, rate=None, host_eu=None, host_us=None):
  import yaml

  with open('deploy/ansible/templates/vars.yml', 'r') as file:
    data = yaml.safe_load(file)

  data['base_dir'] = BASE_DIR
  data['workload_timestamp'] = workload_timestamp if workload_timestamp else None
  data['deployment_type'] = deployment_type if deployment_type else None
  data['deployment_folder'] = deployment_folder
  data['duration'] = duration if duration else None
  data['threads'] = threads if threads else None
  data['rate'] = rate if rate else None
  data['host_eu'] = host_eu if host_eu else None
  data['host_us'] = host_us if host_us else None

  with open('deploy/tmp/ansible-vars.yml', 'w') as file:
    yaml.dump(data, file)

def gen_ansible_inventory_gcp():
  from jinja2 import Environment, FileSystemLoader
  import textwrap

  template = Environment(loader=FileSystemLoader('.')).get_template( "deploy/ansible/templates/inventory.cfg")
  inventory = template.render({
    'username':         GCP_USERNAME,
    'host_db_manager':  get_instance_host(GCP_INSTANCE_DB_MANAGER, GCP_ZONE_MANAGER),
    'host_db_eu':       get_instance_host(GCP_INSTANCE_DB_EU, GCP_ZONE_EU),
    'host_db_us':       get_instance_host(GCP_INSTANCE_DB_US, GCP_ZONE_US),
    'host_app_wrk2':    get_instance_host(GCP_INSTANCE_APP_WRK2, GCP_ZONE_MANAGER),
  })

  filename = "deploy/tmp/ansible-inventory.cfg"
  with open(filename, 'w') as f:
    f.write(textwrap.dedent(inventory))
  print(f"[INFO] generated ansible inventory for GCP at '{filename}'")

def gen_ansible_config():
  from jinja2 import Environment, FileSystemLoader
  from plumbum.cmd import cp
  import textwrap
  import os

  # ensure that public key exists
  ssh_key_path = os.path.expanduser("~/.ssh/google_compute_engine")
  if not os.path.exists(ssh_key_path):
    print(f"[ERROR] google compute engine public key not found at '{ssh_key_path}'")
    exit(-1)

  template = Environment(loader=FileSystemLoader('.')).get_template( "deploy/ansible/templates/ansible.cfg")
  inventory = template.render({
    'gcp_ssh_key_path': ssh_key_path,
  })

  path1 = "deploy/tmp/ansible.cfg"
  with open(path1, 'w') as f:
    f.write(textwrap.dedent(inventory))
  print(f"[INFO] generated ansible config at '{path1}'")

  path2 = os.path.expanduser("~/.ansible.cfg")
  cp[path1, path2] & FG
  print(f"[INFO] copied ansible config to '{path2}'")

# --------------------
# GCP
# --------------------

def gcp_configure():
  from plumbum.cmd import gcloud

  try:
    print("[INFO] configuring firewalls")
    # xcweaver-pn-socialnetwork:
    # tcp ports: 12345, 80
    # xcweaver-pn-storage:
    # tcp ports: 6379,6380,27017,27018,3307,3308,15672,5672,15673,5673
    # xcweaver-pn-swarm:
    # tcp ports: 2376,2377,7946
    # udp ports: 4789,7946
    firewalls = {
      'xcweaver-pn-socialnetwork': 'tcp:12345,tcp:80',
      'xcweaver-pn-storage': 'tcp:6379,tcp:6380,tcp:27017,tcp:27018,tcp:3307,tcp:3308,tcp:15672,tcp:5672,tcp:15673,tcp:5673',
      'xcweaver-pn-swarm': 'tcp:2376,tcp:2377,tcp:7946,udp:4789,udp:7946'
    }

    for name, rules in firewalls.items():
      gcloud['compute', 
            '--project', GCP_PROJECT_ID, 'firewall-rules', 'create', 
            f'{name}',
            '--direction=INGRESS',
            '--priority=100',
            '--network=default',
            '--action=ALLOW',
            f'--rules={rules}',
            '--source-ranges=0.0.0.0/0'] & FG
  except Exception as e:
    print(f"[ERROR] could not configure firewalls: {e}\n\n")

def gcp_deploy():
  from plumbum.cmd import terraform, cp, ansible_playbook

  terraform['-chdir=./deploy/terraform', 'init'] & FG
  terraform['-chdir=./deploy/terraform', 'apply', '-auto-approve'] & FG

  display_progress_bar(30, "waiting for all machines to be ready")

  cp["deploy/ansible/templates/ansible.cfg", os.path.expanduser("~/.ansible.cfg")] & FG
  print("[INFO] copied deploy/ansible/ansible.cfg to ~.ansible.cfg")

  # generate temporary files for this deployment
  os.makedirs("deploy/tmp", exist_ok=True)
  print(f"[INFO] created deploy/tmp/ directory")

  gen_ansible_config()
  # generate ansible inventory with hosts of all gcp machines
  gen_ansible_inventory_gcp()
  # generate ansible inventory with extra variables for current deployment
  gen_ansible_vars()
  
  ansible_playbook["deploy/ansible/playbooks/install-machines.yml", "-i", "deploy/tmp/ansible-inventory.cfg", "--extra-vars", "@deploy/tmp/ansible-vars.yml"] & FG

def gcp_info():
  from plumbum.cmd import gcloud
  gcloud['compute', 'ssh', GCP_INSTANCE_DB_MANAGER, '--command', 'sudo docker node ls'] & FG
  gcloud['compute', 'ssh', GCP_INSTANCE_DB_MANAGER, '--command', 'sudo docker service ls'] & FG

  print("\n--- DATASTORES ---")
  print("storage manager running @", get_instance_host(GCP_INSTANCE_DB_MANAGER, GCP_ZONE_MANAGER))
  print(f"storage in {GCP_ZONE_EU} running @", get_instance_host(GCP_INSTANCE_DB_EU, GCP_ZONE_EU))
  print(f"storage in {GCP_ZONE_US} running @", get_instance_host(GCP_INSTANCE_DB_US, GCP_ZONE_US))
  print("\n--- SERVICES ---")
  print(f"wrk2 in {GCP_ZONE_MANAGER} running @", get_instance_host(GCP_INSTANCE_APP_WRK2, GCP_ZONE_MANAGER))

def gcp_start_datastores():
  from plumbum.cmd import ansible_playbook, gcloud
  gcloud["redis", "instances", "create", "memorystore-primary", "--size=1", "--region=europe-west6", "--tier=STANDARD", "--async"] & FG
  gcloud["redis", "instances", "create", "memorystore-standby", "--size=1", "--region=us-central1", "--tier=STANDARD"] & FG
  ansible_playbook["deploy/ansible/playbooks/start-datastores.yml", "-i", "deploy/tmp/ansible-inventory.cfg", "--extra-vars", "@deploy/tmp/ansible-vars.yml"] & FG
  display_progress_bar(30, "waiting for datastores to initialize")

def gcp_redis_hosts():
  host_eu = get_redis_instance_host(GCP_MEMORY_STORAGE_EU, GCP_REGION_EU)
  host_us = get_redis_instance_host(GCP_MEMORY_STORAGE_US, GCP_REGION_US)

  print(f"memorystore primary host: {host_eu}")
  print(f"memorystore secundary host: {host_us}")

def gcp_update_envoy_file():
  from plumbum.cmd import ansible_playbook
  source = "deploy/memorystorage/envoy.yaml"
  destination = "deploy/tmp/"
  # Copy envoy file to tmp directory
  shutil.copy(source, destination)
  ansible_playbook["deploy/ansible/playbooks/update_envoy_file.yml", "-i", "deploy/tmp/ansible-inventory.cfg", "--extra-vars", "@deploy/tmp/ansible-vars.yml"] & FG

def gcp_replicate_datastores():
  from plumbum.cmd import ansible_playbook
  ansible_playbook["deploy/ansible/playbooks/replicate-databases.yml", "-i", "deploy/tmp/ansible-inventory.cfg", "--extra-vars", "@deploy/tmp/ansible-vars.yml"] & FG
  display_progress_bar(30, "waiting for replication process to be complete")

def gcp_stop_datastores():
  from plumbum.cmd import ansible_playbook, gcloud
  gcloud["redis", "instances", "delete", "memorystore-primary", "--region=europe-west6"] & FG
  gcloud["redis", "instances", "delete", "memorystore-standby", "--region=us-central1"] & FG
  ansible_playbook["deploy/ansible/playbooks/stop-datastores.yml", "-i", "deploy/tmp/ansible-inventory.cfg", "--extra-vars", "@deploy/tmp/ansible-vars.yml"] & FG

def gcp_restart_datastores():
  gcp_stop_datastores()
  gcp_start_datastores()
  
def gcp_clean():
  from plumbum.cmd import terraform
  import shutil

  terraform['-chdir=./deploy/terraform', 'destroy', '-auto-approve'] & FG
  if os.path.exists("deploy/tmp"):
    shutil.rmtree("deploy/tmp")
    print(f"[INFO] removed {BASE_DIR}/deploy/tmp/ directory")
 
def gcp_metrics(timestamp):
  print(f"[INFO] metrics automatically generated at evaluation/gcp after running wrk2")

def write_Ids_to_file(keys):
  filename = 'posIds.txt'
  with open(filename, 'w') as f:
      for key in keys:
          f.write(key + '\n')
  print(f"Post ids written to {filename}.")

def gcp_wrk2_vm(duration, threads, rate, timestamp, host_eu, host_us):
  threads_list = []
  throughputs = []
  latencies = []
  postIds = []
  requestsExecuted = []
  for i in range(int(threads)):
    throughputs.append(0)
    requestsExecuted.append(0)
    latencies.append([])
    postIds.append([])
  for i in range(int(threads)):
    thread = threading.Thread(target=run_workload, args=(duration, f"http://{host_eu}/post_notification", int(int(rate) / int(threads)), throughputs, latencies, postIds, requestsExecuted, i))
    thread.start()
    threads_list.append(thread)

  for thread in threads_list:
    thread.join()

  throughput = 0
  for i in throughputs:
    throughput += i

  throughput_str = f'\n----------------------------------------------------------\n Requests/sec:     {throughput}\n'

  requestsExe = 0
  for i in requestsExecuted:
    requestsExe += i
  
  latencies_results = []
  for i in latencies:
    latencies_results += i

  histogram = hdrh.histogram.HdrHistogram(1, 60*60*1000, 3)
  for latency in latencies_results:
    histogram.record_value(latency)

  histoblob = histogram.encode()
  output = io.BytesIO()
  histogram.dump(histoblob, output)
  median = histogram.get_value_at_percentile(50.0)
  consistencyWindow = gcp_consistency_window(host_us)
  inconsistencies = gcp_inconsistencies(host_us)
  filepath = f"evaluation/gcp/{timestamp}/workload.out"
  os.makedirs(os.path.dirname(filepath), exist_ok=True)
  with open(filepath, "w") as f:
    f.write(output.getvalue().decode() + throughput_str +  f" Latency Median:     {median}\n Consistency Window:     {consistencyWindow}\n Inconsistencies:     {inconsistencies}\n Requests Executed:     {requestsExe}\n")

  postIds_results = []
  for i in postIds:
    postIds_results += i
  write_Ids_to_file(postIds_results)

  print(output.getvalue().decode() + throughput_str + f" Latency Median:     {median}\n Consistency Window:     {consistencyWindow}\n Inconsistencies:     {inconsistencies}\n Requests Executed:     {requestsExe}\n")
  print(f"[INFO] workload results saved at {filepath}")

def gcp_consistency_window(host):
  url = f'http://{host}/consistency_window'
  response = requests.get(url)

  if response.status_code == 200:
      values = json.loads(response.text)

      if values == None or values == []:
        return 0

      # Calculate the median
      median = statistics.median(values)
      print(f"Median: {median}")
      return median
  else:
      print(f"Failed to get consistency window: {response.status_code}")

def gcp_inconsistencies(host):
  url = f'http://{host}/inconsistencies'
  response = requests.get(url)

  if response.status_code == 200:
      inconsistencies = int(response.text)

      if inconsistencies == None:
        return None
      
      return inconsistencies
  else:
      print(f"Failed to get inconsistencies: {response.status_code}")

def gcp_wrk2(duration, threads, rate, host_eu, host_us):
  from plumbum.cmd import ansible_playbook
  timestamp = datetime.datetime.now().strftime("%Y-%m-%d_%H:%M:%S")
  metrics_path = f"{BASE_DIR}/evaluation/gcp/{timestamp}"
  os.makedirs(metrics_path, exist_ok=True)
  gen_ansible_vars(timestamp, 'gcp', None, duration, threads, rate, host_eu, host_us)
  ansible_playbook["deploy/ansible/playbooks/run-wrk2.yml", "-i", "deploy/tmp/ansible-inventory.cfg", "--extra-vars", "@deploy/tmp/ansible-vars.yml"] & FG
  print(f"[INFO] metrics results saved at evaluation/gcp/{timestamp}/ in metrics.yaml file")

def gcp_cluster_eu():
  from plumbum.cmd import gcloud
  gcloud["beta", "container", "--project", GCP_PROJECT_ID, "clusters", "create-auto", "cluster-post-notification-eu", "--region", "europe-west6", "--release-channel", "regular", "--network", f"projects/{GCP_PROJECT_ID}/global/networks/default", "--subnetwork", f"projects/{GCP_PROJECT_ID}/regions/europe-west6/subnetworks/default", "--cluster-ipv4-cidr", "/17", "--binauthz-evaluation-mode=DISABLED"] & FG
  gcloud["container", "clusters", "get-credentials", "cluster-post-notification-eu", "--region", "europe-west6", "--project", GCP_PROJECT_ID] & FG

def gcp_cluster_us():
  from plumbum.cmd import gcloud
  gcloud["beta", "container", "--project", GCP_PROJECT_ID, "clusters", "create-auto", "cluster-post-notification-us", "--region", "us-central1", "--release-channel", "regular", "--network", f"projects/{GCP_PROJECT_ID}/global/networks/default", "--subnetwork", f"projects/{GCP_PROJECT_ID}/regions/us-central1/subnetworks/default", "--cluster-ipv4-cidr", "/17", "--binauthz-evaluation-mode=DISABLED"] & FG
  gcloud["container", "clusters", "get-credentials", "cluster-post-notification-us", "--region", "us-central1", "--project", GCP_PROJECT_ID] & FG

def gcp_posts_sizes():
  host = get_instance_host(GCP_INSTANCE_DB_EU, GCP_ZONE_EU)

  r = redis.StrictRedis(host=host, port=6379, db=0)
  key = r.randomkey()
  size = r.memory_usage(key)
  print(f"Size of the value associated with the key '{key}' in bytes: {size}")

  from pymongo import MongoClient
  import bson

  client = MongoClient('mongodb://localhost:27017/?directConnection=true')  # Update with your MongoDB URI if necessary
  db = client['post-notification']
  collection = db['posts']

  def get_random_document(collection):
      pipeline = [{'$sample': {'size': 1}}]  # Sample one document randomly
      result = list(collection.aggregate(pipeline))
      if result:
          return result[0]
      return None

  random_document = get_random_document(collection)
  print(random_document)
  document_size = len(bson.BSON.encode(random_document))
  print(document_size)

  import mysql.connector

  conn = mysql.connector.connect(
      host='localhost',
      user='root',
      password='root_password',
      database='post_notification',
      port=3307
  )
  cursor = conn.cursor()

  query = "SELECT value, LENGTH(value) as size_in_bytes FROM posts ORDER BY RAND() LIMIT 1"

  # Execute the query
  cursor.execute(query)

  # Fetch the result
  result = cursor.fetchone()
  if result:
      value, size_in_bytes = result
      print(f"Random Value: {value}, Size: {size_in_bytes} bytes")
  else:
      print("No record found.")

  # Close the cursor and connection
  cursor.close()
  conn.close()

# --------------------
# LOCAL
# --------------------

def gcp_consistency_window_kube(host):
  url = f'http://{host}/consistency_window'
  response = requests.get(url)

  if response.status_code == 200:
      values = json.loads(response.text)

      if values == None or values == []:
        return 0

      # Calculate the median
      median = statistics.median(values)
      print(f"Median: {median}")
      return median
  else:
      print(f"Failed to get consistency window: {response.status_code}")

def local_wrk2(duration, threads, rate, host_eu, host_us):
  url = "http://localhost:12345/post_notification"
  threads_list = []
  throughputs = []
  latencies = []
  postIds = []
  requestsExecuted = []
  for i in range(int(threads)):
    throughputs.append(0)
    requestsExecuted.append(0)
    latencies.append([])
    postIds.append([])
  for i in range(int(threads)):
    thread = threading.Thread(target=run_workload, args=(duration, url, int(int(rate) / int(threads)), throughputs, latencies, postIds, requestsExecuted, i))
    thread.start()
    threads_list.append(thread)

  for thread in threads_list:
    thread.join()

  throughput = 0
  for i in throughputs:
    throughput += i

  print(f"throughput: {throughput}")
  throughput_str = f'\n----------------------------------------------------------\n Requests/sec:     {throughput}\n'

  requestsExe = 0
  for i in requestsExecuted:
    requestsExe += i

  latencies_results = []
  for i in latencies:
    latencies_results += i

  histogram = hdrh.histogram.HdrHistogram(1, 60*60*1000, 3)
  for latency in latencies_results:
    histogram.record_value(latency)

  histoblob = histogram.encode()
  output = io.BytesIO()
  histogram.dump(histoblob, output)
  median = histogram.get_value_at_percentile(50.0)
  timestamp = datetime.datetime.now().strftime("%Y-%m-%d_%H:%M:%S")
  filepath = f"evaluation/local/{timestamp}/workload.out"
  os.makedirs(os.path.dirname(filepath), exist_ok=True)
  with open(filepath, "w") as f:
    f.write(output.getvalue().decode() + throughput_str + f"Latency Median:     {median}\n Requests Executed:     {requestsExe}\n")

  postIds_results = []
  for i in postIds:
    postIds_results += i
  write_Ids_to_file(postIds_results)

  print(output.getvalue().decode() + throughput_str + f"Latency Median:     {median}\n Requests Executed:     {requestsExe}\n")
  print(f"[INFO] workload results saved at {filepath}")

  metrics('local', timestamp)

def local_metrics(timestamp):
  metrics('local', timestamp)

def local_consistency_window():
  url = 'http://localhost:12346/consistency_window'
  response = requests.get(url)

  # Check if the request was successful
  if response.status_code == 200:
      # Parse the JSON response
      data = json.loads(response.text)
      response = requests.get(url)
      if response.status_code == 200:
        data2 = json.loads(response.text)
        while data2 == data:
          response = requests.get(url)
          if response.status_code == 200:
            data2 = json.loads(response.text)
      if data == None:
        values = data2
      elif data2 == None:
        values = data
      else:
        values = data + data2

      if values == None:
        return None
      print(values)
      print(len(values))
      average = statistics.mean(values)
      print(f"Average: {average}")

      # Calculate the median
      median = statistics.median(values)
      print(f"Median: {median}")
  else:
      print(f"Failed to get data: {response.status_code}")

def local_storage_deploy():
  print("[INFO] nothing to be done for local")
  exit(0)

def local_storage_build():
  from plumbum.cmd import docker
  docker['build', '-t', 'mongodb-delayed:4.4.6', 'clustering/mongodb-delayed/.'] & FG
  docker['build', '-t', 'mongodb-setup:4.4.6', 'clustering/mongodb-setup/post-storage/.'] & FG
  docker['build', '-t', 'rabbitmq-setup:3.8', 'clustering/rabbitmq-setup/notifier/.'] & FG
  docker['build', '-t', 'redis-setup:latest', 'clustering/redis-im/.'] & FG

def local_storage_run():
  from plumbum.cmd import docker_compose, docker, grep, awk
  docker_compose['up', '-d'] & FG
  print("[INFO] waiting 30 seconds for storages to be ready...")
  for _ in tqdm(range(30)):
      time.sleep(1)
  
  #setup mysql cluster and create database and table inside the mysql cluster
  with open('clustering/mysql/master.sql', 'r') as sql_file:
    sql_content = sql_file.read()
  (docker['exec', '-i', 'mysql1', 'mysql', '-u', 'root', '-proot_password']) << sql_content & FG

  master_status_cmd = docker['exec', 'mysql1', 'mysql', '-u', 'root', '-proot_password', '-e', 'SHOW MASTER STATUS\\G']
  master_status = master_status_cmd()

  master_log_file_cmd = (grep['File:'] << master_status) | (awk['{print $2}'])
  master_log_pos_cmd = (grep['Position:'] << master_status) | (awk['{print $2}'])

  master_log_file = master_log_file_cmd().strip()
  master_log_pos = master_log_pos_cmd().strip()

  with open('clustering/mysql/replica.sql', 'r') as sql_file:
    sql_content = sql_file.read()

  sql_content = sql_content.replace('{{MASTER_LOG_FILE}}', master_log_file)
  sql_content = sql_content.replace('{{MASTER_LOG_POS}}', master_log_pos)

  (docker['exec', '-i', 'mysql2', "mysql", "-u", "root", '-proot_password']) << sql_content & FG


  with open('clustering/mysql/key-value_table.sql', 'r') as sql_file:
    sql_content = sql_file.read()
  (docker['exec', '-i', 'mysql1', 'mysql', '-u', 'root', '-proot_password']) << sql_content & FG

def local_storage_info():
  print("[INFO] nothing to be done for local")
  exit(0)

def local_storage_clean():
  from plumbum.cmd import docker_compose
  docker_compose['down'] & FG

if __name__ == "__main__":
  main_parser = argparse.ArgumentParser()
  command_parser = main_parser.add_subparsers(help='commands', dest='command')

  commands = [
    # gcp
    'configure', 'deploy', 'start-datastores', 'stop-datastores', 'cluster-eu', 'cluster-us',
    'restart-datastores', 'clean', 'info', 'consistency-window', 'consistency-window-kube',
    # datastores
    'storage-run', 'storage-info', 'storage-clean', 'storage-build', 'replicate-datastores', 'update-envoy-file', 'redis-hosts', 'posts-sizes',
    # eval
    'wrk2', 'wrk2-vm', 'metrics',
  ]

  for cmd in commands:
    parser = command_parser.add_parser(cmd)
    parser.add_argument('--local', action='store_true', help="Running in localhost")
    parser.add_argument('--gcp', action='store_true',   help="Running in gcp")
    if cmd == 'wrk2' or cmd =='wrk2-vm':
      parser.add_argument('-d', '--duration', default=30, help="Duration")
      parser.add_argument('-t', '--threads', default=2, help="Number of threads")
      parser.add_argument('-r', '--rate', default=50, help="Number of requests per second")
      parser.add_argument('-hteu', '--host_eu', default="localhost", help="host of the eu load balancer")
      parser.add_argument('-htus', '--host_us', default="localhost", help="host of the us load balancer")
    if cmd == 'wrk2-vm':
      parser.add_argument('-ts', '--timestamp', help="Timestamp of workload")
    if cmd == 'consistency-window-kube':
      parser.add_argument('-ht', '--host', help="host of the us load balancer")
    if cmd == 'metrics':
      parser.add_argument('-t', '--timestamp', help="Timestamp of workload")

  args = vars(main_parser.parse_args())
  command = args.pop('command').replace('-', '_')

  local = args.pop('local')
  gcp = args.pop('gcp')

  if local and gcp or not local and not gcp:
    print("[ERROR] one of --local or --gcp flgs needs to be provided")
    exit(-1)

  if local:
    command = 'local_' + command
  elif gcp:
    load_gcp_profile()
    command = 'gcp_' + command

  print(f"[INFO] ----- {command.upper().replace('_', ' ')} -----\n")
  getattr(sys.modules[__name__], command)(**args)

  print(f"[INFO] done!")
    
