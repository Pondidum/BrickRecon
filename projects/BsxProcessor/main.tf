data "template_file" "bsxprocessor_policy" {
  template = "${file("policies/bsxprocessor-role-policy.json")}"
  vars {
    bucket = "${var.bucket}"
  }
}

resource "aws_iam_role" "bsxprocessor_role" {
  name = "${local.name}_role"
  assume_role_policy = "${file("policies/bsxprocessor-role.json")}"
}

resource "aws_iam_role_policy" "bsxprocessor_role_policy" {
  name = "${local.name}_role_policy"
  role = "${aws_iam_role.bsxprocessor_role.id}"
  policy = "${data.template_file.bsxprocessor_policy.rendered}"
}

data "archive_file" "bsxprocessor_source" {
  type = "zip"
  source_dir = "./src/.build"
  output_path = ".build/bsxprocessor.zip"
}

resource "aws_lambda_function" "bsxprocessor" {
  function_name = "${local.name}"
  role = "${aws_iam_role.bsxprocessor_role.arn}"
  filename = "${data.archive_file.bsxprocessor_source.output_path}"
  handler = "BsxProcessor::BsxProcessor.Handler::Handle"
  runtime = "dotnetcore1.0"
  source_code_hash = "${base64sha256(file("${data.archive_file.bsxprocessor_source.output_path}"))}"
  timeout = 60

  tags = {
    environment = "${var.environment}"
  }
}

resource "aws_s3_bucket_notification" "bsxprocessor_trigger" {
  bucket = "${var.bucket}"

  lambda_function {
    lambda_function_arn = "${aws_lambda_function.bsxprocessor.arn}"
    events = ["s3:ObjectCreated:*"]
    filter_prefix = "upload/"
    filter_suffix = ".bsx"
  }
}
