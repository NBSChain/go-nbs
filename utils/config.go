package utils

import (
	"os"
	"os/user"
	"path/filepath"
	"sync"
)

type Configure struct {
	BaseDir           string
	LevelDBDir        string
	BlocksDir         string
	ShardFun          string
	LogFileName       string
	CmdServicePort    string
	CurrentVersion    string
	SysTimeFormat     string
	BootStrapPeers    []string
	NatServerPort     int
	NatChanSerPort    int
	NatServerIP       []string
	GossipBootStrapIP []string
	GossipCtrlPort    int
}

const (
	cmdServicePort   = "6080"
	currentVersion   = "0.1.4"
	natServerPort    = 11001
	natChanSerPort   = 11002
	gossipCtrlPort   = 12001
	NormalReadBuffer = 1 << 11
)

var (
	natServerIP = []string{
		"192.168.103.101",
		"192.168.103.97",
		"192.168.30.12",
		"52.8.190.235",
		"47.52.172.234",
	}

	gossipContracts = []string{
		"192.168.103.101",
		"192.168.103.97",
		"192.168.30.12",
		"52.8.190.235",
		"47.52.172.234",
	}

	defaultBootstrapAddresses = []string{
		"",
	}

	config   *Configure
	onceConf sync.Once
)

//TODO:: config to local storage.
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

	levelDBDir := filepath.Join(baseDir, string(filepath.Separator), "dataStore")
	blockStoreDir := filepath.Join(baseDir, string(filepath.Separator), "blocks")
	logFileName := filepath.Join(baseDir, string(filepath.Separator), "nbs.log")

	return &Configure{
		BaseDir:           baseDir,
		LevelDBDir:        levelDBDir,
		BlocksDir:         blockStoreDir,
		ShardFun:          "/repo/flatfs/shard/v1/next-to-last/2",
		LogFileName:       logFileName,
		CmdServicePort:    cmdServicePort,
		CurrentVersion:    currentVersion,
		SysTimeFormat:     "2006-01-02 15:04:05",
		BootStrapPeers:    defaultBootstrapAddresses,
		NatServerPort:     natServerPort,
		NatChanSerPort:    natChanSerPort,
		NatServerIP:       natServerIP,
		GossipBootStrapIP: gossipContracts,
		GossipCtrlPort:    gossipCtrlPort,
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
