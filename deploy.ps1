vulcanize web/public/elements/elements.html > web/public/js/elements.html --strip-comments --inline-scripts --inline-css
Copy-Item -Path "web/bower_components/webcomponentsjs/webcomponents.min.js" -Destination "web/public/js/webcomponents.min.js"
gcloud app deploy web/app.yaml web/index.yaml --no-promote -v preview --project fuzzyhipster --quiet