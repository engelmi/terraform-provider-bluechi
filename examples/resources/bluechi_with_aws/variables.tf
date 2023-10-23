/*********************/
/* Variables for AWS */
/*********************/

# tuple(public key path, private key path)
variable "ssh_key_pair" {
  type    = tuple([string, string])
  default = ["~/.ssh/bluechi_aws.pub", "~/.ssh/bluechi_aws"]
}

variable "ssh_user" {
  type    = string
  default = "ec2-user"
}

variable "autosd_ami" {
  type    = string
  default = "ami-0218a73af024c90fa"
}

variable "instance_type" {
  type    = string
  default = "t3a.micro"
}

/*************************/
/* Variables for BlueChi */
/*************************/

variable "use_mock" {
  type    = bool
  default = false
}

variable "bluechi_manager_port" {
  type    = number
  default = 3030
}

variable "bluechi_nodes" {
  type    = list(string)
  default = ["main", "worker1"]
}
