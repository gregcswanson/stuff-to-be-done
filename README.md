# stuff-to-be-done
framework for doing stuff

##run
goapp serve
goapp serve web\

##set go path [this should be done in the system environment]
set GOPATH=C:\Work

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