resource "aws_sns_topic" "kafish_events" {
  name = "${local.name}_events"
}
