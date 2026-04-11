// This is a generated file - do not edit.
//
// Generated from meta.search.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class DateRange extends $pb.GeneratedMessage {
  factory DateRange({
    $core.String? oldest,
    $core.String? newest,
  }) {
    final result = create();
    if (oldest != null) result.oldest = oldest;
    if (newest != null) result.newest = newest;
    return result;
  }

  DateRange._();

  factory DateRange.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DateRange.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DateRange',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'oldest')
    ..aOS(2, _omitFieldNames ? '' : 'newest')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DateRange clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DateRange copyWith(void Function(DateRange) updates) =>
      super.copyWith((message) => updates(message as DateRange)) as DateRange;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DateRange create() => DateRange._();
  @$core.override
  DateRange createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DateRange getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DateRange>(create);
  static DateRange? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get oldest => $_getSZ(0);
  @$pb.TagNumber(1)
  set oldest($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasOldest() => $_has(0);
  @$pb.TagNumber(1)
  void clearOldest() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get newest => $_getSZ(1);
  @$pb.TagNumber(2)
  set newest($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasNewest() => $_has(1);
  @$pb.TagNumber(2)
  void clearNewest() => $_clearField(2);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
