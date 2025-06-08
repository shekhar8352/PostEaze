// file: pkg/config/app.go

package config

import (
	"gorm.io/gorm"
)

type AppContext struct {
	DB *gorm.DB
}

var appCtx *AppContext

func InitAppContext(db *gorm.DB) {
	appCtx = &AppContext{DB: db}
}

func GetAppContext() *AppContext {
	return appCtx
}
