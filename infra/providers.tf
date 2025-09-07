terraform {
  backend "s3" {
    # これらの値は後でterraform initで設定します
    bucket = "kotti-he-oide-terraform-state-1c9b08ff"
    key    = "terraform.tfstate"
    region = "ap-northeast-1"
  }
}

provider "aws" {
  region = var.region
}

provider "random" {}
