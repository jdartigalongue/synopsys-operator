NS=$1

# echo "ns" $NS

oc delete ns $NS

oc delete crd alerts.synopsys.com
oc delete crd hubs.synopsys.com
oc delete crd opssights.synopsys.com


oc delete clusterrolebinding blackduck-operator-admin
#oc delete clusterrolebinding protoform-admin
#oc delete clusterrolebinding blackduck-operator-cluster-admin