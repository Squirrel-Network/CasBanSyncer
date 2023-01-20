package actions

import (
	"CASBanSyncer/src/actions/types"
	"CASBanSyncer/src/consts"
	"CASBanSyncer/src/http"
	"encoding/json"
	"fmt"
	"math/rand"
)

func CheckApi(userId string) (*types.CasResult, error) {
	res := http.ExecuteRequest(
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
	if res.Error != nil {
		return nil, res.Error
	}
	var ban *types.CasResult
	err := json.Unmarshal(res.Read(), &ban)
	if err != nil {
		return nil, err
	}
	return ban, nil
}
