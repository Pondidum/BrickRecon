#! /bin/sh

./build.sh

$variables="../variables.dev.json"

(cd projects/kafish && terraform apply --var-file $variables)
(cd projects/webui && terraform apply --var-file $variables)
(cd projects/ImageCache && terraform apply --var-file $variables)
(cd projects/BsxProcessor && terraform apply --var-file $variables)

(cd projects/webui && ./upload.sh)
