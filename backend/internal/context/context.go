package context

import (
	"backend/constant"
	"backend/logger"
	"fmt"
	"strconv"
)

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

func (ctx *DmContext) RequestId(target string) (int, error) {
	ctx.BackendLogger.CtxLog.Debugf("Requesting new ID for target: %s", target)

	exist, err := ctx.dbContext.Exist(constant.COLL_ID, target)
	if err != nil {
		return -1, fmt.Errorf("failed to check if target id exists: %w", err)
	}

	if exist {
		id, err := ctx.dbContext.Load(constant.COLL_ID, target)
		if err != nil {
			return -1, fmt.Errorf("failed to get %s's id", target)
		}

		numId, err := strconv.Atoi(id)
		if err != nil {
			return -1, fmt.Errorf("failed to convert %s's id number", target)
		}

		numId += 1
		if err := ctx.dbContext.Save(constant.COLL_ID, target, strconv.Itoa(numId)); err != nil {
			return -1, fmt.Errorf("failed to update %s's id", target)
		}

		return numId, nil
	} else {
		numId := 1

		if err := ctx.dbContext.Save(constant.COLL_ID, target, strconv.Itoa(numId)); err != nil {
			return -1, fmt.Errorf("failed to update %s's id", target)
		}

		return numId, nil
	}
}
