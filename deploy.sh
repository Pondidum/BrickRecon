#! /bin/sh

./build.sh

(cd projects/webui && terraform apply)
(cd projects/ImageCache && terraform apply)
(cd projects/BsxProcessor && terraform apply)
