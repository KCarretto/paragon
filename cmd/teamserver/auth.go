// +build !dev

package main

import (
	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/pkg/auth"
	"github.com/kcarretto/paragon/pkg/service"
	"go.uber.org/zap"
)

func getAuthenticator(_ *zap.Logger, graph *ent.Client) service.Authenticator {
	return auth.UserAuthenticator{Graph: graph}
}
