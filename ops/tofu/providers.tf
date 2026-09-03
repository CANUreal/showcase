terraform {
    required_version = ">= 1.6"
    
    required_providers {
        helm = {
            source = "hashicorp/helm"
            version = "~> 2.15"
        }  
        kubernetes = {
            source = "hashicorp/kubernetes"
            version = "~> 2.33"
        }
        random = {
          source = "hashicorp/random"
          version = "~> 3.6"
        }
    } 
}

provider "helm" {
    kubernetes {
        config_path = "~/.kube/config"
    }
    repository_config_path = "~/.config/helm/repositories.yaml"
    repository_cache = "~/.cache/helm"
}

provider "kubernetes" {
    config_path = "~/.kube/config"
}
