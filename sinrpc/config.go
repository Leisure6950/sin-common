package sinrpc

type ProtoType string

const (
	PTHttp ProtoType = "http"
	PTRpc  ProtoType = "rpc"
)

type ClientConfig struct {
	ServiceName    string    `toml:"service_name"`
	Endpoints      []string  `toml:"endpoints"`
	ProtoType      ProtoType `toml:"proto"`
	ConnectTimeout int       `toml:"connnect_timeout"`
	ReadTimeout    int       `toml:"read_timeout"`
	WriteTimeout   int       `toml:"write_timeout"`
	MaxIdleConns   int       `toml:"max_idleconn"`
	RetryTimes     int       `toml:"retry_times"`
	SlowTime       int       `toml:"slow_time"`
	DataCenter     string    `toml:"dc,omitempty"`
}
type ServerConfig struct {
	ServiceName string   `toml:"service_name"`
	Port        int      `toml:"port"`
	Tags        []string `toml:"server_tags"`
}
