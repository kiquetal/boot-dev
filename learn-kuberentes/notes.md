### DEPLOYING AN IMAGE
The kubectl create deployment command will create a "deployment" for us. We'll talk more about the nuances of "deployments" later. But to put it simply, we only need to provide two things:

The name of the deployment (this can be anything, it's used to identify the deployment)
The ID of the Docker image we want to deploy (it would be a full URL if we weren't hosting the image on Docker Hub, which is the default)


kubectl create deployment synergychat-web --image=bootdotdev/synergychat-web:latest



### Expose the service


kubectl port-forward PODNAME 8080:8080




### Expose the dashboard

minikube dashboard


### Create a proxy

minikube proxy

### Enable addons

minikube addons enable metrics-server
