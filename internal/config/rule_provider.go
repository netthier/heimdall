package config

type RuleProviders struct {
	File *FileBasedRuleProviderConfig `koanf:"file"`
}

type FileBasedRuleProviderConfig struct {
	Src   string `koanf:"src"`
	Watch bool   `koanf:"watch"`
}
