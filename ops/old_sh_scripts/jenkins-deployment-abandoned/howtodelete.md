# How to delete this annoying jenkins namespace?

What is this problem?

Thing is when i try and run

```bash
kubectl delete namespace jenkins
```

it stuck on terminating

I ran into this problem it's so overwhelming it nearly got me crazy

So i don't want you or future me to have this problem

To delete that namespace(or any annoying namespace i reckon) run this ->

```bash
NAMESPACE=jenkins
kubectl proxy &
kubectl get namespace $NAMESPACE -o json |jq '.spec = {"finalizers":[]}' >temp.json
curl -k -H "Content-Type: application/json" -X PUT --data-binary @temp.json 127.0.0.1:8001/api/v1/namespaces/$NAMESPACE/finalize
```

> You will reach the zen :) 
