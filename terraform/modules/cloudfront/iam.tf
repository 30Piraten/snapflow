# Description: IAM roles and policies for Cloudfront and WAFV2
resource "aws_iam_role" "cloudfront_role" {
  name = "snapflow-cloudfront-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "cloudfront.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
  
}

# Origin Access Identity (OAI) for Cloudfront
resource "aws_cloudfront_origin_access_identity" "snapflow_oai" {
  comment = "Snapflow Cloudfront OAI"
}

# WAFV2 WebACL permissions
resource "aws_iam_policy" "wafv2_permissions" {
  name = "snapflow-wafv2-permissions"
  description = "IAM policy for WAFV2 permissions"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "wafv2:DescribeWebACL",
          "wafv2:AssociateWebACL",
          "wafv2:DisassociateWebACL"
        ]
        Resource = "${aws_wafv2_web_acl.cloudfront_waf.arn}"
      },
      {
        Effect = "Allow"
        Action = [
          "cloudfront:UpdateDistribution",
          "cloudfront:GetDistributionConfig",
          "cloudfront:DescribeDistribution"
        ]
        Resource = "${aws_cloudfront_distribution.snapflow_cloudfront.arn}"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "wafv2_policy_attachment" {
  role = aws_iam_role.cloudfront_role.name
  policy_arn = aws_iam_policy.wafv2_permissions.arn
}