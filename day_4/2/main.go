//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
// Hint: time.Ticker can be used
// Hint 2: to calculate timediff for Advanced lvl use:
//
//  start := time.Now()
//	// your work
//	t := time.Now()
//	elapsed := t.Sub(start) // 1s or whatever time has passed

package main

import (
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

const LimitProcessingSeconds = int64(10 * time.Second)
const CheckingPeriod = 500 * time.Millisecond

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	ticker := time.NewTicker(CheckingPeriod)
	defer ticker.Stop()

	processEnd := make(chan time.Duration)
	go func(processEnd chan time.Duration) {
		start := time.Now()
		process()
		t := time.Now()
		elapsed := t.Sub(start)
		processEnd <- elapsed
	}(processEnd)

	for {
		select {
		case <-processEnd:
			//User finish current process
			return true
		case <-ticker.C:
			//Each half second increase used time for user and check limit
			u.TimeUsed += int64(CheckingPeriod)
			if !u.IsPremium && u.TimeUsed >= LimitProcessingSeconds {
				return false
			}
		}
	}
}

func main() {
	RunMockServer()
}
