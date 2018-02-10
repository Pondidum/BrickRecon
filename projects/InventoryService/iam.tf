data "template_file" "inventory_role_policy" {
  template = "${file("policies/lambda-role-policy.json")}"
  vars {
    inventory_table = "${aws_dynamodb_table.inventory_sets.arn}"
    boidcache_table = "${aws_dynamodb_table.boid_cache.arn}"
    kafish_lambda = "arn:aws:lambda:${var.region}:${data.aws_caller_identity.current.account_id}:function:${local.kafish}"
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
