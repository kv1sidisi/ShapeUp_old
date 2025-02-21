package operations

import "time"

const (
	refreshOperationType      = "refresh"
	accessOperationType       = "access"
	confirmationOperationType = "confirmation"

	accessTokenExpireTime       = time.Minute * 30
	refreshTokenExpireTime      = time.Hour * 24 * 30
	confirmationTokenExpireTime = time.Hour
)

type OpInfo struct {
	Type       string
	ExpireTime time.Duration
}

func GetOperationInfo(operation string) OpInfo {
	switch operation {
	case refreshOperationType:
		return OpInfo{
			Type:       refreshOperationType,
			ExpireTime: refreshTokenExpireTime,
		}
	case accessOperationType:
		return OpInfo{
			Type:       accessOperationType,
			ExpireTime: accessTokenExpireTime,
		}
	case confirmationOperationType:
		return OpInfo{
			Type:       confirmationOperationType,
			ExpireTime: confirmationTokenExpireTime,
		}
	}

	return OpInfo{}
}
