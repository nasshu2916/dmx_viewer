package infrastructure

import (
	"sync"
	"time"

	"github.com/beevik/ntp"
)

type TimeRepositoryImpl struct {
	mu       sync.RWMutex
	response *ntp.Response
}

func NewTimeRepositoryImpl() *TimeRepositoryImpl {
	return &TimeRepositoryImpl{}
}

func (r *TimeRepositoryImpl) GetTime() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.response == nil {
		return time.Now()
	}
	return time.Now().UTC().Add(r.response.ClockOffset)
}

func (r *TimeRepositoryImpl) ExistsNTPResponse() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.response != nil
}

func (r *TimeRepositoryImpl) SetQueryResponse(response *ntp.Response) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.response = response
}
