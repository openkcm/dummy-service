security {
  scanner = "trivy"

  config = {
    scan-type    = "iac"
    format       = "json"
    severity     = ["MEDIUM", "HIGH", "CRITICAL"]
    exit-code    = 1
    download-db-only = false
  }

  # If you still want tfsec as a fallback:
  overrides = [
    {
      path    = "."
      scanner = "tfsec"
      config = {
        minimum-severity = "MEDIUM"
      }
    }
  ]
}
