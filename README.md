# Getting Started on Okteto with Go

[![Develop on Okteto](https://okteto.com/develop-okteto.svg)](https://cloud.okteto.com/deploy?repository=https://github.com/okteto/go-getting-started)

This example shows how to use the [Okteto CLI](https://github.com/okteto/okteto) to develop a Go Sample App directly in Kubernetes. The Go Sample App is deployed using Kubernetes manifests.

This is the application used for the [Getting Started on Okteto with Go](https://www.okteto.com/docs/samples/golang/) tutorial.


# Run 

## unit tests
```
go test
```

## run the app

```
go run main.go
```


# API
## get pods count
```
https://hello-world-reetasingh.jdm.okteto.net/pods/count
```

## get sorted list of pods 
1. by name ```https://hello-world-reetasingh.jdm.okteto.net/pods/list/name```
2. by age ```https://hello-world-reetasingh.jdm.okteto.net/pods/list/age```
3. by restart count ```https://hello-world-reetasingh.jdm.okteto.net/pods/list/restarts```
