**Project structure**
* proxy - TCP proxy handles sentinel, redis and client connections (written in Golang)
* redis - contains Redis modules
    * BUILDING_MODULE - CHECKS IF VALUE EXCIST AND INCREASES VALUE WHEN TRIGGERED (VALID VALUES 100000000-199999999)
    * UTILITY_BUILDING_MODULE - CHECKS IF VALUE EXCIST AND INCREASES VALUE WHEN TRIGGERED (VALID VALUES 200000000-299999999)
    * PROCEDURE_MODULE - CHECKS IF VALUE EXCIST AND INCREASES VALUE WHEN TRIGGERED (VALID VALUES 0...)
    * YEAR_MODULE - CHECKS IF YEAR MATCHES ELSE RESETS ALL VALUES IN DOCUMENT_MODULE TO 0
    * DOCUMENT_MODULE - CHECKS IF DOCUMENT EXCIST AND INCREASES COUNT (VALID VALUES (YEAR TWO LAST DIGITS)+DOTY + / + COUNT -> EXAMPLE 201111/00001)
* redis-sentinel-proxy-charm - CONTAINS HELM CHART STUFFS





**kubernetes rights:**
```
cat <<EOF | kubectl apply -f -
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: read-pods
  namespace: kube-system
subjects:
  - kind: ServiceAccount
    name: default
    namespace: gitlab-managed-apps
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io
EOF
```