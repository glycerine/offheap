PROJECT = offheap
PACKAGE = github.com/remerge/$(PROJECT)

GOMETALINTER_OPTS = --enable-all --tests --errors -D vet

include Makefile.common
