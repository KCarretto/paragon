package agent

import (
	"context"

	"go.uber.org/zap"
)

func (a *agent) Run(ctx context.Context) {
	a.logger.Info("Running agent")

	// TODO: Panic recover handling
	// TODO: Handle siginterupt and sigterm
	for {
		select {
		case <-ctx.Done():
			a.logger.Info("Agent execution stopping", zap.Error(ctx.Err()))
			a.pending.Wait()
			return
		default:
			a.runOnce(ctx)

			// TODO: Allow ctx to be cancelled even when sleeping
			a.sleep()
		}
	}
}

func (a *agent) runOnce(ctx context.Context) {
	// Collect completed task results
	pickCtx, cancel := context.WithTimeout(ctx, a.pickTimeout)
	defer cancel()
	results := a.pickResults(pickCtx)

	// Report results & fetch new tasks
	if a.reporter == nil {
		// TODO: Error handling
		a.logger.Panic("Reporter was set to nil!")
	}
	tasks := a.reporter.Report(ctx, a, results)

	// Schedule execution of new tasks
	execCtx, cancel := context.WithTimeout(ctx, a.execTimeout)
	defer cancel()
	a.Execute(execCtx, tasks...)
}
