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
