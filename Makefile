OUTPUTDIR  :=  out
MANIFESTDIR := $(OUTPUTDIR)/manifests
TARGETS := $(shell find steps -name "Dockerfile" | sort)
MANIFESTS := $(shell find steps -name "*.yaml" -type f -exec sh -c "head -3 {} | grep -q 'kind: Step' && echo {}" \; | sort)
CONTAINER_REGISTRY ?= us-docker.pkg.dev/stackpulse/public
MANIFEST_PATH ?= gs://stackpulse-steps/
VENDORS_PATH ?= gs://stackpulse-public/
MANIFEST_PARSER ?= gcr.io/stackpulse/step-manifest-parser:prd-21.01.0


TAG ?=
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)

# push latest only for master branch
LATEST ?= $(subst master,latest,$(filter master,$(BRANCH)))
DOCKER_TAGS = $(TAG) $(BRANCH) $(LATEST)

directory-name=$(patsubst ./%,%,$(patsubst %/,%,$(dir $(1))))
remove-step-prefix=$(subst steps/,,$(1))
rulename=$(subst /,-,$(call remove-step-prefix,$(call directory-name, $(1))))

.DELETE_ON_ERROR:
# declare-shortcut(source, dest) // declare a phony rule shortcut from src to dst.
define declare-shortcut
	$(eval .PHONY: $(1))                                                                                                       \
	$(eval                                                                                                                     \
		$(1): $(2)                                                                                                             \
	 )
endef

apps:
	@./scripts/baur_apps.py

gomod:
	@./scripts/go_mod.py

pg:
	@./scripts/local_postgres.sh up

local: pg
	@baur run

clean:
	@./scripts/local_postgres.sh down
	rm -rf $(OUTPUTDIR)

all:
ifeq ("$(CIRCLE_BRANCH)","master")
	@echo "Building master branch"
	@baur run
else
	@echo "Building side branch"
	@baur run --skip-upload
endif

.PHONY: local clean apps gomod all pg

# --- Manifest stuff
# manifest-build-rule(path)
define declare-manifest-build-rule
	$(eval $(1)_packname := $(call rulename, $(1)))                                                                            \
	$(call declare-shortcut,pack-$($(1)_packname),$(MANIFESTDIR)/$($(1)_packname).yml)                                         \
	$(eval                                                                                                                     \
		$(MANIFESTDIR)/$($(1)_packname).yml: $(1) ; $$(info === packing $(1))
			mkdir -p $$(dir $$@);
			grep -q 'imageName:' $(1) && cp $(1) $$@ || docker run -w /root -v $(CURDIR):/root $(MANIFEST_PARSER) image set $(1) $$@ $(CONTAINER_REGISTRY)/$(call remove-step-prefix,$(call directory-name, $(1)))
			docker run -w /root -v $(CURDIR):/root $(MANIFEST_PARSER) validate $$@
	 )
endef

# create pack-$target, packall
$(foreach manifest, $(MANIFESTS), $(call declare-manifest-build-rule,$(manifest)))
packall: $(foreach manifest, $(MANIFESTS),pack-$(call rulename, $(manifest)))


$(call declare-shortcut,indexfile,$(MANIFESTDIR)/indexfile.yml)
$(MANIFESTDIR)/indexfile.yml: packall
	./generate-index-file.py $(MANIFESTDIR) $@

.PHONY: publish-manifests
publish-manifests: packall $(MANIFESTDIR)/indexfile.yml
	gsutil -m cp -r $(MANIFESTDIR)/*.yml $(MANIFEST_PATH)

.PHONY: publish-manifests
publish-manifests-no-deps:
	gsutil -m cp -r $(MANIFESTDIR)/*.yml $(MANIFEST_PATH)

.PHONY: validate-vendors
validate-vendors:
	./validate-vendors.py

.PHONY: publish-vendors
publish-vendors:
	gsutil -m cp -r vendors $(VENDORS_PATH)