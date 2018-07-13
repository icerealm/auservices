package main

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
	t.Logf("info: %+v", cfg)
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
	t.Logf("info: %+v", cfg)
}
