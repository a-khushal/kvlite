package main

import (
	"encoding/json"
	"os"
)

func (kv *KVStore) SaveToFile(filename string) error {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	dataBytes, err := json.MarshalIndent(kv.data, "", "  ")
	if err != nil {
		return err
	}

	tempFile := filename + ".tmp"
	if err := os.WriteFile(tempFile, dataBytes, 0644); err != nil {
		return err
	}

	return os.Rename(tempFile, filename)
}

func (kv *KVStore) LoadFromFile(filename string) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	dataBytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	tempData := make(map[string]string)
	if err := json.Unmarshal(dataBytes, &tempData); err != nil {
		return err
	}

	kv.data = tempData
	return nil
}
