package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type Config struct {
	Server struct {
		HTTP struct {
			Addr    string `yaml:"addr"`
			Timeout string `yaml:"timeout"`
		} `yaml:"http"`
		GRPC struct {
			Addr    string `yaml:"addr"`
			Timeout string `yaml:"timeout"`
		} `yaml:"grpc"`
	} `yaml:"server"`

	Data struct {
		Database struct {
			Driver string `yaml:"driver"`
			Source string `yaml:"source"`
		} `yaml:"database"`
		Redis struct {
			Addr         string `yaml:"addr"`
			Password     string `yaml:"password"`
			DB           int    `yaml:"db"`
			ReadTimeout  string `yaml:"read_timeout"`
			WriteTimeout string `yaml:"write_timeout"`
		} `yaml:"redis"`
	} `yaml:"data"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func connectDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("connect db fail: %w", err))
	}
	return db
}

func main() {
	// 1. 读取配置文件
	cfg, err := LoadConfig("/Users/mengfanxing/Documents/go_kratos_project/review-service/configs/config.yaml")
	if err != nil {
		panic(fmt.Errorf("load config fail: %w", err))
	}

	// 2. 从配置里取 DSN
	dsn := cfg.Data.Database.Source

	// 3. 初始化 gorm gen
	g := gen.NewGenerator(gen.Config{
		OutPath:       "/Users/mengfanxing/Documents/go_kratos_project/review-service/internal/data/query",
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	// 4. 复用数据库连接
	g.UseDB(connectDB(dsn))

	// 5. 生成所有表模型
	g.ApplyBasic(g.GenerateAllTable()...)

	// 6. 执行生成代码
	g.Execute()
}
