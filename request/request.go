package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/afret0/wheel/log"
	"github.com/afret0/wheel/tool"
)

//type RetryOpt struct {
//	retryCount int
//	retryDelay int
//}

//func PostWithRetry(ctx context.Context, ret interface{}, url string, body interface{}, opt *RetryOpt, headers ...http.Header) error {
//	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(opt.retryDelay*opt.retryCount)*time.Second)
//	defer cancel()
//
//	var err error
//
//	for {
//		select {
//		case <-timeoutCtx.Done():
//			return fmt.Errorf("timeout while trying to obtain the lock, err: %+v", err)
//		default:
//
//			err = Post(ctx, ret, url, body, headers...)
//			//lock, err := l.Locker.Obtain(ctx, key, time.Duration(ttl)*time.Second, nil)
//			if err != nil {
//				if errors.Is(err, redislock.ErrNotObtained) {
//					// If the lock is not obtained, wait for a while before trying again
//					time.Sleep(time.Duration(retryDelay) * time.Second)
//					continue
//				}
//				// If there is another error, return it
//				return nil, err
//			}
//
//			// If the lock is obtained, return it
//			return lock, nil
//		}
//	}
//}

func Post(ctx context.Context, ret interface{}, url string, body interface{}, headers ...http.Header) error {
	lg := log.CtxLogger(ctx).WithField("url", url)

	ctx, span := startClientSpan(ctx, "POST", url)

	opId := tool.ConvertOpId(tool.OpId(ctx))
	hd := make(http.Header)
	hd.Add("Content-Type", "application/json")
	hd.Add("opId", opId)

	// 必须在 client span 创建之后注入, 否则透传的是上层 span 而非本次调用
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(hd))

	mergeHeader(hd, headers)

	payloadJson, err := json.Marshal(body)
	if err != nil {
		endClientSpan(span, 0, err)
		return err
	}
	payload := bytes.NewReader(payloadJson)
	req, err := http.NewRequestWithContext(ctx, "POST", url, payload)
	if err != nil {
		endClientSpan(span, 0, err)
		return err
	}
	req.Header = hd
	resp, err := new(http.Client).Do(req)
	if err != nil {
		endClientSpan(span, 0, err)
		return err
	}
	defer func() {
		err = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("statusCode is %v, url: %s", resp.StatusCode, url)
		endClientSpan(span, resp.StatusCode, err)
		return err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		endClientSpan(span, resp.StatusCode, err)
		return err
	}

	lg.Infof("resp: %s", respBody)

	err = json.Unmarshal(respBody, ret)
	if err != nil {
		endClientSpan(span, resp.StatusCode, err)
		return err
	}

	endClientSpan(span, resp.StatusCode, nil)

	return nil
}

func MarshallUrlParams(url string, params map[string]string) string {
	l := make([]string, 0)
	for k, v := range params {
		s := fmt.Sprintf("%s=%s", k, v)
		l = append(l, s)
	}
	s := strings.Join(l, "&")
	return fmt.Sprintf("%s?%s", url, s)
}

func Get(ctx context.Context, ret interface{}, url string, headers ...http.Header) error {
	lg := log.CtxLogger(ctx).WithField("url", url)

	ctx, span := startClientSpan(ctx, "GET", url)

	opId := tool.ConvertOpId(tool.OpId(ctx))
	hd := make(http.Header)
	hd.Add("opId", opId)

	// 必须在 client span 创建之后注入, 否则透传的是上层 span 而非本次调用
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(hd))

	mergeHeader(hd, headers)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		endClientSpan(span, 0, err)
		return err
	}
	req.Header = hd

	resp, err := new(http.Client).Do(req)
	if err != nil {
		endClientSpan(span, 0, err)
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("statusCode is %v, url: %s", resp.StatusCode, url)
		endClientSpan(span, resp.StatusCode, err)
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		endClientSpan(span, resp.StatusCode, err)
		return err
	}

	lg.Infof("resp: %s", body)

	err = json.Unmarshal(body, ret)
	if err != nil {
		endClientSpan(span, resp.StatusCode, err)
		return err
	}

	endClientSpan(span, resp.StatusCode, nil)

	return nil
}
