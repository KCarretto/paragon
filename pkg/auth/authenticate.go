package auth

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/ent/event"
	"github.com/kcarretto/paragon/ent/service"
	"github.com/kcarretto/paragon/ent/tag"
)

const (
	SessionTokenLength = 256
)

type authKey string

var userContextKey authKey = "user"
var serviceContextKey authKey = "service"
var svcTagPrefix string = "svc-"

// GetUser from the context, returns nil for non-user contexts.
func GetUser(ctx context.Context) *ent.User {
	if v := ctx.Value(userContextKey); v != nil {
		if usr, ok := v.(*ent.User); ok {
			return usr
		}
		panic(fmt.Errorf("Received non-user value for user context key: %v", v))
	}

	return nil
}

// GetService from the context, returns nil for non-service contexts.
func GetService(ctx context.Context) *ent.Service {
	if v := ctx.Value(serviceContextKey); v != nil {
		if svc, ok := v.(*ent.Service); ok {
			return svc
		}
		panic(fmt.Errorf("Received non-service value for service context key: %v", v))
	}

	return nil
}

type MultiAuthenticator struct {
	ServiceAuth ServiceAuthenticator
	UserAuth    UserAuthenticator
}

func (auth MultiAuthenticator) Authenticate(w http.ResponseWriter, req *http.Request) (*http.Request, error) {
	req, err := auth.UserAuth.Authenticate(w, req)
	if err != nil {
		return nil, fmt.Errorf("failed user authentication: %w", err)
	}

	req, err = auth.ServiceAuth.Authenticate(w, req)
	if err != nil {
		return nil, fmt.Errorf("failed service authentication: %w", err)
	}

	return req, nil
}

// ServiceAuthenticator parses http requests for service headers and adds service context to the
// request where possible.
type ServiceAuthenticator struct {
	Graph *ent.Client
}

// Authenticate a request by wrapping it's context with the authenticated service identity. It will
// upsert new (unactivated) service identities if the public key is not already registered. Returns
// an error if invalid credentials are provided.
func (auth ServiceAuthenticator) Authenticate(w http.ResponseWriter, req *http.Request) (*http.Request, error) {
	svcName := req.Header.Get(HeaderService)
	pubKeyB64 := req.Header.Get(HeaderIdentity)
	sigB64 := req.Header.Get(HeaderSignature)
	epoch := req.Header.Get(HeaderEpoch)
	if svcName == "" {
		return req, nil
	}

	if pubKeyB64 == "" || sigB64 == "" || epoch == "" {
		return nil, ErrMissingHeaders
	}

	svcQuery := auth.Graph.Service.Query().Where(service.PubKey(pubKeyB64))

	exists, err := svcQuery.Clone().Exist(req.Context())
	if err != nil {
		return nil, err
	}

	// if the service exists
	if exists {
		svc, err := svcQuery.Only(req.Context())
		if err != nil {
			return nil, err
		}
		pubKey, err := base64.StdEncoding.DecodeString(pubKeyB64)
		if err != nil {
			return nil, err
		}
		sig, err := base64.StdEncoding.DecodeString(sigB64)
		if err != nil {
			return nil, err
		}
		valid := ed25519.Verify(pubKey, []byte(epoch), sig)
		if valid {
			return req.WithContext(context.WithValue(req.Context(), serviceContextKey, svc)), nil
		}
		return nil, ErrInvalidSignature
	}

	// now we know the service doesnt exist, but what about the tag?
	tagQuery := auth.Graph.Tag.Query().Where(tag.Name(svcTagPrefix + svcName))

	exists, err = tagQuery.Clone().Exist(req.Context())
	if err != nil {
		return nil, err
	}

	var tag *ent.Tag

	// if a tag exists
	if exists {
		tag, err = tagQuery.Only(req.Context())
		if err != nil {
			return nil, err
		}
		// if not make one
	} else {
		tag, err = auth.Graph.Tag.Create().
			SetName(svcTagPrefix + svcName).
			Save(req.Context())
		if err != nil {
			return nil, err
		}
	}

	svc, err := auth.Graph.Service.Create().
		SetName(svcName).
		SetPubKey(pubKeyB64).
		SetTag(tag).
		Save(req.Context())
	if err != nil {
		return nil, err
	}

	// silent failure is best failure :3
	auth.Graph.Event.Create().
		SetSvcOwner(svc).
		SetKind(event.KindCREATE_SERVICE).
		SetService(svc).
		Save(req.Context())

	return req.WithContext(context.WithValue(req.Context(), serviceContextKey, svc)), nil
}

// UserAuthenticator parses http requests for session cookies and adds user context to the request
// where possible.
type UserAuthenticator struct {
	Graph *ent.Client
}

// Authenticate a request by wrapping it's context with the logged in user. If no user is logged in,
// the original request is returned. Returns an error if it fails to find the logged in user or if
// invalid credentials are provided.
func (auth UserAuthenticator) Authenticate(w http.ResponseWriter, req *http.Request) (*http.Request, error) {
	// Get requested userID, unauthenticated otherwise
	userID, err := parseUserID(req)
	if err != nil {
		return req, nil
	}

	// Get session token, unauthenticated otherwise
	token, err := parseSessionToken(req)
	if err != nil {
		return req, nil
	}

	// Load requested user object, error if no matching user found
	user, err := auth.Graph.User.Get(req.Context(), userID)
	if err != nil {
		ClearUserSession(w)
		return nil, fmt.Errorf("failed to load user: %w", err)
	}
	if user == nil {
		ClearUserSession(w)
		return nil, ErrNotAuthenticated
	}

	// Authenticate as requested user
	if !token.Equals(Secret(user.SessionToken)) {
		ClearUserSession(w)
		return nil, ErrNotAuthenticated
	}

	// User Successfully Authenticated
	return CreateUserSession(w, req, user), nil
}

// parseUserID from the userid cookie.
func parseUserID(req *http.Request) (int, error) {
	userCookie, err := req.Cookie(UserCookieName)
	if err != nil {
		return -1, fmt.Errorf("no user cookie set: %w", err)
	}

	userID, err := strconv.Atoi(userCookie.Value)
	if err != nil {
		return -1, fmt.Errorf("invalid user cookie: %w", err)
	}
	return userID, nil
}

// parseSessionToken from the session token cookie.
func parseSessionToken(req *http.Request) (Secret, error) {
	sessionCookie, err := req.Cookie(SessionCookieName)
	if err != nil {
		return "", fmt.Errorf("no session cookie set: %w", err)
	}

	sessionToken := Secret(sessionCookie.Value)
	if sessionToken == "" {
		return "", fmt.Errorf("invalid session cookie")
	}

	return sessionToken, nil
}
