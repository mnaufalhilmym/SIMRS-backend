package jwt

import "time"

func Renew[T any](data *T, exp *time.Time) (*string, error) {
	if exp.Sub(time.Now()).Seconds() < float64(*conf.duration)/2 {
		return Create(data)
	}
	return nil, nil
}
