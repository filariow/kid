from steps.command import Command
from behave import step


class KId(object):

    kid = "./out/kid"

    def create_identity(self, sa: str, namespace: str) -> (str, int):
        return Command().run(f"{self.kid} create identity {sa} -n {namespace}")


@step(u'Identity "{identity}" is created')
def create_identity(context, identity: str):
    ns = context.namespace
    o, e = KId().create_identity(identity, ns)
    assert e == 0, f'error creating identity {ns}/{identity}: {o}'
