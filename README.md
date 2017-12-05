# stuff-to-be-done
framework for doing stuff

##CloudShell (from app engine)
https://cloud.google.com/shell/docs/quickstart
http://www.ryanday.net/2012/10/01/installing-go-and-gopath/
###Update the .bashrc
```
cd $HOME
vi ./.bashrc
G (move to last line of file)
o (open a new line after)
export GOROOT $HOME/work
export PATH=$PATH:$GOROOT/bin
export GOPATH=$GOPATH:$HOME/stuff-to-be-done/web
[ESC]
:X (quit and save changes)
```

##Setup
Install GO environment for App Engine
https://cloud.google.com/go/quickstarts

Install Node
Install Bower (replace later with webpack)
bower install in the web folder
npm install -g vulcanize
on cloudshell
sudo -i npm install -g vulcanize
? gcloud components install kubectl

be aware of 64bit visual studio code running python causing issues with deployment

windows can have an issue with the gopath if it has a trailing ;
https://github.com/golang/go/issues/21928


##run
goapp serve
goapp serve web\
dev_appserver.py web\app.yaml
dev_appserver.py app.yaml --host $IP --port $PORT

##format source
from root
gofmt -w web

##set go path [this should be done in the system environment]
set GOPATH=C:\Work;C:\Users\Gregory Swanson\Documents\GitHub\stuff-to-be-done\web

[NO] go get github.com/gin-gonic
go get github.com/gin-gonic/gin
sudo -i go get github.com/gin-gonic/gin
export GOPATH=/work;/gae/web
export GOPATH="/home/ubuntu/workspace/work:/home/ubuntu/workspace/jdt/web"

on c9.io
export PATH=$PATH:/home/ubuntu/workspace/go/bin

get the source into /home/ubuntu/workspace/work/src/jdt

Deploy
    // copy all dependancies to the gopath
    del C:\Work\src\stufftobedone
    mkdir C:\Work\src\stufftobedone
    xcopy /E web\stufftobedone C:\Work\src\stufftobedone
    gcloud app deploy app.yaml --no-promote -v preview --project fuzzyhipster
    del C:\Work\src\stufftobedone
    Remove-Item -Recurse -Force C:\Work\src\stufftobedone

Pre build
    vulcanize web/public/elements/elements.html > web/public/js/elements.html

    vulcanize web/public/elements/elements.html > web/public/js/elements.html