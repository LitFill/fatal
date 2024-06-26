VERSION := 0.0.3
.PHONY: release doc help
release: doc ## Create a release on GitHub
	@echo "Creating release $(VERSION) on GitHub"
	@git tag -a v$(VERSION) -m "Version $(VERSION)"
	@git push origin v$(VERSION)
	@gh release create v$(VERSION) --title "$(VERSION)" --generate-notes --notes-from-tag --notes "Release $(VERSION), view changelogs in CHANGELOG.md"

doc: ## Create doc/scc.html
	@echo "Creating scc documentation in html"
	@mkdir -p "doc"
	@touch "doc/scc.html"
	@scc --overhead 1.0 --no-gen -n "scc.html" -s "complexity" -f "html" > doc/scc.html

help: ## Prints help for targets with comments
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
