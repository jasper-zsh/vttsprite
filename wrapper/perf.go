package wrapper

import "time"

type Perf struct {
	startAt              int64
	stopAt               int64
	curMillis            int64
	timingsTotal         map[string]float64
	timingsPeriodTotal   map[string]float64
	timingsCounter       map[string]int64
	timingsPeriodCounter map[string]int64
}

func (p *Perf) Start() {
	p.startAt = time.Now().UnixMilli()
	p.timingsTotal = make(map[string]float64)
	p.timingsPeriodTotal = make(map[string]float64)
	p.timingsCounter = make(map[string]int64)
	p.timingsPeriodCounter = make(map[string]int64)
}

func (p *Perf) GetSpeed() float64 {
	if p.stopAt == 0 {
		return float64(p.curMillis) / float64(time.Now().UnixMilli()-p.startAt)
	} else {
		return float64(p.curMillis) / float64(p.stopAt-p.startAt)
	}
}

func (p *Perf) Stop() {
	p.stopAt = time.Now().UnixMilli()
}

func (p *Perf) Record(curMillis int64) {
	p.curMillis = curMillis
}

func (p *Perf) RecordTiming(key string, value float64) {
	p.timingsTotal[key] += value
	p.timingsPeriodTotal[key] += value
	p.timingsCounter[key] += 1
	p.timingsPeriodCounter[key] += 1
}

func (p *Perf) AvgPeriodTiming(key string) float64 {
	ret := p.timingsPeriodTotal[key] / float64(p.timingsPeriodCounter[key])
	p.timingsPeriodTotal[key] = 0
	p.timingsPeriodCounter[key] = 0
	return ret
}

func (p *Perf) AvgTiming(key string) float64 {
	return p.timingsTotal[key] / float64(p.timingsCounter[key])
}
