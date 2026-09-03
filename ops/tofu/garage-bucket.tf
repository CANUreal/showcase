resource "terraform_data" "garage_s3_key" {
  triggers_replace = [helm_release.garage.id]

  provisioner "local-exec" {
    command = <<-EOT
      set -e 
      kubectl exec -n garage garage-0 -- /garage key create showcase-backend-key 2>&1 | tee /tmp/garage_key_output.txt
      EOT
  }
}
