// This is a generated file - do not edit.
//
// Generated from ddisc/ddisc.discovery.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class AcquisitionState extends $pb.ProtobufEnum {
  static const AcquisitionState Unknown =
      AcquisitionState._(0, _omitEnumNames ? '' : 'Unknown');
  static const AcquisitionState Ephemeral =
      AcquisitionState._(1, _omitEnumNames ? '' : 'Ephemeral');
  static const AcquisitionState Available =
      AcquisitionState._(2, _omitEnumNames ? '' : 'Available');

  /// means it was put into the recommendations queue.
  static const AcquisitionState Downloading =
      AcquisitionState._(3, _omitEnumNames ? '' : 'Downloading');
  static const AcquisitionState Completed =
      AcquisitionState._(4, _omitEnumNames ? '' : 'Completed');

  static const $core.List<AcquisitionState> values = <AcquisitionState>[
    Unknown,
    Ephemeral,
    Available,
    Downloading,
    Completed,
  ];

  static final $core.List<AcquisitionState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static AcquisitionState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AcquisitionState._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
