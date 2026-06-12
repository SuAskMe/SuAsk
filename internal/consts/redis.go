package consts

const (
	RedisSendCodePrefix  = "vCode_"
	RedisCountCodePrefix = "cCode_"
	RedisJWTPrefix       = "jwt_"
)

// for Guest rate limiting

const (
	RedisGuestDevicePrefix = "guest:device:"
	RedisGuestIPPrefix     = "guest:ip:"
	GuestDeviceLimit       = 3
	GuestIPLimit           = 100
	GuestRateLimitTTL      = 3600 // 1 hour in seconds
)
