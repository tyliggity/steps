
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

# declare-docker-rule(directory, tag)
define declare-build-rule
	$(eval $(1)_buildname := $(call rulename, $(1)))                                                                           \
	$(eval $(1)_builddir  := $(call directory-name, $(1)))                                                                     \
	$(eval $(1)_tagname   := $(CONTAINER_REGISTRY)/$(call remove-step-prefix,$(call directory-name, $(1))))                    \
	$(eval .PHONY: build-$($(1)_buildname)-$(2))                                                                               \
	$(eval                                                                                                                     \
		build-$($(1)_buildname)-$(2): ; $$(info === docker building $($(1)_builddir):$(2))
	docker build --build-arg STEP_BASEPATH=$($(1)_builddir) --build-arg BUILD_BRANCH=$(BRANCH) -t $($(1)_tagname):$(2) -f $($(1)_builddir)/Dockerfile .
	 )
endef

# declare-docker-rule(name, tag)
#define declare-push-rule
#	$(eval $(1)_pushname := $(call rulename, $(1)))                                                                            \
#	$(eval $(1)_tagname := $(CONTAINER_REGISTRY)/$(call remove-step-prefix,$(call directory-name, $(1))))                      \
#	$(eval .PHONY: push-$($(1)_pushname)-$(2))                                                                                 \
#	$(eval                                                                                                                     \
#		push-$($(1)_pushname)-$(2): build-$($(1)_pushname)-$(2) ; $$(info === pushing $($(1)_tagname):$(2)...)
#			docker push $($(1)_tagname):$(2)
#	)
#endef

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

all: buildall


buildall: $(foreach target,$(TARGETS),build-$(call rulename, $(target)))
pushall: $(foreach target,$(TARGETS),push-$(call rulename, $(target)))
BUILD_BASES := $(filter %base, $(foreach target,$(TARGETS),build-$(call rulename, $(target))))
buildallbases: $(shell echo $(BUILD_BASES) | sort)

BASES := $(filter %base, $(foreach target,$(TARGETS),push-$(call rulename, $(target))))
pushbase: $(shell echo $(BASES) | sort)

# create the following rules: build-$target, build-$target-$tag, push-$target, push-$target-$tag
$(foreach target, $(TARGETS),                                                                                                                  \
	$(call declare-shortcut,build-$(call rulename, $(target)),$(foreach tag,$(DOCKER_TAGS),build-$(call rulename, $(target))-$(tag)))          \
	$(call declare-shortcut,push-$(call rulename, $(target)),$(foreach tag,$(DOCKER_TAGS),push-$(call rulename, $(target))-$(tag)))            \
	$(foreach tag,$(DOCKER_TAGS),                                                                                                              \
		$(call declare-build-rule,$(target),$(tag))                                                                                            \
		$(call declare-push-rule,$(target),$(tag))                                                                                             \
	)                                                                                                                                          \
)

# create pack-$target, packall
$(foreach manifest, $(MANIFESTS), $(call declare-manifest-build-rule,$(manifest)))
packall: $(foreach manifest, $(MANIFESTS),pack-$(call rulename, $(manifest)))


$(call declare-shortcut,indexfile,$(MANIFESTDIR)/indexfile.yml)
$(MANIFESTDIR)/indexfile.yml: packall 
	./infra/generate-index-file.py $(MANIFESTDIR) $@

.PHONY: publish-manifests  
publish-manifests: packall $(MANIFESTDIR)/indexfile.yml
	gsutil -m cp -r $(MANIFESTDIR)/*.yml $(MANIFEST_PATH)

.PHONY: publish-manifests  
publish-manifests-no-deps:
	gsutil -m cp -r $(MANIFESTDIR)/*.yml $(MANIFEST_PATH)

.PHONY: validate-vendors
validate-vendors:
	./infra/validate-vendors.py

.PHONY: publish-vendors
publish-vendors:
	gsutil -m cp -r vendors $(VENDORS_PATH)

.PHONY: lint-step-manifests
lint-step-manifests:
	docker run --rm -it -v $(shell pwd):/tmp stoplight/spectral --ignore-unknown-format --ruleset=./tmp/infra/lint/spectral.yaml lint '/tmp/steps/**/manifest.yaml' -vv

.PHONY: clean
clean:
	rm -rf $(OUTPUTDIR)
