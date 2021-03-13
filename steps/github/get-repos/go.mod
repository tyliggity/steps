module github.com/stackpulse/public-steps/github/get-repos

go 1.15

require (
	github.com/caarlos0/env/v6 v6.5.0
	github.com/google/go-github/v33 v33.0.0 // indirect
	github.com/shurcooL/githubv4 v0.0.0-20201206200315-234843c633fa
	github.com/shurcooL/graphql v0.0.0-20200928012149-18c5c3165e3a // indirect
	github.com/stackpulse/public-steps/common v0.0.0
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
)

replace github.com/stackpulse/public-steps/common v0.0.0 => ../../common
