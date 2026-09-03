resource "terraform_data" "garage_s3_setup" {
  triggers_replace = [helm_release.garage.id]

  provisioner "local-exec" {
    command = <<-EOT
      set -euo pipefail 
      
      # create key, capture output
      kubectl exec -n garage garage-0 -- /garage key create showcase-backend-key > /tmp/garage_key_raw.txt
      
      KEY_ID=$(grep "Key ID:" /tmp/garage_key_raw.txt | awk '{print $3}')
      SECRET_KEY=$(grep "Secret key:" /tmp/garage_key_raw.txt | awk '{print $3}')

      # create bucket
      kubectl exec -n garage garage-0 -- /garage bucket create showcase-backend

      # allow the key we created for that bucket
      kubectl exec -n garage garage-0 -- /garage bucket allow --read --write --owner showcase-backend --key showcase-backend-key

      # write results as json to tofu back
      cat > ${path.module}/.garage_s3_creds.json <<JSON
      {
        "access_key_id": "$KEY_ID",
        "secret_access_key": "$SECRET_KEY"
      }
      JSON
    EOT
  }
}

data "local_file" "garage_s3_creds" {
  filename    = "${path.module}/.garage_s3_creds.json"
  depends_on = [ terraform_data.garage_s3_setup ]
}

locals {
  garage_s3_creds   = jsondecode(data.local_file.garage_s3_creds.content)
}
