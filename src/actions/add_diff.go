package actions

import (
	"CASBanSyncer/src/actions/types"
	"CASBanSyncer/src/consts"
	dbTypes "CASBanSyncer/src/database/types"
	"CASBanSyncer/src/http"
	"CASBanSyncer/src/utils/concurrency"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"strings"
	"sync"
)

func AddDiff(db *gorm.DB) (int, error) {
	var result []*dbTypes.Superban
	table := db.Table("superban_table")
	table.Find(&result)
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
			res2 := http.ExecuteRequest(
				fmt.Sprintf("https://api.cas.chat/check?user_id=%s", userId),
				http.Retries(3),
				http.Headers(map[string]string{
					// Spoofing the user agent to bypass Cloudflare bot detection
					"User-Agent": fmt.Sprintf(
						"Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
						consts.Devices[rand.Intn(len(consts.Devices))],
					),
				}),
			)
			if res2.Error != nil {
				err = res2.Error
				return
			}
			var ban types.CasResult
			err = json.Unmarshal(res2.Read(), &ban)
			if err != nil {
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
