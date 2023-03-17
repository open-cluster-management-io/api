# v0.6.0

## Changelog since v0.5.0

### New Features 
* Add API AddOnPlacementScore. ([#109](https://github.com/open-cluster-management-io/api/pull/109) [@haoqing0110](https://github.com/haoqing0110))
* Support taint & toleration. ([#115](https://github.com/open-cluster-management-io/api/pull/115) [@elgnay](https://github.com/elgnay))
* Add status feedback fields to work api. ([#108](https://github.com/open-cluster-management-io/api/pull/108) [@qiujian16](https://github.com/qiujian16))

### Added
* Add detached mode spec fields in clusterManager ([#103](https://github.com/open-cluster-management-io/api/pull/103) [@xuezhaojun](https://github.com/xuezhaojun))
* Add detached deploy mode for klusterlet ([#114](https://github.com/open-cluster-management-io/api/pull/114) [@zhujian7](https://github.com/zhujian7))
* Add HealthCheck spec to managedclusteraddon ([#113](https://github.com/open-cluster-management-io/api/pull/113) [@skeeey](https://github.com/skeeey))
* Add const ScoreCoordinate type. ([#126](https://github.com/open-cluster-management-io/api/pull/126) [@haoqing0110](https://github.com/haoqing0110))
* Addon observed generation in managedclusteraddon api. ([#131](https://github.com/open-cluster-management-io/api/pull/131) [@yue9944882](https://github.com/yue9944882))
* Add PlacementConditionMisconfigured. ([#132](https://github.com/open-cluster-management-io/api/pull/132) [@haoqing0110](https://github.com/haoqing0110))

### Changes
* Update docs for cluster join process. ([#101](https://github.com/open-cluster-management-io/api/pull/101) [@zhujian7](https://github.com/zhujian7))
* Set managedcluster leaseduration seconds value with 60. ([#111](https://github.com/open-cluster-management-io/api/pull/111) [@champly](https://github.com/champly))
* Bump the go version to 1.17. ([#116](https://github.com/open-cluster-management-io/api/pull/116) [@zhujian7](https://github.com/zhujian7))
* Fixed issue #112 separated types.go. ([#123](https://github.com/open-cluster-management-io/api/pull/123) [@ilan-pinto](https://github.com/ilan-pinto))
* Update validation rule of taint.timeAdded. ([#125](https://github.com/open-cluster-management-io/api/pull/125) [@elgnay](https://github.com/elgnay))
* Make timeAdded of taint nullable. ([#129](https://github.com/open-cluster-management-io/api/pull/129) [@elgnay](https://github.com/elgnay))

### Bug Fixes
* Update fedora version to 34. ([#105](https://github.com/open-cluster-management-io/api/pull/105) [@xuezhaojun](https://github.com/xuezhaojun))
* Update securitymd to refer to the OCM community security md file. ([#106](https://github.com/open-cluster-management-io/api/pull/106) [@mikeshng](https://github.com/mikesng))
* Upgrade CONTROLLER_GEN_VERSION to v0.6.0. ([#110](https://github.com/open-cluster-management-io/api/pull/110) [@haoqing0110](https://github.com/haoqing0110))
* Fix ci verify. ([#117](https://github.com/open-cluster-management-io/api/pull/117) [@zhujian7](https://github.com/zhujian7))
* Change default value of deleteOption. ([#121](https://github.com/open-cluster-management-io/api/pull/121) [@qiujian16](https://github.com/qiujian16))
* Bugfix: health check mode empty. ([#127](https://github.com/open-cluster-management-io/api/pull/127) [@yue9944882](https://github.com/yue9944882))

### Removed & Deprecated
N/A
