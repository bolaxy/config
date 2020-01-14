package conf

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

const (
	KeystorePath = "keystore"
	PasswordFile = "password"
)

var (
	defaultEthAPIAddr = "0.0.0.0:8080"
	defaultHeartbeat  = 500 * time.Millisecond
	defaultTCPTimeout = 1000 * time.Millisecond
	defaultCacheSize  = 50000
	defaultSyncLimit  = 1000
	defaultMaxPool    = 2
	defaultPath       = "/opt/runbolaxy/bconfig"
	defaultGenesis    = "genesis.toml"
	defaultKeystore   = "keystore"
	defaultPwdFile    = "password"
	defaultDbFile     = "db"

	Key        *ecdsa.PrivateKey
	GensisData *Genesis
	Peers      *PeerSet
	logMux     sync.Mutex
	logger     *logrus.Entry
	Global     *Config
	Logger     *logrus.Entry
)

func DefaultDataConfig() *DataConfig {
	return &DataConfig{
		DataDir:  defaultPath,
		Genesis:  defaultGenesis,
		Keystore: defaultKeystore,
		PwdFile:  defaultPwdFile,
		DbFile:   defaultDbFile,
	}
}

func DefaultNetConfig() *NetConfig {
	return &NetConfig{
		Heartbeat:  defaultHeartbeat,
		TCPTimeout: defaultTCPTimeout,
		MaxPool:    defaultMaxPool,
		EthAPIAddr: defaultEthAPIAddr,
	}
}

func DefaultConfig() *Config {
	return &Config{
		Self:      "",
		Verbose:   true,
		DataCnf:   DefaultDataConfig(),
		NetCnf:    DefaultNetConfig(),
		CacheSize: defaultCacheSize,
		SyncLimit: defaultSyncLimit,
	}
}

type NetConfig struct {
	Heartbeat   time.Duration `mapstructure:"heartbeat"`
	TCPTimeout  time.Duration `mapstructure:"tcp-timeout"`
	JoinTimeout time.Duration `mapstructure:"join_timeout"`
	MaxPool     int           `mapstructure:"max-pool"`
	EthAPIAddr  string        `mapstructure:"listen"`
}

type DataConfig struct {
	DataDir  string `mapstructure:"datadir"`
	Genesis  string `mapstructure:"genesis"`
	Keystore string `mapstructure:"keystore"`
	PwdFile  string `mapstructure:"pwd"`
	DbFile   string `mapstructure:"db"`
}

type Config struct {
	Self      string      `mapstructure:"self"`
	Verbose   bool        `mapstructure:"verbose"`
	DataCnf   *DataConfig `mapstructure:"datacnf"`
	NetCnf    *NetConfig  `mapstructure:"netcnf"`
	Peerlist  []*Peer     `mapstructure:"peerSet"`
	CacheSize int         `mapstructure:"cache-size"`
	SyncLimit int         `mapstructure:"sync-limit"`
}

type PeerList []*Peer

func (list *PeerList) Marshal() ([]byte, error) {
	var b bytes.Buffer

	enc := json.NewEncoder(&b)

	if err := enc.Encode(list); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (list *PeerList) UnMarshal(data []byte) error {
	b := bytes.NewBuffer(data)

	dec := json.NewDecoder(b) // will read from b

	if err := dec.Decode(list); err != nil {
		return err
	}

	return nil
}

func SelfPeer(alias string, list []*Peer) *Peer {
	for _, p := range list {
		if p.Alias == alias {
			return p
		}
	}

	return nil
}

func OtherPeers(excludeAlias string, list *PeerSet) PeerList {
	ps := make([]*Peer, 0, len(list.Peers)-1)
	for _, p := range list.Peers {
		if p.Alias != excludeAlias {
			ps = append(ps, p)
		}
	}

	return ps
}

func (cnf *Config) SelfPeer() *Peer {
	return SelfPeer(cnf.Self, cnf.Peerlist)
}

func (cnf *Config) OtherPeers() PeerList {
	return OtherPeers(cnf.Self, Peers)
}

func (cnf *Config) GetLogger() *logrus.Entry {
	logMux.Lock()
	defer logMux.Unlock()
	if logger == nil {
		logLevel := "info"
		if cnf.Verbose {
			logLevel = "debug"
		}
		logger = newLogger(logLevel)
	}

	return logger
}

func (cnf *Config) GetDBFile() string {
	return filepath.Join(cnf.DataCnf.DataDir, cnf.DataCnf.DbFile)
}

func (cnf *Config) GetGenesis() string {
	return filepath.Join(cnf.DataCnf.DataDir, cnf.DataCnf.Genesis)
}

func (cnf *Config) GetKey() *ecdsa.PrivateKey {
	return Key
}

func (cnf *Config) GetPeers() *PeerSet {
	return Peers
}

func (cnf *Config) GensisData() *Genesis {
	return GensisData
}

func newLogger(lvl string) *logrus.Entry {
	logger := logrus.New()
	logger.Level = LogLevel(lvl)
	logger.Formatter = new(prefixed.TextFormatter)
	return logger.WithField("prefix", "memberlist")
}

// LogLevel ...
func LogLevel(l string) logrus.Level {
	switch l {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel
	}
}

func init() {
	Global = DefaultConfig()
	Logger = Global.GetLogger()
}
