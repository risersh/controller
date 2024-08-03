# Orchestration Controller

The orchestration controller is responsible for orchestrating lifecycles of riser services such as deployments, services, ingresses, etc.

The controller is a fairly lightweight service that primarily listens for events from the REST API by means of a RabbitMQ queue and then creates the necessary Kubernetes resources to run the service.

## Architecture

![diagram](docs/diagram.png)

**Terminology:**

- **Tenant**: A tenant is a top level entity that owns one or more riser deployments.
- **Riser Deployment**: A riser deployment is an entity that contains the dependencies needed to run a docker image build by the [builder](https://github.com/risersh/builder) service.

### Initial Tenant Setup

When a new tenant is created the following resources are requested to be created by the REST API:

- Kubernetes Issuer (cert-manager): Used to issue certificates for deployments for a specific namespace.

### New Riser Deployments

When a new riser deployment is created the following resources are requested to be created by the REST API:

- Kubernetes Deployment: Used to run the docker image built by the [builder](https://github.com/risersh/builder) service.
- Kubernetes Service: Used to expose the deployment to traffic.
- Kubernetes HTTPRoute (istio): Used to route traffic to the deployment.
- Kubernetes Certificate (cert-manager): Used to issue certificates for the new deployment URL(s).

## Development

The controller has the following dependencies:

- Connection to the RabbitMQ server for sending and receiving instructions from the REST API and other services.
- Connection to a Kubernetes cluster with cert-manager and istio installed.

> To start infrastructure services like RabbitMQ see the [infrastructure repository](https://github.com/risersh/infrastructure) readme.

### Setup

#### Installing Dependencies

Download the dependencies using go mod:

```bash
go mod download
```

#### Create `.env.local.yaml`

You can either copy the [.env.local.example.yaml](./.env.local.example.yaml) file to `.env.local.yaml` or create your own.

If you need to create your own config create a `.env.local.yaml` file in the root directory containing something like the following:

```yaml
certificates:
  email: matthew@matthewdavis.io
  server: https://acme-staging-v02.api.letsencrypt.org/directory
```

### Running

In the root directory, run the following command to start the service:

```bash
go run .
```

### Testing

#### Running Single Tests

To run a single test in a specific package, run the following command in the desired directory such as:

```bash
cd kubernetes/resources/cert-manager && go test -v -test.run IssuerSuiteRun
```

This will outut something similar to:

```bash
../controller/kubernetes/resources/cert-manager 🌱 main [!?] ✗ cd resources/cert-manager && go test -v -test.run IssuerSuiteRun
=== RUN   TestIssuerSuiteRun
=== RUN   TestIssuerSuiteRun/Test1NewIssuerWithHTTPSolver
--- PASS: TestIssuerSuiteRun (0.11s)
    --- PASS: TestIssuerSuiteRun/Test1NewIssuerWithHTTPSolver (0.11s)
PASS
ok      github.com/risersh/controller/kubernetes/resources/cert-manager 0.556s
```
