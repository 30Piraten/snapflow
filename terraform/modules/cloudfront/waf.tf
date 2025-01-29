resource "aws_wafv2_web_acl" "cloudfront_waf" {
    name = "cloudfront-web-acl"
    scope = "CLOUDFRONT"
    description = "Web ACL for CloudFront distribution"

    rule {
        name = "AWS-Managed-Rule-Baseline"
        priority = 0 
        override_action {
            none {

            }
        }
        statement {
            managed_rule_group_statement {
                vendor_name = "AWS"
                name = "AWSManagedRulesAdminProtectionRuleSet"
            }
        }
        visibility_config {
            cloudwatch_metrics_enabled = true
            metric_name = "AWS-Managed-Rule-Baseline"
            sampled_requests_enabled = true 
        }
    }

    rule {
        name = "validate-content-type"
        priority = 1
        action {
          block {
            
          }
        }
        statement {
            byte_match_statement {
              search_string = "application/json"
              field_to_match {
                single_header {
                    name = "content-type"
                }
              }
              positional_constraint = "EXACTLY"
              text_transformation {
                priority = 0
                type = "NONE"
              }
            }
        }

        visibility_config {
            cloudwatch_metrics_enabled = true
            metric_name = "validate-content-type"
            sampled_requests_enabled = true
        }
    }

    visibility_config {
        cloudwatch_metrics_enabled = true
        metric_name = "cloudfront-web-acl"
        sampled_requests_enabled = true
    }

    default_action {
        allow {}
    }
}