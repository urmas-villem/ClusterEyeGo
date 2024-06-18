#### Dockerhub link:
https://hub.docker.com/r/huxlee/clustereye   

#### Run the application on kubernetes
### Deployment option 1 (windows environment)
## Prerequisites:
- [ArgoCD cli](https://argo-cd.readthedocs.io/en/stable/cli_installation/)
- [Helm](https://helm.sh/docs/intro/install/)
                                                                                
Run this command in a powershell window.     
```   
Invoke-Expression (Invoke-WebRequest -Uri "https://raw.githubusercontent.com/urmas-villem/ClusterEye/main/kubernetes/ArgoCD/setup.ps1").Content
```
This will:                                              
- setup argocd (on localhost:81)
- setup jenkins(with all of the prerequisites and the pipeline job already configured) on (localhost:8080)
- setup the ClusterEye application (on localhost)
- setup some dummy applications so Clustereye has something to check otherwise the input is empty
