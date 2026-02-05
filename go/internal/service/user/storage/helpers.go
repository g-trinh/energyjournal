package storage

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func isNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
}

func getString(data map[string]interface{}, key string) string {
	if v, ok := data[key].(string); ok {
		return v
	}
	return ""
}

func getTimestamp(data map[string]interface{}, key string) (time.Time, error) {
	v, ok := data[key]
	if !ok || v == nil {
		return time.Time{}, nil
	}

	switch t := v.(type) {
	case time.Time:
		return t, nil
	default:
		return time.Time{}, nil
	}
}
