kubectl create namespace argocd
kubectl create namespace jenkins
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
kubectl patch svc argocd-server -n argocd -p '{\"spec\": {\"type\": \"LoadBalancer\"}}'
kubectl patch svc argocd-server -n argocd --type='json' -p '[{"op": "replace", "path": "/spec/ports/0/port", "value": 81}]'

do {
    try {
        $initial_password = (argocd admin initial-password -n argocd 2>$null).Split("`n")[0].Trim()
        $success = $true
    } catch {
        Write-Host "ArgoCD not ready yet.. trying again.."
        Start-Sleep -Seconds 1
        $success = $false
    }
} while (-not $success)
Write-Host "ArgoCD is ready. Initial password retrieved: "$initial_password -ForegroundColor Green

Write-Host "Trying to login to ArgoCD. (This might take a minute)" -ForegroundColor Yellow
$loginSuccess = $false
do {
    try {
        $argocdOutput = & argocd login localhost:81 --username admin --password $initial_password --insecure 2>&1
    
        if ($argocdOutput -like "*context deadline exceeded*") {
            Write-Host "ArgoCD login failed.. trying again.."
            Start-Sleep -Seconds 1
        }
        else {
            $loginSuccess = $true
        }
    } catch {
        $errorMessage = $_.Exception.Message
        Write-Host "ArgoCD login failed.. trying again.."
    }
} while (-not $loginSuccess)
Write-Host "Logged into ArgoCD." -ForegroundColor Green

argocd account update-password --current-password $initial_password --new-password password
Write-Host "ArgoCD password updated to 'user:admin pw:password'" -ForegroundColor Green
               
argocd app create clustereye --repo https://github.com/urmas-villem/ClusterEyeGo.git --path kubernetes --dest-server https://kubernetes.default.svc --dest-namespace default --sync-policy automated

kubectl apply -n argocd -f https://raw.githubusercontent.com/urmas-villem/ClusterEyeGo/main/kubernetes/Jenkins/jenkins.yaml

Write-Host "Set up dummy software for Clustereye to pick up" -ForegroundColor Yellow
kubectl apply -n argocd -f https://raw.githubusercontent.com/urmas-villem/ClusterEyeGo/main/kubernetes/Testsoftware/argoapp/testsoftware-setup.yaml

$dockerPassword = Read-Host "Enter Docker registry password" -AsSecureString
$plainPassword = [System.Runtime.InteropServices.Marshal]::PtrToStringBSTR([System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($dockerPassword))
kubectl create secret docker-registry docker-credentials --docker-username=huxlee --docker-password=$plainPassword --docker-email=random@random.com --namespace jenkins

#$webhook = Read-Host "Enter slack webhook secret"
#kubectl create secret generic slack-webhook-secret --from-literal=SLACK_WEBHOOK_URL=$webhook

Write-Host "Docker secret created" -ForegroundColor Green
Write-Host "***************************"
Write-Host "**Environment setup ready**"
Write-Host "***************************"  
