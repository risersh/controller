package controller

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/risersh/controller/conf"
	"github.com/risersh/controller/rabbitmq"
)

func init() {
	conf.Init()

	multilog.RegisterLogger(multilog.LogMethod("console"), multilog.NewConsoleLogger(&multilog.NewConsoleLoggerArgs{
		Level:  multilog.DEBUG,
		Format: multilog.FormatText,
		FilterDropPatterns: []*string{
			multilog.PtrString("producer"), // Drop rabbitmq producer logs.
		},
	}))

	multilog.RegisterLogger(multilog.LogMethod("elasticsearch"), multilog.NewElasticsearchLogger(&multilog.NewElasticsearchLoggerArgs{
		Level: multilog.DEBUG,
		Config: elasticsearch.Config{
			Addresses: []string{conf.Config.Elasticsearch.URL},
			Username:  conf.Config.Elasticsearch.Username,
			Password:  conf.Config.Elasticsearch.Password,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		Index: "logs-webrtc",
		Mapping: multilog.PtrString(`
		{
			"mappings": {
				"properties": {
					"time": { "type": "date" },
					"level": { "type": "keyword" },
					"group": { "type": "keyword" },
					"message": { "type": "text" },
					"data": { "type": "object" }
				}
			}
		}`),
		FilterDropPatterns: []*string{
			multilog.PtrString("producer"), // Drop rabbitmq producer logs.
		},
	}))
}

func main() {
	go rabbitmq.Setup()

	ch, err := rabbitmq.StartConsuming[rabbitmq.Message[rabbitmq.MessageTypeNewDeployment]]()
	if err != nil {
		log.Fatalf("Failed to start consuming: %v", err)
	}
	go func() {
		for msg := range ch {
			log.Println(msg)
		}
	}()
}
