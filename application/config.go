package app

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
	"time"
)

type Configuration struct {
	configReader *viper.Viper
}

func NewConfiguration() *Configuration {
	return &Configuration{
		configReader: viper.GetViper(),
	}
}

func (c *Configuration) ConfigPath(path string) {
	dir, file := filepath.Split(path)
	filename := strings.TrimSuffix(file, filepath.Ext(file))
	c.configReader.SetConfigName(filename)
	c.configReader.SetConfigType(strings.TrimPrefix(filepath.Ext(file), "."))
	c.configReader.AddConfigPath(dir)
	fmt.Printf("Reading config, directory: %s, file: %s\n", dir, file)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func (c *Configuration) GetStr(env string) string {
	fmt.Printf("GetConfigStr: %v", env)
	envVar := strings.ToUpper(strings.ReplaceAll(env, ".", "_"))
	c.configReader.SetDefault(env, c.configReader.GetString(env))
	err := c.configReader.BindEnv(env, envVar)
	if err != nil {
		panic(err)
	}
	val := c.configReader.GetString(env)
	fmt.Printf("=%v\n", val)
	return val
}

func (c *Configuration) GetInt(env string) int {
	fmt.Printf("GetConfigInt: %v", env)
	envVar := strings.ToUpper(strings.ReplaceAll(env, ".", "_"))
	c.configReader.SetDefault(env, c.configReader.GetInt(env))
	err := c.configReader.BindEnv(env, envVar)
	if err != nil {
		panic(err)
	}
	val := c.configReader.GetInt(env)
	fmt.Printf("=%v\n", val)
	return val
}

func (c *Configuration) GetInt64(env string) int64 {
	fmt.Printf("GetConfigInt64: %v", env)
	envVar := strings.ToUpper(strings.ReplaceAll(env, ".", "_"))
	c.configReader.SetDefault(env, c.configReader.GetInt64(env))
	err := c.configReader.BindEnv(env, envVar)
	if err != nil {
		panic(err)
	}
	val := c.configReader.GetInt64(env)
	fmt.Printf("=%v\n", val)
	return val
}

func (c *Configuration) GetDuration(env string) time.Duration {
	fmt.Printf("GetConfigDuration: %v", env)
	envVar := strings.ToUpper(strings.ReplaceAll(env, ".", "_"))
	c.configReader.SetDefault(env, c.configReader.GetDuration(env))
	err := c.configReader.BindEnv(env, envVar)
	if err != nil {
		panic(err)
	}
	val := c.configReader.GetDuration(env)
	fmt.Printf("=%v\n", val)
	return val
}

func (c *Configuration) GetBool(env string) bool {
	fmt.Printf("GetConfigBool: %v", env)
	envVar := strings.ToUpper(strings.ReplaceAll(env, ".", "_"))
	c.configReader.SetDefault(env, c.configReader.GetBool(env))
	err := c.configReader.BindEnv(env, envVar)
	if err != nil {
		panic(err)
	}
	val := c.configReader.GetBool(env)
	fmt.Printf("=%v\n", val)
	return val
}
