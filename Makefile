task-cli-add:
	@if [ -z "$(description)" ]; then \
		echo "Error: description is required. Usage: make task-cli-add description=<string>"; \
		exit 1; \
	fi
	go run main.go task-cli add $(description)

task-cli-update:
	@if [ -z "$(id)" ]; then \
		echo "Error: id is required. Usage: make task-cli-update id=<number>"; \
		exit 1; \
	fi

	@if [ -z "$(description)" ]; then \
		echo "Error: description is required. Usage: make task-cli-update id=<number> description=<string>"; \
		exit 1; \
	fi

	go run main.go task-cli update $(id) $(description)

task-cli-delete:
	@if [ -z "$(id)" ]; then \
		echo "Error: id is required. Usage: make task-cli-delete id=<number>"; \
		exit 1; \
	fi

	go run main.go task-cli delete $(id)

task-cli-mip:
	@if [ -z "$(id)" ]; then \
		echo "Error: id is required. Usage: make task-cli-mip id=<number>"; \
		exit 1; \
	fi

	go run main.go task-cli mark-in-progress $(id)

task-cli-done:
	@if [ -z "$(id)" ]; then \
		echo "Error: id is required. Usage: make task-cli-done id=<number>"; \
		exit 1; \
	fi

	go run main.go task-cli mark-done $(id)

task-cli-ls:
	go run main.go task-cli list $(status)
