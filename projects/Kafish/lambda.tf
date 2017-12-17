data "archive_file" "kafish_lambda_source" {
  type = "zip"
  source_dir = "./build/app"
  output_path = "./build/lambda.zip"
}

resource "aws_lambda_function" "kafish_writer" {
  filename = "${data.archive_file.kafish_lambda_source.output_path}"
  function_name = "${local.name}_writer"
  role = "${aws_iam_role.kafish_role.arn}"
  handler = "index.writeHandler"
  runtime = "nodejs6.10"
  source_code_hash = "${base64sha256(file("${data.archive_file.kafish_lambda_source.output_path}"))}"
  publish = true

  environment {
    variables = {
      TABLE_NAME = "${local.table_name}"
    }
  }
}
