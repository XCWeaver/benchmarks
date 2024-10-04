#!/usr/bin/env python3
import argparse
import json
import re
import statistics
import time
import googleapiclient.discovery
from google.oauth2 import service_account
import sys
from plumbum import FG
import requests
import yaml
from time import sleep
from tqdm import tqdm
import datetime
import os
import socket

APP_PORT                    = 12345
NUM_DOCKER_SWARM_SERVICES   = 19
NUM_DOCKER_SWARM_NODES      = 3
BASE_DIR                    = os.path.dirname(os.path.realpath(__file__))

# -----------
# GCP profile
# -----------
# TBD
GCP_PROJECT_ID                  = None
GCP_USERNAME                    = None
GCP_CREDENTIALS                 = None
GCP_COMPUTE                     = None

# ---------------------
# GCP app configuration
# ---------------------
# same as in terraform
APP_FOLDER_NAME           = "socialnetwork"
GCP_INSTANCE_APP_WRK2     = "xcweaver-dsb-app-wrk2"
GCP_INSTANCE_APP_EU       = "xcweaver-dsb-app-eu"
GCP_INSTANCE_APP_US       = "xcweaver-dsb-app-us"
GCP_INSTANCE_DB_MANAGER   = "xcweaver-dsb-db-manager"
GCP_INSTANCE_DB_EU        = "xcweaver-dsb-db-eu"
GCP_INSTANCE_DB_US        = "xcweaver-dsb-db-us"
GCP_ZONE_MANAGER          = "europe-west3-a"
GCP_ZONE_EU               = "europe-west3-a"
GCP_ZONE_US               = "us-central1-a"

# --------------------
# Helpers
# --------------------

def load_gcp_profile():
  import yaml
  global GCP_PROJECT_ID, GCP_USERNAME, GCP_COMPUTE
  try:
    with open('gcp/config.yml', 'r') as file:
      config = yaml.safe_load(file)
      GCP_PROJECT_ID  = str(config['project_id'])
      GCP_USERNAME    = str(config['username'])
    GCP_CREDENTIALS   = service_account.Credentials.from_service_account_file("gcp/credentials.json")
    GCP_COMPUTE = googleapiclient.discovery.build('compute', 'v1', credentials=GCP_CREDENTIALS)
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

def is_port_open(address, port):
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    result = sock.connect_ex((address, port))
    sock.close()
    return result == 0

def run_workload(timestamp, deployment, url, threads, conns, duration, rate):
  import threading

  # verify workload files
  if not os.path.exists(f"{BASE_DIR}/wrk2/wrk"):
    print(f"[ERROR] error running workload: '{BASE_DIR}/wrk2/wrk' file does not exist")
    exit(-1)


  progress_thread = threading.Thread(target=display_progress_bar, args=(duration, "running workload",))
  progress_thread.start()

  from plumbum import local
  with local.env(HOST_EU=url, HOST_US=url):
    wrk2 = local['./wrk2/wrk']
    output = wrk2['-D', 'exp', '-t', str(threads), '-c', str(conns), '-d', str(duration), '-L', '-s', './wrk2/scripts/social-network/compose-post.lua', f'{url}/wrk2-api/post/compose', '-R', str(rate)]()
  
    filepath = f"evaluation/{deployment}/{timestamp}/workload.out"
    os.makedirs(os.path.dirname(filepath), exist_ok=True)
    with open(filepath, "w") as f:
      f.write(output)

    print(output)
    print(f"[INFO] workload results saved at {filepath}")

  progress_thread.join()
  return output

def gen_weaver_config_gcp():
  host_eu = get_instance_host(GCP_INSTANCE_DB_EU, GCP_ZONE_EU)
  host_us = get_instance_host(GCP_INSTANCE_DB_US, GCP_ZONE_US)

#eu
  weaver_eu = "deploy/xcweaver/weaver-template-eu.toml"

  with open(weaver_eu, 'r') as file:
    toml_content = file.read()

  # Replace the placeholder with the actual value
  toml_content = re.sub(r'{{ host_eu }}', host_eu, toml_content)

  filepath_eu = "deploy/tmp/weaver-gcp-eu.toml"
  with open(filepath_eu, 'w') as file:
    file.write(toml_content)

  # us
  weaver_us = "deploy/xcweaver/weaver-template-us.toml"
  with open(weaver_us, 'r') as file:
    toml_content = file.read()

  # Replace the placeholder with the actual value
  toml_content = re.sub(r'{{ host_us }}', host_us, toml_content)
  filepath_us = "deploy/tmp/weaver-gcp-us.toml"
  with open(filepath_us, 'w') as file:
    file.write(toml_content)

  print(f"[INFO] generated app config for GCP at {filepath_eu} and {filepath_us}")

