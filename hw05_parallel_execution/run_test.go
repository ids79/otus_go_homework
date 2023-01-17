package hw05parallelexecution

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 48
		err := Run(tasks, workersCount, maxErrorsCount)

		require.ErrorIs(t, err, ErrErrorsLimitExceeded)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		workersCount := 5
		maxErrorsCount := 1
		var runTasksCount int32
		var startTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&startTasksCount, 1)
				require.Eventually(t,
					func() bool {
						if atomic.LoadInt32(&startTasksCount)-atomic.LoadInt32(&runTasksCount) == int32(workersCount) {
							return true
						} else if i > tasksCount-workersCount {
							return true
						}
						return false
					},
					time.Second, 10*time.Millisecond)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	t.Run("if were errors in all tasks, and m=-1, than started and finished all tasks", func(t *testing.T) {
		tasksCount := 20
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 8
		maxErrorsCount := -1
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Nil(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "all tasks were started and finished")
	})

	t.Run("if n<=0, than all tasks are not executed", func(t *testing.T) {
		tasksCount := 2
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := -1
		maxErrorsCount := 1
		err := Run(tasks, workersCount, maxErrorsCount)

		require.ErrorIs(t, err, ErrNoThreadsForExecute)
	})
}
