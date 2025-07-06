# Example: Fetch a Jamf Compliance Benchmark by ID

data "jamfcompliancebenchmarkengine_benchmark" "example" {
  id = "example-benchmark-id"
}

output "benchmark_title" {
  value = data.jamfcompliancebenchmarkengine_benchmark.example.title
}

output "benchmark_rules" {
  value = [for r in data.jamfcompliancebenchmarkengine_benchmark.example.rules : r.title]
}
