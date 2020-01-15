package conf

import (
	"errors"


	"github.com/bolaxy/common"
	"github.com/spf13/viper"
)

const (
	DefaultHomeBase = ".bolaxy"
	DefaultEnv      = "BOLAXYDIR"
	configName      = "config"
	genesisName     = "genesis"
)

var (
	FileNotFound = errors.New("file not found")
)

// TryLoadConfig 载入配置文件。本函数可以被调用多次。
// 如果传入的文件路径为空，会主动寻找可用配置。查找顺序为
// filePath > ${HOME}/.bolaxy > $ENV[BOLAXYDIR] > WORKDIR > EXE RUN PATH
// 显式提供配置文件路径，会直接处理，非显式提供配置文件路径，直到找到正确的配置文件为止
func TryLoadConfig(filePath string, cfgName ...string) (*Config, error) {
	cname := configName
	if len(cfgName) == 1 {
		cname = cfgName[0]
	}

	if err := genericLoad(filePath, cname, Global); err != nil {
		return nil, err
	}
	Logger = Global.GetLogger()
	return Global, nil
}

// TryLoadGenesis 载入创世配置文件。创世配置文件只在节点初始化时调用一次
// 如果传入的文件路径为空，会主动寻找可用配置。查找顺序为
// filePath > ${HOME}/.bolaxy > $ENV[BOLAXYDIR] > WORKDIR > EXE RUN PATH
// 显式提供配置文件路径，会直接处理，非显式提供配置文件路径，直到找到正确的配置文件为止
func TryLoadGenesis(filePath string, genName ...string) (*Genesis, error) {
	gname := genesisName
	if len(genName) == 1 {
		gname = genName[0]
	}

	return TryLoadGenesisWithName(filePath, gname)
}

func TryLoadGenesisWithName(filePath, name string) (*Genesis, error) {
	var genesis Genesis
	if err := genericLoad(filePath, name, &genesis); err != nil {
		return nil, err
	}

	return &genesis, nil
}

func genericLoad(filePath, fileName string, rawVal interface{}) (err error) {
	if len(filePath) > 0 {
		return load(filePath, fileName, rawVal)
	}
	if err = load(common.Home(DefaultHomeBase), fileName, rawVal); err == nil {
		return nil
	}

	if err = load(common.Env(DefaultEnv), fileName, rawVal); err == nil {
		return nil
	}

	if err = load(common.WorkDir(), fileName, rawVal); err == nil {
		return nil
	}

	if err = load(common.ExeDir(), fileName, rawVal); err == nil {
		return nil
	}

	return err
}

func load(filePath, fileName string, rawVal interface{}) error {
	if len(filePath) == 0 {
		return FileNotFound
	}

	viper.SetConfigName(fileName)
	viper.AddConfigPath(filePath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(rawVal); err != nil {
		return err
	}

	return nil
}