def gen_ansible_vars(workload_timestamp=None, deployment_type=None, duration=None, threads=None, rate=None):
  import yaml

  with open('deploy/ansible/templates/vars.yml', 'r') as file:
    data = yaml.safe_load(file)

  data['base_dir'] = BASE_DIR
  data['workload_timestamp'] = workload_timestamp if workload_timestamp else None
  data['deployment_type'] = deployment_type if deployment_type else None
  data['duration'] = duration if duration else None
  data['threads'] = threads if threads else None
  data['rate'] = rate if rate else None

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
    'host_app_eu':      get_instance_host(GCP_INSTANCE_APP_EU, GCP_ZONE_EU),
    'host_app_us':      get_instance_host(GCP_INSTANCE_APP_US, GCP_ZONE_US),
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


def get_consistency_window(port):
  url = f'http://localhost:{port}/consistency_window'
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
      # Calculate the median
      median = statistics.median(values)
      return median
  else:
      print(f"Failed to get data: {response.status_code}")
      return None

# METRICS FORMAT
#╭────────────────────────────────────────────────────────────────────────╮
#│ // The number of composed posts                                        │
#│ composed_posts: COUNTER                                                │
#├───────────────────┬────────────────────┬───────────────────────┬───────┤
#│ serviceweaver_app │ serviceweaver_node │ serviceweaver_version │ Value │
#├───────────────────┼────────────────────┼───────────────────────┼───────┤
#│ weaver-dsb-db     │ 0932683b           │ 1cd20361              │ 0     │
#│ weaver-dsb-db     │ 1205179c           │ 1cd20361              │ 0     │
#|  ...              | ...                | ...                   | ...   |
#╰───────────────────┴────────────────────┴───────────────────────┴───────╯
#
#╭────────────────────────────────────────────────────────────────────────╮
#│ // The number of times an cross-service inconsistency has occured      │
#│ inconsistencies: COUNTER                                               │
#├───────────────────┬────────────────────┬───────────────────────┬───────┤
#│ serviceweaver_app │ serviceweaver_node │ serviceweaver_version │ Value │
#├───────────────────┼────────────────────┼───────────────────────┼───────┤
#│ weaver-dsb-db     │ 0932683b           │ 1cd20361              │ 0     │
#│ weaver-dsb-db     │ 1205179c           │ 1cd20361              │ 0     │
#|  ...              | ...                | ...                   | ...   |
#╰───────────────────┴────────────────────┴───────────────────────┴───────╯

