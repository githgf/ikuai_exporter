package main

import (
	"encoding/json"
	"github.com/alexflint/go-arg"
	"github.com/githgf/ikuai"
	"github.com/githgf/ikuai-exporter/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Ikuai         string `arg:"env:IK_URL" help:"iKuai URL" default:"http://124.221.237.7:38243"`
	IkuaiUsername string `arg:"env:IK_USER" help:"iKuai username" default:"admin"`
	IkuaiPassword string `arg:"env:IK_PWD" help:"iKuai password" default:"hello12345"`
	Debug         bool   `arg:"env:DEBUG" help:"iKuai 开启 debug 日志" default:"false"`
	InsecureSkip  bool   `arg:"env:SKIP_TLS_VERIFY" help:"是否跳过 iKuai 证书验证" default:"true"`
}

var (
	version   string
	buildTime string
)

func main() {
	config := &Config{}
	arg.MustParse(config)

	i := ikuai.NewIKuai(config.Ikuai, config.IkuaiUsername, config.IkuaiPassword, config.InsecureSkip, true)

	if config.Debug {
		i.Debug()
	}

	pkg.LoadAll(i)

	registry := prometheus.NewRegistry()

	registry.MustRegister(pkg.NewIKuaiExporter(i))

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))
	http.HandleFunc("/allVlan", func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(pkg.GetAllInf())
		w.Write(b)
	})
	go func() {
		time.Sleep(1 * time.Minute)
		pkg.StartLoadIkuaiAsync(i)
	}()
	log.Printf("exporter %v started at :9090", version)

	log.Fatal(http.ListenAndServe(":9090", nil))
}
