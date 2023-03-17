from steps.command import Command
from behave import when, given


class KId(object):

    kid = "./out/kid"

    def create_identity(self, sa: str, namespace: str) -> (str, int):
        return Command().run(f"{self.kid} create identity {sa} -n {namespace}")

    def begin_rotation(self, sa: str, namespace: str) -> (str, int):
        return Command().run(f"{self.kid} begin rotation {sa} -n {namespace}")

    def complete_rotation(self, sa: str, namespace: str) -> (str, int):
        return Command().run(
            f"{self.kid} complete rotation {sa} -n {namespace}")

    def rollback_token(self, identity: str, version: str,
                       namespace: str) -> (str, int):
        return Command().run(
            f'{self.kid} rollback token {identity} {version} -n {namespace}')


@given(u'Identity "{identity}" is created')
@when(u'Identity "{identity}" is created')
def create_identity(context, identity: str):
    ns = context.namespace
    o, e = KId().create_identity(identity, ns)
    assert e == 0, \
        f'error creating identity {ns}/{identity}: {o}'


@given(u'Token rotation begins for identity "{identity}"')
@when(u'Token rotation begins for identity "{identity}"')
def begin_key_rotation(context, identity: str):
    ns = context.namespace
    o, e = KId().begin_rotation(identity, ns)
    assert e == 0, 'error beginning key rotation'\
        f'for identity {ns}/{identity}: {o}'


@given(u'Token rotation completes for identity "{identity}"')
@when(u'Token rotation completes for identity "{identity}"')
def complete_key_rotation(context, identity: str):
    ns = context.namespace
    o, e = KId().complete_rotation(identity, ns)
    assert e == 0, 'error completing key rotation' \
        f'for identity {ns}/{identity}: {o}'


@when(u'Token for identity "{identity}" '
      'with version "{version}" is rolled back')
def rollback_token(context, identity: str, version: str):
    ns = context.namespace
    o, e = KId().rollback_token(identity, version, ns)
    assert e == 0, 'error rolling back token '\
        f'for identity {ns}/{identity}: {o}'
