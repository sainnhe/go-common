This repository contains my personal implementation of some commonly used libraries.

To develop:

1. Use go version specified in `go.mod` and `.tool-versions`. You can use [asdf](https://github.com/asdf-vm/asdf) to manage multiple versions.
2. Install [mockgen](https://github.com/uber-go/mock). If you are using asdf to manage versions, execute `asdf reshim golang` after installing mockgen.
3. Install [golangci-lint](https://golangci-lint.run/welcome/install/).
4. Execute the following commands to install git hooks:
  - Pre-commit hook will run checks and tests before committing: `ln -s ../../githooks/pre-commit .git/hooks/`
  - Commit-msg hook will check where commit message conforms to [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/): `ln -s ../../githooks/commit-msg .git/hooks/`
