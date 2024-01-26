package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

type ScrenTimeConfiguration struct {
	ScreenLimitMin  int    `json:"screenLimitMin"`
	TimeWindowStart string `json:"timeWindowStart"`
	TimeWindowStop  string `json:"timeWindowStop"`
}

func GetScreenTimeConfiguration(user string) (ScrenTimeConfiguration, error) {
	// try to get from api
	apiUrlMap := make(map[string]string)
	err := viper.UnmarshalKey("app.screenTimeApi", &apiUrlMap)
	apiUrl, okAPiUrl := apiUrlMap[user]

	if err == nil && okAPiUrl {
		res, err := http.Get(apiUrl)
		if err == nil {
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err == nil {
				var apiLimitBody ScrenTimeConfiguration
				err = json.Unmarshal(body, &apiLimitBody)
				if err == nil {
					return apiLimitBody, nil
				}
			}
		}
	}

	// get from config with a default
	var screenLimitMin int
	err = viper.UnmarshalKey(fmt.Sprintf("app.screenTimeLimitMinutes.%s", user), &screenLimitMin)

	if err != nil {
		screenLimitMin = 0
	}

	var timeWindowStart string
	err = viper.UnmarshalKey(fmt.Sprintf("app.screenTimeWindow.%s.start", user), &timeWindowStart)

	if err != nil {
		timeWindowStart = "00:00"
	}

	var timeWindowStop string
	err = viper.UnmarshalKey(fmt.Sprintf("app.screenTimeWindow.%s.stop", user), &timeWindowStop)

	if err != nil {
		timeWindowStop = "00:00"
	}

	return ScrenTimeConfiguration{
		ScreenLimitMin:  screenLimitMin,
		TimeWindowStart: timeWindowStart,
		TimeWindowStop:  timeWindowStop,
	}, nil

}
