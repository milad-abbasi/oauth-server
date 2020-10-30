.DEFAULT_GOAL := serve

# .PHONY is used for reserving tasks words
.PHONY: start before stop restart serve

# Variable for filename for store running process id
PID_FILE := /tmp/oauth-server.pid

# We can use such syntax to get main.go and other root Go files.
HTTP_SERVER := ./cmd/http/server.go

# Start task performs "go run main.go" command and writes it's process id to PID_FILE.
start:
	@go run ${HTTP_SERVER} 2>&1 & echo $$! > ${PID_FILE}
# You can also use go build command for start task
# start:
#   go build -o /bin/my-app . && \
#   /bin/my-app & echo $$! > $(PID_FILE)

# Stop task will kill process by ID stored in PID_FILE (and all child processes by pstree).
stop:
	@-kill `pstree -p \`cat ${PID_FILE}\` | tr "\n" " " | sed "s/[^0-9]/ /g" | sed "s/\s\s*/ /g"`

# Before task will only prints message. Actually, it is not necessary. You can remove it, if you want.
before:
	@echo "STOPED my-app" && printf '%*s\n' "40" '' | tr ' ' -

# Restart task will execute stop, before and start tasks in strict order and prints message.
restart: stop before start
	@echo "STARTED my-app" && printf '%*s\n' "40" '' | tr ' ' -

# Serve task will run fswatch monitor and performs restart task if any source file changed. Before serving it will execute start task.
serve: start
	@fswatch -or --event=Updated ~/workspace/oauth-server | xargs -n1 -I {} make restart
