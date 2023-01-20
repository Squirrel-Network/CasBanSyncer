package database

import (
	"CASBanSyncer/src/database/types"
	"gorm.io/gorm"
	"time"
)

func GetAllSuperban(db *gorm.DB) []*types.Superban {
	var result []*types.Superban
	table := db.Table("superban_table")
	table.Where("(motivation_text LIKE ? OR motivation_text LIKE ?) AND user_date < ?", "cas%", "sn ban%", time.Now().AddDate(0, -2, -15)).Find(&result)
	return result
}
