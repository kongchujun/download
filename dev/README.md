# after download project please use mod init command to build it
```
cd Go-download
rm -f go.mod go.sum
go mod init godownload
go mod tidy
```

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

# generate files for test
file_count=30
file_content="This is the contentssssssssssssssssssssdfdsfsdfsdfadsf asdfasdfasdfasdfasdfasdfasdfasdfasdfasdf of the file."

# run loop to create it
for ((i=1; i<=file_count; i++)); do
    echo "$file_content" > "file${i}.txt"
done
```

# set up debug for gin in vscode
```
create .vscode folder in Go_download folder, then create launch.json and add content below in it
{
    "version": "0.2.0",
    "configurations": [
      {
        "name": "Launch Gin",
        "type": "go",
        "request": "launch",
        "mode": "debug",
        "program": "${workspaceFolder}",
        "env": {},
        "args": []
      }
    ]
  }
```