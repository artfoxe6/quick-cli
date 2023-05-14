#!/bin/bash

# 获取当前目录
current_dir=$(pwd)

# 检查参数
if [ $# -eq 0 ]; then
  echo "使用方法：./build.sh <格式>"
  echo "可用格式：linux, mac, windows"
  exit 1
fi

# 设置输出目录
output_dir="${current_dir}/bin"
mkdir -p "${output_dir}"

# 遍历当前目录下的Go程序文件
for file in *.go; do
  filename="${file%.*}"
  output_file="${output_dir}/${filename}"

  # 根据参数选择打包格式
  case "$1" in
    linux)
      GOOS=linux GOARCH=amd64 go build -o "${output_dir}/quick-cli"
      ;;
    mac)
      GOOS=darwin GOARCH=amd64 go build -o "${output_dir}/quick-cli-mac"
      ;;
    windows)
      GOOS=windows GOARCH=amd64 go build -o "${output_dir}/quick-cli.exe"
      ;;
    *)
      echo "不支持的格式：$1"
      exit 1
      ;;
  esac
done

echo "打包完成！"
