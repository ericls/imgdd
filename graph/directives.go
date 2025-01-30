package graph

import (
	"context"
	"fmt"

	"github.com/ericls/imgdd/captcha"
	"github.com/ericls/imgdd/identity"

	"github.com/99designs/gqlgen/graphql"
)

func IsSiteOwner(r *Resolver) func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		currentUser := identity.GetCurrentOrganizationUser(r.ContextUserManager, ctx)
		if currentUser == nil {
			return nil, fmt.Errorf("not authenticated")
		}
		if currentUser.IsSiteOwner() {
			return next(ctx)
		}
		return nil, fmt.Errorf("not site owner")
	}
}

func IsAuthenticated(r *Resolver) func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		currentUser := identity.GetCurrentOrganizationUser(r.ContextUserManager, ctx)
		if currentUser == nil {
			return nil, fmt.Errorf("not authenticated")
		}
		return next(ctx)
	}
}

func CaptchaProtected(r *Resolver) func(ctx context.Context, obj interface{}, next graphql.Resolver, action string) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, action string) (interface{}, error) {
		captchaClient := r.CaptchaClient
		if captchaClient == nil {
			return next(ctx)
		}
		token := captcha.GetToken(ctx)
		if token == "" {
			return nil, fmt.Errorf("captcha token not found")
		}
		fieldContext := graphql.GetFieldContext(ctx)
		resolverLogger.Info().Str("field", fieldContext.Field.Name).Str("action", action).Msgf("Validating captcha token")
		if ok, err := captchaClient.VerifyCaptcha(ctx, token, action); err != nil {
			return nil, err
		} else {
			resolverLogger.Info().Str("field", fieldContext.Field.Name).Str("action", action).Bool("ok", ok).Msgf("captcha token validation result")
			if !ok {
				return nil, fmt.Errorf("captcha verification failed")
			} else {
				return next(ctx)
			}
		}
	}
}
