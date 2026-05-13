package config

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/google/uuid"
)

type Config struct {
	RobotID         string `json:"robot_id"`
	EVMPayeeAddress string `json:"evm_payee_address"`
	Price           string `json:"price"`
	Network         string `json:"network"`
}

var (
	priceRegex   = regexp.MustCompile(`^\$\d+(\.\d+)?$`)
	networkRegex = regexp.MustCompile(`^[a-z0-9]{3,8}:[-_a-zA-Z0-9]{1,32}$`)
)

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if cfg.RobotID == "" {
		cfg.RobotID = uuid.NewString()
	}

	if cfg.Price == "" {
		cfg.Price = "$0.001"
	}
	if !priceRegex.MatchString(cfg.Price) {
		return nil, fmt.Errorf("invalid price format: %q, expected format like $0.001", cfg.Price)
	}

	if cfg.Network == "" {
		cfg.Network = "eip155:8453" // Base mainnet CAIP-2 ID
	}
	if !networkRegex.MatchString(cfg.Network) {
		return nil, fmt.Errorf("invalid network format: %q, expected format like eip155:8453", cfg.Network)
	}

	if cfg.EVMPayeeAddress == "" {
		return nil, fmt.Errorf("evm_payee_address is required")
	}

	return &cfg, nil
}
