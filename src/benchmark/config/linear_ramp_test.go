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
		rpsTarget := ramp.TargetRPS(1)
		rpsNext := ramp.NextValue()

		// then
		assert.Equal(t, -1, rpsTarget)
		assert.Equal(t, -1, rpsNext)
	})

	t.Run("ramp with one value", func(t *testing.T) {
		// given
		ramp := NewLinearRamp([]RequestRamp{
			{Duration: 1, TargetRPS: 1},
		})

		// when
		rpsTarget := ramp.TargetRPS(1)
		rpsNext := ramp.NextValue()

		// then
		assert.Equal(t, 1, rpsTarget)
		assert.Equal(t, 1, rpsNext)
	})

	t.Run("ramp with one value with high RPS", func(t *testing.T) {
		// given
		ramp := NewLinearRamp([]RequestRamp{
			{Duration: 1, TargetRPS: 100},
		})

		// when
		rpsTarget := ramp.TargetRPS(1)
		rpsNext := ramp.NextValue()

		// then
		assert.Equal(t, 100, rpsTarget)
		assert.Equal(t, 100, rpsNext)
	})

	t.Run("ramp with one value with high Duration", func(t *testing.T) {
		// given
		ramp := NewLinearRamp([]RequestRamp{
			{Duration: 100, TargetRPS: 1},
		})

		// when
		for i := 1; i < 100; i++ {
			rpsTarget := ramp.TargetRPS(i)
			rpsNext := ramp.NextValue()

			// then
			assert.Equal(t, 0, rpsTarget)
			assert.Equal(t, 0, rpsNext)
		}

		rpsTarget := ramp.TargetRPS(100)
		rpsNext := ramp.NextValue()
		assert.Equal(t, 1, rpsTarget)
		assert.Equal(t, 1, rpsNext)
	})

	t.Run("ramp with one value with high Duration and high RPS", func(t *testing.T) {
		// given
		ramp := NewLinearRamp([]RequestRamp{
			{Duration: 100, TargetRPS: 100},
		})

		// when
		for i := 1; i < 100; i++ {
			rpsTarget := ramp.TargetRPS(i)
			rpsNext := ramp.NextValue()

			// then
			assert.Equal(t, i, rpsTarget)
			assert.Equal(t, i, rpsNext)
		}

		rpsTarget := ramp.TargetRPS(100)
		rpsNext := ramp.NextValue()
		assert.Equal(t, 100, rpsTarget)
		assert.Equal(t, 100, rpsNext)
	})

	t.Run("ramp with one value with high Duration and high RPS", func(t *testing.T) {
		// given
		ramp := NewLinearRamp([]RequestRamp{
			{Duration: 100, TargetRPS: 50},
		})

		// when
		for i := 1; i < 100; i++ {
			rpsTarget := ramp.TargetRPS(i)
			rpsNext := ramp.NextValue()

			// then
			assert.Equal(t, i/2, rpsTarget)
			assert.Equal(t, i/2, rpsNext)
		}

		rpsTarget := ramp.TargetRPS(100)
		rpsNext := ramp.NextValue()
		assert.Equal(t, 50, rpsTarget)
		assert.Equal(t, 50, rpsNext)
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
			rpsTarget := ramp.TargetRPS(i)
			rpsNext := ramp.NextValue()

			// then
			shouldValue := int((float32(1000) / float32(60)) * float32(i))
			assert.Equal(t, shouldValue, rpsTarget)
			assert.Equal(t, shouldValue, rpsNext)
		}

		for i := 0; i < 60; i++ {
			rpsTarget := ramp.TargetRPS(i + 60)
			rpsNext := ramp.NextValue()

			// then
			assert.Equal(t, 1000, rpsTarget)
			assert.Equal(t, 1000, rpsNext)
		}

		for i := 0; i < 60; i++ {
			rpsTarget := ramp.TargetRPS(i + 120)
			rpsNext := ramp.NextValue()

			// then
			shouldValue := int((float32(2000-1000)/float32(60))*float32(i)) + 1000
			assert.Equal(t, shouldValue, rpsTarget)
			assert.Equal(t, shouldValue, rpsNext)
		}

		for i := 0; i < 60; i++ {
			rpsTarget := ramp.TargetRPS(i + 180)
			rpsNext := ramp.NextValue()

			// then
			shouldValue := 2000 - int((float32(2000)/float32(60))*float32(i))
			assert.Equal(t, shouldValue, rpsTarget)
			assert.Equal(t, shouldValue, rpsNext)
		}

		rpsTarget := ramp.TargetRPS(240)
		rpsNext := ramp.NextValue()
		assert.Equal(t, 0, rpsTarget)
		assert.Equal(t, 0, rpsNext)
	})
}
