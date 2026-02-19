package utils

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
}

func GlobalVoiceKey(name string) string {
	return "voice:global:" + name
}

const GlobalVoiceNamesKey = "voice_names:global"

func GlobalIpaKey(name string) string {
	return "ipa:global:" + name
}

const GlobalIpaNamesKey = "ipa_names:global"
