# Trainer notes

## Setup

For 12 participants, this is the recommended
setup:

```
bastion_count = 1
instance_count = 14
cluster_initial_worker_node_count = 4
cluster_machine_type = "n1-standard-4"
```

Remember afterwards to run the following scripts
in the following order:

- `./wait_bastion_ready.sh && ./create-users.sh`

After this, all training instances should have
kubectl configured.
