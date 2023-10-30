package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	// viper.SetConfigFile("config.yaml") // 配置文件名称

	viper.AddConfigPath("/etc/tgb/")          // 查找配置文件所在的路径
	viper.AddConfigPath("$HOME/.config/tgb/") // 多次调用以添加多个搜索路径
	viper.AddConfigPath(".")                  // 还可以在工作目录中查找配置

	// bind env
	viper.AutomaticEnv()
	viper.SetEnvPrefix("TBD")                              // 设置环境变量的前缀
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 设置环境变量的命名规则

	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {             // 处理读取配置文件的错误
		return fmt.Errorf("Fatal error config file: %s \n", err)
	}
	return nil
}
