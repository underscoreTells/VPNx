package config

var ConfigVersions = map[ConfigVersion]func() any{
	CONFIG_VERSION_ONE: func() any {
		var cfg ConfigVersionOne
		return cfg
	},
}
