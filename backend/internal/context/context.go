package context

import "backend/logger"

type DmContextParams struct {
	DbType string
	DbPath string

	*logger.BackendLogger
}

type DmContext struct {
	*dbContext

	*logger.BackendLogger
}

func NewDmContext(params *DmContextParams) *DmContext {
	dbContext, err := newDbContext(&dbContextIE{
		dbType: params.DbType,
		dbPath: params.DbPath,

		BackendLogger: params.BackendLogger,
	})
	if err != nil {
		params.BackendLogger.CtxLog.Errorf("Failed to create dbContext: %v", err)
		return nil
	}

	return &DmContext{
		dbContext: dbContext,

		BackendLogger: params.BackendLogger,
	}
}

func (ctx *DmContext) Release() {
	ctx.BackendLogger.CtxLog.Infoln("Release DmContext...")

	ctx.dbContext.release()

	ctx.BackendLogger.CtxLog.Infoln("DmContext released")
}

func (ctx *DmContext) Db() *dbContext {
	return ctx.dbContext
}
