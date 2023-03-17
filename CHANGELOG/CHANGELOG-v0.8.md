# v0.8.0

## Changelog since v0.7.0

### New Features 
* Support Registration FeatureGates to pass feature gates ([#149](https://github.com/open-cluster-management-io/api/pull/149) [@ivan-cai](https://github.com/ivan-cai))
* Feat: Adding work subject executor feature related api ([#152](https://github.com/open-cluster-management-io/api/pull/152) [@yue9944882](https://github.com/yue9944882))
* Work Update strategy ([#154](https://github.com/open-cluster-management-io/api/pull/154) [@qiujian16](https://github.com/qiujian16))

### Added
* Add label selector to clusterset spec ([#150](https://github.com/open-cluster-management-io/api/pull/150) [@ldpliu](https://github.com/ldpliu))
* Add HostAlias in KlusterletDeployOption ([#166](https://github.com/open-cluster-management-io/api/pull/166) [@Promacanthus](https://github.com/Promacanthus))

### Changes
* Add help func for cluster and clusterset ([#153](https://github.com/open-cluster-management-io/api/pull/153) [@ldpliu](https://github.com/ldpliu))
* Add status for clustersetbinding ([#158](https://github.com/open-cluster-management-io/api/pull/158) [@qiujian16](https://github.com/qiujian16))
* Add a helper function to get all valid ManagedClusterSetBindings ([#151](https://github.com/open-cluster-management-io/api/pull/151) [@mikeshng](https://github.com/mikeshng))
* move registration feature gates to api ([#162](https://github.com/open-cluster-management-io/api/pull/162) [@ivan-cai](https://github.com/ivan-cai))
* add help func for placement ([#156](https://github.com/open-cluster-management-io/api/pull/156) [@haoqing0110](https://github.com/haoqing0110))
* Register feature "V1beta1CSRAPICompatibility" into FateureGates ([#169](https://github.com/open-cluster-management-io/api/pull/169) [@Promacanthus](https://github.com/Promacanthus))

### Bug Fixes
* Add subresource in binding crd ([#160](https://github.com/open-cluster-management-io/api/pull/160) [@qiujian16](https://github.com/qiujian16))
* update DefaultClusterSet feature gate description ([#163](https://github.com/open-cluster-management-io/api/pull/163) [@ldpliu](https://github.com/ldpliu))
* Add GOPATH, GOHOSTOS and GOHOSTARCH in makefile ([#168](https://github.com/open-cluster-management-io/api/pull/168) [@Promacanthus](https://github.com/Promacanthus))
* add validation rule to klusterlet namespace ([#170](https://github.com/open-cluster-management-io/api/pull/170) [@elgnay](https://github.com/elgnay))

### Removed & Deprecated
N/A
