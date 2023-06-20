#!/bin/bash

# 设置生成文件的数量和内容
file_count=30
file_content="This is the contentssssssssssssssssssssdfdsfsdfsdfadsf asdfasdfasdfasdfasdfasdfasdfasdfasdfasdf of the file."

# 循环生成文件
for ((i=1; i<=file_count; i++)); do
    echo "$file_content" > "file${i}.txt"
done

