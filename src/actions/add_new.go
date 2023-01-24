package actions

import (
	"CASBanSyncer/src/database"
	dbTypes "CASBanSyncer/src/database/types"
	"CASBanSyncer/src/http"
	"CASBanSyncer/src/utils"
	"gorm.io/gorm"
	"regexp"
)

func AddNew(db *gorm.DB) (int, error) {
	var result []*dbTypes.Superban
	db.Table("superban_table").Find(&result)
	rMap := make(map[string]bool)
	for _, superban := range result {
		rMap[superban.UserId] = true
	}
	res := http.ExecuteRequest("https://cas.chat/feed")
	if res.Error != nil {
		return 0, res.Error
	}
	rgx, _ := regexp.Compile(`<a\s+href="/query\?u=[0-9]+"\s+class=".*?">#([0-9]+)</a>.*?([0-9]+|less)\s+([a-z]+)`)
	resMatch := rgx.FindAllStringSubmatch(res.ReadString(), -1)
	var count int
	for _, content := range resMatch {
		userId := content[1]
		if rMap[userId] {
			continue
		}
		dur, err := utils.StringToTime(content[2], content[3])
		if err != nil {
			return 0, nil
		}
		count++
		database.AddSuperBan(db, userId, dur)
	}
	return count, nil
}
