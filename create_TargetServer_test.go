package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/apigee/v1"
)

type MockTargetServersCreateCall struct {
	Parent       string
	TargetServer *apigee.GoogleCloudApigeeV1TargetServer
	Response     *apigee.GoogleCloudApigeeV1TargetServer
	Err          error
}

func (m *MockTargetServersCreateCall) SetParent(parent string) *MockTargetServersCreateCall {
	m.Parent = parent
	return m
}

func (m *MockTargetServersCreateCall) SetTargetServer(targetServer *apigee.GoogleCloudApigeeV1TargetServer) *MockTargetServersCreateCall {
	m.TargetServer = targetServer
	return m
}

func (m *MockTargetServersCreateCall) Do() (*apigee.GoogleCloudApigeeV1TargetServer, error) {
	return m.Response, m.Err
}

type MockService struct {
	Organizations struct {
		Environments struct {
			Targetservers struct {
				CreateCall *MockTargetServersCreateCall
			}
		}
	}
}

func createTargetServer(ctx context.Context, service *MockService, parent string, targetServerConfigJSON []byte) error {
	var targetServerConfig TargetServerConfig
	err := json.Unmarshal(targetServerConfigJSON, &targetServerConfig)
	if err != nil {
		return err
	}

	req := &apigee.GoogleCloudApigeeV1TargetServer{
		Name:      targetServerConfig.Name,
		Host:      targetServerConfig.Host,
		Port:      int64(targetServerConfig.Port),
		Protocol:  "HTTP",
		IsEnabled: true,
		SSLInfo:   &apigee.GoogleCloudApigeeV1TlsInfo{Enabled: targetServerConfig.SSL},
	}

	createRequest := service.Organizations.Environments.Targetservers.CreateCall
	createRequest.SetParent(parent)
	createRequest.SetTargetServer(req)

	resp, err := createRequest.Do()
	if err != nil {
		return err
	}

	_ = resp

	return nil
}

var targetServerConfigJSON []byte

func TestCreateTargetServer(t *testing.T) {
	targetServerConfig := TargetServerConfig{
		// Name: "my-target-exxxx",
		// Host: "www.google.com",
		// Port: 443,
		SSL: true || false,
	}

	portList := []int{443, 8443, 8080, 80}

	for _, port := range portList {
		targetServerConfig.Port = port

		targetServerConfigJSON, _ = json.Marshal(targetServerConfig)
	}

	service := &MockService{}
	createCall := &MockTargetServersCreateCall{}
	service.Organizations.Environments.Targetservers.CreateCall = createCall

	err := createTargetServer(context.Background(), service, "organizations/dock-apigee-nonprod/environments/dev", targetServerConfigJSON)
	if err != nil {
		t.Errorf("Erro ao criar o Target Server: %v", err)
	}

	assert.Equal(t, "organizations/dock-apigee-nonprod/environments/dev", createCall.Parent)
	assert.NotNil(t, createCall.TargetServer)
	assert.Equal(t, targetServerConfig.Name, createCall.TargetServer.Name)
	assert.Equal(t, targetServerConfig.Host, createCall.TargetServer.Host)
	assert.Equal(t, int64(targetServerConfig.Port), createCall.TargetServer.Port)
	assert.Equal(t, true, createCall.TargetServer.IsEnabled)
	assert.NotNil(t, createCall.TargetServer.SSLInfo)
	assert.Equal(t, true, createCall.TargetServer.SSLInfo.Enabled)
}

func TestReadFile(t *testing.T) {
	filename := "target-server-config.json"

	content, err := readFile(filename)
	if err != nil {
		t.Errorf("Erro ao ler o arquivo: %v", err)
	}

	if len(content) == 0 {
		t.Errorf("Conteúdo do arquivo está vazio")
	}
}
