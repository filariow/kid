"""
before_step(context, step), after_step(context, step)
    These run before and after every step.
    The step passed in is an instance of Step.
before_scenario(context, scenario), after_scenario(context, scenario)
    These run before and after each scenario is run.
    The scenario passed in is an instance of Scenario.
before_feature(context, feature), after_feature(context, feature)
    These run before and after each feature file is exercised.
    The feature passed in is an instance of Feature.
before_all(context), after_all(context)
    These run before and after the whole shooting match.
"""


from steps.kubernetescli import Kubernetes
from steps.util import scenario_id


def before_scenario(context, _scenario):
    context.namespace = scenario_id(context)
    Kubernetes().create_namespace(context.namespace)


def after_scenario(context, _scenario):
    Kubernetes().delete_namespace(context.namespace)
