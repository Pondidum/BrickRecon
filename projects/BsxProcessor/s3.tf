resource "aws_s3_bucket_notification" "bsxprocessor_trigger" {
  bucket = "${var.bucket}"

  lambda_function {
    lambda_function_arn = "${aws_lambda_function.bsxprocessor.arn}"
    events = ["s3:ObjectCreated:*"]
    filter_prefix = "upload/"
    filter_suffix = ".bsx"
  }
}
