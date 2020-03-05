// +build !integration all

package account

import (
	"testing"

	"github.com/short-d/short/app/usecase/service"

	"github.com/short-d/app/mdtest"
	kgsentity "github.com/short-d/kgs/app/entity"
	"github.com/short-d/short/app/entity"
	"github.com/short-d/short/app/usecase/keygen"
	"github.com/short-d/short/app/usecase/repository"
)

func TestLinker_IsAccountLinked(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		keys             []string
		users            []entity.User
		mappingUsers     []entity.User
		mappingSSOUsers  []entity.SSOUser
		ssoUser          entity.SSOUser
		expectedIsLinked bool
	}{
		{
			name:            "account not linked",
			keys:            []string{},
			mappingUsers:    []entity.User{},
			mappingSSOUsers: []entity.SSOUser{},
			ssoUser: entity.SSOUser{
				ID:    "alpha",
				Email: "alpha@example.com",
				Name:  "Alpha User",
			},
			expectedIsLinked: false,
		},
		{
			name: "account already linked",
			keys: []string{},
			mappingUsers: []entity.User{
				{
					ID:    "beta",
					Name:  "Beta",
					Email: "beta@example.com",
				},
			},
			mappingSSOUsers: []entity.SSOUser{
				{
					ID:    "alpha",
					Email: "alpha@example.com",
					Name:  "Alpha User",
				},
			},
			ssoUser: entity.SSOUser{
				ID:    "alpha",
				Email: "alpha@example.com",
				Name:  "Alpha User",
			},
			expectedIsLinked: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			keyFetcher := service.NewKeyFetcherFake([]kgsentity.Key{})
			keyGen, err := keygen.NewKeyGenerator(2, &keyFetcher)
			mdtest.Equal(t, nil, err)
			userRepo := repository.NewUserFake(testCase.users)
			accountMappingRepo, err :=
				repository.NewAccountMappingFake(
					testCase.mappingSSOUsers,
					testCase.mappingUsers,
				)
			mdtest.Equal(t, nil, err)

			linker := NewLinker(keyGen, &userRepo, &accountMappingRepo)
			isLinked, err := linker.IsAccountLinked(testCase.ssoUser)
			mdtest.Equal(t, nil, err)
			mdtest.Equal(t, testCase.expectedIsLinked, isLinked)
		})
	}
}

func TestLinker_LinkAccount(t *testing.T) {
	testCases := []struct {
		name            string
		key             string
		mappingUsers    []entity.User
		mappingSSOUsers []entity.SSOUser
		users           []entity.User
		ssoUser         entity.SSOUser
		user            entity.User
		expectedRes     bool
	}{
		{
			name: "account already linked",
			mappingUsers: []entity.User{
				{
					ID: "alpha",
				},
			},
			mappingSSOUsers: []entity.SSOUser{
				{
					ID: "gama",
				},
			},
			users: []entity.User{
				{
					ID: "alpha",
				},
			},
			ssoUser: entity.SSOUser{
				ID: "gama",
			},
			user: entity.User{
				ID: "alpha",
			},
			expectedRes: true,
		},
		{
			name:            "account exists not linked",
			key:             "alpha",
			mappingUsers:    []entity.User{},
			mappingSSOUsers: []entity.SSOUser{},
			users: []entity.User{
				{
					Email: "alpha@example.com",
				},
			},
			ssoUser: entity.SSOUser{
				ID:    "gama",
				Email: "alpha@example.com",
			},
			user: entity.User{
				ID:    "alpha",
				Email: "alpha@example.com",
			},
			expectedRes: false,
		},
		{
			name:            "create new account",
			key:             "alpha",
			mappingUsers:    []entity.User{},
			mappingSSOUsers: []entity.SSOUser{},
			users:           []entity.User{},
			ssoUser: entity.SSOUser{
				ID:    "gama",
				Email: "alpha@example.com",
			},
			user: entity.User{
				ID:    "alpha",
				Email: "alpha@example.com",
			},
			expectedRes: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			keyFetcher := service.NewKeyFetcherFake([]kgsentity.Key{"key", "key2"})
			keyGen, err := keygen.NewKeyGenerator(2, &keyFetcher)
			mdtest.Equal(t, nil, err)
			fakeUserRepo := repository.NewUserFake(testCase.users)
			accountMappingRepo, err :=
				repository.NewAccountMappingFake(
					testCase.mappingSSOUsers,
					testCase.mappingUsers,
				)
			mdtest.Equal(t, nil, err)

			linker := NewLinker(keyGen, &fakeUserRepo, &accountMappingRepo)
			err = linker.CreateAndLinkAccount(testCase.ssoUser)
			mdtest.Equal(t, nil, err)

			gotIsRelationExist := accountMappingRepo.IsRelationExist(testCase.ssoUser, testCase.user)
			mdtest.Equal(t, testCase.expectedRes, gotIsRelationExist)

			gotIsIDExist := fakeUserRepo.IsUserIDExist(testCase.user.ID)
			mdtest.Equal(t, testCase.expectedRes, gotIsIDExist)
		})
	}
}
