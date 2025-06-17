security {
  # 1 Specify which scanner to use
  scanner = "tfsec"

  # 2 Global scanner configuration
  config = {
    # Only show findings of at least this severity (LOW, MEDIUM, HIGH, CRITICAL)
    minimum-severity = "MEDIUM"

    # Exclude specific tfsec checks by their code
    exclude-checks = [
      # "AWS001",  # e.g. public S3 buckets (you may handle differently)
      # "GCP002",  # e.g. open firewall rule
    ]

    # Enable only a subset of checks, rather than excluding many
    # enabled-checks = ["AWS003", "AZU006", ...]

    # Set a timeout (in seconds) for the scan
    timeout = 120

    # Output format for reports (e.g. json, csv)
    output-type = "json"

    # Directory or file globs to scan; defaults to "./"
    path-patterns = ["**/*.tf", "**/*.hcl"]
  }

  # 3 Optional: per-directory overrides
  # If you have multiple sub‑modules with different needs
  overrides = [
    {
      # Limit checks in submodule directory
      path = "modules/networking"
      config = {
        minimum-severity = "HIGH"
      }
    },
    {
      path = "modules/experimental"
      config = {
        # Turn off tfsec for experimental code
        disable = true
      }
    }
  ]

  # 4 (Optional) Custom ignore file for tfsec
  # You can maintain a .tfsecignore file alongside your HCL:
  ignore-file = ".tfsecignore"
}
