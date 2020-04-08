package component

import (
	"fmt"
	"time"

	"github.com/nocmk2/sachima/dur"

	"github.com/allegro/bigcache"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	//DB connection
	DB *gorm.DB
	//GlobalCache cache
	GlobalCache *bigcache.BigCache
)

type RestResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func init() {
	// Connect to DB
	var err error
	DB, err = gorm.Open(dur.ReadDBConnStr("sachima_local"))
	if err != nil {
		panic(fmt.Sprintf("failed to connect to DB: %v", err))
	}

	// Initialize cache
	GlobalCache, err = bigcache.NewBigCache(bigcache.DefaultConfig(30 * time.Minute)) // Set expire time to 30 mins
	if err != nil {
		panic(fmt.Sprintf("failed to initialize cahce: %v", err))
	}
}
