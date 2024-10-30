package grpcRegister

import (
	"context"
	"fmt"
	"github.com/afret0/wheel/frame/router"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"reflect"
	"strings"
)

// GrpcController 接口用于标识实现了 gRPC 服务的 Controller
type GrpcController interface {
	// 可以添加一些通用方法
}

type Option struct {
	PrefixWhiteList []string
}

type Opt = Option

type slot struct {
	GrpcController  GrpcController
	MiddlewareChain []gin.HandlerFunc
}

type GrpcRegister struct {
	e    *gin.Engine
	slot map[string]slot
	opt  *Option
}

func NewGrpcRegister(e *gin.Engine, opts ...*Option) *GrpcRegister {
	opt := new(Option)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}

	return &GrpcRegister{
		e:    e,
		slot: make(map[string]slot),
		opt:  opt,
	}
}

func (g *GrpcRegister) Register(serviceName string, ctrl GrpcController, middlewareChain ...gin.HandlerFunc) {
	//g.slot[serviceName] = ctrl
	g.slot[serviceName] = slot{
		GrpcController:  ctrl,
		MiddlewareChain: middlewareChain,
	}
}

func (g *GrpcRegister) RegisterGrpcControllerToGinRouter() {
	for serviceName, slot := range g.slot {
		g.registerGrpcControllerToGinRouter(serviceName, slot.GrpcController, slot.MiddlewareChain...)
	}
}

func (g *GrpcRegister) registerGrpcControllerToGinRouter(serviceName string, ctrl GrpcController, middlewareChain ...gin.HandlerFunc) {
	// 获取 controller 的反射类型
	ctrlType := reflect.TypeOf(ctrl)
	ctrlValue := reflect.ValueOf(ctrl)

	//R := router.GetRouter(g.e)
	R := router.GetRouter(g.e).Group(fmt.Sprintf("/%s", serviceName))

outerLoop:
	// 遍历所有方法
	for i := 0; i < ctrlType.NumMethod(); i++ {
		method := ctrlType.Method(i)

		// 跳过一些内部方法
		if strings.HasPrefix(method.Name, "mustEmbedUnimplemented") {
			continue
		}

		for _, prefix := range g.opt.PrefixWhiteList {
			if strings.Contains(method.Name, prefix) {
				continue outerLoop
			}
		}

		// 创建对应的 HTTP 处理函数
		handler := g.createHTTPHandler(ctrlValue, method)

		// 注册路由，使用方法名作为路径
		R.POST(fmt.Sprintf("/%s", method.Name), handler, middlewareChain...)
	}
}

// createHTTPHandler 创建处理 HTTP 请求的处理函数
func (g *GrpcRegister) createHTTPHandler(ctrl reflect.Value, method reflect.Method) router.HandleFuncWrap {
	return func(c *gin.Context) (interface{}, error) {
		// 获取方法的参数类型
		methodType := method.Type
		if methodType.NumIn() != 3 { // controller 实例 + context + request
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid method signature"})
			//return
			return nil, fmt.Errorf("invalid method signature")
		}

		// 创建请求参数实例
		reqType := methodType.In(2)
		reqValue := reflect.New(reqType.Elem())
		req := reqValue.Interface().(proto.Message)

		// 从 HTTP 请求体解析参数
		if err := c.ShouldBindJSON(req); err != nil {
			//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			//return
			return nil, err
		}

		// 调用 controller 方法
		results := method.Func.Call([]reflect.Value{
			ctrl,
			reflect.ValueOf(context.Background()),
			reqValue,
		})

		// 检查调用结果
		if len(results) != 2 {
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid method result"})
			//return
			return nil, fmt.Errorf("invalid method result, res: %#v", results)
		}

		// 处理错误
		if !results[1].IsNil() {
			err := results[1].Interface().(error)
			//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			//return
			return nil, err
		}

		// 转换响应为 JSON
		resp := results[0].Interface().(proto.Message)
		//c.JSON(http.StatusOK, resp)
		return resp, nil
	}
}
