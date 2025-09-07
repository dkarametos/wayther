package main

// ConfigProvider is an interface for loading configuration.
type ConfigProvider interface {
	LoadConfig(configPath ConfigPath) (*Config, error)
}

// FileConfigProvider is the real implementation of ConfigProvider that reads from the file system.
type FileConfigProvider struct{}

// LoadConfig loads the application configuration from the default path or a custom path.
// It merges configurations if a custom path is provided.
// It returns a Config struct or an error if loading/creating fails.
func (p *FileConfigProvider) LoadConfig(configPath ConfigPath) (*Config, error) {

	//Load or Create the defaultConfig
	config, err := LoadOrCreateConfig(configPath.DefConf, true)
	if err != nil {
		return nil, err
	}
	config.SetDefaults()

	// If there is a custom config specified
	if configPath.isCustom() {
		customConfig, err := LoadOrCreateConfig(configPath.Custom, false)
		if err != nil {
			return nil, err
		}
		config.MergeConfigs(customConfig)
	}

	return config, nil
}
