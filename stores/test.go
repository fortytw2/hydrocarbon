package stores

import (
	"context"
	"testing"

	"github.com/fortytw2/kiasu"
)

// Test tests a given store for basic compliance
func Test(s kiasu.Store, t *testing.T) {
	t.Parallel()

	var users = []struct {
		Email    string
		Password string
	}{
		{"luke@puke.com", "iamaJedi2319-%"},
		{"george@lucas.org", "IamNotMyFather2320"},
		{"george23@lucas.org", "IamNotMyFatheqwer2320"},
		{"georg432e@lucas.org", "IamNotMyFwadaather2320"},
	}

	for _, u := range users {
		confirm, err := s.CreateUser(context.Background(), nil, u.Email, u.Password)
		if err != nil {
			t.Fatal(err)
		}

		at, err := s.ActivateUser(context.Background(), confirm)
		if err != nil {
			t.Fatal(err)
		}

		seshs, err := s.GetActiveSessions(context.Background(), at, &kiasu.Pagination{
			Page:     1,
			PageSize: 10,
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(seshs) != 1 {
			t.Fatal("more than 1 sesh")
		}
	}
}
