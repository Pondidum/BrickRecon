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

  tags = {
    environment = "${var.environment}"
  }
}

resource "aws_s3_bucket_object" "index" {
  bucket = "${aws_s3_bucket.webui.bucket}"
  key = "index.html"
  source = "../src/webui/index.html"
  content_type = "text/html"
  etag = "${md5(file("../src/webui/index.html"))}"
}
