package zaplogger

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func buildStdOutCore(
	le zap.LevelEnablerFunc,
) zapcore.Core {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	ws := zapcore.Lock(os.Stdout)

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		ws,
		le,
	)
}

func buildFileCore(
	le zap.LevelEnablerFunc,
	c *Config,
) zapcore.Core {
	fw := zapcore.AddSync(&lumberjack.Logger{
		Filename:   c.FileLoggerPath,
		MaxSize:    c.FileLoggerMaxSize,
		MaxBackups: c.FileLoggerMaxBackups,
		MaxAge:     c.FileLoggerMaxAge,
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		fw,
		le,
	)
}

func buildStreamCore(
	le zap.LevelEnablerFunc,
	c *Config,
) zapcore.Core {
	sharedTransport := &kafka.Transport{}

	if c.StreamLoggerUsername != "" &&
		c.StreamLoggerPassword != "" {
		mechanism, err := scram.Mechanism(
			scram.SHA256,
			c.StreamLoggerUsername,
			c.StreamLoggerPassword)
		if err != nil {
			panic(err)
		}
		sharedTransport.SASL = mechanism
	}

	w := &kafka.Writer{
		Addr:      kafka.TCP(strings.Split(c.StreamLoggerAddrs, ",")...),
		Topic:     c.StreamLoggerTopic,
		Balancer:  &kafka.LeastBytes{},
		BatchSize: c.StreamLoggerBatchSize,
		BatchTimeout: time.Duration(c.StreamLoggerBatchTimeout) *
			time.Millisecond,
		Async:     true,
		Transport: sharedTransport,
	}

	kw := kafkaWriter{
		writer: w,
	}

	logTopic := zapcore.AddSync(&kw)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		logTopic,
		le,
	)
}

type kafkaWriter struct {
	writer *kafka.Writer
}

func (w *kafkaWriter) Write(b []byte) (int, error) {
	msg := kafka.Message{
		Value: b,
	}
	err := w.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		return 0, err
	}
	return 0, nil
}
