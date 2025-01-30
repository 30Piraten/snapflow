// make use of s3 bucket config here 

resource "aws_cloudfront_distribution" "snapflow_cloudfront" {
  origin {
    domain_name = var.domain_name
    origin_id   = var.origin_id
    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.snapflow_oai.cloudfront_access_identity_path
    }
  }

  enabled = true
  is_ipv6_enabled = true
  comment = "Snapflow Cloudfront Distribution"
  default_root_object = "index.html"

  # Updated logging configuration to use a dedicated logging bucket
  logging_config {
    include_cookies = false
    bucket         = "${var.logging_bucket}.s3.amazonaws.com"
    prefix         = "cloudfront-logs/"
  }

  # Remove local development alias
  # aliases = ["your-domain.com"]  # Add your production domain here

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = var.origin_id

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }

    # Enforce HTTPS
    viewer_protocol_policy = "redirect-to-https"
    
    # Optimize caching strategy
    min_ttl     = 0
    default_ttl = 60  # 1 minute
    max_ttl     = 300 # 5 Minutes

    # Add security headers
    response_headers_policy_id = aws_cloudfront_response_headers_policy.security_headers.id
  }

  ordered_cache_behavior {
    path_pattern     = "/images/*"
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD", "OPTIONS"]
    target_origin_id = var.origin_id

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }

    min_ttl     = 0
    default_ttl = 60    # 1 minute
    max_ttl     = 300 # 5 Minutes
    viewer_protocol_policy = "redirect-to-https"

    # Add security headers
    response_headers_policy_id = aws_cloudfront_response_headers_policy.security_headers.id
  }

  # Web Application Firewall integration
  web_acl_id = aws_wafv2_web_acl.cloudfront_waf.arn

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  # Use custom SSL certificate
  viewer_certificate {
    # acm_certificate_arn      = aws_acm_certificate.cdn_cert.arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
    cloudfront_default_certificate = true
  }

  tags = {
    Environment = "production"
    Service     = "snapflow"
    Managed_by  = "terraform"
  }
}

# Add security headers
resource "aws_cloudfront_response_headers_policy" "security_headers" {
  name = "security-headers-policy"

  security_headers_config {
    strict_transport_security {
      override = true
      access_control_max_age_sec = 31536000
      include_subdomains = true
      preload = true
    }

    content_security_policy {
    override = true
    content_security_policy = "default-src 'self'; img-src 'self' data:; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'"
  }

    content_type_options {
      override = true
    }

    xss_protection {
      override = true
      protection = true
    }

    referrer_policy {
      override = true
      referrer_policy = "same-origin"
    }
  }
}

