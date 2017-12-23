data "template_file" "kafish_role_policy" {
  template = "${file("policies/lambda-role-policy.json")}"
  vars {
    table_arn = "${aws_dynamodb_table.event_store.arn}"
    sns_arn = "${aws_sns_topic.kafish_events.arn}"
  }
}

resource "aws_iam_role" "kafish_role" {
  name = "${local.name}_role"
  assume_role_policy = "${file("policies/lambda-role.json")}"
}

resource "aws_iam_role_policy" "kafish_role_policy" {
  name = "${local.name}_role_policy"
  role = "${aws_iam_role.kafish_role.id}"
  policy = "${data.template_file.kafish_role_policy.rendered}"
}
