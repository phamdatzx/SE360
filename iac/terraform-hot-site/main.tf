provider "aws" {
  region = "ap-southeast-2"
}

# -----------------------
# VPC
# -----------------------
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true
  tags = {
    Name = "se-360-vpc"
  }
}

# -----------------------
# Internet Gateway
# -----------------------
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.main.id
  tags = {
    Name = "main-igw"
  }
}

# -----------------------
# EIP cho NAT Gateway
# -----------------------
resource "aws_eip" "nat" {
  tags = {
    Name = "nat-eip"
  }
}


# -----------------------
# Route table cho public subnet
# -----------------------
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }
  tags = {
    Name = "public-rt"
  }
}

# -----------------------
# Lấy 3 AZ ở Singapore
# -----------------------
data "aws_availability_zones" "available" {}

# -----------------------
# Public subnet
# -----------------------
resource "aws_subnet" "public" {
  count                   = 3
  vpc_id                  = aws_vpc.main.id
  cidr_block              = cidrsubnet(aws_vpc.main.cidr_block, 4, count.index) # 10.0.0.0/18, 10.0.16.0/18, 10.0.32.0/18
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true
  tags = {
    Name = "public-subnet-${count.index + 1}"
  }
}

# Gắn public subnet vào route table
resource "aws_route_table_association" "public_assoc" {
  count          = 3
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

# -----------------------
# NAT Gateway (1 NAT cho 3 private subnet)
# -----------------------
resource "aws_nat_gateway" "nat" {
  allocation_id = aws_eip.nat.id
  subnet_id     = aws_subnet.public[0].id
  tags = {
    Name = "main-nat"
  }
  depends_on = [aws_internet_gateway.igw]
}

# -----------------------
# Private subnet
# -----------------------
resource "aws_subnet" "private" {
  count             = 3
  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(aws_vpc.main.cidr_block, 4, count.index + 3) # 10.0.48.0/18, 10.0.64.0/18, 10.0.80.0/18
  availability_zone = data.aws_availability_zones.available.names[count.index]
  tags = {
    Name = "private-subnet-${count.index + 1}"
  }
}

# -----------------------
# Route table cho private subnet
# -----------------------
resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id
  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat.id
  }
  tags = {
    Name = "private-rt"
  }
}

# Gắn private subnet vào route table
resource "aws_route_table_association" "private_assoc" {
  count          = 3
  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private.id
}

# -----------------------
# Database subnet (isolated, no internet access)
# -----------------------
resource "aws_subnet" "database" {
  count             = 2
  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(aws_vpc.main.cidr_block, 4, count.index + 6) # 10.0.96.0/20, 10.0.112.0/20
  availability_zone = data.aws_availability_zones.available.names[count.index]
  tags = {
    Name = "database-subnet-${count.index + 1}"
  }
}

# -----------------------
# Route table cho database subnet (no internet access)
# -----------------------
resource "aws_route_table" "database" {
  vpc_id = aws_vpc.main.id
  tags = {
    Name = "database-rt"
  }
}

# Gắn database subnet vào route table
resource "aws_route_table_association" "database_assoc" {
  count          = 2
  subnet_id      = aws_subnet.database[count.index].id
  route_table_id = aws_route_table.database.id
}
