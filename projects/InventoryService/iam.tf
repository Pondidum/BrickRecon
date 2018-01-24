data "template_file" "inventory_role_policy" {
  template = "${file("policies/lambda-role-policy.json")}"
  vars {
    inventory_table = "${aws_dynamodb_table.inventory_sets.arn}"
    boidcache_table = "${aws_dynamodb_table.boid_cache.arn}"
    sns_arn = "${data.aws_sns_topic.kafish.arn}"
  }
}

resource "aws_iam_role" "inventory_role" {
  name = "${local.name}_role"
  assume_role_policy = "${file("policies/lambda-role.json")}"
}

resource "aws_iam_role_policy" "inventory_role_policy" {
  name = "${local.name}_role_policy"
  role = "${aws_iam_role.inventory_role.id}"
  policy = "${data.template_file.inventory_role_policy.rendered}"
}
