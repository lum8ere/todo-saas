package auth

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
)

type KeycloakConfig struct {
	IssuerURL string // https://keycloak/realms/todo-saas
	ClientID  string // audience
}

type KeycloakAuth struct {
	verifier *oidc.IDTokenVerifier
}

func NewKeycloakAuth(ctx context.Context, cfg KeycloakConfig) (*KeycloakAuth, error) {
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.ClientID,
	})

	return &KeycloakAuth{verifier: verifier}, nil
}

type Claims struct {
	Subject  string   `json:"sub"`
	Email    string   `json:"email"`
	TenantID string   `json:"tenant_id"`       // свой кастомный клейм
	Roles    []string `json:"roles,omitempty"` // можно маппить из realm_access/...
}

func (a *KeycloakAuth) VerifyToken(ctx context.Context, rawToken string) (*Claims, error) {
	idToken, err := a.verifier.Verify(ctx, rawToken)
	if err != nil {
		return nil, err
	}
	var c Claims
	if err := idToken.Claims(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
