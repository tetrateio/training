# Edit me!

# Download download the terraform GCP service account key (JSON) for the training-infra-owner project.
credz_file = "/Users/liam/Downloads/training-infra-owner-70f84aa3c556.json"

# Organization ID to place the project folder and projects under
organization_id = ""

# GCP Billing account to bill for the infra
billing_account = ""

# The name of the workshop or conference you'll be delivering at.
# WARNING: This is used in IDs with a count suffix so must be unique across all GCP.
workshop_name = "nist-2020"

# Number of participants in the workshop. This creates projects, service accounts, etc. 
# WARNING: Once set only increase the number unless you definitely don't need the projects as they are much harder to recover.
participant_count = 65

# Number of kube clusters to spin up. Comment out to spin down to 0, otherwise set to same value as participant_count.
# Use this to keep other infra but spin down Kube clusters to save on $$$ when they aren't needed.
cluster_count = 65
