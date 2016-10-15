# terraforce 

SSH and SCP for [terraform](https://www.terraform.io/) infrastructure.

## Notes

- expects `TF_VAR_key_path` and `TF_VAR_key_name` to be set for finding ssh keys
- `terraforce` commands must be executed from an terraform directory.
- `terraforce apply` must output `eip` (ideally elastic ips to use for ssh) - TODO

## TODO

- everything to make this actually generally useable :)
