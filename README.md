# ![stackpulse-logo](vendors/stackpulse.svg) StackPulse Steps

[![Build][badge_ci]][link_circle]
[![Contributors][contributors-shield]][contributors-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][badge_issues]][issues-url]
[![License][badge_license]][link_license]

This repository is the official library of StackPulse steps.

StackPulse steps are containerized applications that can be composed together to form a playbook. Steps can perform any arbitrary computational task but are most often calls to different services and APIs, allowing the playbook composer to interact with their stack in an automated way. This enables faster and more streamlined handling of alerts and remediation of incidents.

Just like any other container, steps can receive input in the form of run arguments or environment variables and have output directed to stdout and stderr. The standard form for a step to issue its output is as a JSON object as defined in the [StackPulse steps Golang SDK](https://github.com/stackpulse/steps-sdk-go).

## Development

...

## Contributing

Pull requests adding new steps and submitting bug fixes are welcome.

Please note that all new steps must use the [StackPulse steps Golang SDK](https://github.com/stackpulse/steps-sdk-go) and be submitted with a step manifest containing example input and output examples.

In the case of bug fixes, where appropriate, please provide an accompanying unit test that validates and verifies the proposed fix.

## Troubleshooting

Please feel free to open an issue for any problem or you may encounter while using one of the steps in the repo. We welcome requests and suggestions for new steps as well.

## License

Distributed under the BSD-3-Clause License. See `LICENSE` for more information.

[badge_ci]:https://circleci.com/gh/stackpulse/steps.svg?style=svg
[contributors-shield]: https://img.shields.io/github/contributors/stackpulse/steps.svg?style=flat-square&maxAge=30
[contributors-url]: https://github.com/stackpulse/steps/graphs/contributors
[badge_issues]:https://img.shields.io/github/issues/stackpulse/steps.svg?style=flat-square&maxAge=30
[issues-url]: https://github.com/stackpulse/steps/issues
[stars-shield]: https://img.shields.io/github/stars/stackpulse/steps.svg?style=flat-square&maxAge=30
[stars-url]: https://github.com/stackpulse/steps/stargazers
[badge_license]:https://img.shields.io/github/license/stackpulse/steps.svg?style=flat-square&maxAge=30
[link_license]:https://github.com/stackpulse/steps/blob/master/LICENSE
[link_circle]:https://circleci.com/gh/stackpulse/steps
