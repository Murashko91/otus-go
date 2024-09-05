package scheduler

import (
	"context"
)

type Scheduler interface {
	Run(context.Context)
	Cancel()
}
