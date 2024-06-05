package env

import (
	"fmt"
	"os"
	"strings"
)

const (
	ModeName       = "ENV"
	ServiceNameKey = "SERVICE_NAME"
)

const (
	DevMode  = "dev"
	TestMode = "test"
	ProdMode = "prod"
)

var environment = DevMode

func init() {
	SetEnv(os.Getenv(ModeName))
}

func SetEnv(value string) {
	switch strings.ToLower(value) {
	case DevMode:
		environment = DevMode
	case TestMode:
		environment = TestMode
	case ProdMode:
		environment = ProdMode
	}
	fmt.Printf("Running in \"%v\" mode. \n", environment)
	fmt.Println("- using env:   export ENV=prod")
	fmt.Println("- using code:  env.SetEnv(env.ProdMode)")
	fmt.Println()
}

func Environment() string {
	return environment
}

func GetServiceName() string {
	return os.Getenv(ServiceNameKey)
}
