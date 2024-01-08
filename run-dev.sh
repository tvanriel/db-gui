#!/bin/bash

DEV=1 $HOME/go/bin/CompileDaemon --command="/tmp/db-gui-backend" -build="go build -o /tmp/db-gui-backend ./cmd/main/main.go" & (cd frontend; npm run dev)