package program_manager

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pefish/fuck-server/pkg/db"
	vm "github.com/pefish/go-jsvm"
	go_logger "github.com/pefish/go-logger"
	go_mysql "github.com/pefish/go-mysql"
)

type ProgramManagerType struct {
	program db.Program

	logs          string
	logsLock      sync.Mutex
	logsMaxLength int

	logger    go_logger.InterfaceLogger
	wrappedVm *vm.WrappedVm
}

func NewProgramManager(logger go_logger.InterfaceLogger, program db.Program) *ProgramManagerType {
	return &ProgramManagerType{
		program:       program,
		logger:        logger,
		wrappedVm:     vm.NewVm(program.Content).SetLogger(logger),
		logsMaxLength: 5000,
	}
}

func (lm *ProgramManagerType) PushLogAndFlush(log string) {
	lm.PushLog(log)
	lm.FlushLogs()
}

func (lm *ProgramManagerType) PushLog(log string) {
	defer lm.logsLock.Unlock()
	lm.logsLock.Lock()
	lm.logs += log + "\n"
	if len(lm.logs) > lm.logsMaxLength {
		lm.logs = lm.logs[len(lm.logs)-lm.logsMaxLength:]
	}
}

func (lm *ProgramManagerType) FlushLogs() {
	defer lm.logsLock.Unlock()
	lm.logsLock.Lock()
	_, err := go_mysql.MysqlInstance.Update(
		&go_mysql.UpdateParams{
			TableName: "program",
			Update: map[string]interface{}{
				"logs": lm.logs,
			},
			Where: map[string]interface{}{
				"id": lm.program.Id,
			},
		},
	)
	if err != nil {
		lm.logger.Error(err)
		return
	}
}

func (lm *ProgramManagerType) watchStatus(ctx context.Context) error {
	timer := time.NewTimer(0)
	for {
		select {
		case <-timer.C:
			var program db.Program
			_, err := go_mysql.MysqlInstance.SelectFirst(
				&program,
				&go_mysql.SelectParams{
					TableName: "program",
					Select:    "*",
					Where:     "id = ?",
				},
				lm.program.Id,
			)
			if err != nil {
				return err
			}
			if program.Status == 2 {
				lm.wrappedVm.Kill()
			}
			timer.Reset(3 * time.Second)
		case <-ctx.Done():
			return nil
		}
	}
}

func (lm *ProgramManagerType) watchLogs(ctx context.Context) error {
	timer := time.NewTimer(0)
	for {
		select {
		case <-timer.C:
			lm.FlushLogs()
			timer.Reset(10 * time.Second)
		case <-ctx.Done():
			return nil
		}
	}
}

func (lm *ProgramManagerType) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := lm.watchLogs(ctx)
		if err != nil {
			lm.logger.Error(err)
			return
		}
	}()

	go func() {
		err := lm.watchStatus(ctx)
		if err != nil {
			lm.logger.Error(err)
			return
		}
	}()

	_, err := go_mysql.MysqlInstance.Update(
		&go_mysql.UpdateParams{
			TableName: "program",
			Update: map[string]interface{}{
				"status": 1,
			},
			Where: map[string]interface{}{
				"id": lm.program.Id,
			},
		},
	)
	if err != nil {
		lm.logger.Error(err)
		return
	}

	_, err = lm.wrappedVm.Run()
	if err != nil {
		lm.PushLogAndFlush(
			fmt.Sprintf(
				"Program <%s> exited with error <%s>.",
				lm.program.Name,
				err.Error(),
			),
		)

		_, err := go_mysql.MysqlInstance.Update(
			&go_mysql.UpdateParams{
				TableName: "program",
				Update: map[string]interface{}{
					"status": 3,
				},
				Where: map[string]interface{}{
					"id": lm.program.Id,
				},
			},
		)
		if err != nil {
			lm.logger.Error(err)
			return
		}
		return
	}

	_, err = go_mysql.MysqlInstance.Update(
		&go_mysql.UpdateParams{
			TableName: "program",
			Update: map[string]interface{}{
				"status": 4,
			},
			Where: map[string]interface{}{
				"id": lm.program.Id,
			},
		},
	)
	if err != nil {
		lm.logger.Error(err)
		return
	}
}
