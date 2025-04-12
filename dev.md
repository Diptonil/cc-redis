# Devlog

## 1: `go.mod`

- This file is used to track project dependencies and module information.
- It stores the go version being used as well as the main place where the module is available (be it a GitHub repo, a website, etc.).
- As dependencies are added in, their versions would get listed down. This is to keep builds reproducibale and consistent.
- It's kind of like `package.json`.
