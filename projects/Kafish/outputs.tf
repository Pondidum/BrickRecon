output "writer.arn" {
  value = "${aws_lambda_function.kafish_writer.arn}"
}

output "reader.arn" {
  value = "${aws_lambda_function.kafish_writer.arn}"
}
