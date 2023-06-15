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