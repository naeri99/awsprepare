resource "aws_eks_node_group" "tfer--test-eks-nodegroup" {
  ami_type        = "AL2_x86_64"
  capacity_type   = "ON_DEMAND"
  cluster_name    = "${var.cluster-name}"
  disk_size       = "20"
  # instance_types  = ["t3a.medium"]
  instance_types  = ["c4.2xlarge"]
  node_group_name = "test-eks-nodegroup"
  node_role_arn   = aws_iam_role.test-iam-role-eks-nodegroup.arn
  release_version = "1.23.17-20230816"

  scaling_config {
    desired_size = "3"
    max_size     = "4"
    min_size     = "1"
  }

  subnet_ids = [aws_subnet.test-private-subnet1.id, aws_subnet.test-private-subnet3.id]

  update_config {
    max_unavailable = "1"
  }

  version = "1.23"

  # launch_template {
  #   name = aws_launch_template.template-one.name
  #   version = aws_launch_template.template-one.default_version 
  # }

  depends_on = [
    aws_iam_role_policy_attachment.test-iam-policy-eks-nodegroup,
    aws_iam_role_policy_attachment.test-iam-policy-eks-nodegroup-cni,
    aws_iam_role_policy_attachment.test-iam-policy-eks-nodegroup-ecr,
    aws_subnet.test-public-subnet1,
    aws_subnet.test-public-subnet3,
    aws_subnet.test-private-subnet1,
    aws_subnet.test-private-subnet3,
    aws_eks_cluster.test-eks-cluster

  ]


}
