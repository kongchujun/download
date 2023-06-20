# install vagrant and virtualbox
```
https://www.vagrantup.com/
https://www.virtualbox.org/wiki/Downloads
```

# how to setup vagrant env
```
open Powershell and navigate folder which contains Vagrantfile
vagrant up 
just in case: vagrant plugin install vagrant-disksize
```

# enter into vagrant env and start docker
```
vagrant ssh
systemctl start docker
```

# when generate the image
change config file path in Go_download/cmd/myapp/main.go

# for test multiprocess
```
#!/bin/bash

# 设置生成文件的数量和内容
file_count=30
file_content="This is the contentssssssssssssssssssssdfdsfsdfsdfadsf asdfasdfasdfasdfasdfasdfasdfasdfasdfasdf of the file."

# 循环生成文件
for ((i=1; i<=file_count; i++)); do
    echo "$file_content" > "file${i}.txt"
done
```