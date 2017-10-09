# stuff-to-be-done
framework for doing stuff

##Setup
Install GO environment for App Engine
https://cloud.google.com/go/quickstarts

Install Node
Install Bower (replace later with webpack)
bower install in the web folder
npm install -g vulcanize
? gcloud components install kubectl

be aware of 64bit visual studio code running python causing issues with deployment

##run
goapp serve
goapp serve web\
dev_appserver.py web\app.yaml

##set go path [this should be done in the system environment]
set GOPATH=C:\Work;C:\Users\Gregory Swanson\Documents\GitHub\stuff-to-be-done\web
https://github.com/golang/go/wiki/Setting-GOPATH#windows
go get github.com/gin-gonic

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