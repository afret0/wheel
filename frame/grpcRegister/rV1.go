package grpcRegister

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	"github.com/afret0/wheel/frame/router"
)

var (
	contextType      = reflect.TypeOf((*context.Context)(nil)).Elem()
	protoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()
	errorType        = reflect.TypeOf((*error)(nil)).Elem()
)

// RegisterWithServiceDesc 使用 proto 生成的 grpc.ServiceDesc 注册 controller,
// 服务名取自 desc.ServiceName(等价于 Register(desc.ServiceName, ctrl, ...)),
// RegisterGrpcControllerToGinRouterV1 会直接使用 desc 中声明的方法列表,
// 不再依赖 proto 全局注册表查找.
func (g *GrpcRegister) RegisterWithServiceDesc(desc *grpc.ServiceDesc, ctrl GrpcController, middlewareChain ...gin.HandlerFunc) {
	g.slot[desc.ServiceName] = slot{
		GrpcController:  ctrl,
		MiddlewareChain: middlewareChain,
		ServiceDesc:     desc,
	}
}

// RegisterGrpcControllerToGinRouterV1 按 proto 中 service 定义的方法注册 gin 路由.
//
// 与 RegisterGrpcControllerToGinRouter 的区别:
//   - V0 遍历 controller 上所有导出方法, 同一个 controller 实现多个 service 时,
//     每个 service 都会被注册上全部方法, 只能靠 Option.PrefixWhiteList 裁剪;
//   - V1 只注册对应 service 在 proto 中声明的一元(unary)方法, 路由路径与 gRPC 的
//     /{package}.{Service}/{Method} 保持一致, streaming 方法会被跳过.
//
// Option.PrefixWhiteList 与 Option.MethodMiddlewareSlot 的语义保持不变.
// 若某个 service 在 proto 注册表中找不到, 或 controller 未实现其声明的方法, 返回 error.
func (g *GrpcRegister) RegisterGrpcControllerToGinRouterV1() error {
	serviceNameL := make([]string, 0, len(g.slot))
	for serviceName := range g.slot {
		serviceNameL = append(serviceNameL, serviceName)
	}
	sort.Strings(serviceNameL)

	for _, serviceName := range serviceNameL {
		s := g.slot[serviceName]

		methodNameL, err := serviceMethodNameL(serviceName, s.ServiceDesc)
		if err != nil {
			return err
		}

		if err = g.registerGrpcControllerToGinRouterV1(serviceName, methodNameL, s.GrpcController, s.MiddlewareChain...); err != nil {
			return err
		}
	}

	return nil
}

func (g *GrpcRegister) registerGrpcControllerToGinRouterV1(serviceName string, methodNameL []string, ctrl GrpcController, middlewareChain ...gin.HandlerFunc) error {
	ctrlType := reflect.TypeOf(ctrl)
	ctrlValue := reflect.ValueOf(ctrl)

	R := router.GetRouter(g.e)
	G := R.Group(fmt.Sprintf("/%s", serviceName))

outerLoop:
	for _, methodName := range methodNameL {
		fullMethodName := fmt.Sprintf("/%s/%s", serviceName, methodName)

		for _, prefix := range g.opt.PrefixWhiteList {
			if strings.Contains(fullMethodName, prefix) {
				continue outerLoop
			}
		}

		method, ok := lookupMethod(ctrlType, methodName)
		if !ok {
			return fmt.Errorf("grpcRegister: controller %s 未实现 service %s 中定义的方法 %s", ctrlType.String(), serviceName, methodName)
		}

		if err := checkMethodSignature(method); err != nil {
			return fmt.Errorf("grpcRegister: %s 方法签名不合法: %w", fullMethodName, err)
		}

		handler := g.createHTTPHandler(ctrlValue, method)

		if slot, ok := g.opt.MethodMiddlewareSlot[fullMethodName]; ok {
			G.POST(fmt.Sprintf("/%s", methodName), handler, slot.MiddlewareChain...)
			continue
		}

		G.POST(fmt.Sprintf("/%s", methodName), handler, middlewareChain...)
	}

	R.RegisterRouter(g.e)

	return nil
}

// serviceMethodNameL 取出 service 中声明的一元方法名, 优先使用 grpc.ServiceDesc,
// 否则从 proto 全局注册表按 service 全名查找
func serviceMethodNameL(serviceName string, desc *grpc.ServiceDesc) ([]string, error) {
	if desc != nil {
		methodNameL := make([]string, 0, len(desc.Methods))
		for _, m := range desc.Methods {
			methodNameL = append(methodNameL, m.MethodName)
		}
		return methodNameL, nil
	}

	d, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(serviceName))
	if err != nil {
		return nil, fmt.Errorf("grpcRegister: 在 proto 注册表中找不到 service %s: %w", serviceName, err)
	}

	sd, ok := d.(protoreflect.ServiceDescriptor)
	if !ok {
		return nil, fmt.Errorf("grpcRegister: %s 不是一个 proto service", serviceName)
	}

	methodL := sd.Methods()
	methodNameL := make([]string, 0, methodL.Len())
	for i := 0; i < methodL.Len(); i++ {
		m := methodL.Get(i)
		if m.IsStreamingClient() || m.IsStreamingServer() {
			continue
		}
		methodNameL = append(methodNameL, string(m.Name()))
	}

	return methodNameL, nil
}

// lookupMethod 先按 proto 中的方法名查找, 找不到时按 protoc-gen-go 的命名规则转换后再查找
func lookupMethod(ctrlType reflect.Type, methodName string) (reflect.Method, bool) {
	if method, ok := ctrlType.MethodByName(methodName); ok {
		return method, true
	}

	return ctrlType.MethodByName(goCamelCase(methodName))
}

func checkMethodSignature(method reflect.Method) error {
	methodType := method.Type
	// controller 实例 + context + request
	if methodType.NumIn() != 3 || !methodType.In(1).Implements(contextType) ||
		methodType.In(2).Kind() != reflect.Ptr || !methodType.In(2).Implements(protoMessageType) {
		return fmt.Errorf("期望 func(context.Context, proto.Message) (proto.Message, error), 实际 %s", methodType.String())
	}

	if methodType.NumOut() != 2 || !methodType.Out(0).Implements(protoMessageType) || !methodType.Out(1).Implements(errorType) {
		return fmt.Errorf("期望 func(context.Context, proto.Message) (proto.Message, error), 实际 %s", methodType.String())
	}

	return nil
}

// goCamelCase 将 proto 中定义的方法名转换为 protoc-gen-go 生成的 Go 方法名
func goCamelCase(s string) string {
	b := new(strings.Builder)
	b.Grow(len(s))

	upper := true
	for _, r := range s {
		switch {
		case r == '_' || r == '.':
			upper = true
		case upper:
			b.WriteRune(unicode.ToUpper(r))
			upper = false
		default:
			b.WriteRune(r)
		}
	}

	return b.String()
}
