package config

import "time"

type RetryConfig struct {
	InitialInterval      time.Duration
	MaxInterval          time.Duration
	MaxElapsedTime       time.Duration
	MultiplicationFactor float64
}

func NewDefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		InitialInterval:      100 * time.Millisecond,
		MaxInterval:          10 * time.Second,
		MaxElapsedTime:       1 * time.Minute,
		MultiplicationFactor: 2.0,
	}
}