def metrics(deployment, timestamp=None):
  from plumbum.cmd import xcweaver
  import re

  if timestamp== None:
    timestamp = datetime.datetime.now().strftime("%Y-%m-%d_%H:%M:%S")

  primary_region = 'europe-west3'
  secondary_region = 'us-central-1' if not local else primary_region

  pattern = re.compile(r'^.*│.*│.*│.*│\s*(\d+\.?\d*)\s*│.*$', re.MULTILINE)

  def get_filter_metrics(metric_name):
    return xcweaver['multi', 'metrics', metric_name]()

  # compose post service
  composed_posts_metrics = get_filter_metrics('sn_composed_posts')
  composed_posts_count = sum(int(value) for value in pattern.findall(composed_posts_metrics))
  sent_notifications_metrics = get_filter_metrics('sn_sent_notifications')
  sent_notifications_count = sum(int(value) for value in pattern.findall(sent_notifications_metrics))
  # post storage service
  write_post_duration_metrics = get_filter_metrics('sn_write_post_duration_ms')
  write_post_duration_metrics_values = pattern.findall(write_post_duration_metrics)
  write_post_duration_avg_ms = sum(float(value) for value in write_post_duration_metrics_values)/len(write_post_duration_metrics_values) if write_post_duration_metrics_values else 0
  # write home timeline service
  read_post_duration_metrics = get_filter_metrics('sn_read_post_duration_ms')
  read_post_duration_metrics_values = pattern.findall(read_post_duration_metrics)
  read_post_duration_avg_ms = sum(float(value) for value in read_post_duration_metrics_values)/len(read_post_duration_metrics_values) if read_post_duration_metrics_values else 0
  consistency_window_metrics = get_filter_metrics('sn_consistency_window_ms')
  consistency_window_metrics_values = pattern.findall(consistency_window_metrics)
  consistency_window_avg_ms = sum(float(value) for value in consistency_window_metrics_values)/len(consistency_window_metrics_values) if consistency_window_metrics_values else 0
  queue_duration_metrics = get_filter_metrics('sn_queue_duration_ms')
  queue_duration_metrics_values = pattern.findall(queue_duration_metrics)
  queue_duration_avg_ms = sum(float(value) for value in queue_duration_metrics_values)/len(queue_duration_metrics_values) if queue_duration_metrics_values else 0
  received_notifications_metrics = get_filter_metrics('sn_received_notifications')
  received_notifications_count = sum(int(value) for value in pattern.findall(received_notifications_metrics))
  inconsitencies_metrics = get_filter_metrics('sn_inconsistencies')
  inconsistencies_count = sum(int(value) for value in pattern.findall(inconsitencies_metrics))

  consistency_window_median_ms = get_consistency_window(12346)
  
  consistency_window_median_ms = "{:.2f}".format(consistency_window_median_ms)
  pc_inconsistencies = "{:.2f}".format((inconsistencies_count / received_notifications_count) * 100) if received_notifications_count != 0 else 0
  pc_received_notifications = "{:.2f}".format((received_notifications_count / sent_notifications_count) * 100) if sent_notifications_count else 0

  write_post_duration_avg_ms = "{:.2f}".format(write_post_duration_avg_ms)
  read_post_duration_avg_ms = "{:.2f}".format(read_post_duration_avg_ms)
  consistency_window_avg_ms = "{:.2f}".format(consistency_window_avg_ms)
  queue_duration_avg_ms = "{:.2f}".format(queue_duration_avg_ms)

  results = {
    'num_composed_posts': int(composed_posts_count),
    'num_received_notifications': int(received_notifications_count),
    'num_notifications_sent': int(sent_notifications_count),
    'num_inconsistencies': int(inconsistencies_count),
    'per_inconsistencies': float(pc_inconsistencies),
    'per_inconsistencies': float(pc_received_notifications),
    'avg_read_post_duration_ms': float(read_post_duration_avg_ms),
    'avg_write_post_duration_ms': float(write_post_duration_avg_ms),
    'consistency_window_ms': float(consistency_window_avg_ms),
    'avg_queue_duration_ms': float(queue_duration_avg_ms),
    'consistency_window_median_ms': float(consistency_window_median_ms),
  }

  # save file if we ran workload
  if timestamp:
    filepath = f"evaluation/{deployment}/{timestamp}/metrics.yml"
    os.makedirs(os.path.dirname(filepath), exist_ok=True)
    with open(filepath, 'w') as outfile:
      yaml.dump(results, outfile, default_flow_style=False)
    print(yaml.dump(results, default_flow_style=False))
    print(f"[INFO] evaluation results saved at {filepath}")

