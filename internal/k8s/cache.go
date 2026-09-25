package k8s

import (
	"context"
	"fmt"
	"sync"

	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	virtv1 "kubevirt.io/api/core/v1"
)

func errUnlessNotFound(err error) error {
	if apierror.IsNotFound(err) {
		return nil
	}
	return err
}

type objectCache struct {
	informer cache.SharedIndexInformer
	cancel   context.CancelFunc
	done     chan struct{}
	mu       sync.Mutex
	changed  chan struct{}
}

var cacheMu sync.Mutex
var objectCaches = make(map[string]*objectCache)
var cacheClient any

func (c *objectCache) notify() {
	c.mu.Lock()
	defer c.mu.Unlock()
	close(c.changed)
	c.changed = make(chan struct{})
}

func (c *objectCache) changes() <-chan struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.changed
}

func Stop() {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	stopCaches()
}

func stopCaches() {
	for _, c := range objectCaches {
		c.cancel()
	}
	for _, c := range objectCaches {
		<-c.done
	}
	objectCaches = make(map[string]*objectCache)
}

// One informer per resource and namespace, shared across all callers. A client
// replacement (in-process restart) discards caches belonging to the old client.
func cachedObjects(ctx context.Context, kind string) (*objectCache, error) {
	cacheMu.Lock()
	if cacheClient != kubeClient {
		stopCaches()
		cacheClient = kubeClient
	}
	key := globalNamespace + "/" + kind
	c := objectCaches[key]
	if c == nil {
		var object runtime.Object
		var list func(context.Context, metav1.ListOptions) (runtime.Object, error)
		var watchFn func(context.Context, metav1.ListOptions) (watch.Interface, error)
		switch kind {
		case "pods":
			client := kubeClient.CoreV1().Pods(globalNamespace)
			object = &corev1.Pod{}
			list = func(ctx context.Context, o metav1.ListOptions) (runtime.Object, error) { return client.List(ctx, o) }
			watchFn = client.Watch
		case "vms":
			client := virtClient.KubevirtV1().VirtualMachines(globalNamespace)
			object = &virtv1.VirtualMachine{}
			list = func(ctx context.Context, o metav1.ListOptions) (runtime.Object, error) { return client.List(ctx, o) }
			watchFn = client.Watch
		default:
			cacheMu.Unlock()
			return nil, fmt.Errorf("unknown informer resource %s", kind)
		}
		lifecycle, cancel := context.WithCancel(context.Background())
		c = &objectCache{cancel: cancel, done: make(chan struct{}), changed: make(chan struct{})}
		c.informer = cache.NewSharedIndexInformer(&cache.ListWatch{ListWithContextFunc: list, WatchFuncWithContext: watchFn}, object, 0, cache.Indexers{})
		_, _ = c.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(any) { c.notify() }, UpdateFunc: func(any, any) { c.notify() }, DeleteFunc: func(any) { c.notify() },
		})
		objectCaches[key] = c
		go func() { defer close(c.done); c.informer.RunWithContext(lifecycle) }()
	}
	cacheMu.Unlock()
	if !cache.WaitForCacheSync(ctx.Done(), c.informer.HasSynced) {
		return nil, fmt.Errorf("sync %s informer: %w", kind, ctx.Err())
	}
	return c, nil
}

func cachedPod(ctx context.Context, name string) (*corev1.Pod, error) {
	c, err := cachedObjects(ctx, "pods")
	if err != nil {
		return nil, err
	}
	obj, exists, err := c.informer.GetStore().GetByKey(globalNamespace + "/" + name)
	if err != nil || !exists {
		return nil, err
	}
	return obj.(*corev1.Pod).DeepCopy(), nil
}

func cachedVM(ctx context.Context, name string) (*virtv1.VirtualMachine, error) {
	c, err := cachedObjects(ctx, "vms")
	if err != nil {
		return nil, err
	}
	obj, exists, err := c.informer.GetStore().GetByKey(globalNamespace + "/" + name)
	if err != nil || !exists {
		return nil, err
	}
	return obj.(*virtv1.VirtualMachine).DeepCopy(), nil
}
