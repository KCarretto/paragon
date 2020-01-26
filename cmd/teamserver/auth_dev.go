// +build dev

package main

import (
	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/pkg/auth"
	"github.com/kcarretto/paragon/pkg/service"

	"go.uber.org/zap"
)

func getAuthenticator(logger *zap.Logger, graph *ent.Client) service.Authenticator {
	logger.Warn("Developer mode, authentication requirement disabled!")
	return auth.DevAuthenticator{Graph: graph}
}
