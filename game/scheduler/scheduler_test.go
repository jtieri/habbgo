package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type MockJob struct {
	t *testing.T

	ctx      context.Context
	cancel   context.CancelFunc
	interval time.Duration
	ticker   *time.Ticker
	running  bool

	// We want to assert that each execution of the job is occurring on the
	// desired interval but the Job interface functions have no return values.
	// We use these variable to assign the number of times the job was executed
	// and if its been stopped, then we make assertions on these values.
	timesExecuted int
	cancelled     bool
}

func NewMockJob(ctx context.Context, cancel context.CancelFunc, interval time.Duration) *MockJob {
	return &MockJob{
		ctx:       ctx,
		cancel:    cancel,
		interval:  interval,
		ticker:    time.NewTicker(interval),
		cancelled: false,
		running:   false,
	}
}

func (m *MockJob) Execute() {
	m.running = true
	defer m.ticker.Stop()
	for {
		select {
		case <-m.ticker.C:
			// duration has passed, execute logic for this job
			m.t.Logf("Executing job - %s\n", time.Now())
			m.timesExecuted += 1
		case <-m.ctx.Done():
			// the job was cancelled, clean up appropriately and stop this job
			m.t.Log("Job was cancelled")
			m.cancelled = true
			m.running = false
			return
		}
	}
}

func (m *MockJob) Stop() {
	m.cancel()
}

func (m *MockJob) Running() bool {
	return m.running
}

func (m *MockJob) Name() string {
	return "Mock Job"
}

func TestCreateAndCancelJob(t *testing.T) {
	jobDuration := 500 * time.Millisecond
	timesToExecute := 10

	ctx, cancel := context.WithCancel(context.Background())

	t.Logf("Creating mock job to run on ticker of %s", jobDuration)
	j := NewMockJob(ctx, cancel, jobDuration)
	j.t = t

	// Call the execute function
	t.Log("Executing job")
	go j.Execute()

	// Let the job execute for the desired amount of times
	for i := 0; i < timesToExecute; i++ {
		time.Sleep(j.interval)
	}

	// Assert that the job actually ran as many times as we expected.
	require.Equal(t, timesToExecute, j.timesExecuted)

	// Cancel the job
	j.Stop()

	// Wait 1 second for the job to be cancelled
	time.Sleep(1 * time.Second)

	// Assert the job was cancelled
	require.Equal(t, true, j.cancelled)
}

func TestTwoJobsCancelOne(t *testing.T) {
	jobDuration := 500 * time.Millisecond
	timesToExecute := 10

	t.Logf("Creating mock jobs to run on ticker of %s", jobDuration)
	ctx, cancel := context.WithCancel(context.Background())
	job1 := NewMockJob(ctx, cancel, jobDuration)
	job1.t = t

	ctx, cancel = context.WithCancel(context.Background())
	job2 := NewMockJob(ctx, cancel, jobDuration)
	job2.t = t

	// Call the execute function
	t.Log("Executing jobs")
	go job1.Execute()
	go job2.Execute()

	// Let the job execute for the desired amount of times
	for i := 0; i < timesToExecute; i++ {
		time.Sleep(job1.interval)
	}

	// Assert that the job actually ran as many times as we expected.
	require.Equal(t, timesToExecute, job1.timesExecuted)

	// Cancel the job
	job1.Stop()

	// Wait 1 second for the job to be cancelled
	time.Sleep(1 * time.Second)

	// Assert that job1 is cancelled and not running but job2 is still running
	require.Equal(t, true, job1.cancelled)
	require.Equal(t, false, job1.running)
	require.Equal(t, false, job2.cancelled)
	require.Equal(t, true, job2.running)
}
