resource "random_password" "postgres" {
    length = 20
    special     = false
}

resource "helm_release" "postgres" {
    name    = "showcase-pg"
    namespace = "postgres"
    create_namespace = true

    repository  = "oci://registry-1.docker.io/bitnamicharts"
    chart       = "postgresql"
    version = "18.8.13"

    set {
        name = "auth.username"
        value = "showcase"
    }
    set {
        name = "auth.database"
        value = "showcase"
    }
    set_sensitive {
        name = "auth.password"
        value = random_password.postgres.result
    }
    set {
        name = "primary.persistence.size"
        value = "6Gi"
    }
    
    wait = true
    timeout = 300
}
