[![Documentation](https://img.shields.io/github/actions/workflow/status/seal-io/docs/pages%2Fpages-build-deployment?label=Documentation)](https://seal-io.github.io/docs/quickstart)
[![Releases](https://img.shields.io/github/v/release/seal-io/walrus)](https://github.com/seal-io/walrus/releases/latest)
[![LICENSE](https://img.shields.io/github/license/seal-io/walrus)](/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/seal-io/walrus)](https://goreportcard.com/report/github.com/seal-io/walrus)
[![Docker Pulls](https://img.shields.io/docker/pulls/sealio/walrus)](https://hub.docker.com/r/sealio/walrus)

<br>

<h1 align="center">Walrus</h1>

<p align="center">
        English&nbsp | &nbsp<a href="docs/README_CN.md">中文</a>&nbsp
</p>
<br>

Walrus is an open-source application management platform that simplifies application deployment and management on any infrastructure. 
It helps platform engineers build golden paths for developers and empowers developers with self-service capabilities.


# Key Features

- <b>Unified Orchestration</b> ：Cloud-Native applications are more than Kubernetes. Walrus orchestrates the entire application system, encompassing both application services and resource dependencies (such as databases, middleware, load balancers, and networks).
- <b>Seperation of Concerns</b> : By leveraging [resource definitions](https://seal-io.github.io/docs/operation/resource-definition) provided by ops team, dev team can define resource type requirements and acheive self-service.
- <b>Polymorphic Resource Management</b> : A single resource type can be translated into polymorphic resources that operate in various modes and environments, from development to production. 
- <b>Single Pane of Glass</b> : Walrus provides a unified view of your entire application system in a [dependency graph](https://seal-io.github.io/docs/application/graph), including all application services, resources, and their sub-components.
- <b>Dynamic Environment Management</b> : Walrus allows you to start or stop application resources or environments as needed, enhancing hardware resource utilization and cost-effectiveness.
- <b>UI Schema</b> : Define how end users interact with and utilize the resources through the [UI schema](you can manage applications and troubleshoot issues through natural language interaction.), without the need for code modification.
- <b>Workflow Engine</b> : Walrus includes a built-in workflow engine with extensible step templates, catering to diverse requirements for complex orchestration and deployment. 
- <b>AI Assistance</b> : With [Appilot](https://github.com/seal-io/appilot) integration, you can manage applications and troubleshoot issues through natural language interaction.

# Quick Start

```shell
sudo docker run -d --privileged --restart=unless-stopped -p 80:80 -p 443:443 sealio/walrus
```

Open your browser to `https://<server-ip-or-domain>`

# Documentation

Please see [the official docs site](https://seal-io.github.io/docs/) for complete documentation.

# Community and Support

If you need any help, please join us at
- [Discord](https://discord.gg/fXZUKK2baF)
- [WeChat](docs/WECHAT_CN.md)

Feel free to [file an issue](https://github.com/seal-io/walrus/issues/new) if you have any feedback or questions.

For security issues, please report by sending an email to <security@seal.io>.

# Contributing

Please read our [contributing guide](./docs/CONTRIBUTING.md) if you're interested in contributing to Walrus.

# License

Copyright (c) 2023 [Seal, Inc.](https://seal.io)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at [LICENSE](./LICENSE) file for details.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
