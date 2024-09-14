#!/bin/bash
echo "*********************************"
echo "Wait for RabbitMQ nodes"
echo "*********************************"
# test connection with
# rabbitmqctl -n rabbit@<server> environment
dockerize -wait tcp://notifier-rabbitmq-us:5672 -wait tcp://notifier-rabbitmq-eu:5672 -wait-retry-interval 10s -timeout 300s echo "[INFO] notifier-rabbitmq ready!"
rabbitmqctl -n rabbit@notifier-rabbitmq-us await_startup
rabbitmqctl -n rabbit@notifier-rabbitmq-eu await_startup
echo "*********************************"
echo "RabbitMQ nodes READY!"
echo "*********************************"

echo "*********************************"
echo "Build RabbitMQ Federation"
echo "*********************************"
# federate notifier-rabbitmq-us to notifier-rabbitmq-eu
config='{"max-hops": 1, "uri": ["amqp://admin:admin@notifier-rabbitmq-eu"]}'
rabbitmqctl -n rabbit@notifier-rabbitmq-us set_parameter federation-upstream cluster1 "${config}"
config='[{"upstream": "cluster1"}]'
rabbitmqctl -n rabbit@notifier-rabbitmq-us set_parameter federation-upstream-set cluster1_federators "${config}"
config='{"federation-upstream-set": "cluster1_federators"}'
rabbitmqctl -n rabbit@notifier-rabbitmq-us set_policy --apply-to exchanges federation_test "notifier\-*" "${config}"
config='{"ha-mode": "all"}'
rabbitmqctl -n rabbit@notifier-rabbitmq-us set_policy ha-federation "^federation:*" "${config}"

# federate notifier-rabbitmq-eu to notifier-rabbitmq-us
config='{"max-hops": 1, "uri": ["amqp://admin:admin@notifier-rabbitmq-us"]}'
rabbitmqctl -n rabbit@notifier-rabbitmq-eu set_parameter federation-upstream cluster2 "${config}"
config='[{"upstream": "cluster2"}]'
rabbitmqctl -n rabbit@notifier-rabbitmq-eu set_parameter federation-upstream-set cluster2_federators "${config}"
config='{"federation-upstream-set": "cluster2_federators"}'
rabbitmqctl -n rabbit@notifier-rabbitmq-eu set_policy --apply-to exchanges federation_test "notifier\-*" "${config}"
config='{"ha-mode": "all"}'
rabbitmqctl -n rabbit@notifier-rabbitmq-eu set_policy ha-federation "^federation:*" "${config}"

echo "*********************************"
echo "RabbitMQ Federation DONE!"
echo "*********************************"


echo "*********************************"
echo "Opening HTTP:8000 server for dockerize coordination"
echo "*********************************"

ran -p 8000 -l -r /tmp/

echo "*********************************"
echo "HTTP:8000 server for dockerize coordination DONE!"
echo "*********************************"
