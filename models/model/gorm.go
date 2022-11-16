package model

import (
	"bookstore_im/common/log"
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func Init(dsn string) error {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.WithContext(context.TODO()).Errorf("open mysql error, err: %v", err)
		return err
	}

	return nil
}
