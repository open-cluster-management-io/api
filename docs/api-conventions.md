\[comment\]: # ( Copyright Contributors to the Open Cluster Management project )
# OCM API Conventions

OCM APIs follow the [Kubernetes API conventions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md)
with some additional guidlines inspired from [OpenShift API Convention](https://github.com/openshift/enhancements/blob/master/dev-guide/api-conventions.md).

## API Guidance

### Configuration vs Workload APIs

A configuration API is one that is typically cluster-scoped, a singleton within the cluster, and managed by a cluster
administrator only. Operator APIs such as `ClusterManager`, `Klusterlet` are configuration APIs.

A workload API typically is namespaced and is used by end users of the cluster. Most of the APIs, such as
`ManagedCluster`, `ManifestWork`, and `Placement` are workload APIs.

#### Defaulting

In workload APIs, we typically default fields on create (or update) when the field isn't set.
This has the effect that changing the default value in the API does not change the value for objects that have
previously been created.

This has the implication that you cannot change the behaviour of a default value once the API is defined as that would
cause the same object to result in different behavior for different versions of the API, which would surprise users and
compromise portability.

To change the default behaviour could constitute a breaking change and disrupt the end user's workload;
the behaviour must remain consistent through the lifetime of the resource.
This also means that defaults cannot be changed without a breaking change to the API.
If a user were to delete their workload API resource and recreate it, the behaviour should remain the same.

With configuration APIs, we typically default fields within the controller and not within the API.

Typically, optional fields on configuration APIs contain a statement within their Godoc to describe
the default behaviour when they are omitted, along with a note that this is subject to change over time.

#### Pointers

In configuration APIs specifically, we advise to avoid making fields pointers unless there is an absolute need to do so.
An absolute need being the need to distinguish between the zero value and a nil value.

Using pointers makes writing code to interact with the API harder and more error prone, and it also harms the
discoverability of the API.

##### Pointers to structs

An exception to this rule is when using a pointer to a struct in combination with an API validation that requires the
field to be unset.

The JSON tag `omitempty` is not compatible with struct references. Meaning any struct will always, when empty, be
marshalled to `{}`. If the struct **must** genuinely be omitted, it must be a pointer.

### Discriminated Unions

A discriminated union is a structure within the API that closely resembles a union type.
A number of fields exist within the structure and we are expecting the user to configure precisely one of
the fields.

In particular, for a discriminated union, an extra field exists which allows the user to declaratively
state which of particular fields they are configuring.

We use discriminated unions in Kubernetes APIs so that we force the user to make a choice.
We do not want our code to guess what the user meant, they should tell us which of the choices they made
using the discriminant.

Important to note:
* All structs within the union **MUST** be pointers
* All structs within the union **MUST** be optional
* The discriminant should be required
* The discriminant **MUST** be a string (or string alias) type
* Discriminant values should be PascalCase and should be equivalent to the camelCase field name (json tag) of one member of the union
* Empty union members (discriminant values without a paired union member) are also permitted

### No annotation-based APIs

Do not use annotations for extending an API. Annotations may seem as a good candidate for introducing experimental/new
API. Nevertheless, migration from annotations to formal schema usually never happens as it requires breaking changes
in customer deployments.

1. Validation does not always come with definition. User set values can be too broad and hard to limit later on.
2. Lack of discoverability. There's no pre-existing schema that can be published.
3. Validation is limited. Certain kinds of validators aren't allowed on annotations, so hooks are more frequently used instead.
4. Hard to extend. An annotation value (a string) can not be extended with additional fields under a parent.
5. Unclear versioning. Annotation keys can omit versioning. Or, there are multiple ways to specify a version.
6. Users can "squat" on annotations by adding an unvalidated annotation value for a key that is used in a future version.

### Use JSON Field Names in Godoc

Ensure that the godoc for a field name matches the JSON name, not the Go name,
in Go definitions for API objects.  In particular, this means that the godoc for
field names should use an initial lower-case letter.

### Use Resource Name Rather Than Kind in Object References

Use resource names rather than kinds for object references.

### Do not use Boolean fields

Many ideas start as a Boolean value, e.g. `FooEnabled: true|false`, but often evolve into needing 3, 4, or even more
states at some point during the API's lifetime.
As a Boolean value can only ever have 2 or in some cases 3 values (`true`, `false`, `omitted` when a pointer), we have
seen examples in which API authors have later added additional fields, paired with a Boolean field, that are only
meaningful when the original field has a certain state. This makes it confusing for an end user as they have to be
aware that the field they are trying to use only has an effect in certain circumstances.
