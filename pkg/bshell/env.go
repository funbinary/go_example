package bshell

import "os"

func GetEnv(key string) string {
	v, _ := os.LookupEnv(key)
	return v
}
