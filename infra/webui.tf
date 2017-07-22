data "template_file" "s3_public_policy" {
  template = "${file("policies/s3-static-site.json")}"
  vars {
    bucket_name = "${var.bucket}-${var.environment}"
  }
}

resource "aws_s3_bucket" "webui" {
  bucket = "${var.bucket}-${var.environment}"

  acl = "public-read" # lock this down to just the website later
  policy = "${data.template_file.s3_public_policy.rendered}"

  website {
    index_document = "index.html"
  }

  cors_rule {
    allowed_headers = [ "*" ]
    allowed_methods = [ "GET" ]
    allowed_origins = [ "http://localhost:3000" ] # todo, allow the webui when hosted by s3 here
    max_age_seconds = 3000
  }

  tags = {
    environment = "${var.environment}"
  }
}
