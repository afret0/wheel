package tool

import "encoding/base64"

// Base64Encode 编码字符串为 base64
func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// Base64Decode 解码 base64 字符串
func Base64Decode(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Base64DecodeWithoutErr 解码 base64 字符串，忽略错误
func Base64DecodeWithoutErr(s string) string {
	result, _ := Base64Decode(s)
	return result
}

// Base64URLEncode URL 安全的 base64 编码
func Base64URLEncode(s string) string {
	return base64.URLEncoding.EncodeToString([]byte(s))
}

// Base64URLDecode URL 安全的 base64 解码
func Base64URLDecode(s string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
