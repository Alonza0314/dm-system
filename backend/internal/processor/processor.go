package processor

import (
	"backend/internal/context"
	"backend/logger"
	"time"
)

type ProcessorParams struct {
	Username string
	Password string

	JwtSecret    string
	JwtExpiresIn time.Duration

	MaxHistory int

	*context.DmContext

	*logger.BackendLogger
}

type Processor struct {
	username string
	password string

	jwtSecret    string
	jwtExpiresIn time.Duration

	maxHistory int

	*context.DmContext

	*logger.BackendLogger
}

func NewProcessor(params *ProcessorParams) *Processor {
	return &Processor{
		username: params.Username,
		password: params.Password,

		jwtSecret:    params.JwtSecret,
		jwtExpiresIn: params.JwtExpiresIn,

		maxHistory: params.MaxHistory,

		DmContext: params.DmContext,

		BackendLogger: params.BackendLogger,
	}
}

func (p *Processor) Release() {
	p.BackendLogger.ProcLog.Infoln("Release Processor...")

	p.DmContext.Release()

	p.BackendLogger.ProcLog.Infoln("Processor released")
}
