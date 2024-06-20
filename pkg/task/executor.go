package task

import (
	"context"
	"time"

	"github.com/pefish/fuck-server/pkg/db"
	program_manager "github.com/pefish/fuck-server/pkg/program-manager"
	go_mysql "github.com/pefish/go-mysql"

	go_logger "github.com/pefish/go-logger"
)

type ExecutorType struct {
	logger go_logger.InterfaceLogger
}

func NewExecutor() *ExecutorType {
	w := &ExecutorType{}
	w.logger = go_logger.Logger.CloneWithPrefix(w.Name())
	return w
}

func (t *ExecutorType) Init(ctx context.Context) error {
	return nil
}

func (t *ExecutorType) Run(ctx context.Context) error {
	programs := make([]db.Program, 0)
	err := go_mysql.MysqlInstance.Select(
		&programs,
		&go_mysql.SelectParams{
			TableName: "program",
			Select:    "*",
			Where:     "status = 0",
		},
	)
	if err != nil {
		return err
	}
	if len(programs) == 0 {
		t.Logger().Debug("No program waiting to run.")
		return nil
	}
	for _, program := range programs {
		t.Logger().InfoF("New program <%s>.", program.Name)
		go func(program db.Program) {
			programManager := program_manager.NewProgramManager(t.Logger(), program)
			programManager.Run(ctx)
		}(program)
	}
	return nil
}

func (t *ExecutorType) Stop() error {
	return nil
}

func (t *ExecutorType) Name() string {
	return "Executor"
}

func (t *ExecutorType) Interval() time.Duration {
	return 3 * time.Second
}

func (t *ExecutorType) Logger() go_logger.InterfaceLogger {
	return t.logger
}
