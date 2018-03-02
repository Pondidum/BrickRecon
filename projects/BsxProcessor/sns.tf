data "aws_sns_topic" "kafish" {
  name = "${local.sns_topic}"
}

resource "aws_sns_topic_subscription" "lambda_trigger" {
  topic_arn = "${data.aws_sns_topic.kafish.arn}"
  protocol  = "lambda"
  endpoint  = "${aws_lambda_function.bsxprocessor.arn}"
}

resource "aws_lambda_permission" "with_sns" {
  statement_id = "AllowExecutionFromSNS"
  action = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.bsxprocessor.arn}"
  principal = "sns.amazonaws.com"
  source_arn = "${data.aws_sns_topic.kafish.arn}"
}