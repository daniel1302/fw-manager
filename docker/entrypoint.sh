#!/bin/bash


	# flag.StringVar(&templateVars.DataCenter, "data-center", "dc1", "Data center")
	# flag.StringVar(&templateVars.DataDir, "data-dir", "/consul-home", "Data dir for consul agent")
	# flag.StringVar(&templateVars.NodeName, "node-name", "some-node", "Node name for consul agent")
	# flag.StringVar(&templateVars.RetryJoinIP, "retry-join-ip", "10.5.0.5", "IP of the consul server")
	# flag.StringVar(&templateVars.NodeEnv, "node-env", "metrics", "Node env for consul node metadata")
	# flag.StringVar(&templateVars.NodeStage, "node-stage", "prod", "Node stage for consul node metadata")

DATA_CENTER="${DATA_CENTER:-dc1}";
DATA_DIR="${DATA_DIR:-/consul/home}";
NODE_NAME="${NODE_NAME:-empty-node-name}";
RETRY_JOIN_IP="${RETRY_JOIN_IP:-10.5.0.5}";
NODE_ENV="${NODE_ENV:-metrics}";
NODE_STAGE="${NODE_STAGE:-prod}";
# generate config
/consul-config-gen \
    --data-center "${DATA_CENTER}" \
    --data-dir "${DATA_DIR}" \
    --node-name "${NODE_NAME}" \
    --retry-join-ip "${RETRY_JOIN_IP}" \
    --node-env "${NODE_ENV}" \
    --node-stage "${NODE_STAGE}" \
| tee /etc/consul-agent.hcl;

# run agent in the background
consul agent --config-file "/etc/consul-agent.hcl" --ui &
sleep 3;
# run service
/consul-svc \
    --name wireguard \
    --env "${NODE_ENV}" \
    --stage "${NODE_STAGE}" \
    --address "localhost" \
    --port 8081 &

echo "Services started"

# wait untill any process exit
wait -n

exit $?