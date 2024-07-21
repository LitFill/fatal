VERSION := 0.2.0
.PHONY: release doc changelog help
release: doc changelog ## Create a release on GitHub
	@echo "Creating release $(VERSION) on GitHub"
	@git tag -a v$(VERSION) -m "Version $(VERSION)"
	@git push origin v$(VERSION)
	@gh release create v$(VERSION) --title "$(VERSION)" --generate-notes --notes "Release $(VERSION), view changelogs in CHANGELOG.md"

doc: ## Create docs/scc.html
	@echo "Creating scc documentation in html"
	@mkdir -p "docs"
	@touch "docs/scc.html"
	@scc --overhead 1.0 --no-gen -n "scc.html" -s "complexity" -f "html" > docs/scc.html

changelog: ## Update CHANGELOG.md
	@echo "Updating CHANGELOG.md"
	@git-chglog --output CHANGELOG.md

help: ## Prints help for targets with comments
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
