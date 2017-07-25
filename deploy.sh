#! /bin/sh

./build.sh

(cd infra && terraform plan -out=infra.tfplan)
read -p "Press enter to confirm"
(cd infra && terraform apply infra.tfplan)
