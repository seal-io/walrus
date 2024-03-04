<br>

<h1 align="center">Walrus</h1>

<p align="center">
        <a href="../README.md">English</a>&nbsp ｜ &nbsp中文&nbsp
</p>
<br>

Walrus是一个开源的应用管理平台，可以简化在任意基础设施之上应用服务和基础设施资源的部署和管理。它帮助平台工程师为开发人员构建黄金路径，并赋予开发人员自服务的能力。

# Why Walrus

- 一个企业级的XaC（一切即代码）平台，赋能DevOps协作与自服务。
- 利用现有的开源工具和生态，包括Terraform和OpenTofu，提供更强大、更灵活的编排和面向应用的抽象。
- 开发人员不需要学习复杂的Kubernetes或基础设施知识，就可以在任意基础设施之上交付和管理整个应用系统。


# 关键特性

- 统一编排：云原生应用不仅仅是Kubernetes。Walrus提供应用系统从上到下的完整编排，包括应用服务和资源依赖项（如数据库、中间件、负载均衡器和网络等等）。
- 关注点分离：Dev团队可以声明资源类型的需求，通过Ops团队提供的[资源定义](https://seal-io.github.io/docs/zh/operation/resource-definition)进行规则匹配并自动化置备资源，实现Dev团队的自助服务。
- 多态资源管理：单一的资源类型可以根据实际从开发到生产中不同的运行模式和运行环境，转换为多态运行的资源。
- 应用系统全局视图：Walrus通过[依赖图](https://seal-io.github.io/docs/zh/application/graph)提供了整个应用系统的统一视图，包括所有应用服务、基础设施资源及它们的底层组件。
- 动态环境管理：Walrus允许您按需启动或停止应用资源或整个应用环境，从而提高资源利用率和成本效益。
- UI Schema：通过[UI Schema](https://seal-io.github.io/docs/zh/operation/template#自定义模板-ui-样式)定义展示给最终用户如何部署资源的UI表单交互，而不需要修改代码。
- 工作流引擎：Walrus内置了工作流引擎和可扩展的步骤模板，可以满足各种复杂的编排和部署需求。
- AI辅助：通过与[Appilot](https://github.com/seal-io/appilot)集成，您可以通过自然语言交互来管理应用程序与进行故障排查诊断。

# 快速入门

```shell
sudo docker run -d --privileged --restart=unless-stopped -p 80:80 -p 443:443 sealio/walrus
```

在浏览器中打开 `https://<服务器IP或域名>`

# 文档

请参阅[官方文档站点](https://seal-io.github.io/docs/zh/)获取完整的文档。

# 社区与支持

如果您需要任何帮助，请加入以下社区：
- [Discord（英文）](https://discord.gg/fXZUKK2baF)
- [微信群](WECHAT_CN.md)

如果您有任何反馈意见或问题，欢迎通过Github Issues[提交问题](https://github.com/seal-io/walrus/issues/new)。

对于安全问题，请通过发送电子邮件至 <security@seal.io> 进行报告。

# 贡献

如果您有兴趣为Walrus做出贡献，请阅读我们的[贡献指南](./CONTRIBUTING.md)。

# 许可证

版权所有 (c) 2023 [数澈软件](https://seal.io)

根据 Apache 许可证 2.0 版（“许可证”）进行许可；
除非符合许可证，否则您不得使用此文件。
您可以在 [LICENSE](../LICENSE) 文件中获取许可证的详细信息。

除非适用法律要求或书面同意，否则在许可下分发的软件是基于“按原样”分发的，
不提供任何形式的保证或条件，无论是明示的还是隐含的。
有关特定语言下的权限和限制，请参阅许可证。
