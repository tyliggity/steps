module github.com/stackpulse/public-steps/steps/kubectl/apply-file

go 1.14

require github.com/stackpulse/public-steps/kubectl/base v0.0.0

replace github.com/stackpulse/public-steps/kubectl/base v0.0.0 => ../base
