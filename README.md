# SE-360 Project

A microservices-based ride-hailing application deployed on AWS EKS with Traefik ingress controller.

## Architecture Overview

This project consists of:
- **Infrastructure**: AWS VPC with public/private subnets managed by Terraform
- **Kubernetes Cluster**: EKS cluster with OIDC provider for IAM integration
- **Ingress**: Traefik ingress controller with CORS and authentication middlewares
- **Load Balancing**: AWS Load Balancer Controller for managing ALB/NLB
- **Services**: User, Driver, and Trip microservices
![alt text](<Infra.jpeg>)
![alt text](<Architect.jpeg>)


## Prerequisites

- AWS CLI configured with appropriate credentials
- Terraform >= 1.0
- kubectl
- eksctl
- Helm 3.x

## Setup Instructions

Follow these steps in order to set up the complete infrastructure and deploy services:

### 1. Create Infrastructure with Terraform

Create the VPC, subnets, NAT gateway, and internet gateway:

```bash
cd terraform
terraform init
terraform plan
terraform apply
```

This will create:
- VPC (`se-360-vpc`) with CIDR `10.0.0.0/16`
- 3 public subnets across 3 availability zones
- 3 private subnets across 3 availability zones
- Internet Gateway for public subnets
- NAT Gateway for private subnets
- Route tables and associations

**Note**: After applying, note down the VPC ID and subnet IDs for the next step.

### 2. Create EKS Cluster

Create the EKS cluster using the configuration in the `eks` folder:

```bash
eksctl create cluster -f eks/cluster.yaml
```

This will:
- Create an EKS cluster named `se-360-vpc` in `ap-southeast-1` region
- Use the existing VPC created by Terraform
- Enable IAM OIDC provider for service accounts
- Create a node group (`ng-main-new`) with 3 t3.medium instances in private subnets
- Configure public and private API endpoints

**Cluster creation takes approximately 15-20 minutes.**

After creation, verify the cluster:

```bash
kubectl get nodes
```

### 3. Create Service Account for AWS Load Balancer Controller

Create an IAM service account for the AWS Load Balancer Controller:

```bash
eksctl create iamserviceaccount \
  --cluster se-360-vpc \
  --namespace kube-system \
  --name aws-load-balancer-controller \
  --attach-policy-arn <your arn policy> \
  --override-existing-serviceaccounts \
  --region ap-southeast-1 \
  --approve
```

**Note**: If the IAM policy doesn't exist, create it first using the [official AWS documentation](https://docs.aws.amazon.com/eks/latest/userguide/aws-load-balancer-controller.html).

Verify the service account:

```bash
kubectl get serviceaccount aws-load-balancer-controller -n kube-system
```

### 4. Install AWS Load Balancer Controller with Helm

Add the EKS Helm repository and install the AWS Load Balancer Controller:

```bash
# Add the EKS chart repository
helm repo add eks https://aws.github.io/eks-charts
helm repo update

# Install the AWS Load Balancer Controller
helm install aws-load-balancer-controller eks/aws-load-balancer-controller \
  -n kube-system \
  -f helm/aws-load-balancer/values-aws-load-balancer-controller.yaml
```

Verify the installation:

```bash
kubectl get deployment -n kube-system aws-load-balancer-controller
kubectl logs -n kube-system deployment/aws-load-balancer-controller
```

### 5. Install Traefik Ingress Controller with Helm

Add the Traefik Helm repository and install Traefik with custom values:

```bash
# Add Traefik Helm repository
helm repo add traefik https://traefik.github.io/charts
helm repo update

# Install Traefik using the installation script
chmod +x scripts/helm-install-traefik.sh
./scripts/helm-install-traefik.sh
```

This script will:
- Install Traefik with custom values from `helm/traefik/values.yaml`
- Apply authentication middleware (`helm/traefik/middlewares-auth.yaml`)
- Apply CORS middleware (`helm/traefik/middlewares-cors.yaml`)

Verify Traefik installation:

```bash
kubectl get pods -l app.kubernetes.io/name=traefik
kubectl get svc traefik
```

### 6. Install Microservices with Helm

Install the microservices in the following order:

#### User Service

```bash
chmod +x scripts/helm-install-user-service.sh
./scripts/helm-install-user-service.sh
```

#### Driver Service

```bash
chmod +x scripts/helm-install-driver-service.sh
./scripts/helm-install-driver-service.sh
```

#### Trip Service

```bash
chmod +x scripts/helm-install-trip-service.sh
./scripts/helm-install-trip-service.sh
```

Verify all services are running:

```bash
kubectl get pods
kubectl get svc
kubectl get ingress
```

## Configuration Files

### Terraform
- `terraform/main.tf` - VPC and networking infrastructure

### EKS
- `eks/cluster.yaml` - EKS cluster configuration with eksctl

### Helm Charts
- `helm/aws-load-balancer/values-aws-load-balancer-controller.yaml` - AWS LB Controller configuration
- `helm/traefik/values.yaml` - Traefik ingress controller configuration
- `helm/traefik/middlewares-auth.yaml` - Authentication middleware
- `helm/traefik/middlewares-cors.yaml` - CORS middleware
- `helm/user-service/` - User service Helm chart
- `helm/driver-service/` - Driver service Helm chart
- `helm/trip-service/` - Trip service Helm chart

### Scripts
- `scripts/helm-install-traefik.sh` - Traefik installation script
- `scripts/helm-install-user-service.sh` - User service installation script
- `scripts/helm-install-driver-service.sh` - Driver service installation script
- `scripts/helm-install-trip-service.sh` - Trip service installation script

## Accessing Services

After deployment, get the Traefik load balancer URL:

```bash
kubectl get svc traefik -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
```

Services will be accessible through the Traefik ingress routes configured in each service's Helm chart.

## Cleanup

To tear down the infrastructure:

```bash
# Uninstall Helm releases
helm uninstall trip-service
helm uninstall driver-service
helm uninstall user-service
helm uninstall traefik
helm uninstall aws-load-balancer-controller -n kube-system

# Delete the EKS cluster
eksctl delete cluster -f eks/cluster.yaml

# Destroy Terraform infrastructure
cd terraform
terraform destroy
```

## Additional Documentation

- [Traefik Local Setup](TRAEFIK-LOCAL-SETUP.md) - Local development setup with Docker Compose
- [Docker Compose CORS Guide](docker-compose-cors-guide.md) - CORS configuration for local development
- [Traefik CORS README](helm/traefik/README-CORS.md) - CORS configuration for Kubernetes

## Troubleshooting

### EKS Cluster Issues
- Ensure VPC ID and subnet IDs in `eks/cluster.yaml` match the Terraform outputs
- Verify IAM OIDC provider is enabled: `eksctl utils associate-iam-oidc-provider --cluster=se-360-vpc --approve`

### AWS Load Balancer Controller Issues
- Check service account annotations: `kubectl describe sa aws-load-balancer-controller -n kube-system`
- Review controller logs: `kubectl logs -n kube-system deployment/aws-load-balancer-controller`

### Traefik Issues
- Check Traefik logs: `kubectl logs -l app.kubernetes.io/name=traefik`
- Verify middlewares are applied: `kubectl get middleware`

### Service Deployment Issues
- Check pod status: `kubectl get pods`
- Review pod logs: `kubectl logs <pod-name>`
- Verify ingress routes: `kubectl get ingressroute`