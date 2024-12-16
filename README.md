# distributed-systems-24-winter

This is a backend application for **Distributed Systems 24 Winter** project.

## Setup Instructions with Docker
1. Open this repository in a new Codespace in GitHub
2. Change the Directory with: cd shoppinglist
3. Open a Terminal and type in the command: docker compose up
4. Then the docker Images get pulled and it runs the Docker Containers (Backend, Database, Frontend)
5. You can now check your endroutes via Swagger: <CodespaceURL>/swagger/index.html / OR the frontend application

## Setup Instructions with Kubernetes
1. Open this repository in a new Codespace in GitHub
2. Open a Terminal type in: "Minikube start"  
3. Change the directory "cd shoppinglist/k8s"
4. Type in the terminal "kubectl apply -f . " 
5. Verify if all pods Running correctly: "kubectl get pods"

