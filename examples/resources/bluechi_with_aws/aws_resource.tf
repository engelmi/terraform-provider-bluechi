resource "aws_vpc" "autosd_demo_vpc" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "autosd-demo-vpc"
  }
}

resource "aws_subnet" "autosd_demo_subnet_public" {
  vpc_id            = aws_vpc.autosd_demo_vpc.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "us-east-2a"

  tags = {
    Name = "autosd-demo-subnet-public"
  }
}

resource "aws_subnet" "autosd_demo_subnet_private" {
  vpc_id            = aws_vpc.autosd_demo_vpc.id
  cidr_block        = "10.0.2.0/24"
  availability_zone = "us-east-2a"

  tags = {
    Name = "autosd-demo-subnet-private"
  }
}

resource "aws_internet_gateway" "autosd_demo_ig" {
  vpc_id = aws_vpc.autosd_demo_vpc.id

  tags = {
    Name = "autosd-demo-ig"
  }
}

resource "aws_route_table" "autosd_demo_rt" {
  vpc_id = aws_vpc.autosd_demo_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.autosd_demo_ig.id
  }

  route {
    ipv6_cidr_block = "::/0"
    gateway_id      = aws_internet_gateway.autosd_demo_ig.id
  }

  tags = {
    Name = "autosd-demo-rt"
  }
}

resource "aws_route_table_association" "public_1_rt_a" {
  subnet_id      = aws_subnet.autosd_demo_subnet_public.id
  route_table_id = aws_route_table.autosd_demo_rt.id
}

resource "aws_security_group" "autosd_demo_sg" {
  name   = "autosd-demo-sg"
  vpc_id = aws_vpc.autosd_demo_vpc.id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = var.bluechi_manager_port
    to_port     = var.bluechi_manager_port
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_key_pair" "bluechi" {
  key_name   = "bluechi-key"
  public_key = file(var.ssh_key_pair[0])
}

resource "aws_instance" "ec2main" {
  ami           = var.autosd_ami
  instance_type = var.instance_type
  key_name      = aws_key_pair.bluechi.id

  subnet_id                   = aws_subnet.autosd_demo_subnet_public.id
  vpc_security_group_ids      = [aws_security_group.autosd_demo_sg.id]
  associate_public_ip_address = true

  tags = {
    Platform = "AutoSD"
  }
}

resource "aws_instance" "ec2worker1" {
  ami           = var.autosd_ami
  instance_type = var.instance_type
  key_name      = aws_key_pair.bluechi.id

  subnet_id                   = aws_subnet.autosd_demo_subnet_public.id
  vpc_security_group_ids      = [aws_security_group.autosd_demo_sg.id]
  associate_public_ip_address = true

  tags = {
    Platform = "AutoSD"
  }
}

output "ec2main_public_ip" {
  value = aws_instance.ec2main.public_ip
}

output "ec2worker1_public_ip" {
  value = aws_instance.ec2worker1.public_ip
}
