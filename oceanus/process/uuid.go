package process

import (
	"encoding/hex"
	uuid "github.com/satori/go.uuid"
)

func NewGlobalUUID() string {
	return hex.EncodeToString(uuid.NewV1().Bytes())
}

func NewNamespaceUUID(name string, space uuid.UUID) string {
	return hex.EncodeToString(uuid.NewV3(namespace, name).Bytes())
}

var namespace = uuid.Must(uuid.FromString("4f31b82c-ca02-432c-afbf-8148c81ccaa2"))
