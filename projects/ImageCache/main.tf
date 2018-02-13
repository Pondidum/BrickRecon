data "template_file" "imagecache_policy" {
  template = "${file("policies/imagecache-role-policy.json")}"
  vars {
    bucket = "${var.bucket}"
  }
}

resource "aws_iam_role" "imagecache_role" {
  name = "${local.name}role"
  assume_role_policy = "${file("policies/imagecache-role.json")}"
}

resource "aws_iam_role_policy" "imagecache_role_policy" {
  name = "${local.name}role_policy"
  role = "${aws_iam_role.imagecache_role.id}"
  policy = "${data.template_file.imagecache_policy.rendered}"
}



data "archive_file" "imagecache_source" {
  type = "zip"
  source_dir = "./build/app"
  output_path = "./build/imagecache.zip"
}

resource "aws_lambda_function" "imagecache" {
  function_name = "${local.name}"
  role = "${aws_iam_role.imagecache_role.arn}"
  filename = "${data.archive_file.imagecache_source.output_path}"
  handler = "index.handler"
  runtime = "nodejs6.10"
  timeout = 20
  source_code_hash = "${base64sha256(file("${data.archive_file.imagecache_source.output_path}"))}"

  tags = {
    environment = "${var.environment}"
  }

  environment {
    variables = {
      bucket = "${var.bucket}"
    }
  }
}
