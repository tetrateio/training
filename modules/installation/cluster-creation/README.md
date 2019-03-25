First, we need to create a cluster:

```
$ gcloud container clusters create --machine-type n1-standard-2 --num-nodes 4 workshop --async
```