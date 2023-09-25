package utils

import "github.com/spf13/viper"

func GetScreenTimeLimitMinutes(user string) (int, error) {
	limits := make(map[string]int)
	err := viper.UnmarshalKey("app.screenTimeLimitMinutes", &limits)
	timeScreenLimit, okTimeScreenLimit := limits[user]

	if err == nil && okTimeScreenLimit {
		return timeScreenLimit, nil
	} else {
		return 0, err
	}
}
