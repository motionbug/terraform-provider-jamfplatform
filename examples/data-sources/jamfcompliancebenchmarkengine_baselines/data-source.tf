# Example: Fetch all available baselines

data "jamfcompliancebenchmarkengine_baselines" "all" {}

output "baseline_titles" {
  value = [for b in data.jamfcompliancebenchmarkengine_baselines.all.baselines : b.title]
}
