package request

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afret0/wheel/tool"
)

func TestPostSendsCallerHeader(t *testing.T) {
	t.Setenv("HOSTNAME", "order-6d8f7b9c4d-x2abc")

	var gotCaller string
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotCaller = r.Header.Get(tool.MetaCaller)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":1}`))
	}))
	defer svr.Close()

	var rpl BaseRpl[int]
	if err := Post(context.Background(), &rpl, svr.URL, map[string]string{"uid": "10086"}); err != nil {
		t.Fatalf("Post: %v", err)
	}

	if gotCaller != "order" {
		t.Fatalf("caller header = %q, want order", gotCaller)
	}
	if rpl.Code != 1 {
		t.Fatalf("code = %d, want 1", rpl.Code)
	}
}

func TestGetSendsCallerHeader(t *testing.T) {
	t.Setenv("HOSTNAME", "order-6d8f7b9c4d-x2abc")

	var gotCaller string
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotCaller = r.Header.Get(tool.MetaCaller)
		_, _ = w.Write([]byte(`{"code":1}`))
	}))
	defer svr.Close()

	var rpl BaseRpl[int]
	if err := Get(context.Background(), &rpl, svr.URL); err != nil {
		t.Fatalf("Get: %v", err)
	}

	if gotCaller != "order" {
		t.Fatalf("caller header = %q, want order", gotCaller)
	}
}

func TestPostCallerHeaderCannotBeOverridden(t *testing.T) {
	t.Setenv("HOSTNAME", "order-6d8f7b9c4d-x2abc")

	var gotCaller, gotOpId []string
	var gotToken string
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotCaller = r.Header.Values(tool.MetaCaller)
		gotOpId = r.Header.Values("opId")
		gotToken = r.Header.Get("token")
		_, _ = w.Write([]byte(`{"code":1}`))
	}))
	defer svr.Close()

	hd := make(http.Header)
	hd.Set(tool.MetaCaller, "custom-caller")
	hd.Set("opId", "custom-opid")
	hd.Set("token", "abc")

	var rpl BaseRpl[int]
	if err := Post(context.Background(), &rpl, svr.URL, nil, hd); err != nil {
		t.Fatalf("Post: %v", err)
	}

	if len(gotCaller) != 1 || gotCaller[0] != "order" {
		t.Fatalf("caller header = %v, want [order]", gotCaller)
	}
	if len(gotOpId) != 1 || gotOpId[0] == "custom-opid" || gotOpId[0] == "" {
		t.Fatalf("opId header = %v, 期望由链路自身生成的单值", gotOpId)
	}
	if gotToken != "abc" {
		t.Fatalf("token header = %q, want abc", gotToken)
	}
}
