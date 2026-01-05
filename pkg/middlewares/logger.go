package middlewares

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type TypeLogger string

const (
	WARN_LOG  TypeLogger = "[WARNING]"
	SUCCESS   TypeLogger = "[SUCCESS]"
	ERROR_LOG TypeLogger = "[ERROR]"
	INFO_LOG  TypeLogger = "[INFO]"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		statusCode := c.Response().StatusCode()
		logType := getLogType(statusCode)
		statusColor := getStatusColor(statusCode)

		log.Printf("%s %s %s | %s%d%s | %v | %s | %s",
			logType,
			time.Now().Format("2006-01-02 15:04:05"),
			c.Method(),
			statusColor,
			statusCode,
			resetColor(),
			duration,
			c.IP(),
			c.Path(),
		)

		if err != nil {
			log.Printf("%s Error occurred: %v", ERROR_LOG, err)
		}

		return err
	}
}

func getLogType(statusCode int) TypeLogger {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return SUCCESS
	case statusCode >= 300 && statusCode < 400:
		return INFO_LOG
	case statusCode >= 400 && statusCode < 500:
		return WARN_LOG
	case statusCode >= 500:
		return ERROR_LOG
	default:
		return INFO_LOG
	}
}

func getStatusColor(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "\033[32m"
	case statusCode >= 300 && statusCode < 400:
		return "\033[36m"
	case statusCode >= 400 && statusCode < 500:
		return "\033[33m"
	case statusCode >= 500:
		return "\033[31m"
	default:
		return "\033[37m"
	}
}

func resetColor() string {
	return "\033[0m"
}

func CustomLogger(config LoggerConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if config.Skip != nil && config.Skip(c) {
			return c.Next()
		}

		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		statusCode := c.Response().StatusCode()
		logType := getLogType(statusCode)

		var logMsg string
		if config.CustomFormat != nil {
			logMsg = config.CustomFormat(c, duration, statusCode)
		} else {
			logMsg = fmt.Sprintf("%s %s | %d | %v | %s | %s",
				logType,
				c.Method(),
				statusCode,
				duration,
				c.IP(),
				c.Path(),
			)
		}

		log.Println(logMsg)

		if err != nil && config.LogErrors {
			log.Printf("%s Error: %v", ERROR_LOG, err)
		}

		return err
	}
}

type LoggerConfig struct {
	Skip         func(*fiber.Ctx) bool
	CustomFormat func(*fiber.Ctx, time.Duration, int) string
	LogErrors    bool
}
