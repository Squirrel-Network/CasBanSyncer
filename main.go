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
		log.Println("Syncing...")
		result, err := actions.RemoveDiff(db)
		log.Println("Removed", result, "expired cas bans")
		if err != nil {
			log.Println(err)
			continue
		}
		result, err = actions.AddDiff(db)
		log.Println("Added", result, "cas bans from export.csv")
		if err != nil {
			log.Println(err)
			continue
		}
		var totalAdd int
		for i := 0; i < 60; i++ {
			result, err = actions.AddNew(db)
			totalAdd += result
			if err != nil {
				log.Println(err)
				continue
			}
			time.Sleep(time.Minute * 1)
		}
		log.Println("Added", result, "cas bans from the feed")
	}
}
