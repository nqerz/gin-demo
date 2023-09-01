package global

import "github.com/znqerz/gin-demo/pkg/setting"

var (
	LogSetting         *setting.LogSettingS
	MemoryCacheSetting *setting.MemoryCacheSettingS
)

func InitSetting(path string) error {
	setting, err := setting.NewSetting(path)
	if err != nil {
		return err
	}

	err = setting.ReadSection("Log", &LogSetting)
	if err != nil {
		return err
	}

	return nil
}
