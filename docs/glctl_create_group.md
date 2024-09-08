## glctl create group
``` yaml
 apiVersion: v1
 kind: Group
 metadata: 
   name: Group1
   parentName: Group
   annotations:
    description: " my first Group"  # 组的简短描述。
    avatar: ""  # 项目头像的 URL。   
 spec:
   visibility: public
   branch:
     defaultName: main
     protection: 
       allowedToPush:
       - developer
       - maintainer
       allowForcePush: false
       allowedToMerge:
       - developer
       - maintainer
       developerCanInitialPush: true
   enabled:
   - autoDevops
   - gitAccessProtocol
   - lfs
   - requestAccess
   disabled:
   - emails
   - mentions
   - shareWithGroupLock
   - requireTwoFactorAuthentication
   - membershipLock
   projectCreationLevel: maintainer
   subgroupCreationLevel: maintainer
   wikiAccessLevel: disabled
   limit:
     extraSharedRunnersMinutes: 60
     sharedRunnersMinutes: 50
```

``` shell
  glctl create -f group.yaml
```