package golog

import (
	"log/slog"
	"net/url"
	"time"
)

// LogKey type for varius additional keys
type LogKey string

const (
	PathKey      LogKey = "path"
	QueryKey     LogKey = "query"
	MethodKey    LogKey = "method"
	ExtraKey     LogKey = "extra"
	UserAgentKey LogKey = "useragent"
	IpKey        LogKey = "ip"
	StatusKey    LogKey = "status"
	LatencyKey   LogKey = "latency"
	LengthKey    LogKey = "length"
)

// Path returns an attribute for PathKey
func Path(path string) slog.Attr {
	return slog.String(string(PathKey), path)
}

// Query returns an attribute for QueryKey
func Query(query url.Values) slog.Attr {
	return slog.Any(string(QueryKey), query)
}

// Method returns an attribute for MethodKey
func Method(method string) slog.Attr {
	return slog.String(string(MethodKey), method)
}

// Extra returns an attribute for ExtraKey
func Extra(value any) slog.Attr {
	return slog.Any(string(ExtraKey), value)
}

// UserAgent returns an attribute for UserAgentKey
func UserAgent(ua string) slog.Attr {
	return slog.String(string(UserAgentKey), ua)
}

// Ip returns an attribute for IpKey
func Ip(ip string) slog.Attr {
	return slog.String(string(IpKey), ip)
}

// Status returns an attribute for StatusKey
func Status(statusCode int) slog.Attr {
	return slog.Int(string(StatusKey), statusCode)
}

// Latency returns an attribute for LatencyKey
func Latency(d time.Duration) slog.Attr {
	return slog.Float64(string(LatencyKey), d.Seconds())
}

// Length returns an attribute for LengthKey
func Length(l int) slog.Attr {
	return slog.Int(string(LengthKey), l)
}
