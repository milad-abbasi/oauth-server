.DEFAULT_GOAL := serve

# .PHONY is used for reserving tasks words
.PHONY: start before stop restart serve

# Filename to store running process id
PID_FILE := /tmp/oauth-server.pid

MIGRATE_DB_URI ?= postgres://postgres:example@localhost:5432/postgres?sslmode=disable
MIGRATE_PATH ?= migrations
MIGRATE_CMD ?= up # up | down | goto
MIGRATE_STEPS ?= 0 # Apply all

migrate:
	@go run cmd/migrate/migrate.go -d ${MIGRATE_DB_URI} -p ${MIGRATE_PATH} ${MIGRATE_CMD} ${MIGRATE_STEPS}

# Start task starts http server up and writes it's process id to PID_FILE.
start:
	@go run cmd/http/server.go 2>&1 & echo $$! > ${PID_FILE}
# You can also use go build command for start task
# start:
#   go build -o /bin/http-server cmd/http/server.go && \
#   /bin/http-server 2>&1 & echo $$! > ${PID_FILE}

# Stop task will kill process by ID stored in PID_FILE (and all child processes by pstree).
stop:
	@-kill `pstree -p \`cat ${PID_FILE}\` | tr "\n" " " | sed "s/[^0-9]/ /g" | sed "s/\s\s*/ /g"`

# Before task will only prints message. Actually, it is not necessary. You can remove it, if you want.
before:
	@echo "STOPED server" && printf '%*s\n' "40" '' | tr ' ' -

# Restart task will execute stop, before and start tasks in strict order and prints message.
restart: stop before start
	@echo "STARTED server" && printf '%*s\n' "40" '' | tr ' ' -

# Serve task will run fswatch monitor and performs restart task if any source file changed. Before serving it will execute start task.
serve: start
	@fswatch -or --event=Updated ~/workspace/oauth-server | xargs -n1 -I {} make restart
