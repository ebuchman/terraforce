# terraforce 

SSH and SCP for [terraform](https://www.terraform.io/) infrastructure.

## Notes

- expects `TF_VAR_key_path` and `TF_VAR_key_name` to be set for finding ssh keys
- `terraforce` commands must be executed from an terraform directory.
- `terraforce apply` must output `eip` and `dns` (its ok for an instance to have dns but not eip)

## TODO

- handle known hosts bit
- everything to make this actually generally useable :)
