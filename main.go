package main

import (
	"CASBanSyncer/src/actions"
	"CASBanSyncer/src/consts"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

func main() {
	InitEnvironment()
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		consts.DBUser,
		consts.DbPassword,
		consts.DbHost,
		consts.DbPort,
		consts.DbName,
		consts.DbCharset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("CASBanSyncer started!")
	for {
		result, err := actions.RemoveDiff(db)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Removed", result, "cas bans")
		result, err = actions.AddDiff(db)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Added", result, "cas bans")
		time.Sleep(time.Hour * 1)
	}
}
