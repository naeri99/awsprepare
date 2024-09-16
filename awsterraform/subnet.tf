resource "aws_subnet" "test-public-subnet1" {
  availability_zone = "ap-northeast-2a"

  assign_ipv6_address_on_creation                = "false"
  cidr_block                                     = "192.168.0.0/20"
  enable_dns64                                   = "false"
  # enable_lni_at_device_index                     = "0"
  enable_resource_name_dns_a_record_on_launch    = "false"
  enable_resource_name_dns_aaaa_record_on_launch = "false"
  ipv6_native                                    = "false"
  # map_customer_owned_ip_on_launch                = "false"
  map_public_ip_on_launch                        = "true"
  private_dns_hostname_type_on_launch            = "ip-name"

  tags = {
    Name                                     = "public-subnet1"
    "kubernetes.io/cluster/test-eks-cluster" = "shared"
    "kubernetes.io/role/elb" = 1 

  }

  tags_all = {
    Name                                     = "public-subnet3"
    "kubernetes.io/cluster/test-eks-cluster" = "shared"
    "kubernetes.io/role/elb" = 1 

  }


  vpc_id = "${aws_vpc.sung-vpc.id}"


}

resource "aws_subnet" "test-public-subnet3" {
  availability_zone = "ap-northeast-2c"

  assign_ipv6_address_on_creation                = "false"
  cidr_block                                     = "192.168.16.0/20"
  enable_dns64                                   = "false"
  # enable_lni_at_device_index                     = "0"
  enable_resource_name_dns_a_record_on_launch    = "false"
  enable_resource_name_dns_aaaa_record_on_launch = "false"
  ipv6_native                                    = "false"
  # map_customer_owned_ip_on_launch                = "false"
  map_public_ip_on_launch                        = "true"
  private_dns_hostname_type_on_launch            = "ip-name"

  tags = {
    Name                                     = "public-subnet3"
    "kubernetes.io/cluster/test-eks-cluster" = "shared"
    "kubernetes.io/role/elb" = 1 

  }

  tags_all = {
    Name                                     = "public-subnet3"
    "kubernetes.io/cluster/test-eks-cluster" = "shared"
    "kubernetes.io/role/elb" = 1 

  }


  vpc_id = "${aws_vpc.sung-vpc.id}"
}


resource "aws_subnet" "test-private-subnet1" {

  depends_on = [
    aws_vpc.sung-vpc  
  ]

  assign_ipv6_address_on_creation                = "false"
  cidr_block                                     = "192.168.32.0/20"
  enable_dns64                                   = "false"
  enable_resource_name_dns_a_record_on_launch    = "false"
  enable_resource_name_dns_aaaa_record_on_launch = "false"
  ipv6_native                                    = "false"
  map_public_ip_on_launch                        = "false"
  private_dns_hostname_type_on_launch            = "ip-name"

  tags = {
    Name                                     = "test-private-subnet1"
    "kubernetes.io/cluster/test-eks-cluster" = "shared"

  }

  tags_all = {
    Name                                     = "private-subnet1"
    "kubernetes.io/cluster/test-eks-cluster" = "shared"
  }


  vpc_id = "${aws_vpc.sung-vpc.id}"
  availability_zone = "ap-northeast-2a"
}

resource "aws_subnet" "test-private-subnet3" {

  depends_on = [
    aws_vpc.sung-vpc  
  ]

  assign_ipv6_address_on_creation                = "false"
  cidr_block                                     = "192.168.48.0/20"
  enable_dns64                                   = "false"
  enable_resource_name_dns_a_record_on_launch    = "false"
  enable_resource_name_dns_aaaa_record_on_launch = "false"
  ipv6_native                                    = "false"
  map_public_ip_on_launch                        = "false"
  private_dns_hostname_type_on_launch            = "ip-name"

  tags = {
    Name                                     = "test-private-subnet3"
    "kubernetes.io/cluster/test-eks-cluster" = "shared"

  }

  tags_all = {
    Name                                     = "private-subnet3"
    "kubernetes.io/cluster/test-eks-cluster" = "shared"
  }


  vpc_id = "${aws_vpc.sung-vpc.id}"
  availability_zone = "ap-northeast-2c"
}
