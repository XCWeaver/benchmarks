def gen_ansible_vars():
  import yaml

  with open('deploy/memorystorage/envoy.yaml', 'r') as file:
    data = yaml.safe_load(file)

  data['static_resources']['clusters'][0]['load_assignment']['endpoints'][0]['lb_endpoints'][0]['endpoint']['address']['socket_address']['address'] = "endpoint"
  data['static_resources']['clusters'][1]['load_assignment']['endpoints'][0]['lb_endpoints'][0]['endpoint']['address']['socket_address']['address'] = "endpoint2"
  
  with open('deploy/tmp/test.yml', 'w') as file:
    yaml.dump(data, file)

gen_ansible_vars()