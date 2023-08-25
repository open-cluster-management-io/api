# api

The canonical location of the Open Cluster Management API definition.

## Generating CRD for the first time

To generate a CRD for the first time, please follow the following steps:
1. Make sure you have `controller-gen` locally installed in `_output/tools/bin` directory. If not, you can use `make ensure-controller-gen` to build one.
2. Run `'_output/tools/bin/controller-gen' crd:preserveUnknownFields=false paths="<API_DIRECTORY>" output:dir="<API_DIRECTORY>"` to generate a CRD.
   - `API_DIRECTORY` is the location of your APIs. For example, addon APIs are in `./addon/v1alpha1` directory. You will need `doc.go` and `types.go` created in your `API_DIRECTORY`.
3. Rename the CRD to follow the naming convention `0000_0x_<GROUP>_<RESOURCE>.crd.yaml`. 
   - `0x` should be a monotonically increasing number for each CRD in the `API_DIRECTORY`, and should start from `00`. 
   - `GROUP` is the group of your API, and `RESOURCE` is the resource name. For example, ManagedCluster is using `0000_00_clusters.open-cluster-management.io_managedclusters.crd.yaml`.
4. Remove `metadata.annotations` in the generated CRD. The generated CRD will have an annotation of the generator's version equals to `(devel)` as showed below, and we can simply remove these two lines:
   ```
    annotations:
      controller-gen.kubebuilder.io/version: (devel)
   ```

## Updating CRD schemas

If you make a change to a CRD type in this repo, calling `make update-codegen-crds` should regenerate all CRDs and update the manifests. If yours is not updated, ensure that the path to its API is included in our [calls to the Makefile targets](https://github.com/openshift/api/blob/release-4.5/Makefile#L17-L29).

To add this generator to another repo:
1. Vendor `github.com/openshift/library-go` (and ensure that the `alpha-build-machinery` subdirectory is also included in your `vendor`)

2. Update your `Makefile` to include the following:
```
include $(addprefix ./vendor/github.com/openshift/library-go/alpha-build-machinery/make/, \
  targets/openshift/crd-schema-gen.mk \
)

$(call add-crd-gen,<TARGET_NAME>,<API_DIRECTORY>,<CRD_MANIFESTS>,<MANIFEST_OUTPUT>)
```
The parameters for the call are:

1. `TARGET_NAME`: The name of your generated Make target. This can be anything, as long as it does not conflict with another make target. Recommended to be your api name.
2. `API_DIRECTORY`: The location of your API. For example if your Go types are located under `pkg/apis/myoperator/v1/types.go`, this should be `./pkg/apis/myoperator/v1`.
3. `CRD_MANIFESTS`: The directory your CRDs are located in. For example, if that is `manifests/my_operator.crd.yaml` then it should be `./manifests`
4. `MANIFEST_OUTPUT`: This should most likely be the same as `CRD_MANIFESTS`, and is only provided for flexibility to output generated code to a different directory.

You can include as many calls to different APIs as necessary, or if you have multiple APIs under the same directory (eg, `v1` and `v2beta1`) you can use 1 call to the parent directory pointing to your API.

After this, calling `make update-codegen-crds` should generate a new structural OpenAPIV3 schema for your CRDs.

**Notes**

- This will not generate entire CRDs, only their OpenAPIV3 schemas. If you do not already have a CRD, you will get no output from the generator.
- Ensure that your API is correctly declared for the generator to pick it up. That means, in your `doc.go`, include the following:
  1. `// +groupName=<API_GROUP_NAME>`, this should match the `group` in your CRD `spec`
  2. `// +kubebuilder:validation:Optional`, this tells the operator that fields should be optional unless explicitly marked with `// +kubebuilder:validation:Required`
- This will not touch any non-schema fields in the CRD. For example, `additionalPrinterColumns` & `preserveUnknownFields`. To regenerate those fields, please follow the instructions in the [Generating CRD for the first time](#Generating-CRD-for-the-first-time) section.

For more information on the API markers to add to your Go types, see the [Kubebuilder book](https://book.kubebuilder.io/reference/markers.html)

## Generating code
To generate `zz_generated.deepcopy.go` & `zz_generated.swagger_doc_generated.go`:
1. You will need to create a `register.go` file in your `API_DIRECTORY`. You can get an example from [/cluster/v1/register.go](/cluster/v1/register.go). Make sure the `package`, `GroupName`, and `GroupVersion` are correct. Make sure the `addKnownTypes()` function adds types you are creating in `types.go`.
2. Add `add-crd-gen` in `Makefile` (step 2 of [Updating CRD schemas](#Updating-CRD-schemas)).
3. Run `make update-scripts` to see `zz_generated.deepcopy.go` & `zz_generated.swagger_doc_generated.go` generated in your API directory, and along with some related scripts generated in `client` directory. If `zz_generated.deepcopy.go` is not generated properly, please run the command in a container with `RUNTIME=docker make update-with-container` (if you are using `podman`, set `RUNTIME=podman`).

## Verify
Before you commit you changes, please run `make verify` locally to make sure you have generated all required files.

## API upgrade flow
OCM community doc [Changing the API](https://github.com/open-cluster-management-io/community/blob/main/sig-architecture/api_changes.md) describes the things to consider when updating the API, suggest you read it before upgrading the API.

Once you decided to introduce a new version API, you will also need to consider deprecate the old version API and then remove it in the future. 
Below is the suggested API migration flow, the example code is to migrate `ManagedClusterSet` API from v1beta1 to v1beta2.

### Relase N (add new version)
- Add new version v1beta2 API to api repo. For example:
  - https://github.com/open-cluster-management-io/api/pull/174

- Convert between old and new versions.

  At this stage, both v1beta1 and v1beta2 are in served versions and storage version is v1beta1. If v1beta2 schema not change or can be converted to v1beta1, the CRD storage version can remain unchanged. If the v1beta2 schema changes and requires custom logic, you need to add conversion webhook to transform CRs between v1beta1 and v1beta2. For example: 
  - https://github.com/open-cluster-management-io/registration/pull/272
  - https://github.com/open-cluster-management-io/registration-operator/pull/279

- The consumers (eg, ui/foundation/submarinar-addon/placement) CAN start to migrate to clusterset api v1beta2.

- See also [How to add a new version](https://github.com/open-cluster-management-io/community/blob/main/sig-architecture/api_changes.md#how-to-add-a-new-version)

### Relase N+1 (deprecated old version and start migration)
- Mark v1beta1 as deprecated, served versions are v1beta1 and v1beta2, set storage: true on v1beta2.

- Migrate stored version of clusterset api to v1beta2 and remove v1beta1 from CRD `status.storedVersions`. For example: 
  - https://github.com/open-cluster-management-io/api/pull/202
  - https://github.com/open-cluster-management-io/registration-operator/pull/315
  - https://github.com/open-cluster-management-io/registration/pull/297

- The consumers (eg, ui/foundation/submarinar-addon/placement) SHOULD migrate to clusterset api v1beta2 before it's removed.

- See also [Migrate stored objects to the new version](https://github.com/open-cluster-management-io/community/blob/main/sig-architecture/api_changes.md#migrate-stored-objects-to-the-new-version)

### Relase N+x (remove old version)
when the new version API becomes mature, and ready to remove old version (for example, the consumers have adopted the new version API), remove the old version related code:
- Delete v1beta1 in clusterset crd. For example:
  - https://github.com/open-cluster-management-io/api/pull/266
- Delete conversion webhook about clusterset v1beta1
- Remove migration files in registration-operator
  - https://github.com/open-cluster-management-io/ocm/pull/257