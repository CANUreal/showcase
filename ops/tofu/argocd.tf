resource "helm_release" "argocd" {
  name    = "argocd"
  namespace = "argocd"
  create_namespace = true

  repository = "https://argoproj.github.io/argo-helm"
  chart   = "argo-cd"
  # this is the latest chart, but we gotta pin it. bcz we don't want no unwanted updates.
  version = "10.7.1"  

  wait = true
  timeout = 300
} 

resource "kubernetes_manifest" "argocd_ingress" {
  manifest = {
    apiVersion = "traefik.io/v1alpha1"
    kind: "IngressRouteTCP" 
    metadata = {
      name  = "argocd-server"
      namespace = "argocd"
    }

    spec = {
      entryPoints = [ "websecure" ]
      routes = [
        {
          match   = "HostSNI(`argocd.local`)"
          services = [
            {
              name = "argocd-server"
              port = 443
            }
          ]
        }
      ]
      tls = {
        passthrough = true
      }
    }
  }

  depends_on = [ helm_release.argocd ]
}
