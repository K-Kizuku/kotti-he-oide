terraform {
  backend "s3" {
    bucket = "kotti-he-oide-terraform-state-1c9b08ff"
    key    = "terraform.tfstate"
    region = "ap-northeast-1"
  }
}

provider "aws" {
  region = var.region
}

provider "random" {}
