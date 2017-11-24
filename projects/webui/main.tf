data "template_file" "s3_public_policy" {
  template = "${file("policies/s3-static-site.json")}"
  vars {
    bucket_name = "${local.bucket}"
  }
}

resource "aws_s3_bucket" "webui" {
  bucket = "${local.bucket}"

  acl = "public-read" # lock this down to just the website later
  policy = "${data.template_file.s3_public_policy.rendered}"

  website {
    index_document = "index.html"
  }

  cors_rule {
    allowed_headers = [ "*" ]
    allowed_methods = [ "GET" ]
    allowed_origins = [
      "http://localhost:3000",
      "http://${local.bucket}.s3-website-${var.region}.amazonaws.com",
      "http://${local.bucket}"
    ]
    max_age_seconds = 3000
  }

  tags = {
    environment = "${var.environment}"
  }
}

output "event_api_url" {
  value = "http://${local.bucket}.s3-website-${var.region}.amazonaws.com"
}
