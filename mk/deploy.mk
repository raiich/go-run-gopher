.PHONY: deploy
deploy: ## deploy
	npx wrangler pages deploy dist/ --project-name=go-run-gopher
