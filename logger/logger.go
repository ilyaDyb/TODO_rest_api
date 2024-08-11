package logger

import (
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

func InitLogger(client *elasticsearch.Client) {
	log := logrus.New()
	log.SetOutput(&lumberjack.Logger{
		Filename: "tinder_clone.log",
		MaxSize: 20,
		MaxBackups: 2,
		MaxAge: 28,
		Compress: true,
	})
	log.SetFormatter(&logrus.JSONFormatter{})
	// log.SetLevel(logrus.InfoLevel)
	// log.SetLevel(logrus.FatalLevel)
	log.SetLevel(logrus.DebugLevel)
	// log.SetLevel(logrus.ErrorLevel)

	hook := NewElasticsearchHook(client, "logs")
	log.AddHook(hook)
	Log = log
}

type ElasticsearchHook struct {
    client *elasticsearch.Client
    index  string
}

func NewElasticsearchHook(client *elasticsearch.Client, index string) *ElasticsearchHook {
    return &ElasticsearchHook{client: client, index: index}
}

func (hook *ElasticsearchHook) Fire(entry *logrus.Entry) error {
	js, err := entry.String()
	if err != nil {
		return err
	}

	req := map[string]interface{}{
		"@timestamp": entry.Time.Format(time.RFC3339),
		"message": js,
		"level": entry.Level.String(),
		"fields": entry.Data,
	}
	_, err = hook.client.Index(hook.index, esutil.NewJSONReader(req))
	return err
}

func (hook *ElasticsearchHook) Levels() []logrus.Level {
	return logrus.AllLevels
}