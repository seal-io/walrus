# Webhook(Configuration) Gen

This directory contains the code generators for the Webhook Configuration API.

# Marker Usage

- Type Markers:
    - `+k8s:webhook-gen:validating:*`: Specify to generate ValidatingWebhookConfiguration.
    - `+k8s:webhook-gen:mutating:*`: Specify to generate MutatingWebhookConfiguration.
    - `+k8s:webhook-gen:*:group="example.domain.io",version="v1",resource="dummies",scope="Namespaced"`: Specify the
      group, version and resource name, which is used for generate the webhook configuration's path.
        - `scope`: Select from `Namespaced`, `Cluster` or `*`.
    - `+k8s:webhook-gen:*:operations=["..."],failurePolicy="Fail",sideEffects="None",matchPolicy="Equivalent",timeoutSeconds=10,reinvocationPolicy="Never"`:
      Specify the configuration of the webhook.
        - `operations`: Specify in the JSON array format. Combine by `CREATE`, `UPDATE`, `DELETE`, `CONNECT`, or `*`.
        - `failurePolicy`: Select from `Fail` or `Ignore`.
        - `sideEffects`: Select from `Unknown`, `None`, `Some` or `NoneOnDryRun`.
        - `matchPolicy`: Select from `Exact` or `Equivalent`.
        - `timeoutSeconds`: Specify the timeout seconds.
        - `reinvocationPolicy`: Select from `Never` or `IfNeeded`, works on MutatingWebhookConfiguration only.
    - `+k8s:webhook-gen:*:namespaceSelector={}`: Specify the namespace selector for the webhook configuration.
    - `+k8s:webhook-gen:*:objectSelector={}`: Specify the object selector for the webhook configuration.
    - `+k8s:webhook-gen:*:matchConditions=[{}]`: Specify the match conditions for the webhook configuration.

# License

Copyright (c) 2024 [Seal, Inc.](https://seal.io)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at [LICENSE](../../LICENSE) file for details.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.