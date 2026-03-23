package config

type ServerConfig struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type StoresConfig struct {
	Servers []ServerConfig `json:"servers"`
}
