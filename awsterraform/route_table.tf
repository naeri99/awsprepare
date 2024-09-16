resource "aws_route_table" "test-route-table-pub-sub1" {

  depends_on = [
    aws_vpc.sung-vpc,
    aws_internet_gateway.sung-igw
  ]

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.sung-igw.id
  }

  tags = {
    Name = "test-route-table-pub-sub1"
  }


  vpc_id = "${aws_vpc.sung-vpc.id}"
}

resource "aws_route_table" "test-route-table-pub-sub3" {

  depends_on = [
    aws_vpc.sung-vpc,
    aws_internet_gateway.sung-igw
  ]

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.sung-igw.id
  }

  tags = {
    Name = "test-route-table-pub-sub3"
  }


  vpc_id = "${aws_vpc.sung-vpc.id}"
}

########## Private Subnet Route Tables ########## 

resource "aws_route_table" "test-route-table-pri-sub1" {

  depends_on = [
    aws_vpc.sung-vpc,
    aws_nat_gateway.test-nat-gateway
  ]

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_nat_gateway.test-nat-gateway.id
  }

  tags = {
    Name = "test-route-table-pri-sub1"
  }



  vpc_id = "${aws_vpc.sung-vpc.id}"
}

resource "aws_route_table" "test-route-table-pri-sub3" {

  depends_on = [
    aws_vpc.sung-vpc,
    aws_nat_gateway.test-nat-gateway
  ]

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_nat_gateway.test-nat-gateway.id
  }

  tags = {
    Name = "test-route-table-pri-sub3"
  }


  vpc_id = "${aws_vpc.sung-vpc.id}"
}
