# proxmox_virtual_environment_backups (Data Source)

Retrieves all Proxmox VE backup schedule jobs.




## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `ids` (List of String) Backup job IDs.
- `schedules` (List of String) Schedule for each job.
- `storages` (List of String) Target storage for each job.
