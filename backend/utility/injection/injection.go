package injection

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/samber/do" // 使用第三方的依赖注入工具
)

var (
	defaultInjector     *do.Injector
	defaultInjectorOnce sync.Once
)

// MustInvoke invokes the function with the default injector and panics if any error occurs.
func MustInvoke[T any]() T {
	return do.MustInvoke[T](defaultInjector)
}

// Invoke invokes the function with the default injector.
func Invoke[T any]() (T, error) {
	return do.Invoke[T](defaultInjector)
}

// SetupDefaultInjector initializes the default injector with the given context.
func SetupDefaultInjector(ctx context.Context) *do.Injector {
	// 确保初始化只执行一次
	defaultInjectorOnce.Do(func() {
		injector := do.NewWithOpts(&do.InjectorOpts{})

		// TODO: 注入服务对象
		injectSnowflake(ctx, injector)
		injectRedis(ctx, injector) // 注入redis client

		defaultInjector = injector
	})
	return defaultInjector
}

// ShutdownDefaultInjector shuts down the default injector.
func ShutdownDefaultInjector() {
	if defaultInjector != nil {
		if err := defaultInjector.Shutdown(); err != nil {
			g.Log().Debugf(context.Background(), "ShutdownDefaultInjector: %+v", err)
		}
		defaultInjector = nil
	}
}

// ShutdownHelper is a helper struct for shutdown.
type ShutdownHelper[T any] struct {
	name       string                // 注入的服务对象的名称
	service    T                     // 注入的服务对象
	onShutdown func(service T) error // 关闭服务时的回调函数
}

// NewShutdownHelper creates a new ShutdownHelper.
func NewShutdownHelper[T any](service T, onShutdown func(service T) error) ShutdownHelper[T] {
	return ShutdownHelper[T]{
		service:    service,
		onShutdown: onShutdown,
	}
}

// NewShutdownHelperNamed creates a new ShutdownHelper with a name.
func NewShutdownHelperNamed[T any](service T, name string, onShutdown func(service T) error) ShutdownHelper[T] {
	return ShutdownHelper[T]{
		name:       name,
		service:    service,
		onShutdown: onShutdown,
	}
}

// Shutdown shuts down the service.
func (h ShutdownHelper[T]) Shutdown() error {
	g.Log().Debugf(
		context.Background(),
		"ShutdownHelper Shutdown: %s, %s",
		reflect.TypeOf(h.service), h.name,
	)
	return h.onShutdown(h.service)
}

// SetupShutdownHelper sets up a shutdown helper.
func SetupShutdownHelper[T any](injector *do.Injector, service T, onShutdown func(service T) error) {
	do.Provide(injector, func(i *do.Injector) (ShutdownHelper[T], error) {
		g.Log().Debugf(context.Background(), "NewShutdownHelper: %s", reflect.TypeOf(service))
		return NewShutdownHelper(service, onShutdown), nil
	})
	do.MustInvoke[ShutdownHelper[T]](injector)
}

// SetupShutdownHelperNamed sets up a shutdown helper with a name.
func SetupShutdownHelperNamed[T any](injector *do.Injector, service T, name string, onShutdown func(service T) error) {
	name = fmt.Sprintf("ShutdownHelper:%s", name)
	do.ProvideNamed(injector, name, func(i *do.Injector) (ShutdownHelper[T], error) {
		g.Log().Debugf(
			context.Background(),
			"NewShutdownHelper: %s, %s",
			reflect.TypeOf(service), name,
		)
		return NewShutdownHelperNamed(service, name, onShutdown), nil
	})
	do.MustInvokeNamed[ShutdownHelper[T]](injector, name)
}
