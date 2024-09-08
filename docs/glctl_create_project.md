## glctl create project

``` yaml
apiVersion: v1
kind: Project
metadata:
  name: "Project1"  # 项目的名称。
  group: Group1
  annotations:
    description: "my first project"  # 项目的简短描述。
    avatar: ""  # 项目头像的 URL。
    markedForDeletion: false  # 项目是否标记为删除。
    deletionDate: null  # 项目计划删除的日期。
  labels:
    tags: [] # 项目的标签列表。
spec:
  visibility: "private"  # 项目的可见性设置：public（公开）、internal（内部）、private（私有）。
  repository:
    initializeWithReadme: true  # 是否初始化项目时自动创建一个 README 文件。
    defaultBranch: "main"  # 默认的分支名称。
    autocloseReferencedIssues: false # 设置是否在默认分支上自动关闭引用的问题
    
    merge:
      enabled: true   # 是否启用合并请求功能。
      method: "merge"  # 合并方法：merge（合并）、rebase_merge（变基合并）、ff（快进合并）。
      approvalsRequired: 1  # 合并请求所需的审批数量。
      removeSourceBranch: true  # 合并后是否删除源分支。
      restrictMergeConditions:
        pipelineSucceeds: false  # 是否仅在流水线成功时允许合并。
        allDiscussionsResolved: false  # 是否仅在所有讨论已解决时允许合并。
        allowMergeOnSkippedPipeline: true  # 是否允许将合并请求与跳过的作业合并。
        onlyAllowMergeIfAllDiscussionsAreResolved: true  # 是否仅在所有讨论解决后合并。
        onlyAllowMergeIfAllStatusChecksPassed: false  # 是否仅在所有状态检查通过后合并请求。

  settings:
    issuesEnabled: true  # 是否启用问题跟踪功能。
    snippetsEnabled: true  # 是否启用代码片段功能。
    packagesEnabled: true  # 是否启用项目包功能。
    securityEnabled: true  # 是否启用安全功能。
    dependencyScanningEnabled: true  # 是否启用依赖项扫描。
    runners:
      groupRunnersEnabled:  true  # 是否启用组级别的 Runner。
      sharedRunnersEnabled: true  # 是否启用共享 Runner。    
    dastEnabled: true  # 是否启用动态应用程序安全测试（DAST）。
    secretDetectionEnabled: true  # 是否启用秘密检测。
    lfsEnabled: true  # 启用大文件存储。
    requestAccessEnabled: true  # 启用请求访问功能。
    emailsDisabled: false  # 是否禁用电子邮件通知。
    approvalsBeforeMerge: 2  # 默认情况下需要的审批者数量
    showDefaultAwardEmojis: true  # 是否显示默认的表情符号回应。
    printingMergeRequestLinkEnabled: true  # 从命令行推送时是否显示创建/查看合并请求的链接。

  cicd:
    enabled: true  # 是否启用 CI/CD 作业。
    publicJobs: true  # CI 作业是否对所有人公开。
    timeout: 3600  # CI 作业的超时时间（秒）。
    gitStrategy: "fetch"  # Git 策略：clone（克隆）、fetch（获取）。
    coverageRegex: ""  # 用于测试覆盖率的正则表达式。
    configPath: ""  # 自定义 CI 配置文件的路径。
    autoCancelPendingPipelines: "enabled"  # 自动取消待处理管道：enabled（启用）或 disabled（禁用）。
    autoDevops:
      enabled: false  # 是否启用 Auto DevOps。
      deployStrategy: "continuous"  # 部署策略：continuous（持续部署）、manual（手动部署）。

  compliance:
    frameworks: []  # 合规框架的列表。
    restrictVariables: false  # 是否限制用户定义的 CI/CD 变量。
    externalAuthorizationLabel: ""  # 外部授权的分类标签。
    licenseScanningEnabled: true  # 是否启用许可证扫描。

  importSettings:
    url: ""  # 用于导入项目的 URL。
    type: ""  # 导入的类型。
    source: ""  # 导入的来源。
    status: ""  # 导入的状态。

  visibilitySettings:
    analyticsAccessLevel: "private"  # 分析功能的可见性级别：private（私有）、enabled（启用）。
    buildsAccessLevel: "private"  # 构建功能的可见性级别：private（私有）、enabled（启用）、public（公开）。
    containerRegistryAccessLevel: "private"  # 容器注册表的可见性级别：private（私有）、enabled（启用）。
    environmentsAccessLevel: "private"  # 环境管理的可见性级别：private（私有）、enabled（启用）。
    featureFlagsAccessLevel: "private"  # 功能标记的可见性级别：private（私有）、enabled（启用）。
    forkingAccessLevel: "private"  # 分叉功能的可见性级别：private（私有）、enabled（启用）。
    infrastructureAccessLevel: "private"  # 基础设施管理的可见性级别：private（私有）、enabled（启用）。
    issuesAccessLevel: "private"  # 问题跟踪的可见性级别：private（私有）、enabled（启用）。
    mergeRequestsAccessLevel: "private"  # 合并请求的可见性级别：private（私有）、enabled（启用）。
    modelExperimentsAccessLevel: "private"  # 模型实验的可见性级别：private（私有）、enabled（启用）。
    modelRegistryAccessLevel: "private"  # 模型注册表的可见性级别：private（私有）、enabled（启用）。
    monitorAccessLevel: "private"  # 监控功能的可见性级别：private（私有）、enabled（启用）。
    pagesAccessLevel: "private"  # 页面功能的可见性级别：private（私有）、enabled（启用）。
    releasesAccessLevel: "private"  # 发布功能的可见性级别：private（私有）、enabled（启用）。
    repositoryAccessLevel: "private"  # 仓库的可见性级别：private（私有）、enabled（启用）。
    requirementsAccessLevel: "private"  # 需求管理的可见性级别：private（私有）、enabled（启用）。
    securityAndComplianceAccessLevel: "private"  # 安全和合规的可见性级别：private（私有）、enabled（启用）。
    snippetsAccessLevel: "private"  # 代码片段的可见性级别：private（私有）、enabled（启用）。
    wikiAccessLevel: "private"  # Wiki 功能的可见性级别：private（私有）、enabled（启用）。

  mirroring:
    enabled: false  # 是否启用镜像功能。
    triggerBuilds: false  # 是否在镜像更新时触发构建。
    protectedBranchesOnly: false  # 是否仅镜像保护分支。
    pullMirror: false  # 是否启用拉取镜像。
    overwriteDivergedBranches: false  # 是否覆盖分叉的分支。
    userId: ""  # 镜像的用户 ID。

  container:  
    registryEnabled: true # 是否启用容器注册表功能。
    scanningEnabled: true  # 是否启用容器扫描。
    expirationPolicyAttributes:
      enabled: false  # 是否启用容器过期策略。
      cadence: "daily"  # 策略应用的频率：daily（每天）、weekly（每周）、monthly（每月）。
      keepN: 10  # 保留的容器镜像数量。
      olderThan: "30d"  # 镜像过期阈值（例如，30d 表示 30 天）。
      nameRegexDelete: ""
      nameRegexKeep: ""
      nameRegex: ""
```

``` shell
  glctl create -f project.yaml
```