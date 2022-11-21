package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

type ViperConfig struct {
	App      App   `mapstructure:"app"`
	Kafka    Kafka `mapstructure:"kafka"`
	Mysql    Mysql `mapstructure:"mysql"`
	Log      Log   `mapstructure:"log"`
	propsMap map[string]map[string]interface{}
	lock     sync.RWMutex // 配置文件读写锁
}

type App struct {
	Address      string `mapstructure:"address"`
	ReadTimeout  int64  `mapstructure:"read_timeout"`
	WriteTimeout int64  `mapstructure:"write_timeout"`
	Static       string `mapstructure:"static"`
	AppLogPath   string `mapstructure:"app_log_path"`
}

type Mysql struct {
	Enable   bool   `mapstructure:"enable"`
	User     string `mapstructure:"user"`
	Host     string `mapstructure:"host"`
	Port     int64  `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type Kafka struct {
	Enable bool     `mapstructure:"enable"`
	Hosts  []string `mapstructure:"hosts"`
}

type Log struct {
	Folder      string `mapstructure:"folder"`
	Driver      string `mapstructure:"driver"`
	Formatter   string `mapstructure:"formatter"`
	DateFormat  string `mapstructure:"date_format"`
	Level       string `mapstructure:"level"`
	File        string `mapstructure:"file"`
	RotateSize  int    `mapstructure:"rotate_size"`
	RotateCount int    `mapstructure:"rotate_count"`
	RotateTime  string `mapstructure:"rotate_time"`
	MaxAge      string `mapstructure:"max_age"`
}

func NewViperConfig(...interface{}) (interface{}, error) {

	var cfg ViperConfig
	viper := viper.New()
	//1.设置配置文件路径
	rootDir, _ := os.Getwd()
	viper.SetConfigFile(rootDir + "/config/config.yml")
	//2.配置读取
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	//3.将配置映射成结构体
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	//4.加载属性列表
	cfg.LoadPropsMap()

	//5. 监听配置文件变动,重新解析配置
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println(e.Name)
		if err := viper.Unmarshal(&cfg); err != nil {
			panic(err)
		}
		// 重新load props
		cfg.LoadPropsMap()
	})

	return &cfg, nil
}

func (v *ViperConfig) LoadPropsMap() {
	propsM := map[string]map[string]interface{}{}

	// 2.遍历parent层
	m := loadPropsMap(*v)

	// 3.处理子层
	for k, value := range m {
		propsM[k] = loadPropsMap(value)
	}

	v.propsMap = propsM
}

func loadPropsMap(v interface{}) map[string]interface{} {
	elems := reflect.TypeOf(v)
	values := reflect.ValueOf(v)
	m := make(map[string]interface{}, elems.NumField())

	for i := 0; i < elems.NumField(); i++ {
		// 取出tag name
		pTagName := elems.Field(i).Tag.Get("mapstructure")
		if pTagName == "" {
			continue
		}

		// 存放父层 tag -> interface{} 对象映射
		m[pTagName] = values.Field(i).Interface()
	}

	return m
}

func (v *ViperConfig) IsExist(key string) bool {
	return v.find(key) != nil
}

func (v *ViperConfig) Get(key string) interface{} {
	return v.find(key)
}

func (v *ViperConfig) GetBool(key string) bool {
	return cast.ToBool(v.find(key))
}

func (v *ViperConfig) GetInt(key string) int {
	return cast.ToInt(v.find(key))
}

func (v *ViperConfig) GetFloat64(key string) float64 {
	return cast.ToFloat64(v.find(key))
}

func (v *ViperConfig) GetTime(key string) time.Time {
	return cast.ToTime(v.find(key))
}

func (v *ViperConfig) GetString(key string) string {
	return cast.ToString(v.find(key))
}

func (v *ViperConfig) GetIntSlice(key string) []int {
	return cast.ToIntSlice(v.find(key))
}

func (v *ViperConfig) GetStringSlice(key string) []string {
	return cast.ToStringSlice(v.find(key))
}

func (v *ViperConfig) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(v.find(key))
}

func (v *ViperConfig) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(v.find(key))
}

func (v *ViperConfig) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(v.find(key))
}

func (v *ViperConfig) Load(key string, val interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml",
		Result:  val,
	})
	if err != nil {
		return err
	}

	return decoder.Decode(v.find(key))
}

// 通过path来获取某个配置项
func (v *ViperConfig) find(key string) interface{} {
	v.lock.RLock()
	defer v.lock.RUnlock()

	m := map[string]interface{}{}
	for k, val := range v.propsMap {
		m[k] = val
	}

	return v.searchMap(m, strings.Split(key, "."))
}

// 查找某个路径的配置项
func (v *ViperConfig) searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	// 判断是否有下个路径
	next, ok := source[path[0]]
	if ok {
		// 判断这个路径是否为1
		if len(path) == 1 {
			return next
		}

		// 判断下一个路径的类型
		switch next.(type) {
		case map[interface{}]interface{}:
			// 如果是interface的map，使用cast进行下value转换
			return v.searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			// 如果是map[string]，直接循环调用
			return v.searchMap(next.(map[string]interface{}), path[1:])
		default:
			// 否则的话，返回nil
			return nil
		}
	}
	return nil
}
