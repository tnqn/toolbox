# toolbox

## K8s development

### Aliases

#### Installation

```bash
curl -L https://raw.githubusercontent.com/tnqn/toolbox/main/k8s/aliases.sh | bash
```

#### Usage

```bash
# Create a pcap Pod on Node kind-worker and attach to the Pod for packet capture.
$ pcap kind-worker
pod/pcap-kind-worker created
pod/pcap-kind-worker condition met
$ tcpdump -i eth0 -n

# Remove the pcap Pod on Node kind-worker
$ unpcap kind-worker
```