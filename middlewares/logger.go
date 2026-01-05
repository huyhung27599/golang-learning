package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func LoggerMiddleware() gin.HandlerFunc {
	logPath := "logs/http.log"
	

logger :=	zerolog.New(&lumberjack.Logger{
    Filename:   logPath,
    MaxSize:    2, // megabytes
    MaxBackups: 5,
    MaxAge:     7, //days
    Compress:   true, // disabled by default
	LocalTime: true,
}).With().Timestamp().Logger()

	return func(c *gin.Context) {
		start := time.Now()
		contentType := c.GetHeader("Content-Type")
		requestBody := make(map[string]any)
		var formFiles []map[string]any

		if strings.HasPrefix(contentType, "multipart/form-data") {
			if err := c.Request.ParseMultipartForm(32<<20); err == nil && c.Request.MultipartForm != nil {
				for key, values := range c.Request.MultipartForm.Value {
					if len(values) == 1 {
						requestBody[key] = values[0]
					} else {
						requestBody[key] = values
					}
				}

				for _, files := range c.Request.MultipartForm.File {
					for _, file := range files {
					formFiles = append(formFiles, map[string]any{
						"filename": file.Filename,
						"size": formatFileSize(file.Size),
						"header": file.Header,
					})
					}
				}

				if len(formFiles) > 0 {
					requestBody["form_files"] = formFiles
				}
			}
			c.Request.Body = io.NopCloser(bytes.NewBufferString(""))
		} else {
			bodyBytes,err := io.ReadAll(c.Request.Body)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to read request body")
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			if strings.HasPrefix(contentType, "application/json") {
				_= json.Unmarshal(bodyBytes, &requestBody)
			} else {
                values, _ := url.ParseQuery(string(bodyBytes))
				for key, vals := range values {
					if len(vals) == 1 {
						requestBody[key] = vals[0]
					} else {
						requestBody[key] = vals
					}
				}
			}

		}


		customWriter := &CustomResponseWriter{
			ResponseWriter: c.Writer,
			body: bytes.NewBufferString(""),
		}



		c.Writer = customWriter
		c.Next()
		duration := time.Since(start)

	
		statusCode := c.Writer.Status()

		responseContentType := c.Writer.Header().Get("Content-Type")
		responseBodyRaw := customWriter.body.String()

		var responseBodyParsed any

		if strings.HasPrefix(responseContentType, "images/") {
                 responseBodyParsed = "Image"
		} else if strings.HasPrefix(responseContentType, "application/json") || strings.HasPrefix(strings.TrimSpace(responseBodyRaw), "{") || strings.HasPrefix(strings.TrimSpace(responseBodyRaw), "[") {
           if err := json.Unmarshal([]byte(responseBodyRaw), &responseBodyParsed); err != nil {
			responseBodyParsed = responseBodyRaw
		   }
		} else {
			responseBodyParsed = responseBodyRaw
		}

		if statusCode >= 500 {
			logEvent := logger.Error()
			logEvent.Str("method", c.Request.Method).Str("path", c.Request.URL.Path).Str("query", c.Request.URL.RawQuery).Str("ip", c.ClientIP()).Str("user_agent", c.Request.UserAgent()).Str("referer", c.Request.Referer()).Str("protocol", c.Request.Proto).Str("host", c.Request.Host).Str("duration", duration.String()).
			Int("status", statusCode).Int("size", c.Writer.Size()).Int("duration", int(duration.Milliseconds())).Msg("Request")
		} else if statusCode >= 400 {
			logEvent := logger.Warn()
			logEvent.Str("method", c.Request.Method).Str("path", c.Request.URL.Path).Str("query", c.Request.URL.RawQuery).Str("ip", c.ClientIP()).Str("user_agent", c.Request.UserAgent()).Str("referer", c.Request.Referer()).Str("protocol", c.Request.Proto).Str("host", c.Request.Host).Str("duration", duration.String()).
			Int("status", statusCode).Int("size", c.Writer.Size()).Int("duration", int(duration.Milliseconds())).Msg("Request")
		} else {
			logEvent := logger.Info()
			logEvent.Str("method", c.Request.Method).Str("path", c.Request.URL.Path).Str("query", c.Request.URL.RawQuery).Str("ip", c.ClientIP()).Str("user_agent", c.Request.UserAgent()).Str("referer", c.Request.Referer()).Str("protocol", c.Request.Proto).Str("host", c.Request.Host).Str("duration", duration.String()).
			Int("status", statusCode).Int("size", c.Writer.Size()).Int("duration", int(duration.Milliseconds())).Interface("request_body", requestBody).Interface("response_body", responseBodyParsed).Msg("Request")
		}
	}
}

func formatFileSize(size int64) string {
	switch {
	case size >= 1<<20:
		return fmt.Sprintf("%.2f MB", float64(size)/(1<<20))
	case size >= 1<<10:
		return fmt.Sprintf("%.2f KB", float64(size)/(1<<10))
	default:
		return fmt.Sprintf("%d bytes", size)
	}
}