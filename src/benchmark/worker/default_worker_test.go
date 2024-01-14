package worker

import (
	"errors"
	"net/url"
	"sync"
	"testing"

	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/benchmark/_mocks"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/benchmark/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDefaultWorker(t *testing.T) {
	ctrl := gomock.NewController(t)

	client := mocks.NewMockClient(ctrl)
	ticksPerSecond := 30

	target, err := url.Parse("http://localhost:8080")
	if err != nil {
		t.Fatal(err)
	}

	targets := []*url.URL{target}
	terminate := make(chan bool)

	t.Run("Empty Ramp", func(t *testing.T) {
		// given
		var wg sync.WaitGroup
		wg.Add(1)
		ramp := config.NewLinearRamp([]config.RequestRamp{})
		worker := NewDefaultWorker(1, &wg, client, ramp, targets, terminate, ticksPerSecond)

		// when
		go worker.Work()

		wg.Wait()

		// then
		select {
		case <-worker.results:
			assert.True(t, false)
		case <-worker.errors:
			assert.True(t, false)
		default:
			assert.True(t, true)
		}
	})

	t.Run("One Request", func(t *testing.T) {
		// given
		var wg sync.WaitGroup
		wg.Add(1)
		ramp := config.NewLinearRamp([]config.RequestRamp{
			{Duration: 1, TargetRPS: 1},
		})
		worker := NewDefaultWorker(1, &wg, client, ramp, targets, terminate, ticksPerSecond)

		// when
		client.EXPECT().Send(target.Host, target.Path).Return(uint64(200), nil).Times(1)
		go worker.Work()

		wg.Wait()

		// then
		select {
		case statusCode := <-worker.results:
			assert.Equal(t, uint64(200), statusCode)
		case <-worker.errors:
			assert.True(t, false)
		default:
			assert.True(t, false)
		}
	})

	t.Run("One Request with Error", func(t *testing.T) {
		// given
		var wg sync.WaitGroup
		wg.Add(1)
		ramp := config.NewLinearRamp([]config.RequestRamp{
			{Duration: 1, TargetRPS: 1},
		})
		worker := NewDefaultWorker(1, &wg, client, ramp, targets, terminate, ticksPerSecond)

		// when
		client.EXPECT().Send(target.Host, target.Path).Return(uint64(500), errors.New("Network Error")).Times(1)
		go worker.Work()

		wg.Wait()

		// then
		select {
		case statusCode := <-worker.results:
			assert.Equal(t, uint64(500), statusCode)
		case err := <-worker.errors:
			assert.Error(t, err)
		default:
			assert.True(t, false)
		}
	})

	t.Run("A Few Requests with errors and successes", func(t *testing.T) {
		t.Skip("Skip this test because it is flaky")
		// given
		var wg sync.WaitGroup
		wg.Add(1)
		ramp := config.NewLinearRamp([]config.RequestRamp{
			{Duration: 0, TargetRPS: 10},
			{Duration: 5, TargetRPS: 10},
		})
		worker := NewDefaultWorker(1, &wg, client, ramp, targets, terminate, ticksPerSecond)

		count := 0
		shouldErrorCount := 0
		shouldSuccessCount := 0

		// when
		client.EXPECT().
			Send(target.Host, target.Path).
			DoAndReturn(func(host, path string) (uint64, error) {
				count++
				if count%2 == 0 {
					shouldSuccessCount++
					return uint64(200), nil
				}
				shouldErrorCount++
				return uint64(500), errors.New("Network Error")
			}).
			Times(50)
		go worker.Work()

		wg.Wait()

		// then
		successes := 0
		errors := 0
		for i := 0; i < count; i++ {
			select {
			case statusCode := <-worker.results:
				if statusCode == 200 {
					successes++
				} else {
					errors++
				}
			default:
				assert.True(t, false)
			}
		}

		assert.Equal(t, shouldSuccessCount, successes)
		assert.Equal(t, shouldErrorCount, errors)

		errors = 0
		for i := 0; i < shouldErrorCount; i++ {
			select {
			case err := <-worker.errors:
				assert.Error(t, err)
				errors++
			default:
				assert.True(t, false)
			}
		}

		assert.Equal(t, shouldErrorCount, errors)
	})
}
