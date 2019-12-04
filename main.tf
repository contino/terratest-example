provider "aws" {
  region = "us-east-1"
}
resource "aws_instance" "web_server" {
  ami                    = "ami-43a15f3e" # Ubuntu 16.04
  instance_type          = "t2.micro"
  vpc_security_group_ids = ["${aws_security_group.web_server.id}"]
  # Run a "Hello, World" web server on port 8080
  user_data = <<-EOF
              #!/bin/bash
              echo ${var.response_text} > index.html
              nohup busybox httpd -f -p 8080 &
              EOF  
}
# Allow the web app to receive requests on port 8080
resource "aws_security_group" "web_server" {
  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
