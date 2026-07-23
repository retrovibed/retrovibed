// This is a generated file - do not edit.
//
// Generated from ddisc.discovery.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'meta.search.pb.dart' as $0;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class Discovery extends $pb.GeneratedMessage {
  factory Discovery({
    $core.String? id,
    $core.List<$core.int>? infohash,
    $core.int? attempts,
    $core.String? nextCheck,
    $core.String? createdAt,
    $core.String? updatedAt,
    $core.String? title,
    $core.String? description,
    $core.int? health,
    $fixnum.Int64? bytes,
    $core.int? policyRank,
    $core.String? knownMediaId,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (infohash != null) result.infohash = infohash;
    if (attempts != null) result.attempts = attempts;
    if (nextCheck != null) result.nextCheck = nextCheck;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (title != null) result.title = title;
    if (description != null) result.description = description;
    if (health != null) result.health = health;
    if (bytes != null) result.bytes = bytes;
    if (policyRank != null) result.policyRank = policyRank;
    if (knownMediaId != null) result.knownMediaId = knownMediaId;
    return result;
  }

  Discovery._();

  factory Discovery.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Discovery.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Discovery',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..a<$core.List<$core.int>>(
        2, _omitFieldNames ? '' : 'infohash', $pb.PbFieldType.OY)
    ..aI(3, _omitFieldNames ? '' : 'attempts', fieldType: $pb.PbFieldType.OU3)
    ..aOS(4, _omitFieldNames ? '' : 'next_check')
    ..aOS(5, _omitFieldNames ? '' : 'created_at')
    ..aOS(6, _omitFieldNames ? '' : 'updated_at')
    ..aOS(7, _omitFieldNames ? '' : 'title')
    ..aOS(8, _omitFieldNames ? '' : 'description')
    ..aI(9, _omitFieldNames ? '' : 'health', fieldType: $pb.PbFieldType.OU3)
    ..a<$fixnum.Int64>(10, _omitFieldNames ? '' : 'bytes', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aI(11, _omitFieldNames ? '' : 'policy_rank',
        fieldType: $pb.PbFieldType.OU3)
    ..aOS(12, _omitFieldNames ? '' : 'known_media_id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Discovery clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Discovery copyWith(void Function(Discovery) updates) =>
      super.copyWith((message) => updates(message as Discovery)) as Discovery;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Discovery create() => Discovery._();
  @$core.override
  Discovery createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Discovery getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Discovery>(create);
  static Discovery? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.List<$core.int> get infohash => $_getN(1);
  @$pb.TagNumber(2)
  set infohash($core.List<$core.int> value) => $_setBytes(1, value);
  @$pb.TagNumber(2)
  $core.bool hasInfohash() => $_has(1);
  @$pb.TagNumber(2)
  void clearInfohash() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.int get attempts => $_getIZ(2);
  @$pb.TagNumber(3)
  set attempts($core.int value) => $_setUnsignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasAttempts() => $_has(2);
  @$pb.TagNumber(3)
  void clearAttempts() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get nextCheck => $_getSZ(3);
  @$pb.TagNumber(4)
  set nextCheck($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasNextCheck() => $_has(3);
  @$pb.TagNumber(4)
  void clearNextCheck() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get createdAt => $_getSZ(4);
  @$pb.TagNumber(5)
  set createdAt($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasCreatedAt() => $_has(4);
  @$pb.TagNumber(5)
  void clearCreatedAt() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get updatedAt => $_getSZ(5);
  @$pb.TagNumber(6)
  set updatedAt($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasUpdatedAt() => $_has(5);
  @$pb.TagNumber(6)
  void clearUpdatedAt() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.String get title => $_getSZ(6);
  @$pb.TagNumber(7)
  set title($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasTitle() => $_has(6);
  @$pb.TagNumber(7)
  void clearTitle() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.String get description => $_getSZ(7);
  @$pb.TagNumber(8)
  set description($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasDescription() => $_has(7);
  @$pb.TagNumber(8)
  void clearDescription() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.int get health => $_getIZ(8);
  @$pb.TagNumber(9)
  set health($core.int value) => $_setUnsignedInt32(8, value);
  @$pb.TagNumber(9)
  $core.bool hasHealth() => $_has(8);
  @$pb.TagNumber(9)
  void clearHealth() => $_clearField(9);

  @$pb.TagNumber(10)
  $fixnum.Int64 get bytes => $_getI64(9);
  @$pb.TagNumber(10)
  set bytes($fixnum.Int64 value) => $_setInt64(9, value);
  @$pb.TagNumber(10)
  $core.bool hasBytes() => $_has(9);
  @$pb.TagNumber(10)
  void clearBytes() => $_clearField(10);

  @$pb.TagNumber(11)
  $core.int get policyRank => $_getIZ(10);
  @$pb.TagNumber(11)
  set policyRank($core.int value) => $_setUnsignedInt32(10, value);
  @$pb.TagNumber(11)
  $core.bool hasPolicyRank() => $_has(10);
  @$pb.TagNumber(11)
  void clearPolicyRank() => $_clearField(11);

  @$pb.TagNumber(12)
  $core.String get knownMediaId => $_getSZ(11);
  @$pb.TagNumber(12)
  set knownMediaId($core.String value) => $_setString(11, value);
  @$pb.TagNumber(12)
  $core.bool hasKnownMediaId() => $_has(11);
  @$pb.TagNumber(12)
  void clearKnownMediaId() => $_clearField(12);
}

class DiscoverySearchRequest extends $pb.GeneratedMessage {
  factory DiscoverySearchRequest({
    $0.DateRange? nextCheck,
    $core.Iterable<$core.String>? id,
    $fixnum.Int64? attemptsMin,
    $fixnum.Int64? attemptsMax,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = create();
    if (nextCheck != null) result.nextCheck = nextCheck;
    if (id != null) result.id.addAll(id);
    if (attemptsMin != null) result.attemptsMin = attemptsMin;
    if (attemptsMax != null) result.attemptsMax = attemptsMax;
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    return result;
  }

  DiscoverySearchRequest._();

  factory DiscoverySearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DiscoverySearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DiscoverySearchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: create)
    ..aOM<$0.DateRange>(1, _omitFieldNames ? '' : 'next_check',
        subBuilder: $0.DateRange.create)
    ..pPS(2, _omitFieldNames ? '' : 'id')
    ..a<$fixnum.Int64>(
        3, _omitFieldNames ? '' : 'attempts_min', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        4, _omitFieldNames ? '' : 'attempts_max', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        1000, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        1001, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoverySearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoverySearchRequest copyWith(
          void Function(DiscoverySearchRequest) updates) =>
      super.copyWith((message) => updates(message as DiscoverySearchRequest))
          as DiscoverySearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DiscoverySearchRequest create() => DiscoverySearchRequest._();
  @$core.override
  DiscoverySearchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DiscoverySearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DiscoverySearchRequest>(create);
  static DiscoverySearchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $0.DateRange get nextCheck => $_getN(0);
  @$pb.TagNumber(1)
  set nextCheck($0.DateRange value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNextCheck() => $_has(0);
  @$pb.TagNumber(1)
  void clearNextCheck() => $_clearField(1);
  @$pb.TagNumber(1)
  $0.DateRange ensureNextCheck() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<$core.String> get id => $_getList(1);

  @$pb.TagNumber(3)
  $fixnum.Int64 get attemptsMin => $_getI64(2);
  @$pb.TagNumber(3)
  set attemptsMin($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasAttemptsMin() => $_has(2);
  @$pb.TagNumber(3)
  void clearAttemptsMin() => $_clearField(3);

  @$pb.TagNumber(4)
  $fixnum.Int64 get attemptsMax => $_getI64(3);
  @$pb.TagNumber(4)
  set attemptsMax($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(4)
  $core.bool hasAttemptsMax() => $_has(3);
  @$pb.TagNumber(4)
  void clearAttemptsMax() => $_clearField(4);

  @$pb.TagNumber(1000)
  $fixnum.Int64 get offset => $_getI64(4);
  @$pb.TagNumber(1000)
  set offset($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(1000)
  $core.bool hasOffset() => $_has(4);
  @$pb.TagNumber(1000)
  void clearOffset() => $_clearField(1000);

  @$pb.TagNumber(1001)
  $fixnum.Int64 get limit => $_getI64(5);
  @$pb.TagNumber(1001)
  set limit($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(1001)
  $core.bool hasLimit() => $_has(5);
  @$pb.TagNumber(1001)
  void clearLimit() => $_clearField(1001);
}

class DiscoverySearchResponse extends $pb.GeneratedMessage {
  factory DiscoverySearchResponse({
    DiscoverySearchRequest? next,
    $core.Iterable<Discovery>? items,
  }) {
    final result = create();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  DiscoverySearchResponse._();

  factory DiscoverySearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DiscoverySearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DiscoverySearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: create)
    ..aOM<DiscoverySearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: DiscoverySearchRequest.create)
    ..pPM<Discovery>(2, _omitFieldNames ? '' : 'items',
        subBuilder: Discovery.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoverySearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoverySearchResponse copyWith(
          void Function(DiscoverySearchResponse) updates) =>
      super.copyWith((message) => updates(message as DiscoverySearchResponse))
          as DiscoverySearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DiscoverySearchResponse create() => DiscoverySearchResponse._();
  @$core.override
  DiscoverySearchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DiscoverySearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DiscoverySearchResponse>(create);
  static DiscoverySearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  DiscoverySearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(DiscoverySearchRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  DiscoverySearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<Discovery> get items => $_getList(1);
}

class DiscoveryCreateRequest extends $pb.GeneratedMessage {
  factory DiscoveryCreateRequest({
    Discovery? discovery,
  }) {
    final result = create();
    if (discovery != null) result.discovery = discovery;
    return result;
  }

  DiscoveryCreateRequest._();

  factory DiscoveryCreateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DiscoveryCreateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DiscoveryCreateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: create)
    ..aOM<Discovery>(1, _omitFieldNames ? '' : 'discovery',
        subBuilder: Discovery.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoveryCreateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoveryCreateRequest copyWith(
          void Function(DiscoveryCreateRequest) updates) =>
      super.copyWith((message) => updates(message as DiscoveryCreateRequest))
          as DiscoveryCreateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DiscoveryCreateRequest create() => DiscoveryCreateRequest._();
  @$core.override
  DiscoveryCreateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DiscoveryCreateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DiscoveryCreateRequest>(create);
  static DiscoveryCreateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Discovery get discovery => $_getN(0);
  @$pb.TagNumber(1)
  set discovery(Discovery value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDiscovery() => $_has(0);
  @$pb.TagNumber(1)
  void clearDiscovery() => $_clearField(1);
  @$pb.TagNumber(1)
  Discovery ensureDiscovery() => $_ensure(0);
}

class DiscoveryCreateResponse extends $pb.GeneratedMessage {
  factory DiscoveryCreateResponse({
    Discovery? discovery,
  }) {
    final result = create();
    if (discovery != null) result.discovery = discovery;
    return result;
  }

  DiscoveryCreateResponse._();

  factory DiscoveryCreateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DiscoveryCreateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DiscoveryCreateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: create)
    ..aOM<Discovery>(1, _omitFieldNames ? '' : 'discovery',
        subBuilder: Discovery.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoveryCreateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoveryCreateResponse copyWith(
          void Function(DiscoveryCreateResponse) updates) =>
      super.copyWith((message) => updates(message as DiscoveryCreateResponse))
          as DiscoveryCreateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DiscoveryCreateResponse create() => DiscoveryCreateResponse._();
  @$core.override
  DiscoveryCreateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DiscoveryCreateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DiscoveryCreateResponse>(create);
  static DiscoveryCreateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Discovery get discovery => $_getN(0);
  @$pb.TagNumber(1)
  set discovery(Discovery value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDiscovery() => $_has(0);
  @$pb.TagNumber(1)
  void clearDiscovery() => $_clearField(1);
  @$pb.TagNumber(1)
  Discovery ensureDiscovery() => $_ensure(0);
}

class DiscoveryDownloadRequest extends $pb.GeneratedMessage {
  factory DiscoveryDownloadRequest() => create();

  DiscoveryDownloadRequest._();

  factory DiscoveryDownloadRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DiscoveryDownloadRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DiscoveryDownloadRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoveryDownloadRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoveryDownloadRequest copyWith(
          void Function(DiscoveryDownloadRequest) updates) =>
      super.copyWith((message) => updates(message as DiscoveryDownloadRequest))
          as DiscoveryDownloadRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DiscoveryDownloadRequest create() => DiscoveryDownloadRequest._();
  @$core.override
  DiscoveryDownloadRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DiscoveryDownloadRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DiscoveryDownloadRequest>(create);
  static DiscoveryDownloadRequest? _defaultInstance;
}

class DiscoveryDownloadResponse extends $pb.GeneratedMessage {
  factory DiscoveryDownloadResponse({
    Discovery? discovery,
  }) {
    final result = create();
    if (discovery != null) result.discovery = discovery;
    return result;
  }

  DiscoveryDownloadResponse._();

  factory DiscoveryDownloadResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DiscoveryDownloadResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DiscoveryDownloadResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: create)
    ..aOM<Discovery>(1, _omitFieldNames ? '' : 'discovery',
        subBuilder: Discovery.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoveryDownloadResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoveryDownloadResponse copyWith(
          void Function(DiscoveryDownloadResponse) updates) =>
      super.copyWith((message) => updates(message as DiscoveryDownloadResponse))
          as DiscoveryDownloadResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DiscoveryDownloadResponse create() => DiscoveryDownloadResponse._();
  @$core.override
  DiscoveryDownloadResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DiscoveryDownloadResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DiscoveryDownloadResponse>(create);
  static DiscoveryDownloadResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Discovery get discovery => $_getN(0);
  @$pb.TagNumber(1)
  set discovery(Discovery value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDiscovery() => $_has(0);
  @$pb.TagNumber(1)
  void clearDiscovery() => $_clearField(1);
  @$pb.TagNumber(1)
  Discovery ensureDiscovery() => $_ensure(0);
}

class DiscoveryDeleteRequest extends $pb.GeneratedMessage {
  factory DiscoveryDeleteRequest() => create();

  DiscoveryDeleteRequest._();

  factory DiscoveryDeleteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DiscoveryDeleteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DiscoveryDeleteRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoveryDeleteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoveryDeleteRequest copyWith(
          void Function(DiscoveryDeleteRequest) updates) =>
      super.copyWith((message) => updates(message as DiscoveryDeleteRequest))
          as DiscoveryDeleteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DiscoveryDeleteRequest create() => DiscoveryDeleteRequest._();
  @$core.override
  DiscoveryDeleteRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DiscoveryDeleteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DiscoveryDeleteRequest>(create);
  static DiscoveryDeleteRequest? _defaultInstance;
}

class DiscoveryDeleteResponse extends $pb.GeneratedMessage {
  factory DiscoveryDeleteResponse({
    Discovery? discovery,
  }) {
    final result = create();
    if (discovery != null) result.discovery = discovery;
    return result;
  }

  DiscoveryDeleteResponse._();

  factory DiscoveryDeleteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DiscoveryDeleteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DiscoveryDeleteResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: create)
    ..aOM<Discovery>(1, _omitFieldNames ? '' : 'discovery',
        subBuilder: Discovery.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoveryDeleteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoveryDeleteResponse copyWith(
          void Function(DiscoveryDeleteResponse) updates) =>
      super.copyWith((message) => updates(message as DiscoveryDeleteResponse))
          as DiscoveryDeleteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DiscoveryDeleteResponse create() => DiscoveryDeleteResponse._();
  @$core.override
  DiscoveryDeleteResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DiscoveryDeleteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DiscoveryDeleteResponse>(create);
  static DiscoveryDeleteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Discovery get discovery => $_getN(0);
  @$pb.TagNumber(1)
  set discovery(Discovery value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDiscovery() => $_has(0);
  @$pb.TagNumber(1)
  void clearDiscovery() => $_clearField(1);
  @$pb.TagNumber(1)
  Discovery ensureDiscovery() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
