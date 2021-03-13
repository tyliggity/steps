module github.com/stackpulse/public-steps/github/get-tags

go 1.15

require (
	github.com/caarlos0/env/v6 v6.5.0
	github.com/google/go-github/v33 v33.0.0 // indirect
	github.com/shurcooL/githubv4 v0.0.0-20201206200315-234843c633fa // indirect
	github.com/shurcooL/graphql v0.0.0-20200928012149-18c5c3165e3a // indirect
	github.com/stackpulse/public-steps/common v0.0.0
)

replace github.com/stackpulse/public-steps/common v0.0.0 => ../../common
