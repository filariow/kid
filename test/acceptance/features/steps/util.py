import os
import time
import yaml
from string import Template
from behave import step
from kubernetes import client, config


def get_api_client_from_kubeconfig(kubeconfig: str) -> client:
    kcd = yaml.safe_load(kubeconfig)
    api_client = config.new_client_from_config_dict(kcd)
    return api_client


def scenario_id(context):
    ssfn = os.path.splitext(context.scenario.filename)[0]
    sf = os.path.basename(ssfn).lower()
    sl = context.scenario.line
    return f"{sf}-{sl}"


def substitute_scenario_id(context, text="$scenario_id"):
    return Template(text).substitute(scenario_id=scenario_id(context))


def get_env(env):
    value = os.getenv(env)
    assert env is not None, f"{env} environment variable needs to be set"
    print(f"{env} = {value}")
    return value


# Behave steps
@step(u'{duration} seconds have passed')
def wait(context, duration):
    time.sleep(float(duration))


@step(u'1 second has passed')
def wait_1_s(context):
    wait(context, 1)
