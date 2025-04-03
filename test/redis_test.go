package test

import (
	"context"
	_ "gorm-gen-skeleton/internal/bootstrap"
	"gorm-gen-skeleton/internal/variable"
	"testing"
)

func TestRedis(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error(err)
		}
	}()
	redisClient := variable.Redis
	t.Log(redisClient.Get(context.Background(), "test").Result())
}
