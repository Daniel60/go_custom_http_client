package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log         *zap.Logger
	LOG_OUTPUT  = "LOG_OUTPUT"
	LOG_LEVEL   = "LOG_LEVEL"
	wordsToMask = []string{"password", "secret", "token", "Authorization"} // Constante com as palavras a serem mascaradas
)

func init() {
	logConfig := zap.Config{
		OutputPaths: []string{getOutputLogs()},
		Level:       zap.NewAtomicLevelAt(getLevelLogs()),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "message",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	log, _ = logConfig.Build()
}

func Info(message string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	log.Error(message, tags...)
	log.Sync()
}

func getOutputLogs() string {
	output := strings.ToLower(strings.TrimSpace(os.Getenv(LOG_OUTPUT)))
	if output == "" {
		return "stdout"
	}
	return output
}

func getLevelLogs() zapcore.Level {
	level := strings.ToLower(strings.TrimSpace(os.Getenv(LOG_LEVEL)))
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

// Função para mascarar palavras sensíveis nas chaves dos headers
func maskSensitiveWords(headers map[string]string, wordsToMask []string) map[string]string {
	maskedHeaders := make(map[string]string)
	for key, value := range headers {
		// Verifica se a chave contém alguma palavra sensível
		for _, word := range wordsToMask {
			if strings.Contains(key, word) {
				// Se a chave contiver a palavra sensível, mascara o valor
				value = "***MASKED***"
				break
			}
		}
		maskedHeaders[key] = value
	}
	return maskedHeaders
}

func CustomGinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Coletar os headers da requisição
		requestHeaders := make(map[string]string)
		for key, values := range c.Request.Header {
			requestHeaders[key] = strings.Join(values, ", ")
		}

		// Mascarar palavras sensíveis nos headers
		maskedHeaders := maskSensitiveWords(requestHeaders, wordsToMask)

		responseBody := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = responseBody
		c.Next()

		end := time.Now()
		latency := end.Sub(start).Milliseconds()
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		var responseBodyMap map[string]interface{}
		responseBodyString := responseBody.body.String()
		if err := json.Unmarshal([]byte(responseBodyString), &responseBodyMap); err != nil {
			// Se não for possível converter, loga como string
			responseBodyMap = map[string]interface{}{
				"raw_response": responseBodyString,
			}
		}

		fields := []zap.Field{
			zap.String("result", "success"),
			zap.String("client_ip", clientIP),
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Int64("latency_ms", latency),
			zap.Any("client_request_headers", maskedHeaders), // Usar os headers mascarados
			zap.Any("response_body", responseBodyMap),
		}

		if reqHeaders, exists := c.Get("external_request_headers"); exists {
			apiHeaders := make(map[string]string)
			for k, v := range reqHeaders.(http.Header) {
				apiHeaders[k] = strings.Join(v, ", ")
			}
			fields = append(fields, zap.Any("external_request_headers", maskSensitiveWords(apiHeaders, wordsToMask)))
		}

		if errorMessage != "" {
			fields = append(fields, zap.String("error", errorMessage))
		}

		log.Info("Request", fields...)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
