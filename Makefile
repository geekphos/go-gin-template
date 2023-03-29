# ==============================================================================
# 定义全局 Makefile 变量方便后面引用

COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# 项目根目录
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))
# 构建产物、临时文件存放目录
OUTPUT_DIR := $(ROOT_DIR)/_output

# ==============================================================================

## 指定应用使用的 version 包，会通过 `-ldflags -X` 向该包中指定的变量注入值
VERSION_PACKAGE=phos.cc/yoo/pkg/version

## 定义 VERSION 语义化版本号
ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --tags --always --match='v*')
endif

## 检查代码仓库是否是 dirty（默认dirty）
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
  GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD)

GO_LDFLAGS += \
  -X $(VERSION_PACKAGE).GitVersion=$(VERSION) \
  -X $(VERSION_PACKAGE).GitCommit=$(GIT_COMMIT) \
  -X $(VERSION_PACKAGE).GitTreeState=$(GIT_TREE_STATE) \
  -X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')


# ==============================================================================
# 定义 Makefile all 伪目标，执行 `make` 时，会默认会执行 all 伪目标
.PHONY: all
all: tidy clean build

# ==============================================================================
# 定义其他需要的伪目标

.PHONY: build
build: tidy # 编译源码，依赖 tidy 目标自动添加/移除依赖包.
	@go build -v -ldflags "$(GO_LDFLAGS)" -o $(OUTPUT_DIR)/yoo $(ROOT_DIR)/cmd/yoo/main.go


.PHONY: tidy
tidy: # 自动添加/移除依赖包.
	@go mod tidy

.PHONY: clean
clean: # 清理构建产物、临时文件等.
	@-rm -vrf $(OUTPUT_DIR)

.PHONY: gen-model
gen-model: # 生成模型代码.
	@go run $(ROOT_DIR)/cmd/gen/main.go

.PHONY: doc
doc: # 生成文档.
	@swag init --parseInternal --pd -d $(ROOT_DIR)/internal/yoo -g yoo.go