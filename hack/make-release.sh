#!/bin/sh

# cat the *.yaml files into the final file: balancer.yaml

DEPLOY_FILE=deploy/balancer.yaml

rm $DEPLOY_FILE 2>/dev/null
touch $DEPLOY_FILE
for file in ./deploy/cluster_role.yaml ./deploy/cluster_role_binding.yaml ./deploy/crds/hliangzhao_v1alpha1_balancer_crd.yaml ./deploy/operator.yaml ./deploy/service_account.yaml
do
    echo "---" >> $DEPLOY_FILE
    cat $file >> $DEPLOY_FILE
done
