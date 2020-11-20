data_excludes = data/*.json data/scripts/*.lua data/object-images/*.png

.PHONY: link-data
link-data:
	ln -s ./example-data ./data

.PHONY: skip-data
skip-data:
	git update-index --skip-worktree $(data_excludes)

.PHONY: no-skip-data
no-skip-data:
	git update-index --no-skip-worktree $(data_excludes)