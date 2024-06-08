package logger_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pansachin/employee-service/pkg/logger"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		Name string
		F    func(buf *bytes.Buffer)
		Want string
	}{
		{
			Name: "NewTextLogger",
			F: func(buf *bytes.Buffer) {
				buf.Reset()
				c := &logger.Config{
					Writer: buf,
					Debug:  true,
				}
				l, _ := logger.NewTextLogger(c, "test", "1.0")
				l.Info("text-formatted-log", "key", "value")
			},
			Want: `level=INFO msg=text-formatted-log service=test version=1.0 key=value`,
		},
		{
			Name: "NewJSONLogger",
			F: func(buf *bytes.Buffer) {
				buf.Reset()
				c := &logger.Config{
					Writer: buf,
					Debug:  true,
					Json:   true,
				}
				l, _ := logger.NewJSONLogger(c, "test", "1.0")
				l.Info("json-formatted-log", "key", "value")
			},
			Want: `"level":"INFO","msg":"json-formatted-log","service":"test","version":"1.0","key":"value"}`,
		},
		{
			Name: "NewLoggerJSONINFO",
			F: func(buf *bytes.Buffer) {
				buf.Reset()
				c := &logger.Config{
					Writer: buf,
					Debug:  false,
					Json:   true,
				}
				l, _ := logger.NewLogger(c, "test", "1.0")
				l.Info("json-formatted-log", "keyinfo", "valueinfo")
				l.Debug("json-formatted-log", "keydebug", "valuedebug")
			},
			Want: `"level":"INFO","msg":"json-formatted-log","service":"test","version":"1.0","keyinfo":"valueinfo"}`,
		},
		{
			Name: "NewLoggerJSONDEBUG",
			F: func(buf *bytes.Buffer) {

				buf.Reset()
				c := &logger.Config{
					Writer: buf,
					Debug:  true,
					Json:   true,
				}
				l, _ := logger.NewLogger(c, "test", "1.0")
				l.Info("json-formatted-log", "keyinfo", "valueinfo")
				l.Debug("json-formatted-log", "keydebug", "valuedebug")
			},
			Want: `"level":"DEBUG","msg":"json-formatted-log","service":"test","version":"1.0","keydebug":"valuedebug"}`,
		},
	}

	var buf bytes.Buffer
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.F(&buf)
			if result := buf.String(); !strings.Contains(result, tt.Want) {
				t.Errorf("Expected log message '%s' to contain '%s'", result, tt.Want)
			}
		})
	}
}
