# APIRegistration Gen

This directory contains the code generators for the APIService API.

# Marker Usage

- Package Markers
    - `+groupName=example.domain.io`: Specify the group name.
    - `+versionName=v1`: Specify the version.
    - `+k8s:apireg-gen:service:insecureSkipTLSVerify=,groupPriorityMinimum=,versionPriority=`: Specify the configuration
      of the APIService.
        - `insecureSkipTLSVerify`: Select from `true` or `false`, default is `false`.
        - `groupPriorityMinimum`: Specify the minimum priority of the group, must between 0 and 2000, default is `100`.
        - `versionPriority`: Specify the priority of the version, must not be negative, default is `100`.
- Type Markers
    - `+k8s:apireg-gen:resource:scope="Namespaced",categories=["..."],shortName=["..."],plural="...",subResources=["status","scale"]`:
      Specify the resource.
        - `scope`: Select from `Namespaced` or `Cluster`.
        - `categories`, `shortName` and `subResources`: Specify in the JSON array format.
        - `subResources` only supports `status` and `scale`.

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