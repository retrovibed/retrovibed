// This is a generated file - do not edit.
//
// Generated from library/library.watch.proto.

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

class WatchHistoryRecord extends $pb.GeneratedMessage {
  factory WatchHistoryRecord({
    $core.String? id,
    $core.String? mediaId,
    $fixnum.Int64? duration,
  }) {
    final result = WatchHistoryRecord._();
    if (id != null) result.id = id;
    if (mediaId != null) result.mediaId = mediaId;
    if (duration != null) result.duration = duration;
    return result;
  }

  WatchHistoryRecord._();

  factory WatchHistoryRecord.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WatchHistoryRecord()..mergeFromBuffer(data, registry);
  factory WatchHistoryRecord.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WatchHistoryRecord()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WatchHistoryRecord',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'library'),
      createEmptyInstance: WatchHistoryRecord.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'media_id')
    ..a<$fixnum.Int64>(
        3, _omitFieldNames ? '' : 'duration', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WatchHistoryRecord clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WatchHistoryRecord copyWith(void Function(WatchHistoryRecord) updates) =>
      super.copyWith((message) => updates(message as WatchHistoryRecord))
          as WatchHistoryRecord;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use WatchHistoryRecord() / WatchHistoryRecord.new instead')
  static WatchHistoryRecord create() => WatchHistoryRecord._();
  static $pb.GeneratedMessage $_createMessage() => WatchHistoryRecord._();
  @$core.override
  WatchHistoryRecord createEmptyInstance() => WatchHistoryRecord._();
  @$core.pragma('dart2js:noInline')
  static WatchHistoryRecord getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WatchHistoryRecord>(
          WatchHistoryRecord.$_createMessage);
  static WatchHistoryRecord? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get mediaId => $_getSZ(1);
  @$pb.TagNumber(2)
  set mediaId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasMediaId() => $_has(1);
  @$pb.TagNumber(2)
  void clearMediaId() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get duration => $_getI64(2);
  @$pb.TagNumber(3)
  set duration($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasDuration() => $_has(2);
  @$pb.TagNumber(3)
  void clearDuration() => $_clearField(3);
}

class WatchHistoryRecordRequest extends $pb.GeneratedMessage {
  factory WatchHistoryRecordRequest({
    WatchHistoryRecord? record,
  }) {
    final result = WatchHistoryRecordRequest._();
    if (record != null) result.record = record;
    return result;
  }

  WatchHistoryRecordRequest._();

  factory WatchHistoryRecordRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WatchHistoryRecordRequest()..mergeFromBuffer(data, registry);
  factory WatchHistoryRecordRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WatchHistoryRecordRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WatchHistoryRecordRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'library'),
      createEmptyInstance: WatchHistoryRecordRequest.$_createMessage)
    ..aOM<WatchHistoryRecord>(1, _omitFieldNames ? '' : 'record',
        subBuilder: WatchHistoryRecord.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WatchHistoryRecordRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WatchHistoryRecordRequest copyWith(
          void Function(WatchHistoryRecordRequest) updates) =>
      super.copyWith((message) => updates(message as WatchHistoryRecordRequest))
          as WatchHistoryRecordRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use WatchHistoryRecordRequest() / WatchHistoryRecordRequest.new instead')
  static WatchHistoryRecordRequest create() => WatchHistoryRecordRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      WatchHistoryRecordRequest._();
  @$core.override
  WatchHistoryRecordRequest createEmptyInstance() =>
      WatchHistoryRecordRequest._();
  @$core.pragma('dart2js:noInline')
  static WatchHistoryRecordRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WatchHistoryRecordRequest>(
          WatchHistoryRecordRequest.$_createMessage);
  static WatchHistoryRecordRequest? _defaultInstance;

  @$pb.TagNumber(1)
  WatchHistoryRecord get record => $_getN(0);
  @$pb.TagNumber(1)
  set record(WatchHistoryRecord value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasRecord() => $_has(0);
  @$pb.TagNumber(1)
  void clearRecord() => $_clearField(1);
  @$pb.TagNumber(1)
  WatchHistoryRecord ensureRecord() => $_ensure(0);
}

class WatchHistoryRecordResponse extends $pb.GeneratedMessage {
  factory WatchHistoryRecordResponse() => WatchHistoryRecordResponse._();

  WatchHistoryRecordResponse._();

  factory WatchHistoryRecordResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WatchHistoryRecordResponse()..mergeFromBuffer(data, registry);
  factory WatchHistoryRecordResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WatchHistoryRecordResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WatchHistoryRecordResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'library'),
      createEmptyInstance: WatchHistoryRecordResponse.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WatchHistoryRecordResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WatchHistoryRecordResponse copyWith(
          void Function(WatchHistoryRecordResponse) updates) =>
      super.copyWith(
              (message) => updates(message as WatchHistoryRecordResponse))
          as WatchHistoryRecordResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use WatchHistoryRecordResponse() / WatchHistoryRecordResponse.new instead')
  static WatchHistoryRecordResponse create() => WatchHistoryRecordResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      WatchHistoryRecordResponse._();
  @$core.override
  WatchHistoryRecordResponse createEmptyInstance() =>
      WatchHistoryRecordResponse._();
  @$core.pragma('dart2js:noInline')
  static WatchHistoryRecordResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WatchHistoryRecordResponse>(
          WatchHistoryRecordResponse.$_createMessage);
  static WatchHistoryRecordResponse? _defaultInstance;
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
