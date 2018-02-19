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

  environment {
    variables = {
      IMAGECACHE_LAMBDA = "${local.imagecache_lambda}"
    }
  }

  tags = {
    environment = "${var.environment}"
  }
}
