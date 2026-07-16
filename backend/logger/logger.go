package logger

import (
	"backend/constant"

	loggergo "github.com/Alonza0314/logger-go/v2"
	loggergoModel "github.com/Alonza0314/logger-go/v2/model"
	loggergoUtil "github.com/Alonza0314/logger-go/v2/util"
)

type BackendLogger struct {
	*loggergo.Logger

	CfgLog  loggergoModel.LoggerInterface
	AccLog  loggergoModel.LoggerInterface
	BckLog  loggergoModel.LoggerInterface
	ProcLog loggergoModel.LoggerInterface
	CtxLog  loggergoModel.LoggerInterface
	DbLog   loggergoModel.LoggerInterface
	CatLog  loggergoModel.LoggerInterface
	DevLog  loggergoModel.LoggerInterface
}

func NewBackendLogger(level loggergoUtil.LogLevelString, writeToFile bool) *BackendLogger {
	logger := loggergo.NewLogger("/tmp/dm-system.log", writeToFile)
	logger.SetLevel(level)

	return &BackendLogger{
		Logger: logger,

		CfgLog:  logger.WithTags(constant.CFG_LOG),
		AccLog:  logger.WithTags(constant.ACC_LOG),
		BckLog:  logger.WithTags(constant.BCK_LOG),
		ProcLog: logger.WithTags(constant.PROC_LOG),
		CtxLog:  logger.WithTags(constant.CTX_LOG),
		DbLog:   logger.WithTags(constant.DB_LOG),
		CatLog:  logger.WithTags(constant.CAT_LOG),
		DevLog:  logger.WithTag(constant.DEV_LOG),
	}
}
