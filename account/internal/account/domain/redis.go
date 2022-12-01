package domain

import "time"

type RedisPort interface {
	Set(key string, value string, expTime time.Duration) bool
	Get(key string) (value string)
	Delete(key string) bool
}
