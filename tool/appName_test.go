package tool

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestAppNameFromHostId(t *testing.T) {
	cases := []struct {
		name   string
		hostId string
		want   string
	}{
		{
			name:   "pod 名去掉最后两段",
			hostId: "order-6d8f7b9c4d-x2abc",
			want:   "order",
		},
		{
			name:   "服务名本身带连字符时只去掉最后两段",
			hostId: "asset-v1-svr-6d8f7b9c4d-x2abc",
			want:   "asset-v1-svr",
		},
		{
			name:   "刚好三段",
			hostId: "order-rs-hash",
			want:   "order",
		},
		{
			name:   "不足三段原样返回",
			hostId: "afretos-mbp",
			want:   "afretos-mbp",
		},
		{
			name:   "空值原样返回",
			hostId: "",
			want:   "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := appNameFromHostId(tc.hostId); got != tc.want {
				t.Fatalf("appNameFromHostId(%q) = %q, want %q", tc.hostId, got, tc.want)
			}
		})
	}
}

func TestIsGeneratedUuid(t *testing.T) {
	cases := []struct {
		name   string
		hostId string
		want   bool
	}{
		{
			name:   "HostId 兜底生成的 32 位 uuid",
			hostId: UUIDWithoutHyphen(),
			want:   true,
		},
		{
			name:   "pod 名",
			hostId: "order-6d8f7b9c4d-x2abc",
			want:   false,
		},
		{
			name:   "长度不足 32",
			hostId: strings.Repeat("a", 31),
			want:   false,
		},
		{
			name:   "长度超过 32",
			hostId: strings.Repeat("a", 33),
			want:   false,
		},
		{
			name:   "含非十六进制字符",
			hostId: "z" + strings.Repeat("a", 31),
			want:   false,
		},
		{
			name:   "大写十六进制",
			hostId: strings.ToUpper(UUIDWithoutHyphen()),
			want:   false,
		},
		{
			name:   "空值",
			hostId: "",
			want:   false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isGeneratedUuid(tc.hostId); got != tc.want {
				t.Fatalf("isGeneratedUuid(%q) = %v, want %v", tc.hostId, got, tc.want)
			}
		})
	}
}

func TestAppNameFromEnvHostname(t *testing.T) {
	t.Setenv("HOSTNAME", "asset-v1-svr-6d8f7b9c4d-x2abc")

	if got := HostId(); got != "asset-v1-svr-6d8f7b9c4d-x2abc" {
		t.Fatalf("HostId() = %q, want 环境变量值", got)
	}
	if got := AppName(); got != "asset-v1-svr" {
		t.Fatalf("AppName() = %q, want asset-v1-svr", got)
	}
}

func TestAppNameIgnoresGeneratedUuidHostId(t *testing.T) {
	t.Setenv("HOSTNAME", "")

	hostId := HostId()
	if !isGeneratedUuid(hostId) {
		t.Skipf("HostId() = %q 不是兜底 uuid, 本用例不适用", hostId)
	}

	app := AppName()
	if app == hostId || strings.Contains(app, hostId) {
		t.Fatalf("AppName() = %q 用了兜底 uuid %q, caller 会变成随机值", app, hostId)
	}

	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		if app != "" {
			t.Fatalf("AppName() = %q, os.Hostname() 也取不到时期望为空", app)
		}

		return
	}

	if want := appNameFromHostId(hostname); app != want {
		t.Fatalf("AppName() = %q, want os.Hostname() 切分结果 %q", app, want)
	}
}

func TestAppNameStableWhenHostIdFallsBackToUuid(t *testing.T) {
	t.Setenv("HOSTNAME", "")

	first, second := AppName(), AppName()
	if first != second {
		t.Fatalf("AppName() 每次不同: %q vs %q, caller 会在日志里漂移", first, second)
	}
	if HostId() == HostId() {
		t.Log("HostId() 兜底值每次不同, 属原有行为, AppName() 不受影响")
	}
}

func TestInjectCallerMD(t *testing.T) {
	t.Setenv("HOSTNAME", "order-6d8f7b9c4d-x2abc")

	md := metadata.Pairs(MetaCaller, "upstream-service")
	InjectCallerMD(md)

	if got := md.Get(MetaCaller); len(got) != 1 || got[0] != "order" {
		t.Fatalf("caller = %v, want [order]", got)
	}
}

func TestInjectCallerHeader(t *testing.T) {
	t.Setenv("HOSTNAME", "order-6d8f7b9c4d-x2abc")

	t.Run("覆盖上游残留的 caller", func(t *testing.T) {
		hd := make(http.Header)
		hd.Set(MetaCaller, "upstream-service")

		InjectCallerHeader(hd)

		if got := hd.Values(MetaCaller); len(got) != 1 || got[0] != "order" {
			t.Fatalf("caller header = %v, want [order]", got)
		}
	})

	t.Run("重复注入只保留一个值", func(t *testing.T) {
		hd := make(http.Header)
		InjectCallerHeader(hd)
		InjectCallerHeader(hd)

		if got := hd.Values(MetaCaller); len(got) != 1 || got[0] != "order" {
			t.Fatalf("caller header = %v, want [order]", got)
		}
	})

	t.Run("header 为 nil 时不 panic", func(t *testing.T) {
		InjectCallerHeader(nil)
	})
}
