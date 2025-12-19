package errors

import "fmt"

// APIError 高德 API 原生错误（包含错误码和描述）
type APIError struct {
	Code string // 错误码（如 10001：Key 无效）
	Info string // 错误描述
}

func (e *APIError) Error() string {
	return fmt.Sprintf("amap api error [code:%s]: %s", e.Code, e.Info)
}

type InvalidConfigError string // 配置错误
type NetworkError string       // 网络错误
type ParseError string         // 响应解析错误

func NewAPIError(code, info string) *APIError { return &APIError{Code: code, Info: info} }

func NewInvalidConfigError(msg string) error { return InvalidConfigError(msg) }
func (e InvalidConfigError) Error() string   { return "invalid config: " + string(e) }

func NewNetworkError(msg string) error { return NetworkError(msg) }
func (e NetworkError) Error() string   { return "network err: " + string(e) }

func NewParseError(msg string) error { return ParseError(msg) }
func (e ParseError) Error() string   { return "parser err: " + string(e) }
