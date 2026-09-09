// This is a generated file - do not edit.
//
// Generated from wireguard/meta.wireguard.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class WireguardNettype extends $pb.ProtobufEnum {
  static const WireguardNettype UNSPECIFIED =
      WireguardNettype._(0, _omitEnumNames ? '' : 'UNSPECIFIED');
  static const WireguardNettype DISTRIBUTION =
      WireguardNettype._(1, _omitEnumNames ? '' : 'DISTRIBUTION');
  static const WireguardNettype SOCIAL =
      WireguardNettype._(2, _omitEnumNames ? '' : 'SOCIAL');

  static const $core.List<WireguardNettype> values = <WireguardNettype>[
    UNSPECIFIED,
    DISTRIBUTION,
    SOCIAL,
  ];

  static final $core.List<WireguardNettype?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static WireguardNettype? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const WireguardNettype._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
