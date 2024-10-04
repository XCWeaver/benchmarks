#!/usr/bin/env python3
import argparse
import datetime
import io
import json
import os
import socket
import sys
import time
from plumbum import FG
from tqdm import tqdm
import requests
import googleapiclient.discovery
from google.oauth2 import service_account
import yaml
import hdrh.dump
import hdrh.histogram
import hdrh
import statistics
import numpy as np

APP_PORT                    = 12345
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
APP_FOLDER_NAME           = "trainticket"
GCP_INSTANCE_WRK2         = "trainticket-wrk2"
GCP_ZONE_EU               = "europe-west6-a"

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
    time.sleep(1)

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

def gen_ansible_vars(workload_timestamp=None, deployment_type=None, duration=None, threads=None, rate=None, host=None):
  import yaml

  with open('deploy/ansible/templates/vars.yml', 'r') as file:
    data = yaml.safe_load(file)

  data['base_dir'] = BASE_DIR
  data['workload_timestamp'] = workload_timestamp if workload_timestamp else None
  data['deployment_type'] = deployment_type if deployment_type else None
  data['duration'] = duration if duration else None
  data['threads'] = threads if threads else None
  data['rate'] = rate if rate else None
  data['host'] = host if host else None

  with open('deploy/tmp/ansible-vars.yml', 'w') as file:
    yaml.dump(data, file)

def gen_ansible_inventory_gcp():
  from jinja2 import Environment, FileSystemLoader
  import textwrap

  host_wrk2   = get_instance_host(GCP_INSTANCE_WRK2, GCP_ZONE_EU)

  template = Environment(loader=FileSystemLoader('.')).get_template( "deploy/ansible/templates/inventory.cfg")
  inventory = template.render({
    'username':         GCP_USERNAME,
    'host_wrk2':   host_wrk2,
  })

  filename = "deploy/tmp/ansible-inventory-gcp.cfg"
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

def get_consistency_window():
  url = f'http://localhost:{APP_PORT}/wrk2-api/user/consistencyWindow'
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
  from plumbum.cmd import weaver
  import re

  if timestamp == None:
    timestamp = datetime.datetime.now().strftime("%Y-%m-%d_%H:%M:%S")

  pattern = re.compile(r'^.*│.*│.*│.*│\s*(\d+\.?\d*)\s*│.*$', re.MULTILINE)

  def get_filter_metrics(metric_name):
    return weaver['multi', 'metrics', metric_name]()

  # wkr2 api
  inconsitencies_metrics = get_filter_metrics('tt_inconsistencies')
  inconsistencies_count = sum(int(value) for value in pattern.findall(inconsitencies_metrics))
  tickets_canceled_metrics = get_filter_metrics('sn_tickets_canceled')
  tickets_canceled = sum(int(value) for value in pattern.findall(tickets_canceled_metrics))
  pc_inconsistencies = "{:.2f}".format((inconsistencies_count / tickets_canceled) * 100)
  consistency_window_metrics = get_filter_metrics('sn_consistency_window_ms')
  consistency_window_metrics_values = pattern.findall(consistency_window_metrics)
  consistency_window_avg_ms = sum(float(value) for value in consistency_window_metrics_values if value != 0)/2 if consistency_window_metrics_values else 0
  
  consistency_window_median_ms = get_consistency_window()

  consistency_window_avg_ms = "{:.2f}".format(consistency_window_avg_ms)
  consistency_window_median_ms = "{:.2f}".format(consistency_window_median_ms)

  results = {
    'num_tickets_canceled': int(tickets_canceled),
    'num_inconsistencies': int(inconsistencies_count),
    'per_inconsistencies': float(pc_inconsistencies),
    'consistency_window_ms': float(consistency_window_avg_ms),
    'consistency_window_median_ms': float(consistency_window_median_ms),
  }

  print(results)

  filepath = f"evaluation/{deployment}/{timestamp}/metrics.yml"
  os.makedirs(os.path.dirname(filepath), exist_ok=True)
  with open(filepath, 'w') as outfile:
    yaml.dump(results, outfile, default_flow_style=False)
  print(yaml.dump(results, default_flow_style=False))
  print(f"[INFO] evaluation results saved at {filepath}")




def run_test(duration, address, rate, id, throughputs, latencies_results, index):

  cancel_url = f"{address}/wrk2-api/user/cancelTicket"

  def execute_request(session, url, params):
    request_time_ms = int(time.time() * 1000)
    try:
      with session.get(url, params=params) as response:
        if response.status_code != 200:
          print(f"Request failed with status code {response.status_code}.")
          return
        latencies.append(int(time.time() * 1000) - request_time_ms)
        requests_counter[0] += 1
    except requests.RequestException as e:
      print(f"Request failed: {e}")
    except json.JSONDecodeError as e:
      print(f"JSON decoding failed: {e}")

  latencies = []
  start_time = time.time()
  interval = 1.0 / int(rate)
  requests_counter  = [0]
  with requests.Session() as session:
    first_request_time = int(time.time() * 1000)
    stopped = False
    #iterates over the orders list
    #if reaches the last order starts again until the time is over
    while True:
      for orderId in range(id, -1, -1):
        params = {
          "token": "token",
          "orderId": str(orderId),
          "loginId": "accountId",
        }

        execute_request(session, cancel_url, params)
        
        elapsed_time = time.time() - start_time
        if elapsed_time >= int(duration):
          stopped = True
          break

        time.sleep(interval - ((time.time() - start_time) % interval))

      if stopped:
        break

  print(f"{requests_counter[0]} requests executed!")
  test_time = int(time.time() * 1000) - first_request_time
  throughput = requests_counter[0] / (test_time / 1000)

  throughputs[index] = throughput
  latencies_results[index] = latencies

