resource "aws_s3_bucket" "happy" {
    bucket = "digitalroute-happy-button-lambda"
    acl    = "private"
}