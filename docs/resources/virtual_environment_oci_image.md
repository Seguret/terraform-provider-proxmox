---
page_title: "proxmox_virtual_environment_oci_image Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Downloads an OCI container image to Proxmox VE storage.
---

# proxmox_virtual_environment_oci_image (Resource)

Downloads an OCI container image (e.g. from Docker Hub) to Proxmox VE storage for use as a container template.

## Example Usage

```terraform
resource "proxmox_virtual_environment_oci_image" "ubuntu" {
  node_name    = "pve"
  datastore_id = "local"
  url          = "docker.io/library/ubuntu:22.04"
}

resource "proxmox_virtual_environment_oci_image" "custom" {
  node_name    = "pve"
  datastore_id = "local"
  url          = "ghcr.io/myorg/myimage:latest"
  file_name    = "myimage-latest.tar"
  pull_method  = "http"
}
```

## Schema

### Required

- `datastore_id` (String) The storage ID on which to store the image. Changing this forces a new resource.
- `node_name` (String) The name of the Proxmox VE node. Changing this forces a new resource.
- `url` (String) The OCI image URL (e.g. `docker.io/library/ubuntu:22.04`). Changing this forces a new resource.

### Optional

- `file_name` (String) Override the target filename on the storage.
- `pull_method` (String) The pull method: `http` or `oci`. Defaults to `oci`.

### Read-Only

- `id` (String) The volume ID of the downloaded image.
- `size` (String) The image size as reported by Proxmox.
