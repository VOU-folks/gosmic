package services

import (
	"gopkg.in/yaml.v3"
	"gosmic/internal/structs"
	"os"
	"sync"
)

type ConfigService struct {
	path string

	config     *structs.Config
	accessLock sync.Mutex
}

func NewConfigService(path string) *ConfigService {
	return &ConfigService{
		path: path,

		accessLock: sync.Mutex{},
		config:     &structs.Config{},
	}
}

func (s *ConfigService) Init() error {
	return s.Load()
}

func (s *ConfigService) Load() error {
	s.accessLock.Lock()
	defer s.accessLock.Unlock()

	buffer, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(buffer, s.config)
	if err != nil {
		return err
	}

	return nil
}

func (s *ConfigService) Path() string {
	return s.path
}

func (s *ConfigService) GetConfig() structs.Config {
	return *s.config
}

func (s *ConfigService) GetAppsConfig() structs.AppsConfig {
	return s.config.Apps
}

func (s *ConfigService) GetApiConfig() structs.ApiConfig {
	return s.config.Apps.Api
}

func (s *ConfigService) GetDatabaseConfig() structs.DatabaseConfig {
	return s.config.Database
}

func (s *ConfigService) GetStorageConfig() structs.StorageConfig {
	return s.config.Storage
}

func (s *ConfigService) GetOsmConfig() structs.OSMConfig {
	return s.config.OSM
}
