package setting

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	sections = make(map[string]interface{})
)

type Setting struct {
	vp *viper.Viper
}

func (s *Setting) ReloadAllSection() error {
	for k, v := range sections {
		if err := s.ReadSection(k, v); err != nil {
			return err
		}
	}

	return nil
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	if err := s.vp.UnmarshalKey(k, v); err != nil {
		return err
	}

	if _, ok := sections[k]; !ok {
		sections[k] = v
	}

	return nil
}

func (s *Setting) WatchSettingChange() {
	go func(st *Setting) {
		st.vp.WatchConfig()
		st.vp.OnConfigChange(func(in fsnotify.Event) {
			_ = s.ReloadAllSection()
		})
	}(s)
}

func rootPath(path string) string {
	if path != "" {
		return path
	}

	path = reflect.TypeOf(Setting{}).PkgPath()
	var (
		wd  string
		err error
	)
	if wd, err = os.Getwd(); err != nil {
		return ""
	}
	splitItems := strings.Split(path, "/")
	splitItemsLen := len(splitItems)
	pkgRootPath := fmt.Sprintf("/%s", strings.Join(splitItems[1:splitItemsLen], "/"))
	return strings.TrimSuffix(wd, pkgRootPath)
}

func NewSetting(path string) (*Setting, error) {
	root := rootPath(path)
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")

	vp.AddConfigPath(fmt.Sprintf("%s", root))
	if err := vp.ReadInConfig(); err != nil {
		// if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// 	return nil, fmt.Errorf("%v", a ...interface{})
		// }
		return nil, err

	}

	s := &Setting{vp}
	s.WatchSettingChange()
	return s, nil
}

type SeverSettingS struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type LogSettingS struct {
	Level      int
	Format     string
	Output     string
	OutputFile string
}

type MemoryCacheSettingS struct {
	Expire           time.Duration
	HardMaxCacheSize int
}
