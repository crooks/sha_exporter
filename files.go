package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"github.com/Masterminds/log-go"
)

func fileHash(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func iterFiles() (success, fail, missing int) {
	for k, v := range cfg.Files {
		fileHash, err := fileHash(v.Path)
		if err != nil {
			log.Warnf("Hashing of %s failed with: %v", k, err)
			prom.fileRead.WithLabelValues(k).Set(0)
			missing++
			continue
		}
		conforms := v.Hash == fileHash
		log.Debugf("Processing file %s with path %s. Collision=%t", k, v.Path, conforms)
		prom.fileRead.WithLabelValues(k).Set(1)
		prom.fileSHA.WithLabelValues(k).Set(bool2Float(conforms))
		if conforms {
			success++
		} else {
			fail++
		}
	}
	return
}
