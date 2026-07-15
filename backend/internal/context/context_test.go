package context_test

import (
	"backend/internal/context"
	"backend/internal/db"
	"backend/logger"
	"fmt"
	"os"
	"testing"

	loggergoUtil "github.com/Alonza0314/logger-go/v2/util"
)

var dmCtx *context.DmContext

func TestMain(m *testing.M) {
	db.DefaultDbPath = "/tmp/nsf_db_context_test.db"

	logger := logger.NewBackendLogger(loggergoUtil.LogLevelString("error"), false)

	dmCtx = context.NewDmContext(&context.DmContextParams{
		DbType: "bbolt",
		DbPath: db.DefaultDbPath,

		BackendLogger: logger,
	})

	exitCode := m.Run()

	dmCtx.Release()

	if err := os.Remove(db.DefaultDbPath); err != nil {
		fmt.Printf("Failed to remove test DB file: %v\n", err)
		exitCode = 1
	}

	os.Exit(exitCode)
}
