data "archive_file" "inventory_lambda_source" {
  type = "zip"
  source_dir = "./build/app"
  output_path = "./build/lambda.zip"
}

resource "aws_lambda_function" "inventory_service" {
  filename = "${data.archive_file.inventory_lambda_source.output_path}"
  function_name = "${local.name}"
  role = "${aws_iam_role.inventory_role.arn}"
  handler = "aws.handler"
  runtime = "nodejs6.10"
  source_code_hash = "${base64sha256(file("${data.archive_file.inventory_lambda_source.output_path}"))}"

  environment {
    variables = {
      TABLE_NAME = "${local.table_name}"
      SNS_TOPIC = "${data.aws_sns_topic.kafish.arn}"
    }
  }
}
