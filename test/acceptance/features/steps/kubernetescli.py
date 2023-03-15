import polling2
from steps.environment import ctx
from steps.command import Command
from behave import then


class Kubernetes(object):

    def __init__(self, kubeconfig=None):
        self.cmd = Command()
        if kubeconfig is not None:
            self.cmd.setenv("KUBECONFIG", kubeconfig)

    def create_namespace(self, namespace):
        cmd = f'{ctx.cli} create ns {namespace}'
        output, exit_status = self.cmd.run(cmd)
        if exit_status == 0:
            return output
        return None

    def delete_namespace(self, namespace):
        cmd = f'{ctx.cli} delete ns {namespace}'
        output, exit_status = self.cmd.run(cmd)
        if exit_status == 0:
            return output
        return None

    def resource_exists(self, resource_type: str,
                        resource_name: str, namespace: str) -> bool:
        cmd = f'{ctx.cli} get {resource_type} {resource_name} -n {namespace}'
        _, exit_code = self.cmd.run(cmd)
        return exit_code == 0

    def apply(self, yaml, namespace=None, user=None):
        if namespace is not None:
            ns_arg = f"-n {namespace}"
        else:
            ns_arg = ""
        if user is not None:
            user_arg = f"--user={user}"
        else:
            user_arg = ""
        (output, exit_code) = self.cmd.run(
            f"{ctx.cli} apply {ns_arg} {user_arg} --validate=false -f -", yaml)
        assert exit_code == 0, \
            f"Non-zero exit code ({exit_code}) while applying a YAML: {output}"
        return output

    def delete_by_name(self, res_type, res_name, namespace=None):
        if namespace is not None:
            ns_arg = f"-n {namespace}"
        else:
            ns_arg = ""
        (output, exit_code) = self.cmd.run(
            f"{ctx.cli} delete {res_type} {res_name} {ns_arg}")
        assert exit_code == 0, \
            f"Non-zero exit code ({exit_code}) "\
            f"while deleting a Custom Resource: {output}"
        return output

    def service_account_exists(self, sa_name, namespace) -> (bool):
        o, ec = self.cmd.run(
            f'{ctx.cli} get serviceaccounts {sa_name} -n {namespace}')
        return ec == 0

    def secret_exists(self, secret, namespace) -> (bool):
        o, ec = self.cmd.run(f'{ctx.cli} get secrets {secret} -n {namespace}')
        return ec == 0


# Behave steps

@then(u'Service Account "{sa_name}" exists')
def service_account_exists(context, sa_name: str):
    k = Kubernetes()
    polling2.poll(
        target=lambda: k.service_account_exists(sa_name, context.namespace),
        step=1,
        timeout=30)


@then(u'Secret "{secret}" exists')
def secret_exists(context, secret: str):
    k = Kubernetes()
    polling2.poll(
        target=lambda: k.secret_exists(secret, context.namespace),
        step=1,
        timeout=30)
