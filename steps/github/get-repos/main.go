package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shurcooL/githubv4"
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
	"golang.org/x/oauth2"
)

type Args struct {
	Owner   string `env:"OWNER,required" envDefault:""`
	Token   string `env:"TOKEN,required" envDefault:""`
	Privacy string `env:"TYPE" envDefault:""`
}

type GHCommand struct {
	args Args
}

type Repository string

type Output struct {
	Repositories []Repository `json:"repositories"`
}

func (output *Output) AddItem(repo Repository) []Repository {
	output.Repositories = append(output.Repositories, repo)
	return output.Repositories
}

func (gh *GHCommand) Init() error {
	return env.Parse(&gh.args)
}

func (gh *GHCommand) Run() (int, []byte, error) {
	var err error
	var typeQuery struct {
		Organization struct {
			Repositories struct {
				Nodes []struct {
					Name string
				}
			} `graphql:"repositories(privacy: $privacy, first: 100, affiliations:  [OWNER])"`
		} `graphql:"organization(login: $login)"`
	}
	var nonTypeQuery struct {
		Organization struct {
			Repositories struct {
				Nodes []struct {
					Name string
				}
			} `graphql:"repositories(first: 100, affiliations:  [OWNER])"`
		} `graphql:"organization(login: $login)"`
	}

	variables := map[string]interface{}{
		"login": githubv4.String(gh.args.Owner),
	}
	if gh.args.Privacy != "" {
		upper := githubv4.RepositoryPrivacy(strings.ToUpper(gh.args.Privacy))
		variables["privacy"] = githubv4.RepositoryPrivacy(upper)
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gh.args.Token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)
	repositories := []Repository{}
	output := Output{repositories}
	if gh.args.Privacy != "" {
		err = client.Query(context.Background(), &typeQuery, variables)
		if err != nil {
			return step.ExitCodeFailure, nil, fmt.Errorf("get-repositories failed with error: error: %w", err)
		}

		for _, node := range typeQuery.Organization.Repositories.Nodes {
			repo := Repository(node.Name)
			output.AddItem(repo)

		}
	} else {
		err = client.Query(context.Background(), &nonTypeQuery, variables)
		if err != nil {
			return step.ExitCodeFailure, nil, fmt.Errorf("get-repositories failed with error: error: %w", err)
		}
		for _, node := range nonTypeQuery.Organization.Repositories.Nodes {
			repo := Repository(node.Name)
			output.AddItem(repo)
		}
	}
	ret, err := json.Marshal(output)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("get-repositories failed with error: error: %w", err)

	}
	return step.ExitCodeOK, ret, err
}

func main() {
	step.Run(&GHCommand{})
}
