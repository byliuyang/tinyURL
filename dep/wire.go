//+build wireinject

package dep

import (
	"database/sql"
	"short/app/adapter/github"
	"short/app/adapter/graphql"
	"short/app/adapter/sqlrepo"
	"short/app/usecase/account"
	"short/app/usecase/keygen"
	"short/app/usecase/repo"
	"short/app/usecase/requester"
	"short/app/usecase/service"
	"short/app/usecase/url"
	"short/dep/provider"
	"time"

	"github.com/byliuyang/app/fw"
	"github.com/byliuyang/app/modern/mdhttp"
	"github.com/byliuyang/app/modern/mdlogger"
	"github.com/byliuyang/app/modern/mdrequest"
	"github.com/byliuyang/app/modern/mdrouting"
	"github.com/byliuyang/app/modern/mdservice"
	"github.com/byliuyang/app/modern/mdtimer"
	"github.com/byliuyang/app/modern/mdtracer"
	"github.com/google/wire"
)

const oneDay = 24 * time.Hour

func InjectGraphQlService(
	name string,
	db *sql.DB,
	graphqlPath provider.GraphQlPath,
	secret provider.ReCaptchaSecret,
	jwtSecret provider.JwtSecret,
) mdservice.Service {
	wire.Build(
		wire.Value(provider.TokenValidDuration(oneDay)),

		wire.Bind(new(fw.GraphQlAPI), new(graphql.Short)),
		wire.Bind(new(url.Retriever), new(url.RetrieverPersist)),
		wire.Bind(new(url.Creator), new(url.CreatorPersist)),
		wire.Bind(new(repo.UserURL), new(sqlrepo.UserURL)),
		wire.Bind(new(repo.URL), new(*sqlrepo.URL)),

		mdservice.New,
		mdlogger.NewLocal,
		mdtracer.NewLocal,
		provider.GraphGophers,
		mdhttp.NewClient,
		mdrequest.NewHTTP,
		mdtimer.NewTimer,
		provider.JwtGo,

		sqlrepo.NewURL,
		sqlrepo.NewUserURL,
		keygen.NewInMemory,
		provider.Authenticator,
		url.NewRetrieverPersist,
		url.NewCreatorPersist,
		provider.ReCaptchaService,
		requester.NewVerifier,
		graphql.NewShort,
	)
	return mdservice.Service{}
}

func InjectRoutingService(
	name string,
	db *sql.DB,
	wwwRoot provider.WwwRoot,
	githubClientID provider.GithubClientID,
	githubClientSecret provider.GithubClientSecret,
	jwtSecret provider.JwtSecret,
) mdservice.Service {
	wire.Build(
		wire.Value(provider.TokenValidDuration(oneDay)),

		wire.Bind(new(service.Account), new(account.RepoService)),
		wire.Bind(new(url.Retriever), new(url.RetrieverPersist)),
		wire.Bind(new(repo.User), new(*(sqlrepo.User))),
		wire.Bind(new(repo.URL), new(*sqlrepo.URL)),

		mdservice.New,
		mdlogger.NewLocal,
		mdtracer.NewLocal,
		mdrouting.NewBuiltIn,
		mdhttp.NewClient,
		mdrequest.NewHTTP,
		mdrequest.NewGraphQl,
		mdtimer.NewTimer,
		provider.JwtGo,

		sqlrepo.NewUser,
		sqlrepo.NewURL,
		url.NewRetrieverPersist,
		account.NewRepoService,
		provider.GithubOAuth,
		github.NewAPI,
		provider.Authenticator,
		provider.ShortRoutes,
	)
	return mdservice.Service{}
}
