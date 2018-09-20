package utils

import (
	"os"
	"os/user"
	"path/filepath"
	"sync"
)

type Configure struct {
	BaseDir        	string
	LevelDBDir     	string
	BlocksDir	string
	LogFileName    	string
	CmdServicePort 	string
	CurrentVersion 	string
}

const cmdServicePort = "6080"
const currentVersion = "0.0.1"

var config *Configure
var onceConf sync.Once

func GetConfig() *Configure {
	onceConf.Do(func() {
		config = initConfig()
	})

	return config
}

func initConfig() *Configure {

	baseDir := getBaseDir()
	if _, ok := FileExists(baseDir); ok == false {
		err := os.Mkdir(baseDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	levelDBDir 	:= filepath.Join(baseDir, string(filepath.Separator), "dataStore")
	blockStoreDir 	:= filepath.Join(baseDir, string(filepath.Separator), "blocks")
	logFileName 	:= filepath.Join(baseDir, string(filepath.Separator), "nbs.log")

	return &Configure{
		BaseDir:        baseDir,
		LevelDBDir:     levelDBDir,
		BlocksDir:     	blockStoreDir,
		LogFileName:    logFileName,
		CmdServicePort:	cmdServicePort,
		CurrentVersion: currentVersion,
	}
}

func getBaseDir() string {

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	baseDir := filepath.Join(usr.HomeDir, string(filepath.Separator), ".nbs")

	return baseDir
}
