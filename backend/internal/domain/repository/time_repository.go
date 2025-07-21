package repository

import (
	"time"

	"github.com/beevik/ntp"
)

type TimeRepository interface {
	GetTime() time.Time
	ExistsNTPResponse() bool
	SetQueryResponse(response *ntp.Response)
}