# --------------------
# GCP
# --------------------

def gcp_configure():
  from plumbum.cmd import gcloud

  try:
    print("[INFO] configuring firewalls")
    # trainticket-app:
    # tcp ports: 12345, 80
    firewalls = {
      'trainticket-app': 'tcp:12345,tcp:80',
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
  # generate ansible inventory with hosts of all gcp machines
  gen_ansible_inventory_gcp()
  # generate ansible inventory with extra variables for current deployment
  gen_ansible_vars()
  
  ansible_playbook["deploy/ansible/playbooks/install-machines.yml", "-i", "deploy/tmp/ansible-inventory-gcp.cfg", "--extra-vars", "@deploy/tmp/ansible-vars.yml"] & FG

def gcp_clean():
  from plumbum.cmd import terraform
  import shutil

  terraform['-chdir=./deploy/terraform', 'destroy', '-auto-approve'] & FG
  if os.path.exists("deploy/tmp"):
    shutil.rmtree("deploy/tmp")
    print(f"[INFO] removed {BASE_DIR}/deploy/tmp/ directory")

def gcp_metrics(timestamp):
  print(f"[INFO] metrics automatically generated at evaluation/gcp after running wrk2")


# --------------------
# KUBERNETES
# --------------------

def kube_consistency_window(host):
  url = f'http://{host}/wrk2-api/user/consistencyWindow'
  response = requests.get(url)

  if response.status_code == 200:
      data = json.loads(response.text)
      
      if data == None:
        return None

      # Calculate the median
      median = np.percentile(data, 90)
      print(f"Median: {median}")
      return median
  else:
      print(f"Failed to get consistency window: {response.status_code}")

def kube_inconsistencies(host):
  url = f'http://{host}/wrk2-api/user/inconsistencies'
  response = requests.get(url)

  if response.status_code == 200:
      inconsistencies = int(response.text)

      return inconsistencies
  else:
      print(f"Failed to get inconsistencies: {response.status_code}")

def gcp_cluster():
  from plumbum.cmd import gcloud
  gcloud["beta", "container", "--project", GCP_PROJECT_ID, "clusters", "create-auto", "cluster-trainticket", "--region", "europe-west6", "--release-channel", "regular", "--network", f"projects/{GCP_PROJECT_ID}/global/networks/default", "--subnetwork", f"projects/{GCP_PROJECT_ID}/regions/europe-west6/subnetworks/default", "--cluster-ipv4-cidr", "/17", "--binauthz-evaluation-mode=DISABLED"] & FG
  gcloud["container", "clusters", "get-credentials", "cluster-trainticket", "--region", "europe-west6", "--project", GCP_PROJECT_ID] & FG


def gcp_wrk2_vm(duration, threads, rate, timestamp, host):
  import threading
  threads_list = []
  individual_rate = int(int(rate) / int(threads))
  last_id = int(int(rate) * int(duration) // int(threads))
  throughputs = []
  latencies = []
  for i in range(int(threads)):
    throughputs.append(0)
    latencies.append([])
  for i in range(int(threads)):
    thread = threading.Thread(target=run_test, args=(duration, f"http://{host}", individual_rate, last_id * (i + 1), throughputs, latencies, i))
    thread.start()
    threads_list.append(thread)

  for thread in threads_list:
    thread.join()

  throughput = 0
  for i in throughputs:
    throughput += i

  throughput_str = f'\n----------------------------------------------------------\n Requests/sec:     {throughput}\n'
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
  consistency_window = kube_consistency_window(host)
  inconsistencies = kube_inconsistencies(host)
  filepath = f"evaluation/gcp/{timestamp}/workload.out"
  os.makedirs(os.path.dirname(filepath), exist_ok=True)
  with open(filepath, "w") as f:
    f.write(output.getvalue().decode() + throughput_str +  f" Latency Median:     {median}\n Consistency Window Median:     {consistency_window}\n  Inconsistencies:     {inconsistencies}\n")

  print(output.getvalue().decode() + throughput_str + f" Latency Median:     {median}\n  Consistency Window Median:     {consistency_window}\n  Inconsistencies:     {inconsistencies}\n")
  print(f"[INFO] workload results saved at {filepath}")

def gcp_wrk2(duration, threads, rate, host):
  from plumbum.cmd import ansible_playbook
  timestamp = datetime.datetime.now().strftime("%Y-%m-%d_%H:%M:%S")
  metrics_path = f"{BASE_DIR}/evaluation/gcp/{timestamp}"
  os.makedirs(metrics_path, exist_ok=True)
  gen_ansible_vars(timestamp, 'gcp', duration, threads, rate, host)
  ansible_playbook["deploy/ansible/playbooks/run-wrk2.yml", "-i", "deploy/tmp/ansible-inventory-gcp.cfg", "--extra-vars", "@deploy/tmp/ansible-vars.yml"] & FG
  print(f"[INFO] workload results saved at evaluation/gcp/{timestamp}/ in workload.yaml")


# --------------------
# LOCAL
# --------------------

def local_wrk2_kube(duration, threads, rate, host):
  import threading
  threads_list = []
  individual_rate = int(int(rate) / int(threads))
  last_id = int(int(rate) * int(duration) // int(threads))
  throughputs = []
  latencies = []
  for i in range(int(threads)):
    throughputs.append(0)
    latencies.append([])
  print(f"[INFO] Running {rate} requests/second for {duration} seconds with {threads} clients!")
  for i in range(int(threads)):
    thread = threading.Thread(target=run_test, args=(duration, f"http://{host}", individual_rate, last_id * (i + 1), throughputs, latencies, i))
    thread.start()
    threads_list.append(thread)

  for thread in threads_list:
    thread.join()

  print("[INFO] Requests finished!")
  print(f"[INFO] Getting workload!")

  throughput = 0
  for i in throughputs:
    throughput += i

  print(f"throughput: {throughput}")
  throughput_str = f'\n----------------------------------------------------------\n Requests/sec:     {throughput}\n'
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
  consistency_window = kube_consistency_window(host)
  inconsistencies = kube_inconsistencies(host)
  timestamp = datetime.datetime.now().strftime("%Y-%m-%d_%H:%M:%S")
  filepath = f"evaluation/local/{timestamp}/workload.out"
  os.makedirs(os.path.dirname(filepath), exist_ok=True)
  with open(filepath, "w") as f:
    f.write(output.getvalue().decode() + throughput_str + f"Latency Median:     {median}\n  Consistency Window Median:     {consistency_window}\n  Inconsistencies:     {inconsistencies}\n")

  print(output.getvalue().decode() + throughput_str + f"Latency Median:     {median}\n  Consistency Window Median:     {consistency_window}\n  Inconsistencies:     {inconsistencies}\n")
  print(f"[INFO] workload results saved at {filepath}")


def local_wrk2(duration, threads, rate):
  import threading
  threads_list = []
  individual_rate = int(int(rate) / int(threads))
  last_id = int(int(rate) * int(duration) // int(threads))
  throughputs = []
  latencies = []
  for i in range(int(threads)):
    throughputs.append(0)
    latencies.append([])
  for i in range(int(threads)):
    thread = threading.Thread(target=run_test, args=(duration, f"http://localhost:{APP_PORT}", individual_rate, last_id * (i + 1), throughputs, latencies, i))
    thread.start()
    threads_list.append(thread)

  for thread in threads_list:
    thread.join()

  throughput = 0
  for i in throughputs:
    throughput += i

  print(f"throughput: {throughput}")
  throughput_str = f'\n----------------------------------------------------------\n Requests/sec:     {throughput}\n'
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
    f.write(output.getvalue().decode() + throughput_str + f"Latency Median:     {median}\n")

  print(output.getvalue().decode() + throughput_str + f"Latency Median:     {median}\n")
  print(f"[INFO] workload results saved at {filepath}")

  metrics('local', timestamp)

def local_metrics(timestamp):
  metrics('local', timestamp)

def local_consistency_window():
  url = 'http://localhost:12345/wrk2-api/user/consistencyWindow'
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

def gcp_tests():
  from plumbum import local
  
  clients = [15]
  manager_script = local["./manager.py"]

  #deploy_command = manager_script["deploy"]["--gcp"]
  #deploy_command()

  #start_command = manager_script["start"]["--gcp"]
  #start_command()
  
  for client in clients:
    rate = str(client * 10)
    clientstr = str(client)
    for i in range(9):
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
    #gcp
     'configure', 'deploy', 'clean', 'consistency-window', 'tests', 'cluster',
    # eval
    'wrk2', 'wrk2-vm', 'metrics', 'wrk2-kube'
  ]

  for cmd in commands:
    parser = command_parser.add_parser(cmd)
    parser.add_argument('--local', action='store_true', help="Running in localhost")
    parser.add_argument('--gcp', action='store_true',   help="Running in gcp")
    if cmd == 'wrk2' or cmd =='wrk2-vm' or cmd =='wrk2-kube':
      parser.add_argument('-d', '--duration', default=30, help="Duration")
      parser.add_argument('-t', '--threads', default=2, help="Number of threads")
      parser.add_argument('-r', '--rate', default=50, help="Number of requests per second")
    if cmd == 'wrk2-vm':
      parser.add_argument('-ts', '--timestamp', help="Timestamp of workload")
    if cmd == 'wrk2-kube' or cmd == 'wrk2-vm' or cmd == 'wrk2':
      parser.add_argument('-ht', '--host', default="localhost", help="Number of requests per second")
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
    
