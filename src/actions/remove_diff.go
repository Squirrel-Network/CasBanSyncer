package actions

import (
	"CASBanSyncer/src/actions/types"
	"CASBanSyncer/src/consts"
	"CASBanSyncer/src/database"
	dbTypes "CASBanSyncer/src/database/types"
	"CASBanSyncer/src/http"
	"CASBanSyncer/src/utils/concurrency"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
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
			res := http.ExecuteRequest(
				fmt.Sprintf("https://api.cas.chat/check?user_id=%s", superban.UserId),
				http.Retries(3),
				http.Headers(map[string]string{
					// Spoofing the user agent to bypass Cloudflare bot detection
					"User-Agent": fmt.Sprintf(
						"Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
						consts.Devices[rand.Intn(len(consts.Devices))],
					),
				}),
			)
			if res.Error != nil {
				fmt.Println(superban.UserId)
				err = res.Error
				return
			}
			var ban types.CasResult
			err = json.Unmarshal(res.Read(), &ban)
			if err != nil {
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
