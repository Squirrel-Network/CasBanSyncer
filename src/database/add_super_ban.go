package database

import (
	dbTypes "CASBanSyncer/src/database/types"
	"gorm.io/gorm"
	"time"
)

func AddSuperBan(db *gorm.DB, userId string, date time.Time) {
	db.Table("superban_table").Create(&dbTypes.Superban{
		UserId:         userId,
		MotivationText: "CAS Ban Import",
		UserDate:       date,
		UserFirstName:  "Unknown",
		IdOperator:     "1065189838",
	})
}
