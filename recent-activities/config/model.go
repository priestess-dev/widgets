package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"time"
)

var (
	reVar = regexp.MustCompile(`\${(\w+)}`)
)

type StringOrEnv string

func (s *StringOrEnv) UnmarshalYAML(value *yaml.Node) error {
	var ss string
	if err := value.Decode(&ss); err != nil {
		return err
	}
	if match := reVar.FindStringSubmatch(ss); len(match) > 0 {
		*s = StringOrEnv(os.Getenv(match[1]))
	} else {
		*s = StringOrEnv(ss)
	}
	return nil
}

type Config struct {
	Host   string `yaml:"host"` // host, default: localhost
	Port   int    `yaml:"port"` // port, default: 8080
	Github *struct {
		ClientID     StringOrEnv   `yaml:"client_id,omitempty"`     // client id of github oauth app
		ClientSecret StringOrEnv   `yaml:"client_secret,omitempty"` // client secret of github oauth app
		RedirectURL  StringOrEnv   `yaml:"redirect_url,omitempty"`  // redirect url of github oauth app
		Scope        []StringOrEnv `yaml:"scope,omitempty"`         // scope of github oauth app
		Username     StringOrEnv   `yaml:"username"`                // github username
	} `yaml:"github"`
	Redis *struct {
		Addr         StringOrEnv   `yaml:"address,omitempty"`   // address of redis server (default: localhost:6379)
		Password     StringOrEnv   `yaml:"password,omitempty"`  // redis password
		DB           int           `yaml:"db,omitempty"`        // redis db (default: 0)
		Prefix       StringOrEnv   `yaml:"prefix,omitempty"`    // redis key prefix (default: "")
		KeySeparator StringOrEnv   `yaml:"separator,omitempty"` // redis key separator (default: "")
		TTL          time.Duration `yaml:"ttl,omitempty"`       // redis key ttl, only int is supported (default: 3600 seconds)
	}
}

// todo: hot load

func Validate(config *Config) error {
	if config.Host == "" {
		config.Host = "localhost"
	}
	if config.Port == 0 {
		config.Port = 8080
	}
	return nil
}

func LoadYAML(path string) (*Config, error) {
	c := &Config{}
	file, err := os.ReadFile(path)
	if err != nil {
		return &Config{}, err
	}
	// marshall yaml to config
	err = yaml.Unmarshal(file, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func NewConfig(configPath string) *Config {
	config, err := LoadYAML(configPath)
	if err != nil {
		panic(err)
	}
	return config
}
