package sender

import (
	"context"
)

type Sender interface {
	Run(context.Context)
	Close()
}
