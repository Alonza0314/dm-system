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

	DbType string
	DbPath string

	*logger.BackendLogger
}

type Processor struct {
	username string
	password string

	jwtSecret    string
	jwtExpiresIn time.Duration

	*context.DmContext

	*logger.BackendLogger
}

func NewProcessor(params *ProcessorParams) *Processor {
	dmCtx := context.NewDmContext(&context.DmContextParams{
		DbType: params.DbType,
		DbPath: params.DbPath,

		BackendLogger: params.BackendLogger,
	})
	if dmCtx == nil {
		params.BackendLogger.ProcLog.Errorf("Failed to create DmContext")
		return nil
	}

	return &Processor{
		username: params.Username,
		password: params.Password,

		jwtSecret:    params.JwtSecret,
		jwtExpiresIn: params.JwtExpiresIn,

		DmContext: dmCtx,

		BackendLogger: params.BackendLogger,
	}
}

func (p *Processor) Release() {
	p.BackendLogger.ProcLog.Infoln("Release Processor...")

	p.DmContext.Release()

	p.BackendLogger.ProcLog.Infoln("Processor released")
}
