package RedisUtils

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"git.paygear.ir/giftino/account/internal/account/domain"
)

type redisService struct {
	client *redis.Client
}

func New(cli *redis.Client) domain.RedisPort {
	return redisService{client: cli}
}

func (r redisService) Set(key string, value string, expTime time.Duration) bool {
	err := r.client.Set(context.TODO(), key, value, expTime).Err()
	if err != nil {
		fmt.Println("error on setting key on redis DB.")
		return false
	} else {
		fmt.Println("successful setting key on redis DB.")
		return true
	}
}

func (r redisService) Get(key string) (value string) {
	redisResult, err := r.client.Get(context.TODO(), key).Result()
	if err != nil {
		fmt.Println("Error on GetAccountConfirmation on RedisDB.Get")
		fmt.Println(err)
		return ""
	} else {
		return redisResult
	}
}

func (r redisService) Delete(key string) bool {
	err := r.client.Del(context.TODO(), key).Err()
	if err != nil {
		fmt.Println("error on deleting key from redis DB.")
		return false
	} else {
		fmt.Println("successful deletion key from redis DB.")
		return true
	}
}
