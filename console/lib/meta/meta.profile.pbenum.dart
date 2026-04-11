// This is a generated file - do not edit.
//
// Generated from meta.profile.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class ProfileStatus extends $pb.ProtobufEnum {
  static const ProfileStatus NONE =
      ProfileStatus._(0, _omitEnumNames ? '' : 'NONE');
  static const ProfileStatus DISABLED =
      ProfileStatus._(1, _omitEnumNames ? '' : 'DISABLED');
  static const ProfileStatus PENDING =
      ProfileStatus._(2, _omitEnumNames ? '' : 'PENDING');
  static const ProfileStatus ENABLED =
      ProfileStatus._(3, _omitEnumNames ? '' : 'ENABLED');

  static const $core.List<ProfileStatus> values = <ProfileStatus>[
    NONE,
    DISABLED,
    PENDING,
    ENABLED,
  ];

  static final $core.List<ProfileStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static ProfileStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ProfileStatus._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
