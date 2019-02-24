# Modules

Each module should contain a README.md at the root and be entirely self contained. This means it should not rely on config in other directories or repos.

Each module must also be standalone. Meaning that it shouldn't rely on other modules being completed for it to work. For example, telemetry modules shouldn't rely on security or traffic routing modules to have been completed already. The only exception to this is that all modules will rely on the installation of Istio and the demo application.
