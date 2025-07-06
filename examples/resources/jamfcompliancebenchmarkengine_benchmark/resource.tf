# Example: Create a Jamf Compliance Benchmark

data "jamfcompliancebenchmarkengine_rules" "cis_lvl1" {
  baseline_id = "cis_lvl1"
}

resource "jamfcompliancebenchmarkengine_benchmark" "example" {
  title              = "Example Benchmark"
  description        = "Created by Terraform"
  source_baseline_id = "cis_lvl1"

  sources = [
    for s in data.jamfcompliancebenchmarkengine_rules.cis_lvl1.sources : {
      branch   = s.branch
      revision = s.revision
    }
  ]

  rules = [
    for r in data.jamfcompliancebenchmarkengine_rules.cis_lvl1.rules : {
      id      = r.id
      enabled = r.enabled
    }
  ]

  target = {
    device_groups = ["example-device-group-id"]
  }
  enforcement_mode = "MONITOR"
}

output "benchmark_id" {
  value = jamfcompliancebenchmarkengine_benchmark.example.id
}
