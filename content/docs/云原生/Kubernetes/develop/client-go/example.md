---
title: "example"
weight: 3
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
bookCollapseSection: false
---
[example code](https://github.com/bsonger/my-controller)
1. 创建CRD 结构

/pkg/apis/example/v1/types.go
```go
type MyResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec MyResourceSpec `json:"spec"`
}

// MyResourceSpec defines the desired state of MyResource
type MyResourceSpec struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}
```
2. 使用code-generator 工具自动生成clientSet、informer、lister
```shell
bash hack/update-codegen.sh
```
3. 定义controller 结构

/pkg/controller/contronller.go
```go
type Controller struct {
	kubeclientset       kubernetes.Interface
	myresourceClientset clientset.Interface
	myresourceLister listers.MyResourceLister
	myresourceSynced cache.InformerSynced
	workqueue workqueue.TypedRateLimitingInterface[string]
}
```
4. 初始化controller

/pkg/controller/contronller.go
```go
func NewController(
	kubeclientset kubernetes.Interface,
	myresourceClientset clientset.Interface,
	myresourceInformer v1.MyResourceInformer) *Controller {

	controller := &Controller{
		kubeclientset:       kubeclientset,
		myresourceClientset: myresourceClientset,
		myresourceLister:    myresourceInformer.Lister(),
		myresourceSynced:    myresourceInformer.Informer().HasSynced,
		workqueue: workqueue.NewTypedRateLimitingQueueWithConfig(
			workqueue.DefaultTypedControllerRateLimiter[string](),
			workqueue.TypedRateLimitingQueueConfig[string]{
				Name: "my-controller",
			},
		),
	}

	myresourceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueMyResource,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueMyResource(new)
		},
		DeleteFunc: controller.enqueueMyResource,
	})
	return controller
}
```
5. set up /metrics and /healthz接口

/pkg/controller/contronller.go

```go
// Set up metrics server
http.Handle("/metrics", promhttp.Handler())
go func() {
    klog.Fatal(http.ListenAndServe(":8080", nil))
}()
```
```go
// Set up health check server
http.HandleFunc("/healthz", health.HealthCheck)
go func() {
    klog.Fatal(http.ListenAndServe(":8081", nil))
}()
```
6. 使用leaderelection 编写高可用controller

/cmd/my-controller.go
```go
lock := &resourcelock.LeaseLock{
    LeaseMeta: metav1.ObjectMeta{
        Name:      leaseLockName,
        Namespace: leaseLockNamespace,
    },
    Client: kubeClient.CoordinationV1(),
    LockConfig: resourcelock.ResourceLockConfig{
        Identity: os.Getenv("POD_NAME"), // Use the pod name as the identity
    },
}

```
```go
leaderelection.RunOrDie(context.TODO(), leaderelection.LeaderElectionConfig{
    Lock:            lock,
    ReleaseOnCancel: true,
    LeaseDuration:   15 * time.Second,
    RenewDeadline:   10 * time.Second,
    RetryPeriod:     2 * time.Second,
    Callbacks: leaderelection.LeaderCallbacks{
        OnStartedLeading: func(ctx context.Context) {
            klog.Info("Started leading")
            informerFactory.Start(ctx.Done()) // Start the controller logic
        },
        OnStoppedLeading: func() {
            klog.Info("Stopped leading")
            le <- struct{}{}
        },
        OnNewLeader: func(identity string) {
            if identity == os.Getenv("POD_NAME") {
                klog.Info("I am the leader")
            } else {
                klog.Infof("New leader elected: %s", identity)
            }
        },
    },
})

```

7. 启动controller

/cmd/my-controller.go
```go
myController.Run(2, ctx.Done())
```

业务逻辑写在 pkg/controller/contronller.go

```go
func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	myresource, err := c.myresourceLister.MyResources(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("myresource '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}

	fmt.Printf("Sync/Add/Update for MyResource %s\n", myresource.GetName())
	return nil
}
```