**Project structure**
* proxy - TCP proxy handles sentinel, redis and client connections (written in Golang)
    * main.go - starts proxy
    * init.go - initialises proxy values
    * proxy.go - handles client connection checks if request is valid and makes redis request and returns it
    * readiness.go - needed for kubernetes if pod is ready
    * redis.go - creates redis pool, handles redis connections and redis init if needed.
    * sentinel.go - tcp connection with sentinel gets new redis master and creates redis pool
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