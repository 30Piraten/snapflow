// IAM role and policy for SQS
resource "aws_iam_role" "sqs_role" {
    name = "sqs_role"
    assume_role_policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                Effect: "Allow"
                Action: "sts:AssumeRole"
                Principal: {
                    Service: "lambda.amazonaws.com"
                }
                Resource = "${aws_sqs_queue.print_queue.arn}"
            }
        ]
    })
}

resource "aws_iam_role_policy" "sqs_role_policy" {
  name = "sqs_role_policy"
    role = aws_iam_role.sqs_role.id
    policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                Effect: "Allow"
                Action: [
                    "sqs:SendMessage",
                    "sqs:ReceiveMessage",
                    "sqs:DeleteMessage"
                ]
                Resource = "${aws_sqs_queue.print_queue.arn}"
            }
        ]
    })
}