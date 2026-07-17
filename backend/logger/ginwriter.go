package logger

import "strings"

// ginWriter adapts an io.Writer so gin's internal logging (gin.DefaultWriter /
// gin.DefaultErrorWriter) is routed through the same tagged logger as the
// rest of the backend, instead of writing straight to os.Stdout.
type ginWriter struct {
	log interface {
		Infoln(args ...interface{})
	}
}

func (w *ginWriter) Write(p []byte) (int, error) {
	w.log.Infoln(strings.TrimRight(string(p), "\n"))
	return len(p), nil
}

// GinWriter returns an io.Writer that forwards everything written to it into
// this logger's GinLog tag. Assign it to gin.DefaultWriter (and, if wanted,
// gin.DefaultErrorWriter) before the gin.Engine is constructed -- gin.Logger()
// captures gin.DefaultWriter at middleware-setup time, not per request.
func (b *BackendLogger) GinWriter() *ginWriter {
	return &ginWriter{log: b.GinLog}
}
