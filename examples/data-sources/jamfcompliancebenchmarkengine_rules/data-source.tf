# Example: Fetch all rules for a given baseline

data "jamfcompliancebenchmarkengine_rules" "cis_lvl1" {
  baseline_id = "cis_lvl1"
}

output "rule_titles" {
  value = [for r in data.jamfcompliancebenchmarkengine_rules.cis_lvl1.rules : r.title]
}
