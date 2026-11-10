package model

import (
	"time"

	"github.com/ericls/imgdd/domainmodels"
)

type APIKey struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"keyPrefix"`
	CreatedAt  time.Time  `json:"createdAt"`
	LastUsedAt *time.Time `json:"lastUsedAt,omitempty"`
}

func FromDomainAPIKey(k *domainmodels.APIKey) *APIKey {
	if k == nil {
		return nil
	}
	return &APIKey{
		ID:         k.Id,
		Name:       k.Name,
		KeyPrefix:  k.KeyPrefix,
		CreatedAt:  k.CreatedAt,
		LastUsedAt: k.LastUsedAt,
	}
}
