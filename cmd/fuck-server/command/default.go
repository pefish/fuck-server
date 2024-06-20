package command

import (
	"github.com/pefish/fuck-server/pkg/global"
	"github.com/pefish/fuck-server/pkg/task"
	go_logger "github.com/pefish/go-logger"

	"github.com/pefish/go-commander"
	go_mysql "github.com/pefish/go-mysql"
	task_driver "github.com/pefish/go-task-driver"
)

type DefaultCommand struct {
}

func NewDefaultCommand() *DefaultCommand {
	return &DefaultCommand{}
}

func (dc *DefaultCommand) Config() interface{} {
	return &global.GlobalConfig
}

func (dc *DefaultCommand) Data() interface{} {
	return nil
}

func (dc *DefaultCommand) Init(command *commander.Commander) error {
	go_mysql.MysqlInstance.SetLogger(go_logger.Logger)
	err := go_mysql.MysqlInstance.ConnectWithConfiguration(go_mysql.Configuration{
		Host:     global.GlobalConfig.Db.Host,
		Username: global.GlobalConfig.Db.User,
		Password: global.GlobalConfig.Db.Pass,
		Database: global.GlobalConfig.Db.Db,
	})
	if err != nil {
		return err
	}

	return nil
}

func (dc *DefaultCommand) OnExited(command *commander.Commander) error {
	go_mysql.MysqlInstance.Close()
	return nil
}

func (dc *DefaultCommand) Start(command *commander.Commander) error {
	taskDriver := task_driver.NewTaskDriver()
	taskDriver.Register(task.NewExecutor())
	taskDriver.RunWait(command.Ctx)
	return nil
}
