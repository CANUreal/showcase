resource "helm_release" "garage" {
    name        =   "garage"
    namespace   =   "garage"
    create_namespace = true

    chart   = "${path.module}/charts/garage"
    # Like it says in the chart from garage writers that i cloned in this dir 
    version = "0.9.3"
    
    values = [
        file("${path.module}/values/garage.yaml")
    ]

    wait = true
    timeout = 300
}
