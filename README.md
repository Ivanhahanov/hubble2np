# hubble2np
Utility to analyze Hubble Flow and convert it to Cilium Network Policies
## Install
You can download `hubble2np` from [Github Releases](https://github.com/Ivanhahanov/hubble2np/releases)
### Macos
Download for MacOS on ARM
```bash
wget -c https://github.com/Ivanhahanov/hubble2np/releases/download/0.1.0/hubble2np_Darwin_arm64.tar.gz -O - | tar -xz
```
### Linux 
Download for Linux on AMD
```bash
wget -c https://github.com/Ivanhahanov/hubble2np/releases/download/0.1.0/hubble2np_Linux_x86_64.tar.gz
 -O - | tar -xz
```
### Move to system `PATH`
Move the `hubble2np` binary to a file location on your system `PATH`
```bash
sudo mv ./hubble2np /usr/local/bin/hubble2np
```
## Usage
`hubble2np` works with hubble flow in `json` format. To generate a policy or graph you need to pass json as input using `stdin`.
```bash
hubble observe -n dev --since 1m -o json | hubble2np
```
> [!TIP]
> Generated policies can be redirected to a file
> ```bash
> ... | hubble2np > policies.json
> ```
To get the flow you can use the [hubble](https://docs.cilium.io/en/stable/observability/hubble/hubble-cli/) utility.
You need port-forward to access the hubble api.
```bash
cilium hubble port-forward&
```
You can also start port-forward with kubectl.

### Read from file 
If there is no direct access to the hubble api, you can read the stream from a prepared file.
```bash
# export flow to json
hubble observe -n dev --since 1m -o json > flow.json
# generate policies
cat flow.json | hubble2np
```
### Show Graph
To view the `graph` you can use the corresponding command
```bash
cat flow.json | hubble2np graph
```
> [!TIP]
> This graph can be used for debugging or checking the correctness of input data for policy generation 

## Options
```
NAME:
   hubble2np - generate Cilium Network Policies from Hubble flow

USAGE:
   hubble2np [global options] [command [command options]]

COMMANDS:
   graph    Show graph
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --nodns      disable dns (default: false)
   --ports, -p  enable ports (default: false)
   --help, -h   show help
```

## Example Policy
```yaml
# cat flow.json | hubble2np -p
---
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  creationTimestamp: null
  name: app
  namespace: test
spec:
  egress:
  - toEndpoints:
    - matchLabels:
        app: app
        io.kubernetes.pod.namespace: dev
    toPorts:
    - ports:
      - port: "8080"
  - toEndpoints:
    - matchLabels:
        k8s-app: kube-dns
        io.kubernetes.pod.namespace: kube-system
    toPorts:
    - ports:
      - port: "53"
        protocol: UDP
      rules:
        dns:
        - matchPattern: '*'
  enableDefaultDeny: {}
  endpointSelector:
    matchLabels:
      app: app
      io.kubernetes.pod.namespace: test
  ingress:
  - fromEndpoints:
    - matchLabels:
        app: api
        io.kubernetes.pod.namespace: dev
    toPorts:
    - ports:
      - port: "8080"
  - fromEndpoints:
    - matchLabels:
        app: app
        io.kubernetes.pod.namespace: dev
    toPorts:
    - ports:
      - port: "8080"
```

## Example graph
```bash
# cat flow.json | hubble2np graph
[test/app] -> dev/app -> [test/app,dev/api]
[dev/app,dev/api] -> test/app -> [dev/app]
[dev/app] -> dev/api -> [test/app]
```