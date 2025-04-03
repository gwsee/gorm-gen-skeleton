package test

import (
	"gorm-gen-gorm-gen-skeleton/internal/config"
	"gorm-gen-gorm-gen-skeleton/internal/config/driver"
	"gorm-gen-gorm-gen-skeleton/internal/variable"
	"testing"
)

func TestConfig(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("TestConfig filed:%v", err)
		}
	}()
	opt := config.Options{
		BasePath: variable.BasePath,
	}
	provider, _ := config.New(driver.New(), opt)
	t.Log(provider.Get("Application"))
	t.Log(provider.GetInt("Application.Int"))
	t.Log(provider.GetStringSlice("Application.Array"))
	t.Log(provider.Get("Database.Mysql.Test").([]string))
}

func TestViper(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("TestConfig filed:%v", err)
		}
	}()
	opt := config.Options{
		BasePath: "../",
	}
	viper := driver.New()
	_ = viper.Apply(opt)
	t.Log(viper.Get("Application"))
}
