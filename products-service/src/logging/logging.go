package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"time"
)

type HTTPLogHandler interface {
	Handle(c context.Context, r slog.Record) error
	WithAttrs(attrs []slog.Attr) slog.Handler
	WithGroup(name string) slog.Handler
	Enabled(_ context.Context, level slog.Level) bool
}

type logging struct {
	client   *http.Client
	endpoint string
	service  string
}

func NewLogging(url string, service string) HTTPLogHandler {
	return &logging{
		client:   &http.Client{Timeout: time.Second * 5},
		endpoint: url,
		service:  service,
	}
}

func (h *logging) Handle(c context.Context, r slog.Record) error {
	// ✅ FIXED: Properly extract source information
	var source map[string]any
	if r.PC != 0 {
		frames := runtime.CallersFrames([]uintptr{r.PC})
		frame, _ := frames.Next()
		source = map[string]any{
			"function": frame.Function,
			"file":     frame.File,
			"line":     frame.Line,
		}
	}

	// Collect all attributes
	attrs := make(map[string]interface{})
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})

	entry := struct {
		Service string                 `json:"service"`
		Time    time.Time              `json:"time"`
		Level   string                 `json:"level"`
		Message string                 `json:"msg"`
		Source  map[string]any         `json:"source,omitempty"`
		Attrs   map[string]interface{} `json:"attrs,omitempty"`
	}{
		Service: h.service,
		Time:    r.Time,
		Level:   r.Level.String(),
		Message: r.Message,
		Source:  source,
		Attrs:   attrs,
	}

	b, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("marshal error:", err)
		return err
	}

	// ✅ Add debug logging to see what's being sent
	fmt.Printf("Sending log to %s: %s\n", h.endpoint, string(b))

	resp, err := h.client.Post(h.endpoint, "application/json", bytes.NewBuffer(b))
	if err != nil {
		fmt.Printf("POST error: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	// ✅ Check response status
	if resp.StatusCode != 200 {
		fmt.Printf("Log service returned status: %d\n", resp.StatusCode)
	} else {
		fmt.Println("✓ Log sent successfully")
	}

	return nil
}
func (l *logging) Enabled(_ context.Context, level slog.Level) bool {
	return true // send all logs; change if needed
}

func (l *logging) WithAttrs(attrs []slog.Attr) slog.Handler { return l }
func (l *logging) WithGroup(name string) slog.Handler       { return l }
