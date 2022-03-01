package runner

import (
	"context"
	"time"

	"github.com/criyle/go-judge/envexec"
)

const minCPUPercent = 40 // 40%
const checkIntervalMS = 50

type waiter struct {
	timeLimit time.Duration
}

func (w *waiter) Wait(ctx context.Context, usg envexec.Process) bool {
	var lastCPUUsage time.Duration
	var totalTime time.Duration

	lastCheckTime := time.Now()
	// wait task done (check each interval)
	ticker := time.NewTicker(checkIntervalMS * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C: // interval
			usage := usg.Usage()
			cpuUsage := usage.Time

			cpuUsageDelta := cpuUsage - lastCPUUsage
			timeDelta := now.Sub(lastCheckTime)

			totalTime += durationMax(cpuUsageDelta, timeDelta*minCPUPercent/100)
			if totalTime > w.timeLimit {
				return true
			}

			lastCheckTime = now
			lastCPUUsage = cpuUsage

		case <-ctx.Done(): // returned
			return false
		}
	}
}

func durationMax(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}
