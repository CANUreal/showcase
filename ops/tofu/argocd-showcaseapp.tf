resource "kubernetes_manifest" "showcase_backend_app" {
 manifest = {
  apiVersion = "argoproj.io/v1alpha1"
    kind       = "Application"
    metadata = {
      name      = "showcase-backend"
      namespace = "argocd"
    }
    spec = {
      project = "default"
      source = {
        repoURL        = "https://github.com/CANUreal/showcase-gitops.git"
        targetRevision = "main"
        path           = "apps/showcase-backend"
      }
      destination = {
        server    = "https://kubernetes.default.svc"
        namespace = "showcase"
      }
      syncPolicy = {
        automated = {
          prune    = true
          selfHeal = true
        }
        syncOptions = ["CreateNamespace=true"]
      }
    }  
  }

  depends_on = [ helm_release.argocd ]
}
