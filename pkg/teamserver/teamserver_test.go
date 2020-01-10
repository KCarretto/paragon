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

func TestTeamserver(t *testing.T) {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	srv := &teamserver.Server{EntClient: client}
	router := teamserver.NewRouter(srv)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/graphql", nil)

	router.ServeHTTP(rr, req)

	// fmt.Printf("%s", rr.Body.String())
	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}

}
