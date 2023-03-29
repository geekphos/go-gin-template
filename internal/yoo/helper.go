package yoo

import (
	"os"
	"phos.cc/yoo/internal/pkg/log"
	"phos.cc/yoo/internal/yoo/store"
	"phos.cc/yoo/pkg/db"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfigName = "yoo.yaml"
)

// initConfig 设置需要读取的配置文件名、环境变量，并读取配置文件内容到 viper 中.
func initConfig() {
	if cfgFile != "" {
		// 从命令行选项指定的配置文件中读取
		viper.SetConfigFile(cfgFile)
	} else {
		// 查找用户主目录
		home, err := os.UserHomeDir()
		// 如果获取用户主目录失败，打印 `'Error: xxx` 错误，并退出程序（退出码为 1）
		cobra.CheckErr(err)

		// 将用 `$HOME` 目录加入到配置文件的搜索路径中
		viper.AddConfigPath(home)

		// 设置配置文件格式为 YAML(YAML格式清晰易读，并且支持复杂的配置结构)
		viper.SetConfigType("yaml")

		// 配置文件名称（没有文件扩展名）
		viper.SetConfigName(defaultConfigName)
	}

	// 读取匹配的环境变量
	viper.AutomaticEnv()

	// 读取配置文件。如果指定了配置文件名，则使用指定的配置文件，否则在注册的搜索路径中搜索
	if err := viper.ReadInConfig(); err != nil {
		log.Errorw("Failed to read config file", "error", err)
	}

	log.Debugw("Using config file", "path", viper.ConfigFileUsed())
}

// logOptions 从 viper 中读取日志配置，构建 `*log.Options` 并返回.
// 注意：`viper.Get<Type>()` 中 key 的名字需要使用 `.` 分割，以跟 YAML 中保持相同的缩进.
func logOptions() *log.Options {
	return &log.Options{
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		Level:             viper.GetString("log.level"),
		Format:            viper.GetString("log.format"),
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
	}
}

func initStore() error {
	dbOptions := &db.MySQLOptions{
		Host:                  viper.GetString("db.host"),
		Username:              viper.GetString("db.username"),
		Password:              viper.GetString("db.password"),
		Database:              viper.GetString("db.database"),
		MaxIdleConnections:    viper.GetInt("db.max-idle-connections"),
		MaxOpenConnections:    viper.GetInt("db.max-open-connections"),
		MaxConnectionLifeTime: viper.GetDuration("db.max-connection-life-time"),
		LogLevel:              viper.GetInt("db.log-level"),
	}

	if ins, err := db.NewMySQL(dbOptions); err != nil {
		return err
	} else {
		_ = store.NewStore(ins)
	}

	return nil
}
