define setup_env
	$(eval ENV_FILE := $(1))
	$(eval include $(1))
	$(eval export)
endef

build-cli:
	go build -o ./cli cmd/*.go

build-push-cli:
	$(call setup_env, server/.env)
	CGO_ENABLED=0 GOOS=linux go build -o ./cli cmd/*.go
	docker build -f Dockerfile -t ${CLI_IMG_TAG} .
	docker push ${CLI_IMG_TAG}

migrate:
	$(call setup_env, server/.env)
	migrate -database ${DATABASE_URL} -path server/migrations up

dump-schema:
	$(call setup_env, server/.env)
	docker run -it --rm \
	-v ./.pgdump:/.pgdump \
	postgres pg_dump -d ${DATABASE_URL} -s -b -v -f .pgdump/pgdump.sql
	mv .pgdump/pgdump.sql .
	rmdir .pgdump

run-http-server:
	$(call setup_env, server/.env)
	./cli run http-server

run-worker:
	$(call setup_env, worker/.env)
	./cli run worker

run-http-server-local:
	$(call setup_env, server/.env.local)
	kubectl port-forward services/temporal-frontend 7233:7233 &
	@$(MAKE) build-cli
	./cli run http-server

run-worker-local:
	$(call setup_env, worker/.env.local)
	kubectl port-forward services/temporal-frontend 7233:7233 &
	@$(MAKE) build-cli
	./cli run worker

deploy-server:
	$(call setup_env, server/.env)
	@$(MAKE) build-push-cli
	kustomize build --load-restrictor=LoadRestrictionsNone server/k8s | \
	sed -e "s;{{DOCKER_REPO}};$(DOCKER_REPO);g" | \
	sed -e "s;{{CLI_IMG_TAG}};$(CLI_IMG_TAG);g" | \
	kubectl apply -f -
	kubectl rollout restart deployment kaggo-backend

deploy-worker:
	$(call setup_env, worker/.env)
	@$(MAKE) build-push-cli
	kustomize build --load-restrictor=LoadRestrictionsNone worker/k8s | \
	sed -e "s;{{DOCKER_REPO}};$(DOCKER_REPO);g" | \
	sed -e "s;{{CLI_IMG_TAG}};$(CLI_IMG_TAG);g" | \
	kubectl apply -f -
	kubectl rollout restart deployment kaggo-treq-worker

initialize-listeners:
	$(call setup_env, server/.env)
	@$(MAKE) build-cli
	./cli admin listener initiate-youtube-listener --endpoint https://api.kaggo.brojonat.com

deploy-all:
	# server
	@$(MAKE) deploy-server
	# worker
	@$(MAKE) deploy-worker
	# kick off long lived processes
	@$(MAKE) initialize-listeners

