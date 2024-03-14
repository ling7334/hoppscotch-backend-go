package strategies

import (
	"os"
	"strconv"
	"time"
)

const defaultAccessExpires = 24 * time.Hour
const defaultRefreshExpires = 7 * 24 * time.Hour

var accessExpires time.Duration
var refreshExpires time.Duration

func init() {
	sec, err := strconv.ParseInt(os.Getenv("ACCESS_TOKEN_VALIDITY"), 10, 0)
	if err != nil {
		accessExpires = defaultAccessExpires
	} else {
		accessExpires = time.Duration(sec) * time.Millisecond
	}
	sec, err = strconv.ParseInt(os.Getenv("REFRESH_TOKEN_VALIDITY"), 10, 0)
	if err != nil {
		refreshExpires = defaultRefreshExpires
	} else {
		refreshExpires = time.Duration(sec) * time.Millisecond
	}
}
