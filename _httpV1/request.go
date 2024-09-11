package _httpV1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/afret0/wheel/log"
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

	hd := make(http.Header)
	hd.Add("Content-Type", "application/json")
	opId := strings.ReplaceAll(uuid.New().String(), "-", "")
	hd.Add("opId", opId)
	if len(headers) != 0 {
		hd = headers[0]
	}
	payloadJson, err := json.Marshal(body)
	if err != nil {
		return err
	}
	payload := bytes.NewReader(payloadJson)
	req, err := http.NewRequestWithContext(ctx, "POST", url, payload)
	if err != nil {
		return err
	}
	req.Header = hd
	resp, err := new(http.Client).Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return fmt.Errorf("statusCode is %v, url: %s", resp.StatusCode, url)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	lg.Infof("resp: %s", respBody)

	err = json.Unmarshal(respBody, ret)
	if err != nil {
		return err
	}

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

	hd := make(http.Header)
	opId := strings.ReplaceAll(uuid.New().String(), "-", "")
	hd.Add("opId", opId)
	if len(headers) != 0 {
		hd = headers[0]
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	req.Header = hd

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return fmt.Errorf("statusCode is %v, url: %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	lg.Infof("resp: %s", body)

	err = json.Unmarshal(body, ret)
	if err != nil {
		return err
	}

	return nil
}
