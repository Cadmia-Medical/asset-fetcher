package main

import (
	"flag"
	"io/ioutil"
	"time"
	"log"
	"fmt"

	"gopkg.in/yaml.v2"
	"asset-fetcher/core"
)

type Config struct {
	AssetNames []string `yaml:"asset_names"`
	Frequency string `yaml:"check_frequency"`
	Bucket string `yaml:"bucket"`
	Region string `yaml:"region"`
	DownloadPath string `yaml:"download_path"`
	SymlinkTarget string `yaml:"symlink_target"`
}

func getConf(path string) *Config {
	var c Config
    yamlFile, err := ioutil.ReadFile(path)
    if err != nil {
        log.Fatalf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &c)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return &c
}


func main() {
	configFilePath := flag.String("config", "config.yml", "path to config file")

	config := getConf(*configFilePath)


	var assetFetchers []*core.AssetRefresher
	for _, name := range config.AssetNames {
		aws, err := core.CreateAwsFetcher(config.Region, config.Bucket)
		if err != nil {
			log.Fatalf("Error creating AWS fetcher")
		}
		assetFetcher := core.AssetRefresher{
			Fetcher: aws,
			AssetName: name,
			LocalTag: "",
			DownloadPath: config.DownloadPath,
			SymlinkTarget: config.SymlinkTarget,
		}
		assetFetchers = append(assetFetchers, &assetFetcher)
	}

	freq, err := time.ParseDuration(config.Frequency)
	if err != nil {
		log.Fatalf("Unable to parse refresh frequency")
	}

	Run(&freq, assetFetchers)

}

func Run(freq *time.Duration, assetFetchers []*core.AssetRefresher) {
	ticker := time.NewTicker(*freq)
	checkAllAssets(assetFetchers)

	fmt.Println("Initializing refreshers")
	for _ = range ticker.C {
		checkAllAssets(assetFetchers)
	}
}

func checkAllAssets(assetFetchers []*core.AssetRefresher) {
	for _, assetChecker := range assetFetchers {
		fmt.Println("Checking asset freshness")
		err := assetChecker.Refresh()
		if err != nil {
			log.Fatalf("Error refreshing assets")
		}
	}
}
