package update

import (
	"Bilibili_Downloader/internal/paths"
	"Bilibili_Downloader/pkg/config"
	"log"
	"os"
)

func HandleUpdateFlag() {
	if len(os.Args) > 1 && os.Args[1] == "--update" {
		if len(os.Args) >= 3 {
			CleanupOldVersion(os.Args[2])
		} else {
			log.Println("更新参数缺失")
		}
	}
}

func EnsureConfigVersion() error {
	currentConfig, err := config.ReadConfig()
	if err != nil {
		return err
	}

	if currentConfig.CurrentVersion != config.CurrentVersion {
		log.Println("内置版本信息与config配置不一致，尝试重建config配置")
		p := paths.Get()
		os.Remove(p.ConfigFilePath())
		if err := config.CreateConfig(); err != nil {
			return err
		}
	}

	return nil
}

func GetUpdateAPIURL() (string, error) {
	cfg, err := config.ReadConfig()
	if err != nil {
		return "", err
	}
	return cfg.VersionUpdateApi, nil
}

func GetCurrentVersion() string {
	return config.CurrentVersion
}
