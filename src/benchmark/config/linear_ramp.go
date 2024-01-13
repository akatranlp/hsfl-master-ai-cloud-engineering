package config

type LinearRamp struct {
	requestRamp []RequestRamp
}

func NewLinearRamp(requestRamp []RequestRamp) *LinearRamp {
	return &LinearRamp{requestRamp}
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
