package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinearRamp(t *testing.T) {

	t.Run("Empty ramp", func(t *testing.T) {
		// given
		ramp := NewLinearRamp([]RequestRamp{})

		// when
		rps := ramp.TargetRPS(1)

		// then
		assert.Equal(t, -1, rps)
	})

	t.Run("ramp with one value", func(t *testing.T) {
		// given
		ramp := NewLinearRamp([]RequestRamp{
			{Duration: 1, TargetRPS: 1},
		})

		// when
		rps := ramp.TargetRPS(1)

		// then
		assert.Equal(t, 1, rps)
	})

	t.Run("ramp with one value with high RPS", func(t *testing.T) {
		// given
		ramp := NewLinearRamp([]RequestRamp{
			{Duration: 1, TargetRPS: 100},
		})

		// when
		rps := ramp.TargetRPS(1)

		// then
		assert.Equal(t, 100, rps)
	})

	t.Run("ramp with one value with high Duration", func(t *testing.T) {
		// given
		ramp := NewLinearRamp([]RequestRamp{
			{Duration: 100, TargetRPS: 1},
		})

		// when
		for i := 1; i < 100; i++ {
			rps := ramp.TargetRPS(i)

			// then
			assert.Equal(t, 0, rps)
		}

		rps := ramp.TargetRPS(100)
		assert.Equal(t, 1, rps)
	})

	t.Run("ramp with one value with high Duration and high RPS", func(t *testing.T) {
		// given
		ramp := NewLinearRamp([]RequestRamp{
			{Duration: 100, TargetRPS: 100},
		})

		// when
		for i := 1; i < 100; i++ {
			rps := ramp.TargetRPS(i)

			// then
			assert.Equal(t, i, rps)
		}

		rps := ramp.TargetRPS(100)
		assert.Equal(t, 100, rps)
	})

	t.Run("ramp with one value with high Duration and high RPS", func(t *testing.T) {
		// given
		ramp := NewLinearRamp([]RequestRamp{
			{Duration: 100, TargetRPS: 50},
		})

		// when
		for i := 1; i < 100; i++ {
			rps := ramp.TargetRPS(i)

			// then
			assert.Equal(t, i/2, rps)
		}

		rps := ramp.TargetRPS(100)
		assert.Equal(t, 50, rps)
	})

	t.Run("ramp with more values", func(t *testing.T) {
		// given
		ramp := NewLinearRamp([]RequestRamp{
			{Duration: 60, TargetRPS: 1000},
			{Duration: 120, TargetRPS: 1000},
			{Duration: 180, TargetRPS: 2000},
			{Duration: 240, TargetRPS: 0},
		})

		// when
		for i := 1; i < 60; i++ {
			rps := ramp.TargetRPS(i)

			// then
			shouldValue := int((float32(1000) / float32(60)) * float32(i))
			assert.Equal(t, shouldValue, rps)
		}

		for i := 0; i < 60; i++ {
			rps := ramp.TargetRPS(i + 60)

			// then
			assert.Equal(t, 1000, rps)
		}

		for i := 0; i < 60; i++ {
			rps := ramp.TargetRPS(i + 120)

			// then
			shouldValue := int((float32(2000-1000)/float32(60))*float32(i)) + 1000
			assert.Equal(t, shouldValue, rps)
		}

		for i := 0; i < 60; i++ {
			rps := ramp.TargetRPS(i + 180)

			// then
			shouldValue := 2000 - int((float32(2000)/float32(60))*float32(i))
			assert.Equal(t, shouldValue, rps)
		}

		rps := ramp.TargetRPS(240)
		assert.Equal(t, 0, rps)
	})
}
