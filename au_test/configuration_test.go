package au_test

import (
	"auservices/utilities"
	"os"
	"testing"
)

func TestDevCreateConfiguration(t *testing.T) {
	//given
	path := "../config/keys.dev.json"
	//when
	_, err := utilities.LoadConfiguration(path)
	if err != nil {
		t.Errorf("Errror: %v", err)
	}
	cfg := utilities.GetConfiguration()
	expectedCfg := utilities.Configuration{
		ApplicationPort:   7777,
		MsgURL:            "nats://localhost:4222",
		MsgClusterID:      "api-cluster",
		DbDriver:          "postgres",
		DbURL:             "postgres://docker:docker@localhost/au?sslmode=disable",
		CategoryChannelID: "category-channel",
		WhoUpdate:         "au-app",
	}
	equals(t, expectedCfg, cfg)
}

func TestProdCreateConfiguration(t *testing.T) {
	//given
	path := "../config/keys.prod.json"
	os.Setenv("AppEnvironment", "PROD")
	os.Setenv("MsgURL", "nats://prodUrl:4223")
	os.Setenv("DbURL", "user:pass@tcp(prodbUrl:3306)/audb")
	//when
	_, err := utilities.LoadConfiguration(path)
	if err != nil {
		t.Errorf("Errror: %v", err)
	}
	cfg := utilities.GetConfiguration()
	expectedCfg := utilities.Configuration{
		ApplicationPort:   7777,
		MsgURL:            "nats://prodUrl:4223",
		MsgClusterID:      "api-cluster",
		DbDriver:          "postgres",
		DbURL:             "user:pass@tcp(prodbUrl:3306)/audb",
		CategoryChannelID: "category-channel",
		WhoUpdate:         "au-app",
	}
	equals(t, expectedCfg, cfg)
}
