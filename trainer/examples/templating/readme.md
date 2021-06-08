This directory contains a barebones helm chart, with nonsense .yaml files that exemplify different templating features.

Each `templates/*.yaml` represents a specific slide.

Use `$ helm template template-example` to render all templates

or `$ helm template template-example --show-only templates/functions.yaml`

To only render a specific template file.
