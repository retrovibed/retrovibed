// This is a generated file - do not edit.
//
// Generated from storage.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class Local extends $pb.GeneratedMessage {
  factory Local({
    $core.bool? reclaim,
    $fixnum.Int64? maximum,
  }) {
    final result = create();
    if (reclaim != null) result.reclaim = reclaim;
    if (maximum != null) result.maximum = maximum;
    return result;
  }

  Local._();

  factory Local.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Local.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Local',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'storage'),
      createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'reclaim')
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'maximum', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Local clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Local copyWith(void Function(Local) updates) =>
      super.copyWith((message) => updates(message as Local)) as Local;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Local create() => Local._();
  @$core.override
  Local createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Local getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Local>(create);
  static Local? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get reclaim => $_getBF(0);
  @$pb.TagNumber(1)
  set reclaim($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasReclaim() => $_has(0);
  @$pb.TagNumber(1)
  void clearReclaim() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get maximum => $_getI64(1);
  @$pb.TagNumber(2)
  set maximum($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasMaximum() => $_has(1);
  @$pb.TagNumber(2)
  void clearMaximum() => $_clearField(2);
}

class StorageSettingsRequest extends $pb.GeneratedMessage {
  factory StorageSettingsRequest({
    Local? local,
  }) {
    final result = create();
    if (local != null) result.local = local;
    return result;
  }

  StorageSettingsRequest._();

  factory StorageSettingsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StorageSettingsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StorageSettingsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'storage'),
      createEmptyInstance: create)
    ..aOM<Local>(1, _omitFieldNames ? '' : 'local', subBuilder: Local.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StorageSettingsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StorageSettingsRequest copyWith(
          void Function(StorageSettingsRequest) updates) =>
      super.copyWith((message) => updates(message as StorageSettingsRequest))
          as StorageSettingsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StorageSettingsRequest create() => StorageSettingsRequest._();
  @$core.override
  StorageSettingsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StorageSettingsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StorageSettingsRequest>(create);
  static StorageSettingsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Local get local => $_getN(0);
  @$pb.TagNumber(1)
  set local(Local value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasLocal() => $_has(0);
  @$pb.TagNumber(1)
  void clearLocal() => $_clearField(1);
  @$pb.TagNumber(1)
  Local ensureLocal() => $_ensure(0);
}

class StorageSettingsResponse extends $pb.GeneratedMessage {
  factory StorageSettingsResponse({
    Local? local,
  }) {
    final result = create();
    if (local != null) result.local = local;
    return result;
  }

  StorageSettingsResponse._();

  factory StorageSettingsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StorageSettingsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StorageSettingsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'storage'),
      createEmptyInstance: create)
    ..aOM<Local>(1, _omitFieldNames ? '' : 'local', subBuilder: Local.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StorageSettingsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StorageSettingsResponse copyWith(
          void Function(StorageSettingsResponse) updates) =>
      super.copyWith((message) => updates(message as StorageSettingsResponse))
          as StorageSettingsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StorageSettingsResponse create() => StorageSettingsResponse._();
  @$core.override
  StorageSettingsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StorageSettingsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StorageSettingsResponse>(create);
  static StorageSettingsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Local get local => $_getN(0);
  @$pb.TagNumber(1)
  set local(Local value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasLocal() => $_has(0);
  @$pb.TagNumber(1)
  void clearLocal() => $_clearField(1);
  @$pb.TagNumber(1)
  Local ensureLocal() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
