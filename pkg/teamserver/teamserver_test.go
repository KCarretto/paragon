package teamserver_test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/pkg/teamserver"

	_ "github.com/mattn/go-sqlite3"
)

// An Error returned by the GraphQL server
type Error struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

func TestTeamserver(t *testing.T) {
	client := testClient()
	defer client.Close()
	svc := &teamserver.Service{
		Graph: client,
	}
	router := http.NewServeMux()
	svc.HTTP(router)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/graphql", nil)

	router.ServeHTTP(rr, req)

	// fmt.Printf("%s", rr.Body.String())
	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}

}

func testClient() *ent.Client {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}
