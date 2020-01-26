package auth

import (
	"context"
	"fmt"

	"github.com/kcarretto/paragon/ent"
)

func GetUser(ctx context.Context) *ent.User {
	if v := ctx.Value(userContextKey); v != nil {
		if usr, ok := v.(*ent.User); ok {
			return usr
		}
		panic(fmt.Errorf("Received non-user value for user context key: %v", v))
	}

	return nil
}

// func CreateUserSession(ctx context.Context, user *ent.User) *Context {
// 	if user.SessionToken == "" {
// 		token := NewSecret(SessionTokenLength)
// 		user = user.Update().SetSessionToken(string(token)).SaveX(ctx)
// 	}
// 	return &Context{
// 		Context:      ctx,
// 		userID:       user.ID,
// 		user:         user,
// 		sessionToken: Secret(user.SessionToken),
// 	}
// }
