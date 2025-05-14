
This is my work from [Roadmap.sh](https://roadmap.sh) on [first task](https://roadmap.sh/projects/task-tracker)  

# Technology
I'm just using Go as programming language.  

# How to run
There are two ways for running this CLI app. First using `go run main.go task-cli <command> <args>` or using my Makefile script `make task-cli-ls`, `make task-cli-add description="..."`, `make task-cli-update id=<number> description="..."`, etc

# Architecture
I just separate the codes by it's meaning like `types.go` for the structs declaration, `util.go` for function that i'm using inside `main.go`, and `main.go` the main executioner of this programs

