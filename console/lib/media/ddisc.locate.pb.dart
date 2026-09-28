// This is a generated file - do not edit.
//
// Generated from media/ddisc.locate.proto.

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

class Locate extends $pb.GeneratedMessage {
  factory Locate({
    $core.String? id,
    $core.String? createdAt,
    $core.String? updatedAt,
    $core.String? knownMediaId,
    $core.String? locatedTorrentId,
    $core.String? query,
    $core.String? mimetype,
    $core.String? tombstonedAt,
    $core.bool? autodownload,
    $core.bool? adult,
    $core.int? attempts,
    $core.String? nextCheckAt,
  }) {
    final result = Locate._();
    if (id != null) result.id = id;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (knownMediaId != null) result.knownMediaId = knownMediaId;
    if (locatedTorrentId != null) result.locatedTorrentId = locatedTorrentId;
    if (query != null) result.query = query;
    if (mimetype != null) result.mimetype = mimetype;
    if (tombstonedAt != null) result.tombstonedAt = tombstonedAt;
    if (autodownload != null) result.autodownload = autodownload;
    if (adult != null) result.adult = adult;
    if (attempts != null) result.attempts = attempts;
    if (nextCheckAt != null) result.nextCheckAt = nextCheckAt;
    return result;
  }

  Locate._();

