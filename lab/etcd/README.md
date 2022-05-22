# 本项目etcd结构

etcd:
```shell
/registry
  - /nods
    - /nodeName - node的信息,数据结构还需要进一步调整
  - /pods
    - /podName - pod的信息,pod.Pod的json格式
/scheduler
  - /pods
    - /podName - 交给scheduler来进行调度的pod
/controller
  - /pods     - 关于pod操作
    - /create - 创建pod的操作
    - /stop   - 停止pod的运行的操作
    - /delete - 删除pod的操作
    - /update - 更新pod的信息的操作
  - /nodes    - 关于nodes的操作
    - /create - 创建node的操作
    - /delete - 删除node的操作,也就是断开连接
    - /update - 更新node的信息的操作
/relations
  - /pods
    - /podName      - 记录了该pod对应的node
  - /nods
    - /nodeName
      - /podName    - 记录了该node下所有的pod
```
