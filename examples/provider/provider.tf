terraform {
  required_providers {
    jamfcompliancebenchmarkengine = {
      source = "Jamf-Concepts/jamfcompliancebenchmarkengine"
    }
  }
}

provider "jamfcompliancebenchmarkengine" {
  region        = "us" # or "eu", "apac"
  client_id     = "example-client-id"
  client_secret = "example-client-secret"
}
