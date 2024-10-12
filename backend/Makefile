# 定义应用程序名称和输出目录
APP_NAME=abandonlgzj
OUTPUT_DIR=./bin

# 创建输出目录
$(OUTPUT_DIR):
	mkdir -p $(OUTPUT_DIR)

# 默认目标，执行全部构建
all: build-macOS-arm build-linux-amd build-linux-arm

.PHONY: build-macOS-arm build-linux-amd build-linux-arm clean

# 构建 macOS 版本的程序
build-macOS-arm: $(OUTPUT_DIR)
	GOOS=darwin GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(APP_NAME)-macOS-arm64 .

# 构建 Linux 版本的程序
build-linux-amd: $(OUTPUT_DIR)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(APP_NAME)-linux-amd64 .

build-linux-arm: $(OUTPUT_DIR)
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(APP_NAME)-linux-arm .

# 清理生成的可执行文件
clean:
	rm -rf $(OUTPUT_DIR)

# 添加一个运行测试的目标
test:
	go test ./...