package database

import (
	dbTypes "CASBanSyncer/src/database/types"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func AddSuperBan(db *gorm.DB, userId string, date time.Time) {
	if userId == "1087968824" {
		return
	}
	db.Table("superban_table").Create(&dbTypes.Superban{
		UserId:         userId,
		MotivationText: "CAS Ban Import",
		UserDate:       date,
		UserFirstName:  fmt.Sprintf("NB%s", userId),
		IdOperator:     "1065189838",
	})
}
