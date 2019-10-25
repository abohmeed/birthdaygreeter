[![CircleCI](https://circleci.com/gh/abohmeed/birthdaygreeter.svg?style=svg)](https://circleci.com/gh/abohmeed/birthdaygreeter) ![GitHub commit activity](https://img.shields.io/github/commit-activity/w/abohmeed/birthdaygreeter?style=plastic)
# Install helm on the target machine as well as on the cluster
Use Ansible to install Helm on the target machine or on localhost
Run the following commands to install Tiller on the server
```bash
kubectl create serviceaccount --namespace kube-system tiller
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
#kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'
helm --service-account tiller
```

# Installing the helm chart

```bash
helm install stable/redis
```

You must create a shared secret on the default namespace. The secret must be a file called `redis-password`. Then the chart can be installed as follows:

```bash
helm install --name bg-redis --set password=mypassword stable/redis
```

