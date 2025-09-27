package usecases

import (
	"share-screen/pkg/domain/entities"
	"share-screen/pkg/domain/interfaces"
)

// ServerInfoUseCase implements the server info use case interface
type ServerInfoUseCase struct {
	networkService interfaces.NetworkService
	stunServer     string
	version        string
}

// NewServerInfoUseCase creates a new server info use case
func NewServerInfoUseCase(networkService interfaces.NetworkService, stunServer, version string) *ServerInfoUseCase {
	return &ServerInfoUseCase{
		networkService: networkService,
		stunServer:     stunServer,
		version:        version,
	}
}

// GetServerInfo returns server information including network details
func (uc *ServerInfoUseCase) GetServerInfo(host string) (*entities.ServerInfo, error) {
	return &entities.ServerInfo{
		Host:       host,
		LANIP:      uc.networkService.GetLANIP(),
		STUNServer: uc.stunServer,
		Version:    uc.version,
	}, nil
}
