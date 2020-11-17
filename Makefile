.PHONY: link-data
link-data:
	ln -s ./example-data ./data

.PHONY: skip-data
skip-data:
	git update-index --skip-worktree data/*.json

.PHONY: no-skip-data
no-skip-data:
	git update-index --no-skip-worktree data/*.json