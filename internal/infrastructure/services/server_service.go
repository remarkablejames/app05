package services

import "time"

type ServerService struct {
	Version     string
	Environment string
	StartTime   int64
}

func NewServerStatusService(version, env string) *ServerService {
	return &ServerService{
		Version:     version,
		Environment: env,
		StartTime:   time.Now().Unix(),
	}
}

func (hs *ServerService) GetHealthStatus() map[string]string {
	uptime := time.Since(time.Unix(hs.StartTime, 0)).String()
	timestamp := time.Now().Format(time.RFC3339)

	return map[string]string{
		"status":      "healthy",
		"uptime":      uptime,
		"version":     hs.Version,
		"environment": hs.Environment,
		"timestamp":   timestamp,
	}
}
