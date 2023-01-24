package actions

import (
	dbTypes "CASBanSyncer/src/database/types"
	"CASBanSyncer/src/http"
	"CASBanSyncer/src/utils/concurrency"
	"gorm.io/gorm"
	"strings"
	"sync"
)

func AddDiff(db *gorm.DB) (int, error) {
	var result []*dbTypes.Superban
	db.Table("superban_table").Find(&result)
	rMap := make(map[string]bool)
	for _, superban := range result {
		rMap[superban.UserId] = true
	}
	res := http.ExecuteRequest("https://api.cas.chat/export.csv")
	if res.Error != nil {
		return 0, res.Error
	}
	users := strings.Split(res.ReadString(), "\n")
	semaphore := concurrency.NewPool[string](100)
	var wg sync.WaitGroup
	var err error
	var count int
	for _, user := range users {
		if user == "" {
			continue
		}
		if rMap[user] {
			continue
		}
		wg.Add(1)
		semaphore.Enqueue(func(params ...string) {
			if err != nil {
				return
			}
			defer wg.Done()
			userId := params[0]
			ban, err2 := CheckApi(userId)
			if err2 != nil {
				err = err2
				return
			}
			if ban.Banned {
				db.Table("superban_table").Create(&dbTypes.Superban{
					UserId:         userId,
					MotivationText: "CAS Ban Import",
					UserDate:       ban.Result.TimeAdded,
					UserFirstName:  "Unknown",
					IdOperator:     "1065189838",
				})
				count++
			}
		}, user)
	}
	wg.Wait()
	return count, nil
}
