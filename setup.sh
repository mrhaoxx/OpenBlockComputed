#!/bin/bash

# run only as root
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

# get os info
os=$(uname -s)


# get cpu arch
cpu=$(uname -m)

# get cpu sku
cpusku=$(lscpu | grep "Model name" | cut -d':' -f2 | xargs)
cputotalcores=$(lscpu | grep "Core(s) per socket" | cut -d':' -f2 | xargs)
cpusockets=$(lscpu | grep "Socket(s)" | cut -d':' -f2 | xargs)

# get gpu info
gpu=$(lspci | grep -i nvidia | cut -d':' -f3 | xargs)

# get ram info
ram=$(free -h | grep Mem | awk '{print $2}')

# get network info
network=$(lspci | grep -i network | cut -d':' -f3 | xargs)


ip=$(hostname -I | cut -d' ' -f1)

# generate json
echo "{
  \"os\": \"$os\",
  \"arch\": \"$cpu\",
  \"cpusku\": \"$cpusku\",
  \"cpucores\": \"$cputotalcores\",
  \"cpusockets\": \"$cpusockets\",
  \"gpu\": \"$gpu\",
  \"ram\": \"$ram\",
  \"network\": \"$network\",
  \"ip\": \"$ip\"
}"


