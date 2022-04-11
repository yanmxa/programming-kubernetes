package pkg

import (
	"context"
	"fmt"
	"reflect"
	"time"

	coreV1 "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apisMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	coreInformer "k8s.io/client-go/informers/core/v1"
	networkingInformer "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	coreLister "k8s.io/client-go/listers/core/v1"
	networkingLister "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const workNum = 5
const maxRetry = 10

type controller struct {
	client        kubernetes.Interface
	ingressLister networkingLister.IngressLister
	serviceLister coreLister.ServiceLister
	queue         workqueue.RateLimitingInterface
}

func NewController(client kubernetes.Interface, serviceInformer coreInformer.ServiceInformer, ingressInformer networkingInformer.IngressInformer) controller {
	c := controller{
		client:        client,
		ingressLister: ingressInformer.Lister(),
		serviceLister: serviceInformer.Lister(),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ingressManager"),
	}

	// 调用informer()方法才会真正的创建informer并进行注册
	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addService,
		UpdateFunc: c.updateService,
	})

	ingressInformer.Informer().AddEventHandler((cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.deleteIngress,
	}))

	return c
}

func (c *controller) Run(stopChan chan struct{}) {
	for i := 0; i < workNum; i++ {
		go wait.Until(c.worker, time.Minute, stopChan)
	}
	<-stopChan
}

func (c *controller) worker() {
	for c.processNextItem() {

	}
}

func (c *controller) processNextItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Done(item)

	key := item.(string)
	err := c.syncService(key)
	if err != nil {
		c.handlerError(key, err)
	}
	return true
}

func (c *controller) syncService(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	// 获取service
	service, err := c.serviceLister.Services(namespace).Get(name)
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	// 新增和删除
	_, ok := service.GetAnnotations()["ingress/http"]
	ingress, err := c.ingressLister.Ingresses(namespace).Get(name)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if ok && errors.IsNotFound(err) {
		// create ingress
		ig := c.constructIngress(service)
		_, err := c.client.NetworkingV1().Ingresses(namespace).Create(context.TODO(), ig, v1.CreateOptions{})
		if err != nil {
			return err
		}
	} else if !ok && ingress != nil {
		// delete ingress
		fmt.Println("Service is not exists and delete ingress")
		err := c.client.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), name, v1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *controller) constructIngress(service *coreV1.Service) *networkingV1.Ingress {
	ingress := networkingV1.Ingress{}
	ingress.Name = service.Name
	ingress.Namespace = service.Namespace
	ingress.ObjectMeta.OwnerReferences = []apisMetaV1.OwnerReference{
		*apisMetaV1.NewControllerRef(service, coreV1.SchemeGroupVersion.WithKind("Service")),
	}
	pathType := networkingV1.PathTypePrefix
	ingressClassName := "nginx"
	ingress.Spec = networkingV1.IngressSpec{
		IngressClassName: &ingressClassName,
		Rules: []networkingV1.IngressRule{
			{
				Host: "example.com",
				IngressRuleValue: networkingV1.IngressRuleValue{
					HTTP: &networkingV1.HTTPIngressRuleValue{
						Paths: []networkingV1.HTTPIngressPath{
							{
								Path:     "/",
								PathType: &pathType,
								Backend: networkingV1.IngressBackend{
									Service: &networkingV1.IngressServiceBackend{
										Name: service.Name,
										Port: networkingV1.ServiceBackendPort{
											Number: 80,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	fmt.Println("Construct Ingress")
	return &ingress
}

func (c *controller) handlerError(key string, err error) {
	if c.queue.NumRequeues(key) <= maxRetry {
		c.queue.AddRateLimited(key)
	}
	runtime.HandleError(err)
	c.queue.Forget(key)
}

func (c *controller) addService(obj interface{}) {
	fmt.Println("Add Service")
	c.enqueue(obj)
}

func (c *controller) updateService(oldObj, newObj interface{}) {
	fmt.Println("Update Service")
	// todo 比较annotation
	if reflect.DeepEqual(oldObj, newObj) {
		return
	}
	fmt.Println("Update Service!")
	c.enqueue(newObj)
}

func (c *controller) deleteIngress(obj interface{}) {
	ingress := obj.(*networkingV1.Ingress)
	ownerReference := apisMetaV1.GetControllerOf(ingress)
	if ownerReference == nil {
		return
	}
	if ownerReference.Kind != "Service" {
		return
	}
	fmt.Println("Delete Ingress")
	c.queue.Add(ingress.Namespace + "/" + ingress.Name)
}

func (c *controller) enqueue(obj interface{}) {
	// only need to add the key
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
	}
	c.queue.Add(key)
}
