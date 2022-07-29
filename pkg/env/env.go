package env

import (
	"os"
	"strings"

	"github.com/funbinary/go_example/pkg/errors"
)

// All
//  @Description:获取环境变量
//  @return []string
//
func All() []string {
	return os.Environ()
}

// Map
//  @Description: 获取环境变量map
//  @return map[string]string
//
func Map() map[string]string {
	m := make(map[string]string)
	i := 0
	for _, s := range os.Environ() {
		i = strings.IndexByte(s, '=')
		m[s[0:i]] = s[i+1:]
	}
	return m
}

// Get
//  @Description: 获取指定key的环境变量值，如果不存在且def存在返回第一个def，如果key和def都不存在返回nil
//  @param key 环境变量key
//  @param def 默认值
//  @return interface{}
//
func Get(key string, def ...interface{}) interface{} {
	v, ok := os.LookupEnv(key)
	if !ok {
		if len(def) > 0 {
			return def[0]
		}
		return nil
	}
	return v
}

// Set
//  @Description: 设置环境变量，仅在当前进程有效
//  @param key
//  @param value
//  @return err
//
func Set(key, value string) (err error) {
	err = os.Setenv(key, value)
	if err != nil {
		err = errors.Wrapf(err, `set environment key-value failed with key "%s", value "%s"`, key, value)
	}
	return
}

// SetMap
//  @Description: 设置环境变量，仅在当前进程有效
//  @param m
//  @return err
//
func SetMap(m map[string]string) (err error) {
	for k, v := range m {
		if err = Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

// Contains
//  @Description: 判断环境变量key是否存在，存在返回true
//  @param key
//  @return bool
//
func Contains(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

// Remove
//  @Description: 从环境变量中移除key
//  @param key
//  @return err
//
func Remove(key ...string) (err error) {
	for _, v := range key {
		if err = os.Unsetenv(v); err != nil {
			err = errors.Wrapf(err, `delete environment key failed with key "%s"`, v)
			return err
		}
	}
	return nil
}
