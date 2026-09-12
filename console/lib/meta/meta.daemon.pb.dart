// This is a generated file - do not edit.
//
// Generated from meta/meta.daemon.proto.

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

class Daemon extends $pb.GeneratedMessage {
  factory Daemon({
    $core.String? id,
    $core.String? createdAt,
    $core.String? updatedAt,
    $core.String? description,
    $core.String? hostname,
    $core.bool? default_100,
    $core.bool? downloads,
  }) {
    final result = Daemon._();
    if (id != null) result.id = id;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (description != null) result.description = description;
    if (hostname != null) result.hostname = hostname;
    if (default_100 != null) result.default_100 = default_100;
    if (downloads != null) result.downloads = downloads;
    return result;
  }

  Daemon._();

  factory Daemon.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Daemon()..mergeFromBuffer(data, registry);
  factory Daemon.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Daemon()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Daemon',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: Daemon.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'created_at')
    ..aOS(3, _omitFieldNames ? '' : 'updated_at')
    ..aOS(4, _omitFieldNames ? '' : 'description')
    ..aOS(5, _omitFieldNames ? '' : 'hostname')
    ..aOB(100, _omitFieldNames ? '' : 'default')
    ..aOB(101, _omitFieldNames ? '' : 'downloads')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Daemon clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Daemon copyWith(void Function(Daemon) updates) =>
      super.copyWith((message) => updates(message as Daemon)) as Daemon;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use Daemon() / Daemon.new instead')
  static Daemon create() => Daemon._();
  static $pb.GeneratedMessage $_createMessage() => Daemon._();
  @$core.override
  Daemon createEmptyInstance() => Daemon._();
  @$core.pragma('dart2js:noInline')
  static Daemon getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Daemon>(Daemon.$_createMessage);
  static Daemon? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get createdAt => $_getSZ(1);
  @$pb.TagNumber(2)
  set createdAt($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasCreatedAt() => $_has(1);
  @$pb.TagNumber(2)
  void clearCreatedAt() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get updatedAt => $_getSZ(2);
  @$pb.TagNumber(3)
  set updatedAt($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasUpdatedAt() => $_has(2);
  @$pb.TagNumber(3)
  void clearUpdatedAt() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(4)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(4)
  void clearDescription() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get hostname => $_getSZ(4);
  @$pb.TagNumber(5)
  set hostname($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasHostname() => $_has(4);
  @$pb.TagNumber(5)
  void clearHostname() => $_clearField(5);

  @$pb.TagNumber(100)
  $core.bool get default_100 => $_getBF(5);
  @$pb.TagNumber(100)
  set default_100($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(100)
  $core.bool hasDefault_100() => $_has(5);
  @$pb.TagNumber(100)
  void clearDefault_100() => $_clearField(100);

  @$pb.TagNumber(101)
  $core.bool get downloads => $_getBF(6);
  @$pb.TagNumber(101)
  set downloads($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(101)
  $core.bool hasDownloads() => $_has(6);
  @$pb.TagNumber(101)
  void clearDownloads() => $_clearField(101);
}

class DaemonSearchRequest extends $pb.GeneratedMessage {
  factory DaemonSearchRequest({
    $core.String? query,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = DaemonSearchRequest._();
    if (query != null) result.query = query;
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    return result;
  }

  DaemonSearchRequest._();

  factory DaemonSearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonSearchRequest()..mergeFromBuffer(data, registry);
  factory DaemonSearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonSearchRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DaemonSearchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: DaemonSearchRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(3, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonSearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonSearchRequest copyWith(void Function(DaemonSearchRequest) updates) =>
      super.copyWith((message) => updates(message as DaemonSearchRequest))
          as DaemonSearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use DaemonSearchRequest() / DaemonSearchRequest.new instead')
  static DaemonSearchRequest create() => DaemonSearchRequest._();
  static $pb.GeneratedMessage $_createMessage() => DaemonSearchRequest._();
  @$core.override
  DaemonSearchRequest createEmptyInstance() => DaemonSearchRequest._();
  @$core.pragma('dart2js:noInline')
  static DaemonSearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DaemonSearchRequest>(
          DaemonSearchRequest.$_createMessage);
  static DaemonSearchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get query => $_getSZ(0);
  @$pb.TagNumber(1)
  set query($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasQuery() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuery() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get offset => $_getI64(1);
  @$pb.TagNumber(2)
  set offset($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasOffset() => $_has(1);
  @$pb.TagNumber(2)
  void clearOffset() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get limit => $_getI64(2);
  @$pb.TagNumber(3)
  set limit($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(3)
  void clearLimit() => $_clearField(3);
}

class DaemonSearchResponse extends $pb.GeneratedMessage {
  factory DaemonSearchResponse({
    DaemonSearchRequest? next,
    $core.Iterable<Daemon>? items,
  }) {
    final result = DaemonSearchResponse._();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  DaemonSearchResponse._();

  factory DaemonSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonSearchResponse()..mergeFromBuffer(data, registry);
  factory DaemonSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonSearchResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DaemonSearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: DaemonSearchResponse.$_createMessage)
    ..aOM<DaemonSearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: DaemonSearchRequest.$_createMessage)
    ..pPM<Daemon>(2, _omitFieldNames ? '' : 'items',
        subBuilder: Daemon.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonSearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonSearchResponse copyWith(void Function(DaemonSearchResponse) updates) =>
      super.copyWith((message) => updates(message as DaemonSearchResponse))
          as DaemonSearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use DaemonSearchResponse() / DaemonSearchResponse.new instead')
  static DaemonSearchResponse create() => DaemonSearchResponse._();
  static $pb.GeneratedMessage $_createMessage() => DaemonSearchResponse._();
  @$core.override
  DaemonSearchResponse createEmptyInstance() => DaemonSearchResponse._();
  @$core.pragma('dart2js:noInline')
  static DaemonSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DaemonSearchResponse>(
          DaemonSearchResponse.$_createMessage);
  static DaemonSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  DaemonSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(DaemonSearchRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  DaemonSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<Daemon> get items => $_getList(1);
}

class DaemonCreateRequest extends $pb.GeneratedMessage {
  factory DaemonCreateRequest({
    Daemon? daemon,
  }) {
    final result = DaemonCreateRequest._();
    if (daemon != null) result.daemon = daemon;
    return result;
  }

  DaemonCreateRequest._();

  factory DaemonCreateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonCreateRequest()..mergeFromBuffer(data, registry);
  factory DaemonCreateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonCreateRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DaemonCreateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: DaemonCreateRequest.$_createMessage)
    ..aOM<Daemon>(1, _omitFieldNames ? '' : 'daemon',
        subBuilder: Daemon.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonCreateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonCreateRequest copyWith(void Function(DaemonCreateRequest) updates) =>
      super.copyWith((message) => updates(message as DaemonCreateRequest))
          as DaemonCreateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use DaemonCreateRequest() / DaemonCreateRequest.new instead')
  static DaemonCreateRequest create() => DaemonCreateRequest._();
  static $pb.GeneratedMessage $_createMessage() => DaemonCreateRequest._();
  @$core.override
  DaemonCreateRequest createEmptyInstance() => DaemonCreateRequest._();
  @$core.pragma('dart2js:noInline')
  static DaemonCreateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DaemonCreateRequest>(
          DaemonCreateRequest.$_createMessage);
  static DaemonCreateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Daemon get daemon => $_getN(0);
  @$pb.TagNumber(1)
  set daemon(Daemon value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDaemon() => $_has(0);
  @$pb.TagNumber(1)
  void clearDaemon() => $_clearField(1);
  @$pb.TagNumber(1)
  Daemon ensureDaemon() => $_ensure(0);
}

class DaemonCreateResponse extends $pb.GeneratedMessage {
  factory DaemonCreateResponse({
    Daemon? daemon,
  }) {
    final result = DaemonCreateResponse._();
    if (daemon != null) result.daemon = daemon;
    return result;
  }

  DaemonCreateResponse._();

  factory DaemonCreateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonCreateResponse()..mergeFromBuffer(data, registry);
  factory DaemonCreateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonCreateResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DaemonCreateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: DaemonCreateResponse.$_createMessage)
    ..aOM<Daemon>(1, _omitFieldNames ? '' : 'daemon',
        subBuilder: Daemon.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonCreateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonCreateResponse copyWith(void Function(DaemonCreateResponse) updates) =>
      super.copyWith((message) => updates(message as DaemonCreateResponse))
          as DaemonCreateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use DaemonCreateResponse() / DaemonCreateResponse.new instead')
  static DaemonCreateResponse create() => DaemonCreateResponse._();
  static $pb.GeneratedMessage $_createMessage() => DaemonCreateResponse._();
  @$core.override
  DaemonCreateResponse createEmptyInstance() => DaemonCreateResponse._();
  @$core.pragma('dart2js:noInline')
  static DaemonCreateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DaemonCreateResponse>(
          DaemonCreateResponse.$_createMessage);
  static DaemonCreateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Daemon get daemon => $_getN(0);
  @$pb.TagNumber(1)
  set daemon(Daemon value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDaemon() => $_has(0);
  @$pb.TagNumber(1)
  void clearDaemon() => $_clearField(1);
  @$pb.TagNumber(1)
  Daemon ensureDaemon() => $_ensure(0);
}

class DaemonDeleteRequest extends $pb.GeneratedMessage {
  factory DaemonDeleteRequest() => DaemonDeleteRequest._();

  DaemonDeleteRequest._();

  factory DaemonDeleteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonDeleteRequest()..mergeFromBuffer(data, registry);
  factory DaemonDeleteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonDeleteRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DaemonDeleteRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: DaemonDeleteRequest.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonDeleteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonDeleteRequest copyWith(void Function(DaemonDeleteRequest) updates) =>
      super.copyWith((message) => updates(message as DaemonDeleteRequest))
          as DaemonDeleteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use DaemonDeleteRequest() / DaemonDeleteRequest.new instead')
  static DaemonDeleteRequest create() => DaemonDeleteRequest._();
  static $pb.GeneratedMessage $_createMessage() => DaemonDeleteRequest._();
  @$core.override
  DaemonDeleteRequest createEmptyInstance() => DaemonDeleteRequest._();
  @$core.pragma('dart2js:noInline')
  static DaemonDeleteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DaemonDeleteRequest>(
          DaemonDeleteRequest.$_createMessage);
  static DaemonDeleteRequest? _defaultInstance;
}

class DaemonDeleteResponse extends $pb.GeneratedMessage {
  factory DaemonDeleteResponse({
    Daemon? daemon,
  }) {
    final result = DaemonDeleteResponse._();
    if (daemon != null) result.daemon = daemon;
    return result;
  }

  DaemonDeleteResponse._();

  factory DaemonDeleteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonDeleteResponse()..mergeFromBuffer(data, registry);
  factory DaemonDeleteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonDeleteResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DaemonDeleteResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: DaemonDeleteResponse.$_createMessage)
    ..aOM<Daemon>(1, _omitFieldNames ? '' : 'daemon',
        subBuilder: Daemon.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonDeleteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonDeleteResponse copyWith(void Function(DaemonDeleteResponse) updates) =>
      super.copyWith((message) => updates(message as DaemonDeleteResponse))
          as DaemonDeleteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use DaemonDeleteResponse() / DaemonDeleteResponse.new instead')
  static DaemonDeleteResponse create() => DaemonDeleteResponse._();
  static $pb.GeneratedMessage $_createMessage() => DaemonDeleteResponse._();
  @$core.override
  DaemonDeleteResponse createEmptyInstance() => DaemonDeleteResponse._();
  @$core.pragma('dart2js:noInline')
  static DaemonDeleteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DaemonDeleteResponse>(
          DaemonDeleteResponse.$_createMessage);
  static DaemonDeleteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Daemon get daemon => $_getN(0);
  @$pb.TagNumber(1)
  set daemon(Daemon value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDaemon() => $_has(0);
  @$pb.TagNumber(1)
  void clearDaemon() => $_clearField(1);
  @$pb.TagNumber(1)
  Daemon ensureDaemon() => $_ensure(0);
}

class DaemonLookupRequest extends $pb.GeneratedMessage {
  factory DaemonLookupRequest() => DaemonLookupRequest._();

  DaemonLookupRequest._();

  factory DaemonLookupRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonLookupRequest()..mergeFromBuffer(data, registry);
  factory DaemonLookupRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonLookupRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DaemonLookupRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: DaemonLookupRequest.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonLookupRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonLookupRequest copyWith(void Function(DaemonLookupRequest) updates) =>
      super.copyWith((message) => updates(message as DaemonLookupRequest))
          as DaemonLookupRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use DaemonLookupRequest() / DaemonLookupRequest.new instead')
  static DaemonLookupRequest create() => DaemonLookupRequest._();
  static $pb.GeneratedMessage $_createMessage() => DaemonLookupRequest._();
  @$core.override
  DaemonLookupRequest createEmptyInstance() => DaemonLookupRequest._();
  @$core.pragma('dart2js:noInline')
  static DaemonLookupRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DaemonLookupRequest>(
          DaemonLookupRequest.$_createMessage);
  static DaemonLookupRequest? _defaultInstance;
}

class DaemonLookupResponse extends $pb.GeneratedMessage {
  factory DaemonLookupResponse({
    Daemon? daemon,
  }) {
    final result = DaemonLookupResponse._();
    if (daemon != null) result.daemon = daemon;
    return result;
  }

  DaemonLookupResponse._();

  factory DaemonLookupResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonLookupResponse()..mergeFromBuffer(data, registry);
  factory DaemonLookupResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonLookupResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DaemonLookupResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: DaemonLookupResponse.$_createMessage)
    ..aOM<Daemon>(1, _omitFieldNames ? '' : 'daemon',
        subBuilder: Daemon.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonLookupResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonLookupResponse copyWith(void Function(DaemonLookupResponse) updates) =>
      super.copyWith((message) => updates(message as DaemonLookupResponse))
          as DaemonLookupResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use DaemonLookupResponse() / DaemonLookupResponse.new instead')
  static DaemonLookupResponse create() => DaemonLookupResponse._();
  static $pb.GeneratedMessage $_createMessage() => DaemonLookupResponse._();
  @$core.override
  DaemonLookupResponse createEmptyInstance() => DaemonLookupResponse._();
  @$core.pragma('dart2js:noInline')
  static DaemonLookupResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DaemonLookupResponse>(
          DaemonLookupResponse.$_createMessage);
  static DaemonLookupResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Daemon get daemon => $_getN(0);
  @$pb.TagNumber(1)
  set daemon(Daemon value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDaemon() => $_has(0);
  @$pb.TagNumber(1)
  void clearDaemon() => $_clearField(1);
  @$pb.TagNumber(1)
  Daemon ensureDaemon() => $_ensure(0);
}

class DaemonUpdateRequest extends $pb.GeneratedMessage {
  factory DaemonUpdateRequest({
    Daemon? daemon,
  }) {
    final result = DaemonUpdateRequest._();
    if (daemon != null) result.daemon = daemon;
    return result;
  }

  DaemonUpdateRequest._();

  factory DaemonUpdateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonUpdateRequest()..mergeFromBuffer(data, registry);
  factory DaemonUpdateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonUpdateRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DaemonUpdateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: DaemonUpdateRequest.$_createMessage)
    ..aOM<Daemon>(1, _omitFieldNames ? '' : 'daemon',
        subBuilder: Daemon.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonUpdateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonUpdateRequest copyWith(void Function(DaemonUpdateRequest) updates) =>
      super.copyWith((message) => updates(message as DaemonUpdateRequest))
          as DaemonUpdateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use DaemonUpdateRequest() / DaemonUpdateRequest.new instead')
  static DaemonUpdateRequest create() => DaemonUpdateRequest._();
  static $pb.GeneratedMessage $_createMessage() => DaemonUpdateRequest._();
  @$core.override
  DaemonUpdateRequest createEmptyInstance() => DaemonUpdateRequest._();
  @$core.pragma('dart2js:noInline')
  static DaemonUpdateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DaemonUpdateRequest>(
          DaemonUpdateRequest.$_createMessage);
  static DaemonUpdateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Daemon get daemon => $_getN(0);
  @$pb.TagNumber(1)
  set daemon(Daemon value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDaemon() => $_has(0);
  @$pb.TagNumber(1)
  void clearDaemon() => $_clearField(1);
  @$pb.TagNumber(1)
  Daemon ensureDaemon() => $_ensure(0);
}

class DaemonUpdateResponse extends $pb.GeneratedMessage {
  factory DaemonUpdateResponse({
    Daemon? daemon,
  }) {
    final result = DaemonUpdateResponse._();
    if (daemon != null) result.daemon = daemon;
    return result;
  }

  DaemonUpdateResponse._();

  factory DaemonUpdateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonUpdateResponse()..mergeFromBuffer(data, registry);
  factory DaemonUpdateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonUpdateResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DaemonUpdateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: DaemonUpdateResponse.$_createMessage)
    ..aOM<Daemon>(1, _omitFieldNames ? '' : 'daemon',
        subBuilder: Daemon.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonUpdateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonUpdateResponse copyWith(void Function(DaemonUpdateResponse) updates) =>
      super.copyWith((message) => updates(message as DaemonUpdateResponse))
          as DaemonUpdateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use DaemonUpdateResponse() / DaemonUpdateResponse.new instead')
  static DaemonUpdateResponse create() => DaemonUpdateResponse._();
  static $pb.GeneratedMessage $_createMessage() => DaemonUpdateResponse._();
  @$core.override
  DaemonUpdateResponse createEmptyInstance() => DaemonUpdateResponse._();
  @$core.pragma('dart2js:noInline')
  static DaemonUpdateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DaemonUpdateResponse>(
          DaemonUpdateResponse.$_createMessage);
  static DaemonUpdateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Daemon get daemon => $_getN(0);
  @$pb.TagNumber(1)
  set daemon(Daemon value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDaemon() => $_has(0);
  @$pb.TagNumber(1)
  void clearDaemon() => $_clearField(1);
  @$pb.TagNumber(1)
  Daemon ensureDaemon() => $_ensure(0);
}

class DaemonDisableRequest extends $pb.GeneratedMessage {
  factory DaemonDisableRequest({
    Daemon? daemon,
  }) {
    final result = DaemonDisableRequest._();
    if (daemon != null) result.daemon = daemon;
    return result;
  }

  DaemonDisableRequest._();

  factory DaemonDisableRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonDisableRequest()..mergeFromBuffer(data, registry);
  factory DaemonDisableRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonDisableRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DaemonDisableRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: DaemonDisableRequest.$_createMessage)
    ..aOM<Daemon>(1, _omitFieldNames ? '' : 'daemon',
        subBuilder: Daemon.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonDisableRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonDisableRequest copyWith(void Function(DaemonDisableRequest) updates) =>
      super.copyWith((message) => updates(message as DaemonDisableRequest))
          as DaemonDisableRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use DaemonDisableRequest() / DaemonDisableRequest.new instead')
  static DaemonDisableRequest create() => DaemonDisableRequest._();
  static $pb.GeneratedMessage $_createMessage() => DaemonDisableRequest._();
  @$core.override
  DaemonDisableRequest createEmptyInstance() => DaemonDisableRequest._();
  @$core.pragma('dart2js:noInline')
  static DaemonDisableRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DaemonDisableRequest>(
          DaemonDisableRequest.$_createMessage);
  static DaemonDisableRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Daemon get daemon => $_getN(0);
  @$pb.TagNumber(1)
  set daemon(Daemon value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDaemon() => $_has(0);
  @$pb.TagNumber(1)
  void clearDaemon() => $_clearField(1);
  @$pb.TagNumber(1)
  Daemon ensureDaemon() => $_ensure(0);
}

class DaemonDisableResponse extends $pb.GeneratedMessage {
  factory DaemonDisableResponse({
    Daemon? daemon,
  }) {
    final result = DaemonDisableResponse._();
    if (daemon != null) result.daemon = daemon;
    return result;
  }

  DaemonDisableResponse._();

  factory DaemonDisableResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonDisableResponse()..mergeFromBuffer(data, registry);
  factory DaemonDisableResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DaemonDisableResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DaemonDisableResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: DaemonDisableResponse.$_createMessage)
    ..aOM<Daemon>(1, _omitFieldNames ? '' : 'daemon',
        subBuilder: Daemon.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonDisableResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DaemonDisableResponse copyWith(
          void Function(DaemonDisableResponse) updates) =>
      super.copyWith((message) => updates(message as DaemonDisableResponse))
          as DaemonDisableResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use DaemonDisableResponse() / DaemonDisableResponse.new instead')
  static DaemonDisableResponse create() => DaemonDisableResponse._();
  static $pb.GeneratedMessage $_createMessage() => DaemonDisableResponse._();
  @$core.override
  DaemonDisableResponse createEmptyInstance() => DaemonDisableResponse._();
  @$core.pragma('dart2js:noInline')
  static DaemonDisableResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DaemonDisableResponse>(
          DaemonDisableResponse.$_createMessage);
  static DaemonDisableResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Daemon get daemon => $_getN(0);
  @$pb.TagNumber(1)
  set daemon(Daemon value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDaemon() => $_has(0);
  @$pb.TagNumber(1)
  void clearDaemon() => $_clearField(1);
  @$pb.TagNumber(1)
  Daemon ensureDaemon() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
