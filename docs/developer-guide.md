# Guida per Sviluppatori

Questa guida descrive l'architettura interna del provider Terraform per Proxmox VE e fornisce istruzioni dettagliate per estenderlo con nuove risorse, data source, modelli API e metodi client.

---

## Indice

1. [Panoramica dell'Architettura](#1-panoramica-dellarchitettura)
2. [Aggiungere un Nuovo Data Source](#2-aggiungere-un-nuovo-data-source)
3. [Aggiungere una Nuova Risorsa](#3-aggiungere-una-nuova-risorsa)
4. [Aggiungere Modelli API](#4-aggiungere-modelli-api)
5. [Pattern del Client HTTP](#5-pattern-del-client-http)
6. [Gestione degli Errori](#6-gestione-degli-errori)
7. [Operazioni Asincrone e Polling UPID](#7-operazioni-asincrone-e-polling-upid)
8. [Eseguire i Test](#8-eseguire-i-test)
9. [Convenzioni di Denominazione](#9-convenzioni-di-denominazione)
10. [Errori Comuni da Evitare](#10-errori-comuni-da-evitare)

---

## 1. Panoramica dell'Architettura

Il provider e' strutturato come un modulo Go standard. Ogni livello ha una singola responsabilita'.

### Struttura dei Package

```
terraform-provider-proxmox/
├── main.go                              # Entry point; chiama providerserver.Serve
├── GNUmakefile                          # Target: build, install, test, lint, fmt, clean
├── go.mod                               # Modulo: github.com/Seguret/terraform-provider-proxmox
│
└── internal/
    ├── provider/
    │   └── provider.go                  # ProxmoxProvider: Schema, Configure, Resources, DataSources
    │
    ├── client/
    │   ├── client.go                    # Client struct, metodi HTTP (Get/Post/Put/Delete),
    │   │                                # WaitForTask, WaitForTaskWithTimeout, authenticate
    │   ├── access.go                    # Metodi client: Utenti, Gruppi, Ruoli, ACL, Pool
    │   ├── vm.go                        # Metodi client: VM QEMU (Create/Get/Update/Delete/Clone/...)
    │   ├── container.go                 # Metodi client: Container LXC
    │   ├── storage.go                   # Metodi client: Definizioni storage
    │   ├── network.go                   # Metodi client: Interfacce di rete
    │   ├── firewall.go                  # Metodi client: Regole e opzioni firewall
    │   ├── errors.go                    # APIError, TaskError, IsNotFound()
    │   └── models/
    │       ├── common.go                # APIResponse[T], TaskStatus, Version
    │       ├── node.go                  # NodeListEntry, NodeStatus, CPUInfo, MemoryInfo, ...
    │       ├── vm.go                    # VMConfig, VMCreateRequest, VMCloneRequest, VMStatus, ...
    │       ├── container.go             # ContainerConfig, ContainerCreateRequest, ...
    │       ├── user.go                  # User, UserCreateRequest, UserUpdateRequest
    │       ├── group.go                 # Group, GroupCreateRequest, GroupUpdateRequest
    │       ├── role.go                  # Role, RoleCreateRequest, RoleUpdateRequest
    │       ├── acl.go                   # ACLEntry, ACLUpdateRequest
    │       ├── pool.go                  # Pool, PoolCreateRequest, PoolUpdateRequest
    │       ├── storage.go               # StorageListEntry (a livello di nodo)
    │       ├── storage_config.go        # StorageConfig, StorageCreateRequest, StorageUpdateRequest
    │       ├── network.go               # NetworkInterface, NetworkInterfaceCreateRequest
    │       └── firewall.go              # FirewallRule, FirewallOptions, ...
    │
    ├── resources/
    │   ├── vm/resource.go               # proxmox_virtual_environment_vm
    │   ├── container/resource.go        # proxmox_virtual_environment_container
    │   ├── user/resource.go             # proxmox_virtual_environment_user
    │   ├── group/resource.go            # proxmox_virtual_environment_group
    │   ├── role/resource.go             # proxmox_virtual_environment_role
    │   ├── acl/resource.go              # proxmox_virtual_environment_acl
    │   ├── pool/resource.go             # proxmox_virtual_environment_pool
    │   ├── storage/resource.go          # proxmox_virtual_environment_storage
    │   ├── network_interface/resource.go # proxmox_virtual_environment_network_interface
    │   ├── firewall_rule/resource.go    # proxmox_virtual_environment_firewall_rule
    │   └── firewall_options/resource.go # proxmox_virtual_environment_firewall_options
    │
    └── datasources/
        ├── version/datasource.go        # proxmox_virtual_environment_version
        ├── nodes/datasource.go          # proxmox_virtual_environment_nodes
        ├── node/datasource.go           # proxmox_virtual_environment_node
        ├── datastores/datasource.go     # proxmox_virtual_environment_datastores
        ├── vms/datasource.go            # proxmox_virtual_environment_vms
        └── containers/datasource.go     # proxmox_virtual_environment_containers
```

### Flusso dei Dati

```
Terraform Core
     |
     v
provider.Configure()
     |  crea *client.Client, lo memorizza in resp.ResourceData / resp.DataSourceData
     v
Resource.Configure() / DataSource.Configure()
     |  riceve *client.Client da req.ProviderData
     v
Resource.Create/Read/Update/Delete() / DataSource.Read()
     |  chiama i metodi client (es. client.GetVMs, client.CreateVM)
     v
client.Get() / client.Post() / client.Put() / client.Delete()
     |  avvolge DoRequest(), imposta gli header di autenticazione
     v
REST API di Proxmox VE  (https://<host>:8006/api2/json/...)
```

### Decisioni Architetturali Principali

**Singola istanza client per sessione del provider.** Il `*client.Client` viene creato una volta in `provider.Configure` e iniettato in ogni risorsa e data source tramite il meccanismo `ProviderData`.

**Wrapper generico per le risposte API.** Tutte le risposte delle API Proxmox seguono il formato `{"data": <payload>}`. Il tipo generico `models.APIResponse[T]` gestisce questo uniformemente in tutto il codice.

**Interfaccia sincrona per API asincrone.** Quando Proxmox restituisce un UPID (task identifier), il provider chiama sempre `client.WaitForTask` prima di ritornare, in modo che lo stato di Terraform venga scritto solo dopo che l'operazione e' completata sul cluster.

**Terraform Plugin Framework (non SDKv2).** Tutte le risorse e i data source implementano le interfacce del Plugin Framework (`resource.Resource`, `datasource.DataSource`). Non usare il legacy `terraform-plugin-sdk/v2`.

**Asserzioni di interfaccia a compile-time.** Ogni risorsa e data source include:

```go
var _ resource.Resource = &MiaRisorsa{}
var _ resource.ResourceWithConfigure = &MiaRisorsa{}
var _ resource.ResourceWithImportState = &MiaRisorsa{}
```

Queste righe producono un errore di compilazione se manca un metodo richiesto, rendendo l'errore visibile immediatamente invece di fallire a runtime.

---

## 2. Aggiungere un Nuovo Data Source

Un data source legge l'infrastruttura esistente ed espone i suoi attributi come valori calcolati. Non crea, aggiorna o elimina nulla.

### Passo 1: Creare il package e il file

```
internal/datasources/<nome>/datasource.go
```

Seguire la convenzione di denominazione: un data source per "snapshot" vive in `internal/datasources/snapshots/datasource.go`.

### Passo 2: Implementare il template completo

Ogni data source deve implementare `datasource.DataSource` e `datasource.DataSourceWithConfigure`. Il template completo e':

```go
package snapshots

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"

    "github.com/Seguret/terraform-provider-proxmox/internal/client"
)

// Asserzioni di interfaccia a compile-time. Producono un errore di build se
// manca un metodo richiesto, rendendo l'errore visibile immediatamente.
var _ datasource.DataSource = &SnapshotsDataSource{}
var _ datasource.DataSourceWithConfigure = &SnapshotsDataSource{}

// SnapshotsDataSource mantiene il client del provider.
type SnapshotsDataSource struct {
    client *client.Client
}

// SnapshotsDataSourceModel e' il modello di stato Terraform. Ogni campo
// corrisponde a un attributo dello schema tramite il tag struct `tfsdk`.
type SnapshotsDataSourceModel struct {
    ID       types.String   `tfsdk:"id"`
    NodeName types.String   `tfsdk:"node_name"`
    VMID     types.Int64    `tfsdk:"vmid"`
    Names    []types.String `tfsdk:"names"`
}

// NewDataSource e' il costruttore registrato in provider.go.
func NewDataSource() datasource.DataSource {
    return &SnapshotsDataSource{}
}

// Metadata imposta il nome del tipo usato nella configurazione Terraform.
// Convenzione: <nome_tipo_provider>_virtual_environment_<risorsa>
func (d *SnapshotsDataSource) Metadata(
    _ context.Context,
    req datasource.MetadataRequest,
    resp *datasource.MetadataResponse,
) {
    resp.TypeName = req.ProviderTypeName + "_virtual_environment_snapshots"
}

// Schema dichiara ogni attributo esposto dal data source.
func (d *SnapshotsDataSource) Schema(
    _ context.Context,
    _ datasource.SchemaRequest,
    resp *datasource.SchemaResponse,
) {
    resp.Schema = schema.Schema{
        Description: "Recupera la lista degli snapshot per una VM Proxmox VE.",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Description: "Identificatore placeholder.",
                Computed:    true,
            },
            "node_name": schema.StringAttribute{
                Description: "Il nome del nodo.",
                Required:    true,
            },
            "vmid": schema.Int64Attribute{
                Description: "L'ID della VM.",
                Required:    true,
            },
            "names": schema.ListAttribute{
                Description: "Nomi degli snapshot.",
                Computed:    true,
                ElementType: types.StringType,
            },
        },
    }
}

// Configure riceve il *client.Client dal provider.
// Verificare sempre che ProviderData non sia nil; e' nil durante la validazione.
func (d *SnapshotsDataSource) Configure(
    _ context.Context,
    req datasource.ConfigureRequest,
    resp *datasource.ConfigureResponse,
) {
    if req.ProviderData == nil {
        return
    }
    c, ok := req.ProviderData.(*client.Client)
    if !ok {
        resp.Diagnostics.AddError(
            "Tipo Configure del Data Source Inatteso",
            fmt.Sprintf("Atteso *client.Client, ottenuto: %T", req.ProviderData),
        )
        return
    }
    d.client = c
}

// Read recupera i dati dall'API e li scrive nello stato.
func (d *SnapshotsDataSource) Read(
    ctx context.Context,
    req datasource.ReadRequest,
    resp *datasource.ReadResponse,
) {
    // 1. Leggere la configurazione (attributi di input forniti dall'utente).
    var config SnapshotsDataSourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }

    nodeName := config.NodeName.ValueString()
    vmid := int(config.VMID.ValueInt64())

    tflog.Debug(ctx, "Lettura snapshot VM", map[string]any{
        "node": nodeName,
        "vmid": vmid,
    })

    // 2. Chiamare il client API.
    snapshots, err := d.client.GetVMSnapshots(ctx, nodeName, vmid)
    if err != nil {
        resp.Diagnostics.AddError(
            "Impossibile Leggere gli Snapshot della VM",
            fmt.Sprintf("Errore nella lettura degli snapshot per VM %d sul nodo %s: %s", vmid, nodeName, err),
        )
        return
    }

    // 3. Mappare la risposta API nel modello di stato Terraform.
    state := SnapshotsDataSourceModel{
        ID:       types.StringValue(fmt.Sprintf("snapshots/%s/%d", nodeName, vmid)),
        NodeName: types.StringValue(nodeName),
        VMID:     config.VMID,
        Names:    make([]types.String, len(snapshots)),
    }
    for i, s := range snapshots {
        state.Names[i] = types.StringValue(s.Name)
    }

    // 4. Scrivere lo stato finale.
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
```

### Passo 3: Registrare in provider.go

Aprire `internal/provider/provider.go` e aggiungere l'import e il costruttore allo slice `DataSources`:

```go
import (
    // import esistenti ...
    "github.com/Seguret/terraform-provider-proxmox/internal/datasources/snapshots"
)

func (p *ProxmoxProvider) DataSources(_ context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource{
        // voci esistenti ...
        snapshots.NewDataSource,
    }
}
```

---

## 3. Aggiungere una Nuova Risorsa

Una risorsa gestisce il ciclo di vita completo di un oggetto: Create, Read, Update e Delete. La maggior parte delle risorse implementa anche `ImportState`.

### Passo 1: Creare il package e il file

```
internal/resources/<nome>/resource.go
```

### Passo 2: Implementare il template completo

Ogni risorsa deve implementare `resource.Resource`, `resource.ResourceWithConfigure` e (di norma) `resource.ResourceWithImportState`.

```go
package snapshot

import (
    "context"
    "fmt"
    "strings"

    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"

    "github.com/Seguret/terraform-provider-proxmox/internal/client"
    "github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// Asserzioni di interfaccia a compile-time.
var _ resource.Resource = &SnapshotResource{}
var _ resource.ResourceWithConfigure = &SnapshotResource{}
var _ resource.ResourceWithImportState = &SnapshotResource{}

type SnapshotResource struct {
    client *client.Client
}

// SnapshotResourceModel e' la rappresentazione dello stato Terraform.
// I campi Computed+Optional mantengono valori riletti dall'API.
// I campi con RequiresReplace forzano la ricreazione quando vengono modificati.
type SnapshotResourceModel struct {
    ID          types.String `tfsdk:"id"`
    NodeName    types.String `tfsdk:"node_name"`
    VMID        types.Int64  `tfsdk:"vmid"`
    SnapName    types.String `tfsdk:"snap_name"`
    Description types.String `tfsdk:"description"`
}

func NewResource() resource.Resource {
    return &SnapshotResource{}
}

func (r *SnapshotResource) Metadata(
    _ context.Context,
    req resource.MetadataRequest,
    resp *resource.MetadataResponse,
) {
    resp.TypeName = req.ProviderTypeName + "_virtual_environment_snapshot"
}

func (r *SnapshotResource) Schema(
    _ context.Context,
    _ resource.SchemaRequest,
    resp *resource.SchemaResponse,
) {
    resp.Schema = schema.Schema{
        Description: "Gestisce uno snapshot di una VM Proxmox VE.",
        Attributes: map[string]schema.Attribute{
            // id e' sempre Computed; UseStateForUnknown previene diff inutili.
            "id": schema.StringAttribute{
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            // node_name e vmid insieme identificano la VM padre.
            // RequiresReplace significa che cambiarne uno forza un nuovo snapshot.
            "node_name": schema.StringAttribute{
                Description: "Il nome del nodo.",
                Required:    true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "vmid": schema.Int64Attribute{
                Description: "L'ID della VM.",
                Required:    true,
            },
            "snap_name": schema.StringAttribute{
                Description: "Il nome dello snapshot.",
                Required:    true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "description": schema.StringAttribute{
                Description: "Descrizione dello snapshot.",
                Optional:    true,
                Computed:    true,
            },
        },
    }
}

func (r *SnapshotResource) Configure(
    _ context.Context,
    req resource.ConfigureRequest,
    resp *resource.ConfigureResponse,
) {
    if req.ProviderData == nil {
        return
    }
    c, ok := req.ProviderData.(*client.Client)
    if !ok {
        resp.Diagnostics.AddError(
            "Tipo Configure della Risorsa Inatteso",
            fmt.Sprintf("Atteso *client.Client, ottenuto: %T", req.ProviderData),
        )
        return
    }
    r.client = c
}

// Create provisiona la risorsa e scrive lo stato.
func (r *SnapshotResource) Create(
    ctx context.Context,
    req resource.CreateRequest,
    resp *resource.CreateResponse,
) {
    // 1. Leggere i valori pianificati.
    var plan SnapshotResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    node := plan.NodeName.ValueString()
    vmid := int(plan.VMID.ValueInt64())

    tflog.Debug(ctx, "Creazione snapshot VM", map[string]any{
        "node":     node,
        "vmid":     vmid,
        "snapname": plan.SnapName.ValueString(),
    })

    // 2. Chiamare l'API. Molte operazioni di creazione Proxmox restituiscono un UPID.
    createReq := &models.VMSnapshotCreateRequest{
        SnapName:    plan.SnapName.ValueString(),
        Description: plan.Description.ValueString(),
    }
    upid, err := r.client.CreateVMSnapshot(ctx, node, vmid, createReq)
    if err != nil {
        resp.Diagnostics.AddError("Errore nella creazione dello snapshot VM", err.Error())
        return
    }

    // 3. Attendere il completamento del task asincrono prima di scrivere lo stato.
    if upid != "" {
        if err := r.client.WaitForTask(ctx, node, upid); err != nil {
            resp.Diagnostics.AddError("Errore nell'attesa della creazione snapshot", err.Error())
            return
        }
    }

    // 4. Impostare un ID deterministico, poi rileggere dall'API per popolare
    //    accuratamente tutti gli attributi Computed.
    plan.ID = types.StringValue(fmt.Sprintf("%s/%d/%s", node, vmid, plan.SnapName.ValueString()))
    r.readIntoModel(ctx, &plan, &resp.Diagnostics)
    if resp.Diagnostics.HasError() {
        return
    }

    // 5. Commit dello stato.
    resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read aggiorna lo stato dall'API. Deve chiamare resp.State.RemoveResource
// se la risorsa non esiste piu' (404).
func (r *SnapshotResource) Read(
    ctx context.Context,
    req resource.ReadRequest,
    resp *resource.ReadResponse,
) {
    var state SnapshotResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    r.readIntoModel(ctx, &state, &resp.Diagnostics)
    if resp.Diagnostics.HasError() {
        return
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update applica le modifiche agli attributi mutabili.
func (r *SnapshotResource) Update(
    ctx context.Context,
    req resource.UpdateRequest,
    resp *resource.UpdateResponse,
) {
    var plan SnapshotResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    node := plan.NodeName.ValueString()
    vmid := int(plan.VMID.ValueInt64())
    snapName := plan.SnapName.ValueString()

    updateReq := &models.VMSnapshotUpdateRequest{
        Description: plan.Description.ValueString(),
    }
    if err := r.client.UpdateVMSnapshot(ctx, node, vmid, snapName, updateReq); err != nil {
        resp.Diagnostics.AddError("Errore nell'aggiornamento dello snapshot VM", err.Error())
        return
    }

    plan.ID = types.StringValue(fmt.Sprintf("%s/%d/%s", node, vmid, snapName))
    r.readIntoModel(ctx, &plan, &resp.Diagnostics)
    resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete rimuove la risorsa dall'API.
func (r *SnapshotResource) Delete(
    ctx context.Context,
    req resource.DeleteRequest,
    resp *resource.DeleteResponse,
) {
    var state SnapshotResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    node := state.NodeName.ValueString()
    vmid := int(state.VMID.ValueInt64())
    snapName := state.SnapName.ValueString()

    upid, err := r.client.DeleteVMSnapshot(ctx, node, vmid, snapName)
    if err != nil {
        resp.Diagnostics.AddError("Errore nell'eliminazione dello snapshot VM", err.Error())
        return
    }
    if upid != "" {
        if err := r.client.WaitForTask(ctx, node, upid); err != nil {
            resp.Diagnostics.AddError("Errore nell'attesa dell'eliminazione snapshot", err.Error())
        }
    }
}

// ImportState ricostruisce lo stato da una stringa ID fornita dall'utente.
// Formato: <nodo>/<vmid>/<nome_snapshot>
func (r *SnapshotResource) ImportState(
    ctx context.Context,
    req resource.ImportStateRequest,
    resp *resource.ImportStateResponse,
) {
    parts := strings.SplitN(req.ID, "/", 3)
    if len(parts) != 3 {
        resp.Diagnostics.AddError(
            "ID di importazione non valido",
            "L'ID di importazione deve essere nel formato 'nome_nodo/vmid/nome_snapshot'",
        )
        return
    }

    var vmid int
    if _, err := fmt.Sscan(parts[1], &vmid); err != nil {
        resp.Diagnostics.AddError("VMID non valido", err.Error())
        return
    }

    state := SnapshotResourceModel{
        ID:       types.StringValue(req.ID),
        NodeName: types.StringValue(parts[0]),
        VMID:     types.Int64Value(int64(vmid)),
        SnapName: types.StringValue(parts[2]),
    }
    r.readIntoModel(ctx, &state, &resp.Diagnostics)
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// readIntoModel e' un helper che recupera lo stato corrente dall'API e lo scrive
// nel puntatore al modello fornito. Chiamarlo alla fine di Create, Read,
// Update e ImportState. Gestisce il caso 404 tramite il percorso diagnostics
// (i chiamanti verificano resp.Diagnostics.HasError()).
func (r *SnapshotResource) readIntoModel(
    ctx context.Context,
    model *SnapshotResourceModel,
    diagnostics *diag.Diagnostics,
) {
    node := model.NodeName.ValueString()
    vmid := int(model.VMID.ValueInt64())
    snapName := model.SnapName.ValueString()

    snap, err := r.client.GetVMSnapshot(ctx, node, vmid, snapName)
    if err != nil {
        // Il caso 404 viene gestito in Read dal chiamante che verifica i diagnostics.
        diagnostics.AddError("Errore nella lettura dello snapshot VM", err.Error())
        return
    }

    model.Description = types.StringValue(snap.Description)
}
```

### Passo 3: Registrare in provider.go

```go
import (
    // import esistenti ...
    "github.com/Seguret/terraform-provider-proxmox/internal/resources/snapshot"
)

func (p *ProxmoxProvider) Resources(_ context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        // voci esistenti ...
        snapshot.NewResource,
    }
}
```

---

## 4. Aggiungere Modelli API

Tutti gli struct di request e response vivono in `internal/client/models/`. Ogni file corrisponde a un dominio API di Proxmox (vm, user, storage, ecc.).

### Convenzioni di Denominazione degli Struct

| Scopo | Pattern di denominazione | Esempio |
|-------|------------------------|---------|
| Corpo risposta GET | `<Entita'>` | `VMSnapshot` |
| Corpo richiesta POST | `<Entita'>CreateRequest` | `VMSnapshotCreateRequest` |
| Corpo richiesta PUT | `<Entita'>UpdateRequest` | `VMSnapshotUpdateRequest` |
| Voce lista (es. `GET /qemu`) | `<Entita'>ListEntry` | `VMListEntry` |

### Tag JSON

- Usare sempre `json:"<nome_campo>"` corrispondente all'esatto nome del campo nell'API Proxmox.
- Usare `omitempty` sui campi opzionali di request create/update in modo che i valori zero non vengano inviati all'API.
- Usare un puntatore (`*int`, `*string`) per i campi opzionali sia nelle request che nelle response e dove il valore zero ha significato (es. `onboot=0` significa disabilitato, non assente).

**Esempio di modelli per snapshot VM:**

```go
// internal/client/models/vm.go (aggiunta)

// VMSnapshot rappresenta una voce snapshot da
// GET /nodes/{node}/qemu/{vmid}/snapshot.
type VMSnapshot struct {
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    SnapTime    int64  `json:"snaptime,omitempty"`
    VMState     *int   `json:"vmstate,omitempty"`
}

// VMSnapshotCreateRequest rappresenta il corpo per
// POST /nodes/{node}/qemu/{vmid}/snapshot.
type VMSnapshotCreateRequest struct {
    SnapName    string `json:"snapname"`
    Description string `json:"description,omitempty"`
    VMState     *int   `json:"vmstate,omitempty"`
}

// VMSnapshotUpdateRequest rappresenta il corpo per
// PUT /nodes/{node}/qemu/{vmid}/snapshot/{snapname}/config.
type VMSnapshotUpdateRequest struct {
    Description string `json:"description,omitempty"`
}
```

---

## 5. Pattern del Client HTTP

### Metodi Helper HTTP

Lo struct `Client` in `internal/client/client.go` fornisce cinque metodi helper HTTP:

| Metodo | Firma | Quando usarlo |
|--------|-------|--------------|
| `Get` | `Get(ctx, path, target) error` | Qualsiasi richiesta `GET`. Decodifica `{"data": ...}` in `target`. |
| `Post` | `Post(ctx, path, body, target) error` | `POST` quando e' necessario decodificare il corpo della risposta. Usare per operazioni che restituiscono un UPID o un nuovo ID risorsa. |
| `PostNoResponse` | `PostNoResponse(ctx, path, body) error` | `POST` quando l'API restituisce `{"data": null}` (operazioni di creazione senza valore di ritorno significativo). |
| `Put` | `Put(ctx, path, body) error` | `PUT` per le operazioni di aggiornamento. Proxmox restituisce `{"data": null}` per la maggior parte degli aggiornamenti. |
| `Delete` | `Delete(ctx, path) error` | `DELETE` per eliminazioni sincrone. Per eliminazioni asincrone, usare `DoRequest` direttamente per catturare l'UPID dalla risposta. |

### Aggiungere un Metodo Client

Aggiungere metodi al file che corrisponde al dominio API:

- `internal/client/access.go` - Utenti, Gruppi, Ruoli, ACL, Pool
- `internal/client/vm.go` - VM QEMU
- `internal/client/container.go` - Container LXC
- `internal/client/storage.go` - Definizioni storage
- `internal/client/network.go` - Interfacce di rete
- `internal/client/firewall.go` - Regole e opzioni firewall
- Creare un nuovo file per un nuovo dominio

**Template per un nuovo metodo client:**

```go
// GetVMSnapshots recupera la lista degli snapshot per una VM
// (GET /nodes/{node}/qemu/{vmid}/snapshot).
func (c *Client) GetVMSnapshots(ctx context.Context, node string, vmid int) ([]models.VMSnapshot, error) {
    path := fmt.Sprintf("/nodes/%s/qemu/%d/snapshot", url.PathEscape(node), vmid)
    var result models.APIResponse[[]models.VMSnapshot]
    if err := c.Get(ctx, path, &result); err != nil {
        return nil, err
    }
    return result.Data, nil
}

// GetVMSnapshot recupera un singolo snapshot per nome
// (GET /nodes/{node}/qemu/{vmid}/snapshot/{snapname}/config).
func (c *Client) GetVMSnapshot(ctx context.Context, node string, vmid int, snapName string) (*models.VMSnapshot, error) {
    path := fmt.Sprintf("/nodes/%s/qemu/%d/snapshot/%s/config",
        url.PathEscape(node), vmid, url.PathEscape(snapName))
    var result models.APIResponse[models.VMSnapshot]
    if err := c.Get(ctx, path, &result); err != nil {
        return nil, err
    }
    return &result.Data, nil
}

// CreateVMSnapshot crea uno snapshot e restituisce l'UPID del task asincrono
// (POST /nodes/{node}/qemu/{vmid}/snapshot).
func (c *Client) CreateVMSnapshot(
    ctx context.Context, node string, vmid int, req *models.VMSnapshotCreateRequest,
) (string, error) {
    path := fmt.Sprintf("/nodes/%s/qemu/%d/snapshot", url.PathEscape(node), vmid)
    body, err := json.Marshal(req)
    if err != nil {
        return "", err
    }
    var result models.APIResponse[string]
    if err := c.Post(ctx, path, bytes.NewReader(body), &result); err != nil {
        return "", err
    }
    return result.Data, nil
}

// UpdateVMSnapshot aggiorna la descrizione di uno snapshot
// (PUT /nodes/{node}/qemu/{vmid}/snapshot/{snapname}/config).
func (c *Client) UpdateVMSnapshot(
    ctx context.Context, node string, vmid int, snapName string,
    req *models.VMSnapshotUpdateRequest,
) error {
    path := fmt.Sprintf("/nodes/%s/qemu/%d/snapshot/%s/config",
        url.PathEscape(node), vmid, url.PathEscape(snapName))
    body, err := json.Marshal(req)
    if err != nil {
        return err
    }
    return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteVMSnapshot elimina uno snapshot e restituisce l'UPID del task asincrono
// (DELETE /nodes/{node}/qemu/{vmid}/snapshot/{snapname}).
func (c *Client) DeleteVMSnapshot(ctx context.Context, node string, vmid int, snapName string) (string, error) {
    path := fmt.Sprintf("/nodes/%s/qemu/%d/snapshot/%s",
        url.PathEscape(node), vmid, url.PathEscape(snapName))
    resp, err := c.DoRequest(ctx, "DELETE", path, nil)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        return "", c.parseError(resp)
    }
    var result models.APIResponse[string]
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", fmt.Errorf("impossibile decodificare la risposta di eliminazione: %w", err)
    }
    return result.Data, nil
}
```

---

## 6. Gestione degli Errori

### APIError

`client.parseError` viene chiamato ogni volta che lo stato HTTP della risposta non e' 200. Produce un `*APIError`:

```go
type APIError struct {
    StatusCode int
    Status     string
    Message    string
    Errors     map[string]string
}
```

Il metodo `Error()` formatta il messaggio come: `proxmox API error <codice>: <messaggio>`.

### IsNotFound

Usare `IsNotFound()` per rilevare se una risorsa e' stata eliminata al di fuori di Terraform:

```go
func (r *SnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state SnapshotResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    snap, err := r.client.GetVMSnapshot(
        ctx,
        state.NodeName.ValueString(),
        int(state.VMID.ValueInt64()),
        state.SnapName.ValueString(),
    )
    if err != nil {
        // Se l'API ha restituito 404, la risorsa e' stata eliminata fuori da Terraform.
        // Rimuoverla dallo stato in modo che il prossimo plan la mostri come da ricreare.
        if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
            resp.State.RemoveResource(ctx)
            return
        }
        resp.Diagnostics.AddError("Errore nella lettura dello snapshot VM", err.Error())
        return
    }

    state.Description = types.StringValue(snap.Description)
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
```

### TaskError

Quando `WaitForTask` determina che un task si e' fermato con uno stato di uscita non-OK, restituisce un `*TaskError`:

```go
type TaskError struct {
    UPID       string
    Node       string
    ExitStatus string
}
```

Lo stato di uscita e' una stringa leggibile dall'uomo proveniente dal log dei task Proxmox. Esporlo in un errore diagnostico fornisce informazioni utili all'utente:

```go
upid, err := r.client.CreateVMSnapshot(ctx, node, vmid, req)
if err != nil {
    resp.Diagnostics.AddError("Errore nella creazione snapshot", err.Error())
    return
}
if err := r.client.WaitForTask(ctx, node, upid); err != nil {
    resp.Diagnostics.AddError("Creazione snapshot fallita", err.Error())
    return
}
```

---

## 7. Operazioni Asincrone e Polling UPID

L'API Proxmox VE e' asincrona per la maggior parte delle operazioni di scrittura. Quando un'operazione avvia un job in background, l'API restituisce immediatamente un UPID (Universal Process Identifier) invece di attendere il completamento.

### Quando Aspettarsi un UPID

| Operazione | Restituisce UPID |
|-----------|----------------|
| `CreateVM` | Si |
| `DeleteVM` | Si |
| `StartVM` | Si |
| `StopVM` | Si |
| `ShutdownVM` | Si |
| `CloneVM` | Si |
| `CreateVMSnapshot` | Si |
| `DeleteVMSnapshot` | Si |
| `CreateContainer` | Si |
| `DeleteContainer` | Si |
| `StartContainer` | Si |
| `StopContainer` | Si |
| `CreateUser` | No (restituisce null) |
| `UpdateUser` | No (restituisce null) |
| `CreateStorage` | No (restituisce null) |
| `CreateRole` | No (restituisce null) |

### Attendere Sempre Prima di Leggere lo Stato

Dopo qualsiasi operazione che restituisce un UPID, chiamare `WaitForTask` prima di rileggere la risorsa dall'API. Saltare questo passaggio porta a leggere dati obsoleti: la configurazione o lo stato della VM potrebbe non riflettere ancora l'operazione completata sul cluster.

```go
upid, err := r.client.CreateVM(ctx, node, createReq)
if err != nil {
    resp.Diagnostics.AddError("Errore nella creazione della VM", err.Error())
    return
}
// L'UPID puo' essere una stringa vuota per alcune operazioni; verificarlo.
if upid != "" {
    if err := r.client.WaitForTask(ctx, node, upid); err != nil {
        resp.Diagnostics.AddError("Errore nell'attesa della creazione VM", err.Error())
        return
    }
}
// Solo ora e' sicuro leggere la configurazione VM.
r.readIntoModel(ctx, &plan, &resp.Diagnostics)
```

### Parametri di Polling

| Costante | Valore | Descrizione |
|----------|--------|-------------|
| `taskPollInterval` | 2 secondi | Intervallo tra i poll dello stato |
| `taskDefaultTimeout` | 30 minuti | Tempo massimo di attesa (`WaitForTask`) |

Usare `WaitForTaskWithTimeout` per operazioni con una durata nota limitata:

```go
if err := r.client.WaitForTaskWithTimeout(ctx, node, upid, 5*time.Minute); err != nil {
    resp.Diagnostics.AddError("Timeout snapshot", err.Error())
    return
}
```

---

## 8. Eseguire i Test

### Test Unitari

I test unitari non richiedono un cluster attivo. Posizionarli accanto al codice che testano, seguendo la convenzione Go:

```
internal/client/client_test.go
internal/resources/vm/resource_test.go
```

Eseguire con:

```shell
make test
# oppure direttamente:
go test ./... -v
```

### Test di Accettazione

I test di accettazione (file `*_test.go` contenenti `resource.Test(t, ...)`) richiedono un'istanza Proxmox VE attiva. Sono protetti dalla variabile d'ambiente `TF_ACC=1`.

**Variabili d'ambiente richieste:**

| Variabile | Descrizione |
|-----------|-------------|
| `TF_ACC` | Impostare a `1` per abilitare i test di accettazione |
| `PROXMOX_VE_ENDPOINT` | URL dell'API Proxmox VE |
| `PROXMOX_VE_API_TOKEN` | Token API con privilegi sufficienti |
| `PROXMOX_VE_INSECURE` | Impostare a `true` per usare un certificato auto-firmato |

Eseguire con:

```shell
export TF_ACC=1
export PROXMOX_VE_ENDPOINT="https://pve.lab:8006"
export PROXMOX_VE_API_TOKEN="root@pam!test=<uuid>"
export PROXMOX_VE_INSECURE="true"

make testacc
# oppure con un test specifico:
TF_ACC=1 go test ./internal/resources/vm/... -run TestAccVM -v -timeout 30m
```

### Struttura di un Test di Accettazione

```go
package vm_test

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-framework/providerserver"
    "github.com/hashicorp/terraform-plugin-go/tfprotov6"
    "github.com/hashicorp/terraform-plugin-testing/helper/resource"

    "github.com/Seguret/terraform-provider-proxmox/internal/provider"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
    "proxmox": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func TestAccVM_base(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: `
                    resource "proxmox_virtual_environment_vm" "test" {
                        node_name = "pve"
                        name      = "acc-test-vm"
                        memory    = 512
                        started   = false
                    }
                `,
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr(
                        "proxmox_virtual_environment_vm.test", "name", "acc-test-vm",
                    ),
                    resource.TestCheckResourceAttr(
                        "proxmox_virtual_environment_vm.test", "memory", "512",
                    ),
                ),
            },
        },
    })
}
```

---

## 9. Convenzioni di Denominazione

### Nomi dei Tipi di Risorsa

Tutte le risorse e i data source usano il prefisso `proxmox_virtual_environment_` seguito dal nome dell'entita' in snake_case.

| Nome tipo Go | Tipo risorsa Terraform |
|--------------|----------------------|
| `VMResource` | `proxmox_virtual_environment_vm` |
| `ContainerResource` | `proxmox_virtual_environment_container` |
| `UserResource` | `proxmox_virtual_environment_user` |
| `GroupResource` | `proxmox_virtual_environment_group` |
| `RoleResource` | `proxmox_virtual_environment_role` |
| `ACLResource` | `proxmox_virtual_environment_acl` |
| `PoolResource` | `proxmox_virtual_environment_pool` |
| `StorageResource` | `proxmox_virtual_environment_storage` |
| `NetworkInterfaceResource` | `proxmox_virtual_environment_network_interface` |
| `FirewallRuleResource` | `proxmox_virtual_environment_firewall_rule` |
| `FirewallOptionsResource` | `proxmox_virtual_environment_firewall_options` |

### Layout dei File Rispecchia l'Albero API

La directory per risorse e data source dovrebbe riflettere il segmento del percorso API Proxmox dove vive la risorsa:

| Percorso API | Posizione package |
|-------------|-----------------|
| `/access/users` | `internal/resources/user/` |
| `/access/groups` | `internal/resources/group/` |
| `/access/roles` | `internal/resources/role/` |
| `/access/acl` | `internal/resources/acl/` |
| `/pools` | `internal/resources/pool/` |
| `/storage` | `internal/resources/storage/` |
| `/nodes/{node}/qemu` | `internal/resources/vm/` |
| `/nodes/{node}/lxc` | `internal/resources/container/` |
| `/nodes/{node}/network` | `internal/resources/network_interface/` |
| `/cluster/firewall/rules` | `internal/resources/firewall_rule/` |
| `/nodes` | `internal/datasources/nodes/` |
| `/nodes/{node}/storage` | `internal/datasources/datastores/` |

### Nomi dei Metodi Client

I metodi client seguono il pattern `<Verbo><Entita'>` dove il verbo e' uno tra: `Get`, `Create`, `Update`, `Delete`, `Clone`, `Start`, `Stop`, `Shutdown`, `Resize`, `Apply`.

### Nomi degli Struct Modello

Vedere la tabella nella [Sezione 4](#4-aggiungere-modelli-api).

---

## 10. Errori Comuni da Evitare

### `int` vs `*int` per i Campi dell'API Proxmox

L'API Proxmox usa flag interi per campi simil-booleani (es. `onboot=1` significa abilitato). Nei modelli Go:

- Usare `int` (non puntatore) quando `0` e' un valore valido e il campo e' sempre presente nella risposta (es. `sockets int`).
- Usare `*int` quando il campo puo' essere assente dalla risposta e occorre distinguere "assente" da "impostato a zero" (es. `onboot *int`).

In `readIntoModel`, verificare sempre nil prima di dereferenziare un campo puntatore:

```go
if cfg.OnBoot != nil {
    model.OnBoot = types.BoolValue(*cfg.OnBoot == 1)
}
```

Dimenticare di gestire i puntatori nil provoca un panic a runtime, non un errore a compile-time.

### Chiavi Dinamiche per Dischi e Schede di Rete

La configurazione delle VM Proxmox usa chiavi stringa dinamiche (`scsi0`, `scsi1`, `net0`, `net1`, ecc.) piuttosto che oggetti annidati strutturati. Il provider rappresenta questi come attributi `types.String` individuali invece di una lista o set.

Quando si implementa `buildConfigMap` per gli aggiornamenti VM, includere una chiave nella mappa solo se il campo del modello corrispondente non e' vuoto. Se si include una stringa vuota, Proxmox potrebbe interpretarla come rimozione del dispositivo:

```go
// Corretto: inviare solo se il valore e' impostato
if val := plan.SCSI0.ValueString(); val != "" {
    m["scsi0"] = val
}
// Sbagliato: invia "scsi0": "" che potrebbe rimuovere il disco
m["scsi0"] = plan.SCSI0.ValueString()
```

### Il VMID Deve Essere Impostato Prima di readIntoModel

Se `vmid` non e' specificato nella configurazione Terraform, il provider chiama `GetNextVMID` per ottenerne uno. Dopo la creazione, impostare sempre `plan.VMID` con il valore assegnato prima di chiamare `readIntoModel`:

```go
vmid := int(plan.VMID.ValueInt64())
if plan.VMID.IsNull() || plan.VMID.IsUnknown() || vmid == 0 {
    nextID, err := r.client.GetNextVMID(ctx)
    if err != nil {
        resp.Diagnostics.AddError("Errore nel recupero del prossimo VMID", err.Error())
        return
    }
    vmid = nextID
}

// ... creare la VM ...

plan.VMID = types.Int64Value(int64(vmid))  // deve essere impostato prima di readIntoModel
plan.ID = types.StringValue(fmt.Sprintf("%s/%d", node, vmid))
r.readIntoModel(ctx, &plan, &resp.Diagnostics)
```

Omettere `plan.VMID = ...` fa si' che `readIntoModel` interroghi VMID 0, che non esiste, causando un errore 404.

### Non Mescolare `api_token` con `username/password`

La logica di autenticazione in `client.New` verifica prima `api_token`. Se entrambi sono forniti, `username` e `password` vengono silenziosamente ignorati. Il provider valida questo al momento della configurazione:

```go
if apiToken == "" && (username == "" || password == "") {
    resp.Diagnostics.AddError("Autenticazione Mancante", "...")
}
```

Questo significa che entrambi i metodi possono essere forniti tramite variabili d'ambiente senza errori, finche' `api_token` e' impostato.

### Usare `RequiresReplace` Correttamente

Marcare un attributo con `stringplanmodifier.RequiresReplace()` (o gli equivalenti integer/bool) quando l'API Proxmox non supporta aggiornamenti in-place per quel campo. Dimenticare questo fa si' che Terraform tenti un Update che verra' silenziosamente ignorato dall'API o restituira' un errore.

Esempi di campi che richiedono replace: `user_id`, `group_id`, `role_id`, `pool_id`, `storage` (ID storage), `type` (tipo storage), `node_name` (per VM e container), `snap_name`, `iface` (interfaccia di rete), `scope` (firewall).

### Verifica TLS in Sviluppo

Quando `insecure = true` (o `PROXMOX_VE_INSECURE=true`), il provider salta la validazione del certificato TLS. Questo e' accettabile in un laboratorio locale ma non deve mai essere usato in produzione. Il commento `//nolint:gosec` sull'assegnazione `InsecureSkipVerify` e' intenzionale e documenta questo compromesso.

### Terraform Plugin Framework vs. SDKv2

Questo provider usa `terraform-plugin-framework`, non `terraform-plugin-sdk/v2`. Il framework usa tipi di interfaccia diversi, pattern di accesso allo stato e diagnostics. Non mescolare i due. Differenze chiave:

| Aspetto | Plugin Framework | SDKv2 |
|---------|-----------------|-------|
| Accesso allo stato | `req.State.Get(ctx, &model)` | `d.Get("campo")` |
| Errori | `resp.Diagnostics.AddError(...)` | `return diag.Errorf(...)` |
| Schema | `schema.StringAttribute{...}` | `schema.Schema{Type: schema.TypeString, ...}` |
| Interfaccia | `resource.Resource` | `schema.Resource` |
| Tipo schema DS | `datasource.DataSource` | `schema.Resource` con ReadContext |

### La Cancellazione di `firewall_options` Disabilita, Non Elimina

La risorsa `proxmox_virtual_environment_firewall_options` ha un comportamento speciale nella cancellazione: invece di eliminare la risorsa (che non e' possibile via API), disabilita semplicemente il firewall impostando `enable=0`. Questo e' il comportamento corretto perche' le opzioni firewall sono un'entita' singleton per ogni scope.

```go
func (r *FirewallOptionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Disabilita il firewall invece di eliminare la risorsa.
    var state FirewallOptionsResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }
    disabled := 0
    pathPrefix := r.scopeToPath(state.Scope.ValueString())
    _ = r.client.UpdateFirewallOptions(ctx, pathPrefix, &models.FirewallOptions{Enable: &disabled})
}
```

### Il Pattern `scopeToPath` per il Firewall

Le risorse firewall usano un campo `scope` per determinare su quale oggetto si applica la regola. La funzione `scopeToPath` converte lo scope nel prefisso del percorso API corrispondente:

| Scope | Percorso API |
|-------|------------|
| `cluster` | `/cluster/firewall` |
| `node/pve1` | `/nodes/pve1/firewall` |
| `vm/pve1/100` | `/nodes/pve1/qemu/100/firewall` |
| `ct/pve1/300` | `/nodes/pve1/lxc/300/firewall` |

---

*Ultimo aggiornamento: 2026-02-18*
