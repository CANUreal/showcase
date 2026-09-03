resource "kubernetes_namespace" "showcase" {
 metadata {
   name = "showcase"
 } 
}

resource "kubernetes_secret" "showcase_backend" {
  metadata {
    name = "showcase-backend-secrets"
    namespace = kubernetes_namespace.showcase.metadata[0].name
  }

  data = {
    AWS_ACCESS_KEY_ID   = local.garage_s3_creds.access_key_id
    AWS_SECRET_ACCESS_KEY = local.garage_s3_creds.secret_access_key
    DATABASE_URL =  "postgres://showcase:${random_password.postgres.result}@showcase-pg-postgresql.postgres.svc.cluster.local:5432/showcase?sslmode=disable"
  }

  type = "Opaque"
}
