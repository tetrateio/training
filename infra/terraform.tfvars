# Edit me!

# Download download the terraform GCP service account key (JSON) for the training-infra-owner project.
credz_file = "/Users/liam/Downloads/training-infra-owner-70f84aa3c556.json"

# The name of the workshop or conference you'll be delivering at.
# WARNING: This is used in IDs with a count suffix so must be unique across all GCP.
workshop_name = "nist-2020"

# Number of participants in the workshop. This creates projects, service accounts, etc. 
# WARNING: Once set only increase the number unless you definitely don't need the projects as they are much harder to recover.
participant_count = 65

# Number of kube clusters to spin up. If unset will spin down to 0.
# Use this to keep other infra but spin down Kube clusters to save on $$$ when they aren't needed.
cluster_count = 1
