package teamserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kcarretto/paragon/ent"
	entCred "github.com/kcarretto/paragon/ent/credential"
	entTar "github.com/kcarretto/paragon/ent/target"
	"github.com/kcarretto/paragon/pkg/auth"
	"github.com/kcarretto/paragon/pkg/service"
	"go.uber.org/zap"
)

// InitService provides HTTP handlers for the Init.
type InitService struct {
	Log   *zap.Logger
	Graph *ent.Client
	Auth  service.Authenticator
}

type credential struct {
	Kind      string
	Principal string
	Secret    string
}

type target struct {
	Name        string
	PrimaryIP   string
	PublicIP    string
	OS          string `json:"os"`
	Tags        []string
	Credentials []credential
}

// HTTP registers http handlers for the Init.
func (svc *InitService) HTTP(router *http.ServeMux) {
	upload := &service.Endpoint{
		Log:           svc.Log.Named("schema"),
		Authenticator: svc.Auth,
		Authorizer:    auth.NewAuthorizer().IsActivated(),
		Handler:       service.HandlerFn(svc.HandleSchemaUpload),
	}

	router.Handle("/init/", upload)
}

// HandleSchemaUpload is an http.HandlerFunc which parses JSON and creates the ents.
func (svc InitService) HandleSchemaUpload(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var schema []target
	err = json.Unmarshal(body, &schema)
	if err != nil {
		return err
	}

	tagSet := make(map[string]int)
	// embedded map for same cred per target
	credentialSet := make(map[credential]map[string]int)

	// first pass for credentials and tags
	for _, tar := range schema {
		// go through tags
		for _, tag := range tar.Tags {
			if _, exists := tagSet[tag]; !exists {
				t, err := svc.Graph.Tag.Create().
					SetName(tag).
					Save(ctx)
				if err != nil {
					return err
				}
				tagSet[tag] = t.ID
			}
		}

		// go through credentials
		for _, cred := range tar.Credentials {
			if _, exists := credentialSet[cred]; !exists {
				credentialSet[cred] = make(map[string]int)
			}
			if _, exists := credentialSet[cred][tar.Name]; !exists {
				c, err := svc.Graph.Credential.Create().
					SetKind(entCred.KindPassword). //TODO: Make work for other creds
					SetPrincipal(cred.Principal).
					SetSecret(cred.Secret).
					Save(ctx)
				if err != nil {
					return err
				}
				credentialSet[cred][tar.Name] = c.ID
			}
		}
	}

	// second pass for targets
	for _, tar := range schema {
		var creds []int
		var tags []int

		for _, t := range tar.Tags {
			tags = append(tags, tagSet[t])
		}
		for _, c := range tar.Credentials {
			creds = append(creds, credentialSet[c][tar.Name])
		}

		var os entTar.OS

		switch tar.OS {
		case entTar.OSLINUX.String():
			os = entTar.OSLINUX
			break
		case entTar.OSWINDOWS.String():
			os = entTar.OSWINDOWS
			break
		case entTar.OSBSD.String():
			os = entTar.OSBSD
			break
		case entTar.OSMACOS.String():
			os = entTar.OSMACOS
			break
		default:
			return fmt.Errorf("OS string passed did not conform to OS ENUM")
		}

		_, err := svc.Graph.Target.Create().
			SetName(tar.Name).
			SetPrimaryIP(tar.PrimaryIP).
			SetPublicIP(tar.PublicIP).
			SetOS(os).AddTagIDs(tags...).
			AddCredentialIDs(creds...).
			Save(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
