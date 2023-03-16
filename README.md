# KId - Kubernetes Identity

Command Line Application to manage Service Account based Identities.

## How to use it

The command line application get the kubeconfig by looking for the KUBECONFIG environment variable and, if not found, for the default `$HOME/.kube/config` file.

To set the namespace, you can use the `-n` or `--namespace` argument.

## Porcelain commands

### Create an Identity

To create an Identity you can use the following command:

```console
kid create identity "IDENTITY_NAME"
```
As a result, the following resources will be created:
- A Service Account with the name `IDENTITY_NAME`
- A Secret with the name `IDENTITY_NAME-secret-1` and type `kubernetes.io/service-account-token`

### Read the JWT Token for an identity

```console
kid get token "IDENTITY_NAME"
```

As a result it will print in json format the following information:
- CA Certificate
- Namespace
- JWT Token

### Get kubeconfig for an identity

```console
kid get kubeconfig "IDENTITY_NAME"
```

As a result it will print a kubeconfig valid for authenticating as the given Identity.

The following parameters may be overwritten:
- Server URL
- Context's namespace
- Context's username

### Rotate Identity's Token

Key rotation is performed in two steps.
In the first step you will create a new key, in the second you will delete the old one.

```console
kid begin rotation "IDENTITY_NAME"
```

If the last secret for Identity with name `IDENTITY_NAME` is `IDENTITY_NAME-secret-<n>`, a new `IDENTITY_NAME-secret-<n+1>` is created.
You have time to now spread the `IDENTITY_NAME-secret-<n+1>` among the services using that identity.

Once you are done, you can delete the old secret with the following command:

```console
kid complete rotation "IDENTITY_NAME"
```

### Rollback Identity's Token

If you need to resume a deleted token, you can simply recreate the version using the following command:

```console
kid rollback token "IDENTITY_NAME" "VERSION"
```

This command will recreate the token with version `VERSION` for Service Account `IDENTITY_NAME`.
Provided version must be lower than higher existing.

## Plumbing commands 

### Create a new Token Version

To create a new Token version you can use the following command:

```console
kid create token "IDENTITY_NAME"
```

If the last secret for Identity with name `IDENTITY_NAME` is `IDENTITY_NAME-secret-<n>`, a new `IDENTITY_NAME-secret-<n+1>` is created.


### Revoke Identity's Token

To revoke a token version, you can use the following command:

```console
kid revoke token "IDENTITY_NAME" "VERSION"
```

This command will delete the token with version `VERSION` for Service Account `IDENTITY_NAME`.
> Before revoking the last version of a token, please do generate a new one.
> If you revoke a token and then create a new one, the same token you revoked will be created again.
