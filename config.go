package zaplogger

type Config struct {
	ServiceID                string
	StdOutLoggerEnabled      bool
	FileLoggerEnabled        bool
	FileLoggerPath           string
	FileLoggerMaxSize        int // megabytes
	FileLoggerMaxBackups     int // amount
	FileLoggerMaxAge         int // days
	StreamLoggerEnabled      bool
	StreamLoggerAddrs        string
	StreamLoggerUsername     string
	StreamLoggerPassword     string
	StreamLoggerTopic        string
	StreamLoggerBatchSize    int // messages
	StreamLoggerBatchTimeout int // millisexonds
}
