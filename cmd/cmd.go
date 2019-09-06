package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"short/dep"
	"short/dep/provider"
	"strconv"

	"github.com/spf13/cobra"
)

func Execute(
	host string,
	portStr string,
	user string,
	password string,
	dbName string,
	recaptchaSecret string,
	githubClientID string,
	githubClientSecret string,
	jwtSecret string,
) {
	var migrationRoot string
	var wwwRoot string

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start service",
		Run: func(cmd *cobra.Command, args []string) {
			port := MustInt(portStr)

			start(
				host,
				port,
				user,
				password,
				dbName,
				migrationRoot,
				wwwRoot,
				recaptchaSecret,
				githubClientID,
				githubClientSecret,
				jwtSecret,
			)
		},
	}

	startCmd.Flags().StringVar(&migrationRoot, "migration", "app/adapter/migration", "migration migrations root directory")
	startCmd.Flags().StringVar(&wwwRoot, "www", "public", "www root directory")

	rootCmd := &cobra.Command{Use: "short"}
	rootCmd.AddCommand(startCmd)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func MustInt(numStr string) int {
	num, err := strconv.Atoi(numStr)

	if err != nil {
		panic(err)
	}

	return num
}

func start(
	host string,
	port int,
	user string,
	password string,
	dbName string,
	migrationRoot string,
	wwwRoot string,
	recaptchaSecret string,
	githubClientID string,
	githubClientSecret string,
	jwtSecret string,
) {
	provider.DB(host, port, user, password, dbName, migrationRoot, func(db *sql.DB) {
		service := dep.InjectGraphQlService(
			"GraphQL API",
			db,
			"/graphql",
			provider.ReCaptchaSecret(recaptchaSecret),
			provider.JwtSecret(jwtSecret),
		)
		service.Start(8080)

		service = dep.InjectRoutingService(
			"Routing API",
			db,
			provider.WwwRoot(wwwRoot),
			provider.GithubClientID(githubClientID),
			provider.GithubClientSecret(githubClientSecret),
			provider.JwtSecret(jwtSecret),
		)
		service.StartAndWait(80)
	})
}
