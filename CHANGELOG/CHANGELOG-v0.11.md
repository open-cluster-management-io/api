# v0.11.0

## Changelog since v0.11.0

### New Features
* Add installStrategy in addon v1alpha1 API ([#207](https://github.com/open-cluster-management-io/api/pull/207) [@qiujian16](https://github.com/qiujian16))
* add cluster autoapproval feature gate ([#210](https://github.com/open-cluster-management-io/api/pull/210) [@skeeey](https://github.com/skeeey))
* Add ManifestWorkReplicaSet Feature ([#218](https://github.com/open-cluster-management-io/api/pull/218) [@serngawy](https://github.com/serngawy))
* addon rollout strategy ([#208](https://github.com/open-cluster-management-io/api/pull/208) [@haoqing0110](https://github.com/haoqing0110))

### Added
* add evictionStartTime for appliedmanifestwork ([#214](https://github.com/open-cluster-management-io/api/pull/214) [@skeeey](https://github.com/skeeey))
* add addon manager ([#211](https://github.com/open-cluster-management-io/api/pull/211) [@haoqing0110](https://github.com/haoqing0110))
* add addon condition types and reasons ([#213](https://github.com/open-cluster-management-io/api/pull/213) [@zhiweiyin318](https://github.com/zhiweiyin318))
* Add addon lifecycle annotation definition ([#229](https://github.com/open-cluster-management-io/api/pull/229) [@qiujian16](https://github.com/qiujian16))
* add reasons for condition progressing ([#230](https://github.com/open-cluster-management-io/api/pull/230) [@haoqing0110](https://github.com/haoqing0110))
* add addon condition RegistrationApplied ([#231](https://github.com/open-cluster-management-io/api/pull/231) [@zhiweiyin318](https://github.com/zhiweiyin318))
* Add registration configuration ([#237](https://github.com/open-cluster-management-io/api/pull/237) [@qiujian16](https://github.com/qiujian16))
* Add registries for addon deployment config ([#238](https://github.com/open-cluster-management-io/api/pull/238) [@zhujian7](https://github.com/zhujian7))
* add short names for addon CRDs ([#240](https://github.com/open-cluster-management-io/api/pull/240) [@zhiweiyin318](https://github.com/zhiweiyin318))
* Return json raw string in feedback result ([#235](https://github.com/open-cluster-management-io/api/pull/235) [@qiujian16](https://github.com/qiujian16))

### Changes
* clusterset v1beta1 migration ([#202](https://github.com/open-cluster-management-io/api/pull/202) [@ldpliu](https://github.com/ldpliu))
* Set omitempty on the ManifestWork UpdateStrategy field ([#206](https://github.com/open-cluster-management-io/api/pull/206) [@mprahl](https://github.com/mprahl))
* Rename placemanifestwork to manifestworkset ([#205](https://github.com/open-cluster-management-io/api/pull/205) [@qiujian16](https://github.com/qiujian16))
* modify InstallStrategy to structure ([#222](https://github.com/open-cluster-management-io/api/pull/222) [@haoqing0110](https://github.com/haoqing0110))
* add maxLength and pattern for clustername ([#217](https://github.com/open-cluster-management-io/api/pull/217) [@ycyaoxdu](https://github.com/ycyaoxdu))
* upgrade sets to use generic and add cilint ([#227](https://github.com/open-cluster-management-io/api/pull/227) [@ycyaoxdu](https://github.com/ycyaoxdu))
* update work annotations ([#232](https://github.com/open-cluster-management-io/api/pull/232) [@haoqing0110](https://github.com/haoqing0110))
* check work annotation update in ut ([#236](https://github.com/open-cluster-management-io/api/pull/236) [@haoqing0110](https://github.com/haoqing0110))

### Bug Fixes
* Fix inline bug in placement install strategy ([#209](https://github.com/open-cluster-management-io/api/pull/209) [@qiujian16](https://github.com/qiujian16))
* fix work hash confilict and compare mainfests ([#233](https://github.com/open-cluster-management-io/api/pull/233) [@zhiweiyin318](https://github.com/zhiweiyin318))
* skip IsAlreadyExists during create work ([#234](https://github.com/open-cluster-management-io/api/pull/234) [@zhiweiyin318](https://github.com/zhiweiyin318))

### Removed & Deprecated
* Remove Enable field in ClusterManager API ([#228](https://github.com/open-cluster-management-io/api/pull/228) [@qiujian16](https://github.com/qiujian16))
