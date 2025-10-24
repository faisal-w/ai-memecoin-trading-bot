package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Set test environment variables
	os.Setenv("DRY_RUN", "true")
	os.Setenv("AUTO_EXECUTE", "false")
	os.Setenv("WIN_PROBABILITY_THRESHOLD", "0.85")
	
	cfg := LoadConfig()
	
	if !cfg.DryRun {
		t.Error("Expected DryRun to be true")
	}
	
	if cfg.AutoExecute {
		t.Error("Expected AutoExecute to be false")
	}
	
	if cfg.WinProbabilityThreshold != 0.85 {
		t.Errorf("Expected WinProbabilityThreshold to be 0.85, got %f", cfg.WinProbabilityThreshold)
	}
	
	// Check defaults
	if cfg.MaxHoneypotScore != 0.2 {
		t.Errorf("Expected default MaxHoneypotScore to be 0.2, got %f", cfg.MaxHoneypotScore)
	}
	
	// Clean up
	os.Unsetenv("DRY_RUN")
	os.Unsetenv("AUTO_EXECUTE")
	os.Unsetenv("WIN_PROBABILITY_THRESHOLD")
}

func TestGetEnvDefaults(t *testing.T) {
	// Test that defaults work when env vars not set
	value := getEnv("NONEXISTENT_VAR", "default_value")
	if value != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", value)
	}
	
	intValue := getEnvInt("NONEXISTENT_INT", 42)
	if intValue != 42 {
		t.Errorf("Expected 42, got %d", intValue)
	}
	
	floatValue := getEnvFloat("NONEXISTENT_FLOAT", 3.14)
	if floatValue != 3.14 {
		t.Errorf("Expected 3.14, got %f", floatValue)
	}
	
	boolValue := getEnvBool("NONEXISTENT_BOOL", true)
	if !boolValue {
		t.Error("Expected true")
	}
}
