resource "kubernetes_manifest" "showcase-ingress" {
  manifest = {
    apiVersion = "networking.k8s.io/v1"
    kind = "Ingress"

    metadata = {
      name = "showcase-ingress"
      namespace = "showcase"
    }

    spec = {
      ingressClassName = "traefik"
      
      rules = [
        {
          host = "showcase.local"

          http = {
            paths = [
              {
                path = "/"
                pathType = "Prefix"
                

                backend = {
                  service = {
                    name = "showcase-backend"

                    port =  {
                      number = 80
                    }
                  }
                }
              }
            ]
          }
        }
      ]
    }
  }
}
