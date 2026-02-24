---
page_title: "proxmox_virtual_environment_file Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a file uploaded to Proxmox VE storage.
---

# proxmox_virtual_environment_file (Resource)

Manages a file uploaded to Proxmox VE storage (ISOs, container templates, snippets, etc.).

## Example Usage

```terraform
resource "proxmox_virtual_environment_file" "ubuntu_iso" {
  node_name    = "pve"
  datastore_id = "local"
  source_file  = "https://releases.ubuntu.com/22.04/ubuntu-22.04-live-server-amd64.iso"
  content_type = "iso"
}

resource "proxmox_virtual_environment_file" "snippet" {
  node_name    = "pve"
  datastore_id = "local"
  source_file  = "https://example.com/cloud-init.yaml"
  content_type = "snippets"
  file_name    = "cloud-init.yaml"
}
```

## Schema

### Required

- `datastore_id` (String) The storage ID on which to store the file (e.g. `local`). Changing this forces a new resource.
- `node_name` (String) The name of the Proxmox VE node. Changing this forces a new resource.
- `source_file` (String) The source URL to download the file from. Changing this forces a new resource.

### Optional

- `checksum` (String) Expected checksum of the file for verification. Changing this forces a new resource.
- `checksum_algorithm` (String) The checksum algorithm (`md5`, `sha1`, `sha256`, `sha512`). Changing this forces a new resource.
- `content_type` (String) The content type: `iso`, `vztmpl`, `snippets`, `import`. Inferred from URL if omitted. Changing this forces a new resource.
- `file_name` (String) Override the target filename on the storage. Changing this forces a new resource.
- `overwrite` (Boolean) Whether to overwrite an existing file with the same name. Defaults to `true`.
- `upload_timeout` (Number) The upload timeout in seconds. Defaults to `1800`.

### Read-Only

- `id` (String) The volume ID of the file (e.g. `local:iso/ubuntu.iso`).
- `size` (String) The file size as reported by Proxmox.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_file.example pve/local/local:iso/ubuntu.iso
```
