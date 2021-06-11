package config

import (
	"flag"
	"os"

	"github.com/0chain/bandwidth_marketplace/code/core/errors"
	"github.com/0chain/bandwidth_marketplace/code/core/provider"
	"gopkg.in/yaml.v2"
)

type (
	// Config represents configs stores in Path.
	Config struct {
		ServerChain ServerChain `yaml:"server_chain"`
		Handler     Handler     `yaml:"handler"`
		Database    Database    `yaml:"db"`
		Logging     Logging     `yaml:"logging"`
		Proxy       Proxy       `yaml:"proxy"`

		MagmaAddress      string `yaml:"magma_address"`       // represents address of Magma service
		GRPCServerTimeout int    `yaml:"grpc_server_timeout"` // in seconds

		CLConfig CommandLineConfig
	}

	// ServerChain represents config options described in "server_chain" section of the config yaml file.
	// ServerChain must be a field of Config struct
	ServerChain struct {
		ID              string `yaml:"id"`
		OwnerID         string `yaml:"owner_id"`
		BlockWorker     string `yaml:"block_worker"`
		SignatureScheme string `yaml:"signature_scheme"`
	}

	// Handler represents config options described in "handler" section of the config yaml file.
	// Handler must be a field of Config struct
	Handler struct {
		RateLimit float64 `yaml:"rate_limit"` // per second
	}

	// Database represents config options described in "db" section of the config yaml file.
	// Database must be a field of Config struct
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		UserName string `yaml:"user"`
		Name     string `yaml:"name"`
		Password string `yaml:"password"`
	}

	// Logging represents config options described in "logging" section of the config yaml file.
	// Database must be a field of Config struct
	Logging struct {
		Level string `yaml:"level"`
	}

	// Proxy represents config options described in "proxy" section of the config yaml file.
	// Proxy must be a field of Config struct
	Proxy struct {
		Terms Terms `yaml:"terms"`
	}

	// Terms represents quality of services options described in "proxy.terms" section of the config yaml file.
	Terms struct {
		MaxPrice int64 `yaml:"max_price"`
		QoS      QoS   `yaml:"qos"`
	}

	// QoS represents quality of services options described in "proxy.terms.qos" section of the config yaml file.
	QoS struct {
		MinDownloadMBPS int64 `yaml:"min_download_mbps"`
		MinUploadMBPS   int64 `yaml:"min_upload_mbps"`
	}

	// CommandLineConfig represents config options passed from the command line.
	CommandLineConfig struct {
		Host           string
		Port           int
		GrpcPort       int
		DeploymentMode Deployment
		KeysFile       string
		LogDir         string
		DBDir          string
	}
)

const (
	// Path is a constant stores path to config file from root application directory.
	Path = "./config/cons-config.yaml"
)

// Read reads configs from config file existing in Path.
//
// If an error occurs during execution, the program terminates with code 2 and the error will be written in os.Stderr.
//
// Read should be used only once while application is starting.
func Read() *Config {
	f, err := os.Open(Path)
	if err != nil {
		errors.ExitErr("err while open config file", err, 2)
	}
	defer func(f *os.File) { _ = f.Close() }(f)

	decoder := yaml.NewDecoder(f)
	cfg := new(Config)
	err = decoder.Decode(cfg)
	if err != nil {
		errors.ExitErr("err while decoding config file", err, 2)
	}

	return cfg
}

// Development returns true if Config.ClConfig.DeploymentMode == Development. Either returns false.
func (cfg Config) Development() bool {
	return cfg.CLConfig.DeploymentMode.Development()
}

// ExtractCommandLineConfigs parses passed command line flags, validates required flags and returns CommandLineConfig
// contains necessary flags.
//
// If an error occurs during execution, the program terminates with code 2 and the error will be written in os.Stderr.
//
// ExtractCommandLineConfigs should be used only once while application is starting.
func ExtractCommandLineConfigs() CommandLineConfig {
	deploymentMode := flag.Int("deployment_mode", 2, "apps mode")
	keysFile := flag.String("keys_file", "", "keys file path")
	logDir := flag.String("log_dir", "", "logs directory")
	port := flag.Int("port", 8080, "port")
	hostname := flag.String("hostname", "", "hostname")
	grpcPort := flag.Int("grpc_port", 8081, "grpc port")
	dbDir := flag.String("db_dir", "", "db_dir")

	flag.Parse()
	validateFlags()

	return CommandLineConfig{
		Host:           *hostname,
		Port:           *port,
		DeploymentMode: Deployment(*deploymentMode),
		KeysFile:       *keysFile,
		LogDir:         *logDir,
		GrpcPort:       *grpcPort,
		DBDir:          *dbDir,
	}
}

// validateFlags validates required flags.
//
// If an error occurs during execution, the program terminates with code 2 and the error will be written in os.Stderr.
func validateFlags() {
	required := []string{
		"hostname",
		"port",
		"db_dir",
		"grpc_port",
	}

	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			errors.ExitMsg("missing required -"+req+" argument/flag", 2)
		}
	}
}

func (db Database) String() string {
	return "host=" + db.Host +
		" port=" + db.Port +
		" user=" + db.UserName +
		" dbname=" + db.Name +
		" password=" + db.Password +
		" sslmode=disable"
}

// Validate checks provider.Terms for compliance with configuration of Terms.
func (tCfg Terms) Validate(terms *provider.Terms) bool {
	switch {
	case terms.Price > tCfg.MaxPrice:
	case terms.QoS.DownloadMBPS < tCfg.QoS.MinDownloadMBPS:
	case terms.QoS.UploadMBPS < tCfg.QoS.MinUploadMBPS:

	default:
		return true
	}

	return false
}
