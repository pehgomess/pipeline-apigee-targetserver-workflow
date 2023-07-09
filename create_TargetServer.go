package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/apigee/v1"
	"google.golang.org/api/option"
)

func readFile(targetServerConfigFile string) ([]byte, error) {
	content, err := os.ReadFile(targetServerConfigFile)
	if err != nil {
		return nil, err
	}

	return content, nil
}

type TargetServerConfig struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
	SSL  bool   `json:"SSL"`
}

func main() {
	serviceAccountFile := "admapi.json"

	ctx := context.Background()

	serviceAccountJSON, err := os.ReadFile(serviceAccountFile)
	if err != nil {
		log.Fatalf("Erro ao carregar as credenciais de Service Account %v", err)
	}

	credentials, err := google.CredentialsFromJSON(ctx, serviceAccountJSON, apigee.CloudPlatformScope)
	if err != nil {
		log.Fatalf("Erro ao carregar as credenciais da Service Account: %v", err)
	}

	service, err := apigee.NewService(ctx, option.WithCredentials(credentials))
	if err != nil {
		log.Fatalf("Erro ao criar o cliente do Apigee: %v", err)
	}

	targetServerConfigFile := "target-server-config.json"
	targetServerConfigJSON, err := readFile(targetServerConfigFile)
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo: %v", err)
	}

	var targetServerConfig TargetServerConfig
	err = json.Unmarshal(targetServerConfigJSON, &targetServerConfig)
	if err != nil {
		log.Fatalf("Erro ao decodificar o arquivo de configuração do Target Server: %v", err)
	}

	req := &apigee.GoogleCloudApigeeV1TargetServer{
		Name:      targetServerConfig.Name,
		Host:      targetServerConfig.Host,
		Port:      int64(targetServerConfig.Port),
		Protocol:  "HTTP",
		IsEnabled: true,
		SSLInfo:   &apigee.GoogleCloudApigeeV1TlsInfo{Enabled: targetServerConfig.SSL},
	}

	createRequest := service.Organizations.Environments.Targetservers.Create("organizations/dock-apigee-nonprod/environments/dev", req)
	resp, err := createRequest.Do()
	if err != nil {
		log.Fatalf("Erro ao criar o Target Server: %v", err)
	}

	responseJSON, _ := json.MarshalIndent(resp, "", " ")
	fmt.Printf("Target Server criado com sucesso:\n%s\n", string(responseJSON))

}
