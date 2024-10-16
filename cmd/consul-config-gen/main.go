package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
)

const configTemplate = `
datacenter = "{{ .DataCenter }}"
data_dir = "{{ .DataDir }}"
log_level = "INFO"
node_name = "{{ .NodeName }}"
retry_join = [{{ .RetryJoinIP }}]

node_meta {
  env = "{{ .NodeEnv }}"
  stage = "{{ .NodeStage }}"
}
`

type TemplateVars struct {
	DataCenter  string
	DataDir     string
	NodeName    string
	RetryJoinIP string
	NodeEnv     string
	NodeStage   string
}

var (
	templateVars     TemplateVars = TemplateVars{}
	retryJoinIPsList string
)

func init() {
	flag.StringVar(&templateVars.DataCenter, "data-center", "dc1", "Data center")
	flag.StringVar(&templateVars.DataDir, "data-dir", "/consul-home", "Data dir for consul agent")
	flag.StringVar(&templateVars.NodeName, "node-name", "some-node", "Node name for consul agent")
	flag.StringVar(&retryJoinIPsList, "retry-join-ip", "10.10.0.5", "IP of the consul server")
	flag.StringVar(&templateVars.NodeEnv, "node-env", "metrics", "Node env for consul node metadata")
	flag.StringVar(&templateVars.NodeStage, "node-stage", "prod", "Node stage for consul node metadata")

	flag.Parse()
}

func main() {
	tmpl, err := template.New("config.hcl").Parse(configTemplate)
	if err != nil {
		panic(err)
	}

	retryIps := []string{}
	for _, ip := range strings.Split(retryJoinIPsList, ",") {
		retryIps = append(retryIps, fmt.Sprintf("\"%s\"", ip))
	}

	templateVars.RetryJoinIP = strings.Join(retryIps, ",")
	if err := tmpl.Execute(os.Stdout, templateVars); err != nil {
		panic(err)
	}
}
