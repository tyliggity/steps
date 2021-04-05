package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shurcooL/githubv4"
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
	"golang.org/x/oauth2"
)

type Args struct {
	Owner      string `env:"OWNER,required" envDefault:""`
	Token      string `env:"TOKEN,required" envDefault:""`
	Repository string `env:"REPOSITORY,required" envDefault:""`
	NumOfTags  int    `env:"NUM_OF_LAST_TAGS" envDefault:"5"`
}

type GHCommand struct {
	args Args
}

type Tag struct {
	Name    string
	Message string
}

type Output struct {
	Tags []Tag
}

func (output *Output) AddItem(tag Tag) []Tag {
	output.Tags = append(output.Tags, tag)
	return output.Tags
}
func (gh *GHCommand) Init() error {
	return env.Parse(&gh.args)
}

func (gh *GHCommand) Run() (int, []byte, error) {
	var query struct {
		Repository struct {
			Description string
			Refs        struct {
				Edges []struct {
					Node struct {
						Target struct {
							Tag struct {
								Name    string
								Message string
							} `graphql:"... on Tag"`
						}
					}
				}
			} `graphql:"refs(refPrefix: \"refs/tags/\", orderBy: {direction: DESC, field: TAG_COMMIT_DATE}, first: $num_of_tags)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":       githubv4.String(gh.args.Owner),
		"name":        githubv4.String(gh.args.Repository),
		"num_of_tags": githubv4.Int(gh.args.NumOfTags),
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gh.args.Token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)
	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("get-steps failed with error: error: %w", err)
	}
	tags := []Tag{}
	output := Output{tags}
	for _, edge := range query.Repository.Refs.Edges {
		tag := Tag{edge.Node.Target.Tag.Name, edge.Node.Target.Tag.Message}
		output.AddItem(tag)
	}

	ret, err := json.Marshal(output)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("get-steps failed with error: error: %w", err)

	}
	return step.ExitCodeOK, ret, err
}

func main() {
	step.Run(&GHCommand{})
}
