# Orchestration Controller

The orchestration controller is responsible for orchestrating lifecycles of riser services such as deployments, services, ingresses, etc.

## Dependencies

The controller is fairly lightweight and does have a few dependencies:

- Connection to the RabbitMQ server for sending and receiving instructions from the REST API and other services.
- Connection to a Kubernetes cluster with cert-manager and istio installed.

> To start infrastructure services like RabbitMQ see the [infrastructure repository](https://github.com/risersh/infrastructure) readme.

## Development

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
