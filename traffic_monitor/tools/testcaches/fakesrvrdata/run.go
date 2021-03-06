package fakesrvrdata

import (
	"errors"
	"math/rand"
	"time"
)

type BytesPerSec struct { // TODO change to PerMin? PerHour? (to allow, e.g. one 5xx per hour)
	Min FakeRemap
	Max FakeRemap
}

// runValidate verifies the FakeServerData Remaps match the RemapIncrements
func runValidate(s *FakeServerData, remapIncrements map[string]BytesPerSec) error {
	for r, _ := range s.ATS.Remaps {
		if _, ok := remapIncrements[r]; !ok {
			return errors.New("remap increments missing server remap '" + r + "'")
		}
	}
	for r, bps := range remapIncrements {
		if _, ok := s.ATS.Remaps[r]; !ok {
			return errors.New("remap increments has remap not in server '" + r + "'")
		}
		if bps.Min.InBytes > bps.Max.InBytes || bps.Min.InBytes < 0 {
			return errors.New("invalid remap increments InBytes: must be Max >= Min >= 0)")
		}
		if bps.Min.OutBytes > bps.Max.OutBytes || bps.Min.OutBytes < 0 {
			return errors.New("invalid remap increments OutBytes: must be Max >= Min >= 0)")
		}
		if bps.Min.Status2xx > bps.Max.Status2xx || bps.Min.Status2xx < 0 {
			return errors.New("invalid remap increments Status2xx: must be Max >= Min >= 0)")
		}
		if bps.Min.Status3xx > bps.Max.Status3xx || bps.Min.Status3xx < 0 {
			return errors.New("invalid remap increments Status3xx: must be Max >= Min >= 0)")
		}
		if bps.Min.Status4xx > bps.Max.Status4xx || bps.Min.Status4xx < 0 {
			return errors.New("invalid remap increments Status4xx: must be Max >= Min >= 0)")
		}
		if bps.Min.Status5xx > bps.Max.Status5xx || bps.Min.Status5xx < 0 {
			return errors.New("invalid remap increments Status5xx: must be Max >= Min >= 0)")
		}
	}
	return nil
}

// Run takes a FakeServerData and a config, and starts running it, incrementing stats per the config. Returns a Threadsafe accessor to the running FakeServerData pointer, whose value may be accessed, but MUST NOT be modified.
// TODO add increments for Rcv,SndPackets, ProcLoadAvg variance, ConfigReloads
func Run(s FakeServerData, remapIncrements map[string]BytesPerSec) (Ths, error) {
	// TODO seed rand? Param?
	if err := runValidate(&s, remapIncrements); err != nil {
		return Ths{}, errors.New("invalid configuration: " + err.Error())
	}
	ths := NewThs()
	ths.Set(&s)
	go run(ths, remapIncrements)
	return ths, nil
}

// run starts a goroutine incrementing the FakeServerData's values according to the remapIncrements. Never returns.
func run(srvrThs Ths, remapIncrements map[string]BytesPerSec) {
	tickSecs := uint64(1) // adjustable for performance (i.e. a higher number is less CPU work)
	for {
		time.Sleep(time.Second * time.Duration(tickSecs))
		srvr := srvrThs.Get()
		newRemaps := copyRemaps(srvr.ATS.Remaps)
		for remap, increments := range remapIncrements {
			srvrRemap := newRemaps[remap]
			if increments.Min.InBytes != increments.Min.InBytes {
				i := uint64(rand.Int63n(int64((increments.Max.InBytes-increments.Min.InBytes)*tickSecs))) + (increments.Min.InBytes * tickSecs)
				srvrRemap.InBytes += i
				srvr.System.ProcNetDev.RcvBytes += i
			}
			if increments.Min.OutBytes != increments.Max.OutBytes {
				i := uint64(rand.Int63n(int64((increments.Max.OutBytes-increments.Min.OutBytes)*tickSecs))) + (increments.Min.OutBytes * tickSecs)
				srvrRemap.OutBytes += i
				srvr.System.ProcNetDev.SndBytes += i
			}
			if increments.Min.Status2xx != increments.Max.Status2xx {
				srvrRemap.Status2xx += uint64(rand.Int63n(int64((increments.Max.Status2xx-increments.Min.Status2xx)*tickSecs))) + (increments.Min.Status2xx * tickSecs)
			}
			if increments.Min.Status3xx != increments.Max.Status3xx {
				srvrRemap.Status3xx += uint64(rand.Int63n(int64((increments.Max.Status3xx-increments.Min.Status3xx)*tickSecs))) + (increments.Min.Status3xx * tickSecs)
			}
			if increments.Min.Status4xx != increments.Max.Status4xx {
				srvrRemap.Status4xx += uint64(rand.Int63n(int64((increments.Max.Status4xx-increments.Min.Status4xx)*tickSecs))) + (increments.Min.Status4xx * tickSecs)
			}
			if increments.Min.Status5xx != increments.Max.Status5xx {
				srvrRemap.Status5xx += uint64(rand.Int63n(int64((increments.Max.Status5xx-increments.Min.Status5xx)*tickSecs))) + (increments.Min.Status5xx * tickSecs)
			}
			newRemaps[remap] = srvrRemap
		}
		srvr.ATS.Remaps = newRemaps
		srvrThs.Set(srvr)
	}
}