def update_metrics(deployment_type='gke', timestamp=None, local=True):
  from plumbum.cmd import xcweaver, grep
  import re

  primary_region = 'europe-west3'
  secondary_region = 'us-central-1' if not local else primary_region

  pattern = re.compile(r'^.*│.*│.*│.*│\s*(\d+\.?\d*)\s*│.*$', re.MULTILINE)

  def get_filter_metrics(deployment_type, metric_name, region):
    #return (weaver[deployment_type, 'metrics', metric_name] | grep[region])()
    return xcweaver[deployment_type, 'metrics', metric_name]()

  # wkr2 api
  update_post_duration_metrics = get_filter_metrics(deployment_type, 'sn_update_post_duration_ms', primary_region)
  update_post_duration_metrics_values = pattern.findall(update_post_duration_metrics)
  update_post_duration_avg_ms = sum(float(value) for value in update_post_duration_metrics_values)/len(update_post_duration_metrics_values)
  # compose post service
  updated_posts_metrics = get_filter_metrics(deployment_type, 'sn_updated_posts', primary_region)
  updated_posts_count = sum(int(value) for value in pattern.findall(updated_posts_metrics))
  # post storage service
  update_post_operation_duration_metrics = get_filter_metrics(deployment_type, 'sn_update_post_operation_duration_ms', primary_region)
  update_post_operation_duration_metrics_values = pattern.findall(update_post_operation_duration_metrics)
  update_post_operation_duration_avg_ms = sum(float(value) for value in update_post_operation_duration_metrics_values)/len(update_post_operation_duration_metrics_values)
  # update home timeline service
  queue_duration_metrics = get_filter_metrics(deployment_type, 'sn_queue_duration_ms', secondary_region)
  queue_duration_metrics_values = pattern.findall(queue_duration_metrics)
  queue_duration_avg_ms = sum(float(value) for value in queue_duration_metrics_values)/len(queue_duration_metrics_values)
  read_post_duration_metrics = get_filter_metrics(deployment_type, 'sn_read_post_duration_ms', primary_region)
  read_post_duration_metrics_values = pattern.findall(read_post_duration_metrics)
  read_post_duration_avg_ms = sum(float(value) for value in read_post_duration_metrics_values)/len(read_post_duration_metrics_values)
  received_notifications_metrics = get_filter_metrics(deployment_type, 'sn_received_notifications', secondary_region)
  received_notifications_count = sum(int(value) for value in pattern.findall(received_notifications_metrics))
  inconsitencies_metrics = get_filter_metrics(deployment_type, 'sn_update_inconsistencies', secondary_region)
  inconsistencies_count = sum(int(value) for value in pattern.findall(inconsitencies_metrics))

  pc_inconsistencies = "{:.2f}".format((inconsistencies_count / updated_posts_count) * 100)
  pc_received_notifications = "{:.2f}".format((received_notifications_count / updated_posts_count) * 100)
  update_post_duration_avg_ms = "{:.2f}".format(update_post_duration_avg_ms)
  update_post_operation_duration_avg_ms = "{:.2f}".format(update_post_operation_duration_avg_ms)
  read_post_duration_avg_ms = "{:.2f}".format(read_post_duration_avg_ms)
  queue_duration_avg_ms = "{:.2f}".format(queue_duration_avg_ms)

  results = f"""
    # updated posts:\t\t\t{updated_posts_count}
    # received notifications @ US:\t{received_notifications_count} ({pc_received_notifications}%)
    # inconsistencies @ US:\t\t{inconsistencies_count}
    % inconsistencies @ US:\t\t{pc_inconsistencies}%
    > avg. update post duration:\t{update_post_duration_avg_ms}ms
    > avg. update post duration:\t\t{update_post_operation_duration_avg_ms}ms
    > avg. read post duration:\t\t{read_post_duration_avg_ms}ms
    > avg. queue duration @ US:\t\t{queue_duration_avg_ms}ms
  """
  print(results)

  # save file if we ran workload
  if timestamp:
    eval_folder = 'local' if deployment_type == 'multi' else 'gke'
    filepath = f"evaluation/{eval_folder}/{timestamp}_update_metrics.txt"
    with open(filepath, "w") as f:
      f.write(results)
    print(f"[INFO] evaluation results saved at {filepath}")


def merge_metrics_gcp(timestamp):
  filepatheu = f"evaluation/gcp/{timestamp}/metrics-eu.yml"
  filepathus = f"evaluation/gcp/{timestamp}/metrics-us.yml"

  with open(filepatheu, 'r') as outfile_eu, open(filepathus, 'r') as outfile_us:
    metrics_eu = yaml.safe_load(outfile_eu)
    metrics_us = yaml.safe_load(outfile_us)

  merged_data = {**metrics_eu, **metrics_us}

  notifications_received = merged_data['num_received_notifications']
  notifications_sent = merged_data['num_notifications_sent']

  per_notifications_received = "{:.2f}".format((notifications_received / notifications_sent) * 100)

  percentages = {
    'per_notifications_received': float(per_notifications_received),
  }

  # Add percentages to metrics
  merged_data.update(percentages)
  
  filepathfinal = f"evaluation/gcp/{timestamp}/metrics.yml"
  with open(filepathfinal, 'w') as out_file:
      yaml.safe_dump(merged_data, out_file)

# --------------------
# GCP
# --------------------

