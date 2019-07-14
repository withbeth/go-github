/*
 * refer :
 * https://github.com/github/platform-samples/blob/master/api/golang/basics-of-authentication/server.go
 * https://mingrammer.com/getting-started-with-oauth2-in-go/
 */

package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"log"
	"os"
	"time"
)

var userName = "withbeth"
var oAuth2AccessToken = os.Getenv("GH_PERSONAL_ACCESS_TOKEN")

func main() {

	/* Authentication :
	 * go-github library does not directly handle authentication.
	 * Instead, when creating a new client, pass an http.Client that can handle authentication for you.
	 */
	if oAuth2AccessToken == "" {
		log.Fatal("Unauthorized: No token present")
	}

	ctx := context.Background()
	cli := getGitHubClient(ctx, oAuth2AccessToken)

	/** 1. Get all repositories */
	repos, _, err := cli.Repositories.List(ctx, userName, nil)
	if err != nil {
		log.Fatal(err)
	}

	/** 2. Check updated repo since yesterday */
	sumOfRepoNotifications := sumRepoNotifications(cli, ctx, userName, repos)
	fmt.Println(sumOfRepoNotifications)
	if (sumOfRepoNotifications < 1) {
		log.Fatal("Oh goat.... you failed AGAIN.")
	}

}

func sumRepoNotifications(cli *github.Client, ctx context.Context, owner string, repos []*github.Repository) (result int) {
	for _, repo := range repos {
		notifications, _, err := cli.Activity.ListRepositoryNotifications(
			ctx,
			owner,
			repo.GetName(),
			&github.NotificationListOptions{Since: time.Now().AddDate(0, 0, -1)})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(repo.GetName())
		fmt.Println(notifications)
		result += len(notifications)
	}
	return
}

// Authenticates GitHub Client with provided OAuth access token
func getGitHubClient(ctx context.Context, accessToken string) *github.Client {
	/* Authentication :
	 * go-github library does not directly handle authentication.
	 * Instead, when creating a new client, pass an http.Client that can handle authentication for you.
	 */
	if accessToken == "" {
		log.Fatal("Unauthorized: No token present")
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	httpClient := oauth2.NewClient(ctx, tokenSource)

	return github.NewClient(httpClient)
}


