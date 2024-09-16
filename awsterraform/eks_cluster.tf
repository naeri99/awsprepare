resource "aws_eks_cluster" "test-eks-cluster" {

  depends_on = [
    aws_iam_role_policy_attachment.test-iam-policy-eks-cluster,
    aws_iam_role_policy_attachment.test-iam-policy-eks-cluster-vpc,
    aws_iam_role_policy_attachment.test-iam-policy-eks-nodegroup,
    aws_iam_role_policy_attachment.test-iam-policy-eks-nodegroup-cni,
    aws_iam_role_policy_attachment.test-iam-policy-eks-nodegroup-ecr,
    aws_subnet.test-public-subnet1,
    aws_subnet.test-public-subnet3

  ]

  name     = "${var.cluster-name}"
  role_arn = aws_iam_role.test-iam-role-eks-cluster.arn
  version  = "1.23"

  enabled_cluster_log_types = ["api", "audit", "authenticator", "controllerManager", "scheduler"]


  vpc_config {
    endpoint_public_access = true
    security_group_ids      = [aws_security_group.test-sg-eks-cluster.id]
    subnet_ids              = [aws_subnet.test-public-subnet1.id, aws_subnet.test-public-subnet3.id]
  }
}
