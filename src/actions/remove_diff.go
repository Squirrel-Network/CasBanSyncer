package actions

import (
	"CASBanSyncer/src/database"
	dbTypes "CASBanSyncer/src/database/types"
	"CASBanSyncer/src/utils/concurrency"
	"gorm.io/gorm"
	"sync"
)

func RemoveDiff(db *gorm.DB) (int, error) {
	result := database.GetAllSuperban(db)
	var wg sync.WaitGroup
	var listAllSuperban []string
	var removedCount int
	var err error
	semaphore := concurrency.NewPool[any](100)
	for _, v := range result {
		wg.Add(1)
		semaphore.Enqueue(func(params ...any) {
			defer wg.Done()
			if err != nil {
				return
			}
			superban := params[0].(*dbTypes.Superban)
			ban, err2 := CheckApi(superban.UserId)
			if err2 != nil {
				err = err2
				return
			}
			if ban.Banned {
				if ban.Result.TimeAdded.Unix() != superban.UserDate.Unix() {
					db.Table("superban_table").Where("user_id = ?", superban.UserId).Update("user_date", ban.Result.TimeAdded)
				}
			} else {
				listAllSuperban = append(listAllSuperban, superban.UserId)
				if len(listAllSuperban) == 500 {
					db.Table("superban_table").Where("user_id IN ?", listAllSuperban).Delete(&dbTypes.Superban{})
					removedCount += 500
					listAllSuperban = []string{}
				}
			}
		}, v)
	}
	wg.Wait()
	if err != nil {
		return 0, err
	}
	if len(listAllSuperban) > 0 {
		db.Table("superban_table").Where("user_id IN ?", listAllSuperban).Delete(&dbTypes.Superban{})
		removedCount += len(listAllSuperban)
	}
	return removedCount, nil
}
