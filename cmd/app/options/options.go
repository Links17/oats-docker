package options

import (
	"fmt"
	"github.com/bndr/gojenkins"
	oatsConfig "github.com/caoyingjunz/pixiulib/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"oats-docker/cmd/app/config"
	"oats-docker/pkg/db"
	"oats-docker/pkg/log"
	"oats-docker/pkg/util"
	"os"
	"strings"
)

const (
	maxIdleConns = 10
	maxOpenConns = 100

	defaultConfigFile = ""
)

// Options has all the params needed to run a oats
type Options struct {
	// The default values.
	ComponentConfig config.Config
	GinEngine       *gin.Engine

	DB      *gorm.DB
	Factory db.ShareDaoFactory // 数据库接口

	// CICD 的驱动接口
	CicdDriver *gojenkins.Jenkins

	// ConfigFile is the location of the oats server's configuration file.
	ConfigFile string
}

func NewOptions() (*Options, error) {
	return &Options{
		ConfigFile: defaultConfigFile,
	}, nil
}

// Complete completes all the required options
func (o *Options) Complete() error {
	// 配置文件优先级: 默认配置，环境变量，命令行
	if len(o.ConfigFile) == 0 {
		// Try to read config file path from env.
		if cfgFile := os.Getenv("ConfigFile"); cfgFile != "" {
			o.ConfigFile = cfgFile
		} else {
			o.ConfigFile = defaultConfigFile
		}
	}

	c := oatsConfig.New()
	c.SetConfigFile(o.ConfigFile)
	c.SetConfigType("yaml")
	if err := c.Binding(&o.ComponentConfig); err != nil {
		return err
	}

	// 初始化默认 api 路由
	o.GinEngine = gin.Default()

	// 注册依赖组件
	if err := o.register(); err != nil {
		return err
	}
	return nil
}

// BindFlags binds the oats Configuration struct fields
func (o *Options) BindFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.ConfigFile, "config", "", "The location of the oats configuration file")
}

func (o *Options) register() error {
	if err := o.registerLogger(); err != nil { // 注册日志
		return err
	}
	if err := o.registerDatabase(); err != nil { // 注册数据库
		return err
	}

	return nil
}

func (o *Options) registerLogger() error {
	logType := strings.ToLower(o.ComponentConfig.Default.LogType)
	if logType == "file" {
		// 判断文件夹是否存在，不存在则创建
		if err := util.EnsureDirectoryExists(o.ComponentConfig.Default.LogDir); err != nil {
			return err
		}
	}
	// 注册日志
	log.Register(logType, o.ComponentConfig.Default.LogDir, o.ComponentConfig.Default.LogLevel)

	return nil
}

func (o *Options) registerDatabase() error {
	sqlConfig := o.ComponentConfig.Pgsql
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		sqlConfig.Host,
		sqlConfig.User,
		sqlConfig.Password,
		sqlConfig.Name,
		sqlConfig.Port,
	)
	var err error
	if o.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		return err
	}
	// 设置数据库连接池
	sqlDB, err := o.DB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)

	o.Factory = db.NewDaoFactory(o.DB)

	return nil
}

// Validate validates all the required options.
func (o *Options) Validate() error {
	return nil
}
