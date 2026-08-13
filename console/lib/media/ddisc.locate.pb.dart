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
  }) {
    final result = create();
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
    return result;
  }

  Locate._();

  factory Locate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Locate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Locate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: create)
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
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Locate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Locate copyWith(void Function(Locate) updates) =>
      super.copyWith((message) => updates(message as Locate)) as Locate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Locate create() => Locate._();
  @$core.override
  Locate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Locate getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Locate>(create);
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
}

class LocateSearchRequest extends $pb.GeneratedMessage {
  factory LocateSearchRequest({
    $core.String? query,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = create();
    if (query != null) result.query = query;
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    return result;
  }

  LocateSearchRequest._();

  factory LocateSearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LocateSearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocateSearchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
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
  static LocateSearchRequest create() => LocateSearchRequest._();
  @$core.override
  LocateSearchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LocateSearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocateSearchRequest>(create);
  static LocateSearchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get query => $_getSZ(0);
  @$pb.TagNumber(1)
  set query($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasQuery() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuery() => $_clearField(1);

  @$pb.TagNumber(900)
  $fixnum.Int64 get offset => $_getI64(1);
  @$pb.TagNumber(900)
  set offset($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(900)
  $core.bool hasOffset() => $_has(1);
  @$pb.TagNumber(900)
  void clearOffset() => $_clearField(900);

  @$pb.TagNumber(901)
  $fixnum.Int64 get limit => $_getI64(2);
  @$pb.TagNumber(901)
  set limit($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(901)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(901)
  void clearLimit() => $_clearField(901);
}

class LocateSearchResponse extends $pb.GeneratedMessage {
  factory LocateSearchResponse({
    LocateSearchRequest? next,
    $core.Iterable<Locate>? items,
  }) {
    final result = create();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  LocateSearchResponse._();

  factory LocateSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LocateSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocateSearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: create)
    ..aOM<LocateSearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: LocateSearchRequest.create)
    ..pPM<Locate>(2, _omitFieldNames ? '' : 'items', subBuilder: Locate.create)
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
  static LocateSearchResponse create() => LocateSearchResponse._();
  @$core.override
  LocateSearchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LocateSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocateSearchResponse>(create);
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
  factory LocateLookupRequest() => create();

  LocateLookupRequest._();

  factory LocateLookupRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LocateLookupRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocateLookupRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: create)
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
  static LocateLookupRequest create() => LocateLookupRequest._();
  @$core.override
  LocateLookupRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LocateLookupRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocateLookupRequest>(create);
  static LocateLookupRequest? _defaultInstance;
}

class LocateLookupResponse extends $pb.GeneratedMessage {
  factory LocateLookupResponse({
    Locate? locate,
  }) {
    final result = create();
    if (locate != null) result.locate = locate;
    return result;
  }

  LocateLookupResponse._();

  factory LocateLookupResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LocateLookupResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocateLookupResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: create)
    ..aOM<Locate>(1, _omitFieldNames ? '' : 'locate', subBuilder: Locate.create)
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
  static LocateLookupResponse create() => LocateLookupResponse._();
  @$core.override
  LocateLookupResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LocateLookupResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocateLookupResponse>(create);
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
    final result = create();
    if (locate != null) result.locate = locate;
    return result;
  }

  LocateCreateRequest._();

  factory LocateCreateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LocateCreateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocateCreateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: create)
    ..aOM<Locate>(1, _omitFieldNames ? '' : 'locate', subBuilder: Locate.create)
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
  static LocateCreateRequest create() => LocateCreateRequest._();
  @$core.override
  LocateCreateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LocateCreateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocateCreateRequest>(create);
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
    final result = create();
    if (locate != null) result.locate = locate;
    return result;
  }

  LocateCreateResponse._();

  factory LocateCreateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LocateCreateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocateCreateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'ddisc'),
      createEmptyInstance: create)
    ..aOM<Locate>(1, _omitFieldNames ? '' : 'locate', subBuilder: Locate.create)
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
  static LocateCreateResponse create() => LocateCreateResponse._();
  @$core.override
  LocateCreateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LocateCreateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocateCreateResponse>(create);
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

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
