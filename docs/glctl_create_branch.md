## glctl create branch

``` yaml
apiVersion: v1
kind: Branch
metadata:
  name: "develop"  # 分支名称。
ref:
  branch: main
subjects: 
  group: Group1
  name: Project1
```

``` shell
  glctl create -f project-branch.yaml
```