  factory Locate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Locate()..mergeFromBuffer(data, registry);
  factory Locate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Locate()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Locate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: Locate.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'created_at')
    ..aOS(3, _omitFieldNames ? '' : 'updated_at')
    ..aOS(4, _omitFieldNames ? '' : 'known_media_id')
    ..aOS(5, _omitFieldNames ? '' : 'located_torrent_id')
    ..aOS(6, _omitFieldNames ? '' : 'query')
    ..aOS(7, _omitFieldNames ? '' : 'mimetype')
    ..aOS(8, _omitFieldNames ? '' : 'tombstoned_at')
    ..aOB(9, _omitFieldNames ? '' : 'autodownload')
    ..aOB(10, _omitFieldNames ? '' : 'adult')
    ..aI(11, _omitFieldNames ? '' : 'attempts', fieldType: $pb.PbFieldType.OU3)
    ..aOS(12, _omitFieldNames ? '' : 'next_check_at')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Locate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Locate copyWith(void Function(Locate) updates) =>
      super.copyWith((message) => updates(message as Locate)) as Locate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use Locate() / Locate.new instead')
  static Locate create() => Locate._();
  static $pb.GeneratedMessage $_createMessage() => Locate._();
  @$core.override
  Locate createEmptyInstance() => Locate._();
  @$core.pragma('dart2js:noInline')
  static Locate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Locate>(Locate.$_createMessage);
  static Locate? _defaultInstance;

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
  $core.String get knownMediaId => $_getSZ(3);
  @$pb.TagNumber(4)
  set knownMediaId($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasKnownMediaId() => $_has(3);
  @$pb.TagNumber(4)
  void clearKnownMediaId() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get locatedTorrentId => $_getSZ(4);
  @$pb.TagNumber(5)
  set locatedTorrentId($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasLocatedTorrentId() => $_has(4);
  @$pb.TagNumber(5)
  void clearLocatedTorrentId() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get query => $_getSZ(5);
  @$pb.TagNumber(6)
  set query($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasQuery() => $_has(5);
  @$pb.TagNumber(6)
  void clearQuery() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.String get mimetype => $_getSZ(6);
  @$pb.TagNumber(7)
  set mimetype($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasMimetype() => $_has(6);
  @$pb.TagNumber(7)
  void clearMimetype() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.String get tombstonedAt => $_getSZ(7);
  @$pb.TagNumber(8)
  set tombstonedAt($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasTombstonedAt() => $_has(7);
  @$pb.TagNumber(8)
  void clearTombstonedAt() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.bool get autodownload => $_getBF(8);
  @$pb.TagNumber(9)
  set autodownload($core.bool value) => $_setBool(8, value);
  @$pb.TagNumber(9)
  $core.bool hasAutodownload() => $_has(8);
  @$pb.TagNumber(9)
  void clearAutodownload() => $_clearField(9);

  @$pb.TagNumber(10)
  $core.bool get adult => $_getBF(9);
  @$pb.TagNumber(10)
  set adult($core.bool value) => $_setBool(9, value);
  @$pb.TagNumber(10)
  $core.bool hasAdult() => $_has(9);
  @$pb.TagNumber(10)
  void clearAdult() => $_clearField(10);

  @$pb.TagNumber(11)
  $core.int get attempts => $_getIZ(10);
  @$pb.TagNumber(11)
  set attempts($core.int value) => $_setUnsignedInt32(10, value);
  @$pb.TagNumber(11)
  $core.bool hasAttempts() => $_has(10);
  @$pb.TagNumber(11)
  void clearAttempts() => $_clearField(11);

  @$pb.TagNumber(12)
  $core.String get nextCheckAt => $_getSZ(11);
  @$pb.TagNumber(12)
  set nextCheckAt($core.String value) => $_setString(11, value);
  @$pb.TagNumber(12)
  $core.bool hasNextCheckAt() => $_has(11);
  @$pb.TagNumber(12)
  void clearNextCheckAt() => $_clearField(12);
}

class LocateSearchRequest extends $pb.GeneratedMessage {
  factory LocateSearchRequest({
    $core.String? query,
    $core.Iterable<$core.String>? id,
    $fixnum.Int64? attemptsMin,
    $fixnum.Int64? attemptsMax,
    $core.bool? pending,
    $core.bool? completed,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = LocateSearchRequest._();
    if (query != null) result.query = query;
    if (id != null) result.id.addAll(id);
    if (attemptsMin != null) result.attemptsMin = attemptsMin;
    if (attemptsMax != null) result.attemptsMax = attemptsMax;
    if (pending != null) result.pending = pending;
    if (completed != null) result.completed = completed;
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    return result;
  }

  LocateSearchRequest._();

  factory LocateSearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocateSearchRequest()..mergeFromBuffer(data, registry);
  factory LocateSearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocateSearchRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocateSearchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: LocateSearchRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..pPS(2, _omitFieldNames ? '' : 'id')
    ..a<$fixnum.Int64>(
        3, _omitFieldNames ? '' : 'attempts_min', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        4, _omitFieldNames ? '' : 'attempts_max', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOB(5, _omitFieldNames ? '' : 'pending')
    ..aOB(6, _omitFieldNames ? '' : 'completed')
    ..a<$fixnum.Int64>(
        900, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(901, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocateSearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocateSearchRequest copyWith(void Function(LocateSearchRequest) updates) =>
      super.copyWith((message) => updates(message as LocateSearchRequest))
          as LocateSearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use LocateSearchRequest() / LocateSearchRequest.new instead')
  static LocateSearchRequest create() => LocateSearchRequest._();
  static $pb.GeneratedMessage $_createMessage() => LocateSearchRequest._();
  @$core.override
  LocateSearchRequest createEmptyInstance() => LocateSearchRequest._();
  @$core.pragma('dart2js:noInline')
  static LocateSearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocateSearchRequest>(
          LocateSearchRequest.$_createMessage);
  static LocateSearchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get query => $_getSZ(0);
  @$pb.TagNumber(1)
  set query($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasQuery() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuery() => $_clearField(1);

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

  @$pb.TagNumber(5)
  $core.bool get pending => $_getBF(4);
  @$pb.TagNumber(5)
  set pending($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(5)
  $core.bool hasPending() => $_has(4);
  @$pb.TagNumber(5)
  void clearPending() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.bool get completed => $_getBF(5);
  @$pb.TagNumber(6)
  set completed($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(6)
  $core.bool hasCompleted() => $_has(5);
  @$pb.TagNumber(6)
  void clearCompleted() => $_clearField(6);

  @$pb.TagNumber(900)
  $fixnum.Int64 get offset => $_getI64(6);
  @$pb.TagNumber(900)
  set offset($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(900)
  $core.bool hasOffset() => $_has(6);
  @$pb.TagNumber(900)
  void clearOffset() => $_clearField(900);

  @$pb.TagNumber(901)
  $fixnum.Int64 get limit => $_getI64(7);
  @$pb.TagNumber(901)
  set limit($fixnum.Int64 value) => $_setInt64(7, value);
  @$pb.TagNumber(901)
  $core.bool hasLimit() => $_has(7);
  @$pb.TagNumber(901)
  void clearLimit() => $_clearField(901);
}

class LocateSearchResponse extends $pb.GeneratedMessage {
  factory LocateSearchResponse({
    LocateSearchRequest? next,
    $core.Iterable<Locate>? items,
  }) {
    final result = LocateSearchResponse._();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  LocateSearchResponse._();

  factory LocateSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocateSearchResponse()..mergeFromBuffer(data, registry);
  factory LocateSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocateSearchResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocateSearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: LocateSearchResponse.$_createMessage)
    ..aOM<LocateSearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: LocateSearchRequest.$_createMessage)
    ..pPM<Locate>(2, _omitFieldNames ? '' : 'items',
        subBuilder: Locate.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocateSearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocateSearchResponse copyWith(void Function(LocateSearchResponse) updates) =>
      super.copyWith((message) => updates(message as LocateSearchResponse))
          as LocateSearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use LocateSearchResponse() / LocateSearchResponse.new instead')
  static LocateSearchResponse create() => LocateSearchResponse._();
  static $pb.GeneratedMessage $_createMessage() => LocateSearchResponse._();
  @$core.override
  LocateSearchResponse createEmptyInstance() => LocateSearchResponse._();
  @$core.pragma('dart2js:noInline')
  static LocateSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocateSearchResponse>(
          LocateSearchResponse.$_createMessage);
  static LocateSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  LocateSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(LocateSearchRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  LocateSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<Locate> get items => $_getList(1);
}

class LocateLookupRequest extends $pb.GeneratedMessage {
  factory LocateLookupRequest() => LocateLookupRequest._();

  LocateLookupRequest._();

  factory LocateLookupRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocateLookupRequest()..mergeFromBuffer(data, registry);
  factory LocateLookupRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocateLookupRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocateLookupRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: LocateLookupRequest.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocateLookupRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocateLookupRequest copyWith(void Function(LocateLookupRequest) updates) =>
      super.copyWith((message) => updates(message as LocateLookupRequest))
          as LocateLookupRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use LocateLookupRequest() / LocateLookupRequest.new instead')
  static LocateLookupRequest create() => LocateLookupRequest._();
  static $pb.GeneratedMessage $_createMessage() => LocateLookupRequest._();
  @$core.override
  LocateLookupRequest createEmptyInstance() => LocateLookupRequest._();
  @$core.pragma('dart2js:noInline')
  static LocateLookupRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocateLookupRequest>(
          LocateLookupRequest.$_createMessage);
  static LocateLookupRequest? _defaultInstance;
}

class LocateLookupResponse extends $pb.GeneratedMessage {
  factory LocateLookupResponse({
    Locate? locate,
  }) {
    final result = LocateLookupResponse._();
    if (locate != null) result.locate = locate;
    return result;
  }

  LocateLookupResponse._();

  factory LocateLookupResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocateLookupResponse()..mergeFromBuffer(data, registry);
  factory LocateLookupResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocateLookupResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocateLookupResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: LocateLookupResponse.$_createMessage)
    ..aOM<Locate>(1, _omitFieldNames ? '' : 'locate',
        subBuilder: Locate.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocateLookupResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocateLookupResponse copyWith(void Function(LocateLookupResponse) updates) =>
      super.copyWith((message) => updates(message as LocateLookupResponse))
          as LocateLookupResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use LocateLookupResponse() / LocateLookupResponse.new instead')
  static LocateLookupResponse create() => LocateLookupResponse._();
  static $pb.GeneratedMessage $_createMessage() => LocateLookupResponse._();
  @$core.override
  LocateLookupResponse createEmptyInstance() => LocateLookupResponse._();
  @$core.pragma('dart2js:noInline')
  static LocateLookupResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocateLookupResponse>(
          LocateLookupResponse.$_createMessage);
  static LocateLookupResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Locate get locate => $_getN(0);
  @$pb.TagNumber(1)
  set locate(Locate value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasLocate() => $_has(0);
  @$pb.TagNumber(1)
  void clearLocate() => $_clearField(1);
  @$pb.TagNumber(1)
  Locate ensureLocate() => $_ensure(0);
}

class LocateCreateRequest extends $pb.GeneratedMessage {
  factory LocateCreateRequest({
    Locate? locate,
  }) {
    final result = LocateCreateRequest._();
    if (locate != null) result.locate = locate;
    return result;
  }

  LocateCreateRequest._();

  factory LocateCreateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocateCreateRequest()..mergeFromBuffer(data, registry);
  factory LocateCreateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocateCreateRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocateCreateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: LocateCreateRequest.$_createMessage)
    ..aOM<Locate>(1, _omitFieldNames ? '' : 'locate',
        subBuilder: Locate.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocateCreateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocateCreateRequest copyWith(void Function(LocateCreateRequest) updates) =>
      super.copyWith((message) => updates(message as LocateCreateRequest))
          as LocateCreateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use LocateCreateRequest() / LocateCreateRequest.new instead')
  static LocateCreateRequest create() => LocateCreateRequest._();
  static $pb.GeneratedMessage $_createMessage() => LocateCreateRequest._();
  @$core.override
  LocateCreateRequest createEmptyInstance() => LocateCreateRequest._();
  @$core.pragma('dart2js:noInline')
  static LocateCreateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocateCreateRequest>(
          LocateCreateRequest.$_createMessage);
  static LocateCreateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Locate get locate => $_getN(0);
  @$pb.TagNumber(1)
  set locate(Locate value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasLocate() => $_has(0);
  @$pb.TagNumber(1)
  void clearLocate() => $_clearField(1);
  @$pb.TagNumber(1)
  Locate ensureLocate() => $_ensure(0);
}

class LocateCreateResponse extends $pb.GeneratedMessage {
  factory LocateCreateResponse({
    Locate? locate,
  }) {
    final result = LocateCreateResponse._();
    if (locate != null) result.locate = locate;
    return result;
  }

  LocateCreateResponse._();

  factory LocateCreateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocateCreateResponse()..mergeFromBuffer(data, registry);
  factory LocateCreateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocateCreateResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocateCreateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: LocateCreateResponse.$_createMessage)
    ..aOM<Locate>(1, _omitFieldNames ? '' : 'locate',
        subBuilder: Locate.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocateCreateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocateCreateResponse copyWith(void Function(LocateCreateResponse) updates) =>
      super.copyWith((message) => updates(message as LocateCreateResponse))
          as LocateCreateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use LocateCreateResponse() / LocateCreateResponse.new instead')
  static LocateCreateResponse create() => LocateCreateResponse._();
  static $pb.GeneratedMessage $_createMessage() => LocateCreateResponse._();
  @$core.override
  LocateCreateResponse createEmptyInstance() => LocateCreateResponse._();
  @$core.pragma('dart2js:noInline')
  static LocateCreateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocateCreateResponse>(
          LocateCreateResponse.$_createMessage);
  static LocateCreateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Locate get locate => $_getN(0);
  @$pb.TagNumber(1)
  set locate(Locate value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasLocate() => $_has(0);
  @$pb.TagNumber(1)
  void clearLocate() => $_clearField(1);
  @$pb.TagNumber(1)
  Locate ensureLocate() => $_ensure(0);
}

class LocateDeleteResponse extends $pb.GeneratedMessage {
  factory LocateDeleteResponse({
    Locate? locate,
  }) {
    final result = LocateDeleteResponse._();
    if (locate != null) result.locate = locate;
    return result;
  }

  LocateDeleteResponse._();

  factory LocateDeleteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocateDeleteResponse()..mergeFromBuffer(data, registry);
  factory LocateDeleteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocateDeleteResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocateDeleteResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: LocateDeleteResponse.$_createMessage)
    ..aOM<Locate>(1, _omitFieldNames ? '' : 'locate',
        subBuilder: Locate.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocateDeleteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocateDeleteResponse copyWith(void Function(LocateDeleteResponse) updates) =>
      super.copyWith((message) => updates(message as LocateDeleteResponse))
          as LocateDeleteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use LocateDeleteResponse() / LocateDeleteResponse.new instead')
  static LocateDeleteResponse create() => LocateDeleteResponse._();
  static $pb.GeneratedMessage $_createMessage() => LocateDeleteResponse._();
  @$core.override
  LocateDeleteResponse createEmptyInstance() => LocateDeleteResponse._();
  @$core.pragma('dart2js:noInline')
  static LocateDeleteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocateDeleteResponse>(
          LocateDeleteResponse.$_createMessage);
  static LocateDeleteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Locate get locate => $_getN(0);
  @$pb.TagNumber(1)
  set locate(Locate value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasLocate() => $_has(0);
  @$pb.TagNumber(1)
  void clearLocate() => $_clearField(1);
  @$pb.TagNumber(1)
  Locate ensureLocate() => $_ensure(0);
}

class LocateRetryRequest extends $pb.GeneratedMessage {
  factory LocateRetryRequest({
    $core.bool? resetAttempts,
  }) {
    final result = LocateRetryRequest._();
    if (resetAttempts != null) result.resetAttempts = resetAttempts;
    return result;
  }

  LocateRetryRequest._();

  factory LocateRetryRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocateRetryRequest()..mergeFromBuffer(data, registry);
  factory LocateRetryRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocateRetryRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocateRetryRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: LocateRetryRequest.$_createMessage)
    ..aOB(1, _omitFieldNames ? '' : 'reset_attempts')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocateRetryRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocateRetryRequest copyWith(void Function(LocateRetryRequest) updates) =>
      super.copyWith((message) => updates(message as LocateRetryRequest))
          as LocateRetryRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use LocateRetryRequest() / LocateRetryRequest.new instead')
  static LocateRetryRequest create() => LocateRetryRequest._();
  static $pb.GeneratedMessage $_createMessage() => LocateRetryRequest._();
  @$core.override
  LocateRetryRequest createEmptyInstance() => LocateRetryRequest._();
  @$core.pragma('dart2js:noInline')
  static LocateRetryRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocateRetryRequest>(
          LocateRetryRequest.$_createMessage);
  static LocateRetryRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get resetAttempts => $_getBF(0);
  @$pb.TagNumber(1)
  set resetAttempts($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasResetAttempts() => $_has(0);
  @$pb.TagNumber(1)
  void clearResetAttempts() => $_clearField(1);
}

class LocateRetryResponse extends $pb.GeneratedMessage {
  factory LocateRetryResponse({
    Locate? locate,
  }) {
    final result = LocateRetryResponse._();
    if (locate != null) result.locate = locate;
    return result;
  }

  LocateRetryResponse._();

  factory LocateRetryResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocateRetryResponse()..mergeFromBuffer(data, registry);
  factory LocateRetryResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocateRetryResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocateRetryResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: LocateRetryResponse.$_createMessage)
    ..aOM<Locate>(1, _omitFieldNames ? '' : 'locate',
        subBuilder: Locate.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocateRetryResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocateRetryResponse copyWith(void Function(LocateRetryResponse) updates) =>
      super.copyWith((message) => updates(message as LocateRetryResponse))
          as LocateRetryResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use LocateRetryResponse() / LocateRetryResponse.new instead')
  static LocateRetryResponse create() => LocateRetryResponse._();
  static $pb.GeneratedMessage $_createMessage() => LocateRetryResponse._();
  @$core.override
  LocateRetryResponse createEmptyInstance() => LocateRetryResponse._();
  @$core.pragma('dart2js:noInline')
  static LocateRetryResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocateRetryResponse>(
          LocateRetryResponse.$_createMessage);
  static LocateRetryResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Locate get locate => $_getN(0);
  @$pb.TagNumber(1)
  set locate(Locate value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasLocate() => $_has(0);
  @$pb.TagNumber(1)
  void clearLocate() => $_clearField(1);
  @$pb.TagNumber(1)
  Locate ensureLocate() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
