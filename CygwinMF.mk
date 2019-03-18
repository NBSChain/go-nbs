ifeq ($(PLATFORM), Cygwin)
    INCLUDE := ${shell echo "$(INCLUDE)" | sed -e 's/://g'}
    INCLUDE := /cygdrive/$(INCLUDE)
endif

