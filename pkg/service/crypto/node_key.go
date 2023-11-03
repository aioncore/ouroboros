package crypto

import (
	"crypto/ed25519"
	"encoding/json"
	serviceed25519 "github.com/aioncore/ouroboros/pkg/service/crypto/ed25519"
	"github.com/aioncore/ouroboros/pkg/service/utils"
)

type NodeKey struct {
	PrivateKey *ed25519.PrivateKey `json:"private_key"`
}

func LoadNodeKey(filePath string) (*NodeKey, error) {
	jsonBytes, err := utils.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	nodeKey := &NodeKey{}
	err = json.Unmarshal(jsonBytes, nodeKey)
	if err != nil {
		return nil, err
	}
	return nodeKey, nil
}

func SaveNodeKey(nodeKey *NodeKey, filePath string) error {
	jsonBytes, err := json.Marshal(nodeKey)
	if err != nil {
		return err
	}
	err = utils.WriteFile(filePath, jsonBytes, 0666)
	if err != nil {
		return err
	}
	return nil
}

func LoadOrGenerateNodeKey(filePath string) (*NodeKey, error) {
	exist, err := utils.PathExist(filePath)
	if err != nil {
		return nil, err
	}
	if exist {
		return LoadNodeKey(filePath)
	}
	privateKey, err := serviceed25519.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	nodeKey := &NodeKey{
		PrivateKey: &privateKey,
	}
	err = SaveNodeKey(nodeKey, filePath)
	if err != nil {
		return nil, err
	}
	return nodeKey, nil
}
