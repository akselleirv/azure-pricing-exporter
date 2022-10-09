.PHONY: openapi gen dev

openapi:
	cd api/spec  \
	&& cue export --out openapi types.cue | cue export - routes.cue --out yaml  > openapi.yaml

gen: openapi
	mkdir -p api/models
	oapi-codegen --config api/config/server.cfg.yaml api/spec/openapi.yaml \
	&& oapi-codegen --config api/config/models.cfg.yaml api/spec/openapi.yaml

dev: gen
	tilt up