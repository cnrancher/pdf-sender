package apis

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var CacheClient *cache.Cache

func CreateCacheClient() {

	CacheClient = cache.New(10*time.Minute, 20*time.Minute)

}
