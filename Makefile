.PHONY: clean
clean: 
	rm -rf app/static/
	rm db-gui || :
	rm -rf frontend/src/components/monaco-editor-theme/onedark/OnedarkTextmate.tmtheme.ts

.PHONY: theme
theme:
	cd frontend &&\
	node src/static/mktmTheme.js

.PHONY: frontend
frontend: theme
	cd frontend &&\
	npm ci &&\
	npm run build

.PHONY: backend
backend: frontend
	go build -o db-gui cmd/main/main.go

.PHONY: all
all: clean backend