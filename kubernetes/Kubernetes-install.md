# Installation notes

## Setup VM

### Step 1: Load kernel modules

```bash
sudo modprobe overlay
sudo modprobe br_netfilter

sudo vim /etc/modules
# Add the following lines
# overlay
# br_netfilter
```

### Step 2: Load iptables

```bash
sudo vim /etc/sysctl.d/k8s.conf
# Add the following lines
# net.bridge.bridge-nf-call-iptables = 1
# net.ipv4.ip_forward = 1
# net.bridge.bridge-nf-call-ip6tables = 1

sudo sysctl --system
```

### Step 3: Disable swap

```bash
sudo vim /etc/fstab
# Remove the following line
# /swap.img none swap sw 0 0

sudo swapoff -a
sudo rm -f /swap.img
```

## Install container runtime

### Step 1: Install containerd

```bash
wget https://github.com/containerd/containerd/releases/download/v1.6.25/containerd-1.6.25-linux-amd64.tar.gz
sudo tar Cxzvf /usr/local containerd-1.6.25-linux-amd64.tar.gz
```

### Step 2: Enable containerd service

```bash
wget https://raw.githubusercontent.com/containerd/containerd/main/containerd.service
sudo mv containerd.service /usr/lib/systemd/system/
sudo systemctl enable --now containerd
```

### Step 3: Install runc

```bash
wget https://github.com/opencontainers/runc/releases/download/v1.1.10/runc.amd64
sudo install -m 755 runc.amd64 /usr/local/sbin/runc
```

### Step 4: Reload daemon

```bash
	sudo systemctl daemon-reload
```

## Install Kubernetes

```bash
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl gpg
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.28/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.28/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
```

## Create Kubernetes Control Panel Node

### Initialize Node as Control Panel Node

```bash
sudo kubeadm init
```

### Create configuration

```bash
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

### Setup Container Network Interface

```bash
wget https://docs.projectcalico.org/manifests/calico.yaml
kubectl apply -f calico.yaml
```

## Create Kubernetes Worker Node

On control plane:

```bash
sudo kubeadm token create
openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //'
```

On worker nodes:

```bash
kubeadm join --token <token> <control-plane-host>:<control-plane-port> --discovery-token-ca-cert-hash sha256:<hash>
```
