package configreader

// Config config using on service
var Config *ConfigData

func init() {
	Config = &ConfigData{}
}

type ConfigData struct {
	Logger struct {
		Mode      string `mapstructure:"mode"`
		Debug     bool   `mapstructure:"debug"`
		Sensitive bool   `mapstructure:"sensitive"`
	} `mapstructure:"logger"`
	HTTP struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"http"`
	Runtime struct {
		MaxProcs int `mapstructure:"maxProcs"`
	} `mapstructure:"runtime"`
	Postgresql struct {
		Host         string `mapstructure:"host"`
		Port         int    `mapstructure:"port"`
		Username     string `mapstructure:"username"`
		Password     string `mapstructure:"password"`
		DatabaseName string `mapstructure:"databaseName"`
		Schema       string `mapstructure:"schema"`
		MaxIdleConns int    `mapstructure:"maxIdleConns"`
		MaxOpenConns int    `mapstructure:"maxOpenConns"`
		MaxLifetime  string `mapstructure:"maxLifetime"`
		Parameters   string `mapstructure:"parameters"`
		LogSQL       bool   `mapstructure:"logSql"`
		AutoMigrate  bool   `mapstructure:"autoMigrate"`
	} `mapstructure:"postgresql"`
}