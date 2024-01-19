package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

type ApiLimitBody struct {
	LimitMin int `json:"limitMin"`
}

func GetScreenTimeLimitMinutes(user string) (int, error) {
	apiLimits := make(map[string]string)
	err := viper.UnmarshalKey("app.screenTimeLimitApiMinutes", &apiLimits)
	timeScreenLimitApi, okTimeScreenLimitApi := apiLimits[user]

	if err == nil && okTimeScreenLimitApi {
		res, err := http.Get(timeScreenLimitApi)
		if err == nil {
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err == nil {
				var apiLimitBody ApiLimitBody
				err = json.Unmarshal(body, &apiLimitBody)
				if err == nil {
					return apiLimitBody.LimitMin, err
				}
			}
		}
	}

	limits := make(map[string]int)
	err = viper.UnmarshalKey("app.screenTimeLimitMinutes", &limits)
	timeScreenLimit, okTimeScreenLimit := limits[user]

	if err == nil && okTimeScreenLimit {
		return timeScreenLimit, nil
	} else {
		// try api limimts
		apiLimits := make(map[string]string)
		err := viper.UnmarshalKey("app.screenTimeLimitMinutes", &apiLimits)
		return 0, err
	}
}
