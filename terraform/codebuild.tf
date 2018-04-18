module "github_ci" {
  source  = "github.com/squidfunk/terraform-aws-github-ci"
  version = "0.3.0"

  github_owner       = "hussfelt"
  github_repository  = "dr-happy-push"
  github_oauth_token = "${var.oauth_token}"
  namespace          = "digitalroute-happy-button"
  codebuild_image    = "aws/codebuild/golang:1.10"
}

data "aws_iam_policy_document" "additional_buckets" {
  statement {
    sid = "AllowAdditionalBuckets"
    effect = "Allow"
    actions = [
      "s3:GetBucketVersioning",
      "s3:GetObject",
      "s3:GetObjectVersion",
      "s3:PutObject"
    ]
    resources = [
      "${aws_s3_bucket.happy.arn}/*",
    ]
  }

  statement {
    sid = "AllowCloudFormationDeploy"
    effect = "Allow"
    actions = [
      "cloudformation:Describe*",
      "cloudformation:List*",
      "cloudformation:Get*"
    ]
    resources = [
      "*",
    ]
  }
}

# aws_iam_policy.codebuild 
resource "aws_iam_policy" "codebuild_additional_buckets" { 
  name = "CodebuildAdditionalBuckets" 
  path = "/dr-happy-push/codebuild/" 
  
  policy = "${data.aws_iam_policy_document.additional_buckets.json}" 
}

# aws_iam_policy_attachment.codebuild 
resource "aws_iam_policy_attachment" "codebuild" { 
  name = "HappyButtonCodebuildAdditionalBuckets"

  policy_arn = "${aws_iam_policy.codebuild_additional_buckets.arn}" 
  roles      = ["${module.github_ci.codebuild_service_role_name}"] 
}