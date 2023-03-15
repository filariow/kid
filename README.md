# KSA - Kubernetes Service Account

Library to extract data needed to authenticate as a Service Account.
This data can be used to build a kubeconfig.


## How to use it

The command line application get the kubeconfig by looking for the KUBECONFIG environment variable and, if not found, for the default `$HOME/.kube/config` file.

To set the namespace, you can use the `-n` or `--namespace` argument.


## Create an Identity

To create an Identity you can use the following command:

```console
ksa create identity "IDENTITY_NAME"
```
As a result, the following resources will be created:
- A Service Account with the name `IDENTITY_NAME`
- A Secret with the name `IDENTITY_NAME-secret-1` and type `kubernetes.io/service-account-token`

## Read the JWT Token for an identity

```console
ksa get token "IDENTITY_NAME"
```

As a result it will print in json format the following information:
- CA Certificate
- Namespace
- JWT Token

## Get kubeconfig for an identity

```console
ksa get kubeconfig "IDENTITY_NAME"
```

As a result it will print a kubeconfig valid for authenticating as the given Identity.

The following parameters may be overwritten:
- Server URL
- Context's namespace
- Context's username
