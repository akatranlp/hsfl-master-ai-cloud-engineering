package config

type LinearRamp struct {
	requestRamp     []RequestRamp
	currentIndex    int
	currentDuration int
}

func NewLinearRamp(requestRamp []RequestRamp) *LinearRamp {
	return &LinearRamp{requestRamp, -1, 0}
}

func (r *LinearRamp) TargetRPS(duration int) int {
	currentRPS := 0
	currentDuration := 0
	for _, ramp := range r.requestRamp {
		if ramp.Duration < duration {
			currentRPS = ramp.TargetRPS
			currentDuration = ramp.Duration
		} else if ramp.Duration == duration {
			return ramp.TargetRPS
		} else {
			newStep := ramp.TargetRPS
			stepRate := float32(newStep-currentRPS) / float32(ramp.Duration-currentDuration)
			y := currentRPS + int(stepRate*float32(duration-currentDuration))
			return y
		}
	}
	return -1
}

func (r *LinearRamp) NextValue() int {
	r.currentDuration++
	var rampDuration int
	var trampTargetRPS int

	if len(r.requestRamp) == 0 {
		return -1
	}

	if r.currentIndex == -1 {
		rampDuration = 0
		trampTargetRPS = 0
	} else if r.currentIndex == len(r.requestRamp)-1 {
		return -1
	} else {
		rampDuration = r.requestRamp[r.currentIndex].Duration
		trampTargetRPS = r.requestRamp[r.currentIndex].TargetRPS
	}

	nextRamp := r.requestRamp[r.currentIndex+1]

	if r.currentDuration == nextRamp.Duration {
		r.currentIndex++
		return nextRamp.TargetRPS
	} else if r.currentDuration < nextRamp.Duration {
		stepRate := float32(nextRamp.TargetRPS-trampTargetRPS) / float32(nextRamp.Duration-rampDuration)
		y := trampTargetRPS + int(stepRate*float32(r.currentDuration-rampDuration))
		return y
	}

	return -1
}
