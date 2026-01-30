BRANCH ?= main
FRONT_BRANCH ?= main

worktree:
	git worktree add trees/$(NEW_BRANCH) $(BRANCH) -B $(NEW_BRANCH) && cd trees/$(NEW_BRANCH) && git submodule update && $(MAKE) rebase-front

rebase-front:
	cd front && git fetch && git reset --hard origin/main && git rebase origin/$(FRONT_BRANCH)

build-front:
	cd front && npm ci && npm run build

deploy: rebase-front build-front
	cd go/infra && cdk deploy -y --progress --all

deploy-no-rebase: build-front
	cd go/infra && cdk deploy -y --progress --all