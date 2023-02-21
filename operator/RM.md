# Kubernetes Operator
    本文为Kubernetes Operator简单demo。业务相关逻辑为：PodsManage作为pod的管理，负责将PodsManage转换为对应的
    kubernetes内置资源。
    operator-sdk api参考：https://sdk.operatorframework.io/docs/building-operators/golang/references/client/
## 1、创建项目
    go mod init pods-manager
## 2、使用kubebuilder生成代码框架
    创建项目:
        kubebuilder init --domain my.domain 
    创建一个API:
        kubebuilder create api --group podsmanager --version v1 --kind PodsManage
    生成配置文件:
        make manifests
    安装CRD:
        make install #需要连接外网，可能会失败
    运行控制器:
    make run
    安装CR:
        kubectl apply -f config/samples/
    卸载CRD:
        make uninstall
    卸载控制器:
        make undeploy