provider "aws" {
    region = "us-east-1"
}
terraform {
    required_providers {
        aws = {
            version = "3.19.0"
        }
    }
}

resource "aws_s3_bucket" "bucket-with-env-tag" {
    bucket = "bucketwithenvtag"

    tags = {
        Name        = "bucket_with_env_tag"
        Environment = "Dev"
    }
}

resource "aws_s3_bucket" "bucket-without-tags" {
    bucket = "bucketwithouttags"
}

resource "aws_iam_user_policy" "lb_ec2_ro" {
    name = "lb_ec2_ro"

    user = aws_iam_user.loadbalancer.name

    # Terraform's "jsonencode" function converts a
    # Terraform expression result to valid JSON syntax.
    policy = jsonencode({
        Version   = "2012-10-17"
        Statement = [
            {
                Action = [
                    "ec2:Describe*",
                ]
                Effect   = "Allow"
                Resource = "*"
            },
        ]
    })
}

resource "aws_iam_user" "loadbalancer" {
    name = "loadbalancer"
    path = "/system/"
}

resource "aws_iam_access_key" "lb" {
    user = aws_iam_user.loadbalancer.name
}

resource "aws_iam_user_policy" "storageservice_s3_admin" {
    name = "storageservice_s3_admin"
    user = aws_iam_user.storageservice.name

    # Terraform's "jsonencode" function converts a
    # Terraform expression result to valid JSON syntax.
    policy = jsonencode({
        Version   = "2012-10-17"
        Statement = [
            {
                "Effect" : "Allow",
                "Action" : ["s3:ListBucket"],
                Resource = "*"
            },
            {
                "Effect" : "Allow",
                "Action" : ["s3:ListBucket"],
                Resource = "*"
            },
        ]
    })
}

resource "aws_iam_user" "storageservice" {
    name = "storageservice"
    path = "/system/"
}

resource "aws_iam_access_key" "storageservice" {
    user = aws_iam_user.storageservice.name
}

resource "aws_iam_policy" "sqshandlingpolicy" {
    name = "sqshandlingpolicy"

    policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "sqs:ReceiveMessage",
        "sqs:DeleteMessage",
        "sqs:GetQueueAttributes"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "s3lstpol" {
    name = "s3listingpolicy"

    policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
            "Effect": "Allow",
            "Action": ["s3:ListBucket"],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_policy_attachment" "sqspolicyattachment" {
    name       = "sqspolicyattachment"
    users      = [aws_iam_user.loadbalancer.name, aws_iam_user.storageservice.name]
    policy_arn = aws_iam_policy.sqshandlingpolicy.arn
}

resource "aws_iam_policy_attachment" "s3listingpolicyattachment" {
    name       = "s3listingpolicyattachment"
    users      = [aws_iam_user.loadbalancer.name]
    policy_arn = aws_iam_policy.s3lstpol.arn
}
