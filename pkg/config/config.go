package config

import (
	"Bilibili_Downloader/internal/paths"
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	VersionUpdateApi string `json:"VersionUpdateApi"`
	CurrentVersion   string `json:"CurrentVersion"`
}

const CurrentVersion = `v1.4.2`

func configFilePath() string {
	return paths.Get().ConfigFilePath()
}

func CreateConfig() error {
	p := paths.Get()
	if err := os.MkdirAll(p.ConfigDir, 0755); err != nil {
		return err
	}

	config := Config{
		VersionUpdateApi: "https://api.github.com/repos/KUAPT/Bilibili_Downloader/releases/latest",
		CurrentVersion:   CurrentVersion,
	}

	file, err := os.Create(configFilePath())
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("关闭文件失败：", err)
		}
	}()

	return json.NewEncoder(file).Encode(config)
}

func ReadConfig() (Config, error) {
	var config Config
	file, err := os.Open(configFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			if err := CreateConfig(); err != nil {
				return config, err
			}
			file, err = os.Open(configFilePath())
			if err != nil {
				return config, err
			}
		} else {
			return config, err
		}
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("关闭文件失败：", err)
		}
	}()

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}
