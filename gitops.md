# GitOps

XXX things to do:

- Script to install argoCD in the default namespace, and add a nodeport to it

- Fork the repo
- Access the argo frontend
    - Get password: `kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d`
    - kubectl get svc -n argocd
    - kubectl get nodes -o wide
- login
- create a new rollout
- test that it is working
- change something in the helm charts
- see that it syncs, and replaces
- try to change something in the cluster
- see that Argo goes in and changes things back again.