def gcp_configure():
  from plumbum.cmd import gcloud

  try:
    print("[INFO] configuring firewalls")
    # xcweaver-dsb-socialnetwork:
    # tcp ports: 12345, 14318
    # xcweaver-dsb-storage:
    # tcp ports: 27017,27018,15672,15673,5672,5673,6381,6382,6383,6384,6385,6386,6387,6388,11212,11213,11214,11215,11216,11217
    # xcweaver-dsb-swarm:
    # tcp ports: 2376,2377,7946
    # udp ports: 4789,7946
    firewalls = {
      'xcweaver-dsb-socialnetwork': 'tcp:12345,tcp:14318',
      'xcweaver-dsb-storage': 'tcp:27017,tcp:27018,tcp:15672,tcp:15673,tcp:5672,tcp:5673,tcp:6381,tcp:6382,tcp:6383,tcp:6384,tcp:6385,tcp:6386,tcp:6387,tcp:6388,tcp:11212,tcp:11213,tcp:11214,tcp:11215,tcp:11216,tcp:11217',
      'xcweaver-dsb-swarm': 'tcp:2376,tcp:2377,tcp:7946,udp:4789,udp:7946'
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
  from plumbum.cmd import terraform, ansible_playbook

  terraform['-chdir=./deploy/terraform', 'init'] & FG
  terraform['-chdir=./deploy/terraform', 'apply', '-auto-approve'] & FG

  display_progress_bar(30, "waiting for all machines to be ready")

  # generate temporary files for this deployment
  os.makedirs("deploy/tmp", exist_ok=True)
  print(f"[INFO] created deploy/tmp/ directory")

  gen_ansible_config()
  # generate weaver config with hosts of datastores in gcp machines
  gen_weaver_config_gcp()
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
  print(f"services in {GCP_ZONE_EU} running @", get_instance_host(GCP_INSTANCE_APP_EU, GCP_ZONE_EU))
  print(f"services in {GCP_ZONE_US} running @\n\n", get_instance_host(GCP_INSTANCE_APP_US, GCP_ZONE_US))
  
def gcp_start():
  from plumbum.cmd import ansible_playbook
  ansible_playbook["deploy/ansible/playbooks/start-datastores.yml", "-i", "deploy/tmp/ansible-inventory.cfg", "--extra-vars", "@deploy/tmp/ansible-vars.yml"] & FG
  display_progress_bar(30, "waiting for datastores to initialize")
  ansible_playbook["deploy/ansible/playbooks/start-app.yml", "-i", "deploy/tmp/ansible-inventory.cfg", "--extra-vars", "@deploy/tmp/ansible-vars.yml"] & FG

def gcp_stop():
  from plumbum.cmd import ansible_playbook
  ansible_playbook["deploy/ansible/playbooks/stop-datastores.yml", "-i", "deploy/tmp/ansible-inventory.cfg", "--extra-vars", "@deploy/tmp/ansible-vars.yml"] & FG
  ansible_playbook["deploy/ansible/playbooks/stop-app.yml", "-i", "deploy/tmp/ansible-inventory.cfg", "--extra-vars", "@deploy/tmp/ansible-vars.yml"] & FG

def gcp_restart():
  gcp_stop()
  gcp_start()
  
def gcp_clean():
  from plumbum.cmd import terraform
  import shutil

  terraform['-chdir=./deploy/terraform', 'destroy', '-auto-approve'] & FG
  if os.path.exists("deploy/tmp"):
    shutil.rmtree("deploy/tmp")
    print(f"[INFO] removed {BASE_DIR}/deploy/tmp/ directory")
 
def gcp_init_social_graph():
  print("[INFO] nothing to be done for gcp")
  exit(0)

def gcp_metrics(timestamp):
  from plumbum.cmd import ansible_playbook
  if timestamp== None:
    timestamp = datetime.datetime.now().strftime("%Y-%m-%d_%H:%M:%S")
  gen_ansible_vars(timestamp, 'gcp')
  metrics_path = f"{BASE_DIR}/evaluation/gcp/{timestamp}"
  os.makedirs(metrics_path, exist_ok=True)
  ansible_playbook["deploy/ansible/playbooks/gather-metrics.yml", "-i", "deploy/tmp/ansible-inventory.cfg", "--extra-vars", "@deploy/tmp/ansible-vars.yml"] & FG
  merge_metrics_gcp(timestamp)
  print(f"[INFO] metrics results saved at evaluation/gcp/{timestamp}/ in metrics.yaml")

def gcp_wrk2_vm(conns, duration, threads, rate, timestamp):
  host_eu = get_instance_host(GCP_INSTANCE_APP_EU, GCP_ZONE_EU)
  if timestamp == None:
    timestamp = datetime.datetime.now().strftime("%Y-%m-%d_%H:%M:%S")
  run_workload(timestamp, 'gcp', f"http://{host_eu}:{APP_PORT}", threads, conns, duration, rate)

def gcp_consistency_window():
  host = get_instance_host(GCP_INSTANCE_APP_US, GCP_ZONE_US)
  url = f'http://{host}:{APP_PORT}/consistency_window'
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


def gcp_wrk2(threads, conns, duration, rate):
  from plumbum.cmd import ansible_playbook
  timestamp = datetime.datetime.now().strftime("%Y-%m-%d_%H:%M:%S")
  metrics_path = f"{BASE_DIR}/evaluation/gcp/{timestamp}"
  os.makedirs(metrics_path, exist_ok=True)
  gen_ansible_vars(timestamp, 'gcp', duration, threads, rate)
  ansible_playbook["deploy/ansible/playbooks/run-wrk2.yml", "-i", "deploy/tmp/ansible-inventory.cfg", "--extra-vars", "@deploy/tmp/ansible-vars.yml"] & FG
  sleep(3)
  ansible_playbook["deploy/ansible/playbooks/gather-metrics.yml", "-i", "deploy/tmp/ansible-inventory.cfg", "--extra-vars", "@deploy/tmp/ansible-vars.yml"] & FG
  merge_metrics_gcp(timestamp)
  print(f"[INFO] metrics results saved at evaluation/gcp/{timestamp}/ in metrics-eu.yaml and metrics-us.yaml files")

# --------------------
# LOCAL
# --------------------

def local_init_social_graph():
  from plumbum import local
  with local.env(HOST_EU=f"http://127.0.0.1:{APP_PORT}", HOST_US=f"http://127.0.0.1:{APP_PORT}"):
    local['./social-graph/init_social_graph.py'] & FG

def local_wrk2(threads, conns, duration, rate):
  timestamp = datetime.datetime.now().strftime("%Y-%m-%d_%H:%M:%S")
  run_workload(timestamp, 'local', f"http://127.0.0.1:{APP_PORT}", threads, conns, duration, rate)
  metrics('local', timestamp)

def local_wrk2_update(address, duration):
  import threading
  
  timestamp = datetime.datetime.now().strftime("%Y-%m-%d_%H:%M:%S")

  url_register = f"http://{address}:{APP_PORT}/wrk2-api/user/register"
  params = {
    "user_id": "0",
    "username": "tiago",
    "first_name": "tiago",
    "last_name": "malhadas",
    "password": "qsgechgj345Fwr_cc",
  }
  response = requests.get(url_register, params=params)
  if response.status_code != 200:
    print(f"register request failed with status code {response.status_code}.")
    return
  
  print("[INFO] user tiago successfully registered!")

  sleep(0.5)

  def tqdm_progress(duration):
    print(f"[INFO] composing posts for {duration} seconds...")
    for _ in tqdm(range(int(duration))):
        time.sleep(1)

  progress_thread = threading.Thread(target=tqdm_progress, args=(duration,))
  progress_thread.start()

  start_time = time.time()

  postIds = []
  url_compose = f"http://{address}:{APP_PORT}/wrk2-api/post/compose"
  params = {
    "user_id": "0",
    "username": "tiago",
    "text": "old_post",
    "post_type": "0",
    "media_types": "['png']",
    "media_ids": "[0]",
  }
  while True:
    response = requests.get(url_compose, params=params)
    if response.status_code != 200:
      print(f"compose request failed with status code {response.status_code}.")
      return
    postIds.append(response.json()["PostId"])

    elapsed_time = time.time() - start_time
    if elapsed_time >= duration:
        break
  
  print(f"[INFO] Compose posts done!")
  progress_thread.join()

  print(f"[INFO] updating all posts...")
  url_update = f"http://{address}:{APP_PORT}/wrk2-api/post/update"
  for postId in postIds:
    params = {
      "post_id": postId,
      "user_id": "0",
      "username": "tiago",
      "text": "updatedText",
      "post_type": "0",
      "media_types": "['png']",
      "media_ids": "[0]",
    }
    response = requests.get(url_update, params=params)
    if response.status_code != 200:
      print(f"update request failed with status code {response.status_code}.")
      return
  update_metrics('multi', timestamp, True)

def local_metrics(timestamp):
  metrics('multi', timestamp)

def local_metrics_update():
  update_metrics('multi', None, True)

def local_metrics_eu(timestamp):
  import yaml
  from plumbum.cmd import xcweaver
  import re

  if timestamp== None:
    timestamp = datetime.datetime.now().strftime("%Y-%m-%d_%H:%M:%S")

  pattern = re.compile(r'^.*│.*│.*│.*│\s*(\d+\.?\d*)\s*│.*$', re.MULTILINE)

  def get_filter_metrics(metric_name):
    return xcweaver['multi', 'metrics', metric_name]()

  # Steps:
  # 1. get desired metrics that will be listed for each process using get_filter_metrics
  # 2. then average the values of all processes which is ok since weaver metrics are limited to averages

  
  # compose post service
  composed_posts_metrics = get_filter_metrics('sn_composed_posts')
  composed_posts_count = sum(int(value) for value in pattern.findall(composed_posts_metrics))
  sent_notifications_metrics = get_filter_metrics('sn_sent_notifications')
  sent_notifications_count = sum(int(value) for value in pattern.findall(sent_notifications_metrics))
  # post storage service
  write_post_duration_metrics = get_filter_metrics('sn_write_post_duration_ms')
  write_post_duration_metrics_values = pattern.findall(write_post_duration_metrics)
  write_post_duration_avg_ms = sum(float(value) for value in write_post_duration_metrics_values)/len(write_post_duration_metrics_values) if write_post_duration_metrics_values else 0

  write_post_duration_avg_ms = "{:.2f}".format(write_post_duration_avg_ms)

  results = {
    'num_composed_posts': int(composed_posts_count),
    'num_notifications_sent': int(sent_notifications_count),
    'avg_write_post_duration_ms': float(write_post_duration_avg_ms),
  }

  # save file if we ran workload
  if timestamp:
    filepath = f"evaluation/local/{timestamp}/metrics-eu.yml"
    os.makedirs(os.path.dirname(filepath), exist_ok=True)
    with open(filepath, 'w') as outfile:
      yaml.dump(results, outfile, default_flow_style=False)
    print(yaml.dump(results, default_flow_style=False))
    print(f"[INFO] evaluation results saved at {filepath}")

def local_metrics_us(timestamp):
  import yaml
  from plumbum.cmd import xcweaver
  import re

  if timestamp== None:
    timestamp = datetime.datetime.now().strftime("%Y-%m-%d_%H:%M:%S")

  pattern = re.compile(r'^.*│.*│.*│.*│\s*(\d+\.?\d*)\s*│.*$', re.MULTILINE)

  def get_filter_metrics(metric_name):
    return xcweaver['multi', 'metrics', metric_name]()

  # Steps:
  # 1. get desired metrics that will be listed for each process using get_filter_metrics
  # 2. then average the values of all processes which is ok since weaver metrics are limited to averages

  # write home timeline service
  read_post_duration_metrics = get_filter_metrics('sn_read_post_duration_ms')
  read_post_duration_metrics_values = pattern.findall(read_post_duration_metrics)
  read_post_duration_avg_ms = sum(float(value) for value in read_post_duration_metrics_values)/len(read_post_duration_metrics_values) if read_post_duration_metrics_values else 0
  consistency_window_metrics = get_filter_metrics('sn_consistency_window_ms')
  consistency_window_metrics_values = pattern.findall(consistency_window_metrics)
  consistency_window_avg_ms = sum(float(value) for value in consistency_window_metrics_values)/len(consistency_window_metrics_values) if consistency_window_metrics_values else 0
  queue_duration_metrics = get_filter_metrics('sn_queue_duration_ms')
  queue_duration_metrics_values = pattern.findall(queue_duration_metrics)
  queue_duration_avg_ms = sum(float(value) for value in queue_duration_metrics_values)/len(queue_duration_metrics_values) if queue_duration_metrics_values else 0
  received_notifications_metrics = get_filter_metrics('sn_received_notifications')
  received_notifications_count = sum(int(value) for value in pattern.findall(received_notifications_metrics))
  inconsitencies_metrics = get_filter_metrics('sn_inconsistencies')
  inconsistencies_count = sum(int(value) for value in pattern.findall(inconsitencies_metrics))

  consistency_window_median_ms = get_consistency_window(APP_PORT)
  
  consistency_window_median_ms = "{:.2f}".format(consistency_window_median_ms)
  pc_inconsistencies = "{:.2f}".format((inconsistencies_count / received_notifications_count) * 100) if received_notifications_count != 0 else 0

  queue_duration_avg_ms = "{:.2f}".format(queue_duration_avg_ms)
  read_post_duration_avg_ms = "{:.2f}".format(read_post_duration_avg_ms)
  consistency_window_avg_ms = "{:.2f}".format(consistency_window_avg_ms)

  results = {
    'num_received_notifications': int(received_notifications_count),
    'num_inconsistencies': int(inconsistencies_count),
    'per_inconsistencies': float(pc_inconsistencies),
    'avg_queue_duration_ms': float(queue_duration_avg_ms),
    'avg_read_post_duration_ms': float(read_post_duration_avg_ms),
    'consistecy_window_ms': float(consistency_window_avg_ms),
    'consistency_window_median_ms': float(consistency_window_median_ms),
  }

  # save file if we ran workload
  if timestamp:
    filepath = f"evaluation/local/{timestamp}/metrics-us.yml"
    os.makedirs(os.path.dirname(filepath), exist_ok=True)
    with open(filepath, 'w') as outfile:
      yaml.dump(results, outfile, default_flow_style=False)
    print(yaml.dump(results, default_flow_style=False))
    print(f"[INFO] evaluation results saved at {filepath}")

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
  docker['build', '-t', 'mongodb-delayed:4.4.6', 'docker/mongodb-delayed/.'] & FG
  docker['build', '-t', 'mongodb-setup:4.4.6', 'docker/mongodb-setup/post-storage/.'] & FG
  docker['build', '-t', 'rabbitmq-setup:3.8', 'docker/rabbitmq-setup/write-home-timeline/.'] & FG

def local_storage_run():
  from plumbum.cmd import docker_compose
  docker_compose['up', '-d'] & FG
  display_progress_bar(30, "waiting for storages to be ready")

def local_storage_info():
  print("[INFO] nothing to be done for local")
  exit(0)

def local_storage_clean():
  from plumbum.cmd import docker_compose
  docker_compose['down'] & FG

def gcp_tests():
  from plumbum import local
  
  clients = [1]
  manager_script = local["./manager.py"]

  #deploy_command = manager_script["deploy"]["--gcp"]
  #deploy_command()

  start_command = manager_script["start"]["--gcp"]
  start_command()
  
  for client in clients:
    rate = str(client * 10)
    clientstr = str(client)
    for i in range(10):
      wrk2_command = manager_script["wrk2"]["--gcp"]["-d", "300"]["-r", rate]["-t", clientstr]
      wrk2_command()
      
      # Run the restart command
      restart_command = manager_script["restart"]["--gcp"]
      restart_command()
      print(f"[TEST REPORT] Client {client} already executed {i+1} times!!!")

if __name__ == "__main__":
  main_parser = argparse.ArgumentParser()
  command_parser = main_parser.add_subparsers(help='commands', dest='command')

  commands = [
    # gcp
    'configure', 'deploy', 'start', 'stop', 'restart', 'clean', 'info', 'consistency-window', 'tests',
    # datastores
    'storage-build', 'storage-deploy', 'storage-run', 'storage-info', 'storage-clean',
    # eval
    'init-social-graph', 'wrk2', 'wrk2-vm', 'wrk2-compose', 'wrk2-update', 'metrics', 'metrics-eu', 'metrics-us',
  ]
  
  for cmd in commands:
    parser = command_parser.add_parser(cmd)
    parser.add_argument('--local', action='store_true', help="Running in localhost")
    parser.add_argument('--gcp', action='store_true',   help="Running in gcp")
    if cmd == 'wrk2-compose' or cmd == 'wrk2' or cmd == 'wrk2-vm':
      parser.add_argument('-t', '--threads', default=2, help="Number of threads")
      parser.add_argument('-c', '--conns', default=2, help="Number of connections")
      parser.add_argument('-d', '--duration', default=30, help="Duration")
      parser.add_argument('-r', '--rate', default=50, help="Number of requests per second")
    if cmd == 'wrk2-update':
      parser.add_argument('-d', '--duration', default=30, help="Duration")
    if cmd == 'wrk2-vm':
      parser.add_argument('-ts', '--timestamp', help="Timestamp of workload")
    if cmd == 'metrics' or cmd == 'metrics-eu' or cmd == 'metrics-us':
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
