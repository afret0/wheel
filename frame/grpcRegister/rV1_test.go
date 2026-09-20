package grpcRegister

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type testCtrl struct{}

func (t *testCtrl) Get(_ context.Context, req *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	return wrapperspb.String("get:" + req.GetValue()), nil
}

func (t *testCtrl) Add(_ context.Context, req *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	return wrapperspb.String("add:" + req.GetValue()), nil
}

// InternalHelper 不在 proto 中声明, V1 不应该注册它
func (t *testCtrl) InternalHelper(_ context.Context, req *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	return req, nil
}

func routeSet(e *gin.Engine) map[string]bool {
	m := make(map[string]bool)
	for _, r := range e.Routes() {
		m[r.Method+" "+r.Path] = true
	}
	return m
}

func newEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestRegisterGrpcControllerToGinRouterV1WithServiceDesc(t *testing.T) {
	e := newEngine()

	desc := &grpc.ServiceDesc{
		ServiceName: "descTest.DescSvc",
		Methods: []grpc.MethodDesc{
			{MethodName: "Get"},
			{MethodName: "Add"},
		},
	}

	g := NewGrpcRegister(e, &Option{PrefixWhiteList: []string{"/descTest.DescSvc/Add"}})
	g.RegisterWithServiceDesc(desc, new(testCtrl))

	if err := g.RegisterGrpcControllerToGinRouterV1(); err != nil {
		t.Fatalf("RegisterGrpcControllerToGinRouterV1 err: %v", err)
	}

	routes := routeSet(e)
	if !routes["POST /descTest.DescSvc/Get"] {
		t.Fatalf("期望注册 /descTest.DescSvc/Get, 实际路由: %v", routes)
	}
	if routes["POST /descTest.DescSvc/Add"] {
		t.Fatalf("白名单方法不应该被注册, 实际路由: %v", routes)
	}
	if routes["POST /descTest.DescSvc/InternalHelper"] {
		t.Fatalf("proto 未声明的方法不应该被注册, 实际路由: %v", routes)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/descTest.DescSvc/Get", strings.NewReader(`{"value":"a"}`))
	req.Header.Set("Content-Type", "application/json")
	e.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望 200, 实际 %d, body: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "get:a") {
		t.Fatalf("响应不符合预期: %s", w.Body.String())
	}
}

func TestRegisterGrpcControllerToGinRouterV1FromProtoRegistry(t *testing.T) {
	registerTestProtoService(t)

	e := newEngine()

	g := NewGrpcRegister(e)
	g.Register("wheelTest.RegistrySvc", new(testCtrl))

	if err := g.RegisterGrpcControllerToGinRouterV1(); err != nil {
		t.Fatalf("RegisterGrpcControllerToGinRouterV1 err: %v", err)
	}

	routes := routeSet(e)
	if !routes["POST /wheelTest.RegistrySvc/Get"] || !routes["POST /wheelTest.RegistrySvc/Add"] {
		t.Fatalf("期望注册 proto 中声明的 Get/Add, 实际路由: %v", routes)
	}
	if routes["POST /wheelTest.RegistrySvc/InternalHelper"] {
		t.Fatalf("proto 未声明的方法不应该被注册, 实际路由: %v", routes)
	}
	if routes["POST /wheelTest.RegistrySvc/Watch"] {
		t.Fatalf("streaming 方法不应该被注册, 实际路由: %v", routes)
	}
}

func TestRegisterGrpcControllerToGinRouterV1Error(t *testing.T) {
	e := newEngine()

	g := NewGrpcRegister(e)
	g.RegisterWithServiceDesc(&grpc.ServiceDesc{
		ServiceName: "missTest.MissSvc",
		Methods:     []grpc.MethodDesc{{MethodName: "NotImplemented"}},
	}, new(testCtrl))

	if err := g.RegisterGrpcControllerToGinRouterV1(); err == nil {
		t.Fatal("controller 未实现 proto 声明的方法时期望返回 error")
	}

	g = NewGrpcRegister(e)
	g.Register("notExist.NotExistSvc", new(testCtrl))

	if err := g.RegisterGrpcControllerToGinRouterV1(); err == nil {
		t.Fatal("proto 注册表中找不到 service 时期望返回 error")
	}
}

// TestRegisterGrpcControllerToGinRouterV0 保证旧逻辑不变: 注册 controller 上的全部导出方法
func TestRegisterGrpcControllerToGinRouterV0(t *testing.T) {
	e := newEngine()

	g := NewGrpcRegister(e, &Option{PrefixWhiteList: []string{"/v0Test.V0Svc/Add"}})
	g.Register("v0Test.V0Svc", new(testCtrl))
	g.RegisterGrpcControllerToGinRouter()

	routes := routeSet(e)
	if !routes["POST /v0Test.V0Svc/Get"] || !routes["POST /v0Test.V0Svc/InternalHelper"] {
		t.Fatalf("旧逻辑应该注册全部导出方法, 实际路由: %v", routes)
	}
	if routes["POST /v0Test.V0Svc/Add"] {
		t.Fatalf("白名单方法不应该被注册, 实际路由: %v", routes)
	}
}

func TestGoCamelCase(t *testing.T) {
	caseM := map[string]string{
		"Get":              "Get",
		"get_multi_by_uid": "GetMultiByUid",
		"getRecordList":    "GetRecordList",
	}

	for in, want := range caseM {
		if got := goCamelCase(in); got != want {
			t.Fatalf("goCamelCase(%q) = %q, want %q", in, got, want)
		}
	}
}

func registerTestProtoService(t *testing.T) {
	t.Helper()

	strValue := ".google.protobuf.StringValue"
	fdp := &descriptorpb.FileDescriptorProto{
		Name:       proto.String("wheel/frame/grpcRegister/registry_test.proto"),
		Package:    proto.String("wheelTest"),
		Syntax:     proto.String("proto3"),
		Dependency: []string{"google/protobuf/wrappers.proto"},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: proto.String("RegistrySvc"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: proto.String("Get"), InputType: &strValue, OutputType: &strValue},
					{Name: proto.String("Add"), InputType: &strValue, OutputType: &strValue},
					{
						Name:            proto.String("Watch"),
						InputType:       &strValue,
						OutputType:      &strValue,
						ServerStreaming: proto.Bool(true),
					},
				},
			},
		},
	}

	fd, err := protodesc.NewFile(fdp, protoregistry.GlobalFiles)
	if err != nil {
		t.Fatalf("protodesc.NewFile err: %v", err)
	}

	if err = protoregistry.GlobalFiles.RegisterFile(fd); err != nil {
		t.Fatalf("RegisterFile err: %v", err)
	}
}
