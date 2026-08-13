// This is a generated file - do not edit.
//
// Generated from rss/rss.proto.

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

class Feed extends $pb.GeneratedMessage {
  factory Feed({
    $core.String? id,
    $core.String? createdAt,
    $core.String? updatedAt,
    $core.String? nextCheck,
    $core.String? description,
    $core.String? url,
    $core.bool? autodownload,
    $core.bool? autoarchive,
    $core.bool? contributing,
    $core.String? disabledAt,
    $fixnum.Int64? ttlMinimum,
    $core.String? encryptionSeed,
    $core.String? digest,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (nextCheck != null) result.nextCheck = nextCheck;
    if (description != null) result.description = description;
    if (url != null) result.url = url;
    if (autodownload != null) result.autodownload = autodownload;
    if (autoarchive != null) result.autoarchive = autoarchive;
    if (contributing != null) result.contributing = contributing;
    if (disabledAt != null) result.disabledAt = disabledAt;
    if (ttlMinimum != null) result.ttlMinimum = ttlMinimum;
    if (encryptionSeed != null) result.encryptionSeed = encryptionSeed;
    if (digest != null) result.digest = digest;
    return result;
  }

  Feed._();

  factory Feed.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Feed.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Feed',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'rss'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'created_at')
    ..aOS(3, _omitFieldNames ? '' : 'updated_at')
    ..aOS(4, _omitFieldNames ? '' : 'next_check')
    ..aOS(5, _omitFieldNames ? '' : 'description')
    ..aOS(6, _omitFieldNames ? '' : 'url')
    ..aOB(7, _omitFieldNames ? '' : 'autodownload')
    ..aOB(8, _omitFieldNames ? '' : 'autoarchive')
    ..aOB(9, _omitFieldNames ? '' : 'contributing')
    ..aOS(10, _omitFieldNames ? '' : 'disabled_at')
    ..a<$fixnum.Int64>(
        11, _omitFieldNames ? '' : 'ttl_minimum', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOS(12, _omitFieldNames ? '' : 'encryption_seed')
    ..aOS(13, _omitFieldNames ? '' : 'digest')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Feed clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Feed copyWith(void Function(Feed) updates) =>
      super.copyWith((message) => updates(message as Feed)) as Feed;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Feed create() => Feed._();
  @$core.override
  Feed createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Feed getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Feed>(create);
  static Feed? _defaultInstance;

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
  $core.String get nextCheck => $_getSZ(3);
  @$pb.TagNumber(4)
  set nextCheck($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasNextCheck() => $_has(3);
  @$pb.TagNumber(4)
  void clearNextCheck() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get description => $_getSZ(4);
  @$pb.TagNumber(5)
  set description($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasDescription() => $_has(4);
  @$pb.TagNumber(5)
  void clearDescription() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get url => $_getSZ(5);
  @$pb.TagNumber(6)
  set url($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasUrl() => $_has(5);
  @$pb.TagNumber(6)
  void clearUrl() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.bool get autodownload => $_getBF(6);
  @$pb.TagNumber(7)
  set autodownload($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(7)
  $core.bool hasAutodownload() => $_has(6);
  @$pb.TagNumber(7)
  void clearAutodownload() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.bool get autoarchive => $_getBF(7);
  @$pb.TagNumber(8)
  set autoarchive($core.bool value) => $_setBool(7, value);
  @$pb.TagNumber(8)
  $core.bool hasAutoarchive() => $_has(7);
  @$pb.TagNumber(8)
  void clearAutoarchive() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.bool get contributing => $_getBF(8);
  @$pb.TagNumber(9)
  set contributing($core.bool value) => $_setBool(8, value);
  @$pb.TagNumber(9)
  $core.bool hasContributing() => $_has(8);
  @$pb.TagNumber(9)
  void clearContributing() => $_clearField(9);

  @$pb.TagNumber(10)
  $core.String get disabledAt => $_getSZ(9);
  @$pb.TagNumber(10)
  set disabledAt($core.String value) => $_setString(9, value);
  @$pb.TagNumber(10)
  $core.bool hasDisabledAt() => $_has(9);
  @$pb.TagNumber(10)
  void clearDisabledAt() => $_clearField(10);

  @$pb.TagNumber(11)
  $fixnum.Int64 get ttlMinimum => $_getI64(10);
  @$pb.TagNumber(11)
  set ttlMinimum($fixnum.Int64 value) => $_setInt64(10, value);
  @$pb.TagNumber(11)
  $core.bool hasTtlMinimum() => $_has(10);
  @$pb.TagNumber(11)
  void clearTtlMinimum() => $_clearField(11);

  @$pb.TagNumber(12)
  $core.String get encryptionSeed => $_getSZ(11);
  @$pb.TagNumber(12)
  set encryptionSeed($core.String value) => $_setString(11, value);
  @$pb.TagNumber(12)
  $core.bool hasEncryptionSeed() => $_has(11);
  @$pb.TagNumber(12)
  void clearEncryptionSeed() => $_clearField(12);

  @$pb.TagNumber(13)
  $core.String get digest => $_getSZ(12);
  @$pb.TagNumber(13)
  set digest($core.String value) => $_setString(12, value);
  @$pb.TagNumber(13)
  $core.bool hasDigest() => $_has(12);
  @$pb.TagNumber(13)
  void clearDigest() => $_clearField(13);
}

class FeedSearchRequest extends $pb.GeneratedMessage {
  factory FeedSearchRequest({
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

  FeedSearchRequest._();

  factory FeedSearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FeedSearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FeedSearchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'rss'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..a<$fixnum.Int64>(
        900, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(901, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FeedSearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FeedSearchRequest copyWith(void Function(FeedSearchRequest) updates) =>
      super.copyWith((message) => updates(message as FeedSearchRequest))
          as FeedSearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FeedSearchRequest create() => FeedSearchRequest._();
  @$core.override
  FeedSearchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FeedSearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FeedSearchRequest>(create);
  static FeedSearchRequest? _defaultInstance;

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

class FeedSearchResponse extends $pb.GeneratedMessage {
  factory FeedSearchResponse({
    FeedSearchRequest? next,
    $core.Iterable<Feed>? items,
  }) {
    final result = create();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  FeedSearchResponse._();

  factory FeedSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FeedSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FeedSearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'rss'),
      createEmptyInstance: create)
    ..aOM<FeedSearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: FeedSearchRequest.create)
    ..pPM<Feed>(2, _omitFieldNames ? '' : 'items', subBuilder: Feed.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FeedSearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FeedSearchResponse copyWith(void Function(FeedSearchResponse) updates) =>
      super.copyWith((message) => updates(message as FeedSearchResponse))
          as FeedSearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FeedSearchResponse create() => FeedSearchResponse._();
  @$core.override
  FeedSearchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FeedSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FeedSearchResponse>(create);
  static FeedSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  FeedSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(FeedSearchRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  FeedSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<Feed> get items => $_getList(1);
}

class FeedCreateRequest extends $pb.GeneratedMessage {
  factory FeedCreateRequest({
    Feed? feed,
  }) {
    final result = create();
    if (feed != null) result.feed = feed;
    return result;
  }

  FeedCreateRequest._();

  factory FeedCreateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FeedCreateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FeedCreateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'rss'),
      createEmptyInstance: create)
    ..aOM<Feed>(1, _omitFieldNames ? '' : 'feed', subBuilder: Feed.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FeedCreateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FeedCreateRequest copyWith(void Function(FeedCreateRequest) updates) =>
      super.copyWith((message) => updates(message as FeedCreateRequest))
          as FeedCreateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FeedCreateRequest create() => FeedCreateRequest._();
  @$core.override
  FeedCreateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FeedCreateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FeedCreateRequest>(create);
  static FeedCreateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Feed get feed => $_getN(0);
  @$pb.TagNumber(1)
  set feed(Feed value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasFeed() => $_has(0);
  @$pb.TagNumber(1)
  void clearFeed() => $_clearField(1);
  @$pb.TagNumber(1)
  Feed ensureFeed() => $_ensure(0);
}

class FeedCreateResponse extends $pb.GeneratedMessage {
  factory FeedCreateResponse({
    Feed? feed,
  }) {
    final result = create();
    if (feed != null) result.feed = feed;
    return result;
  }

  FeedCreateResponse._();

  factory FeedCreateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FeedCreateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FeedCreateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'rss'),
      createEmptyInstance: create)
    ..aOM<Feed>(1, _omitFieldNames ? '' : 'feed', subBuilder: Feed.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FeedCreateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FeedCreateResponse copyWith(void Function(FeedCreateResponse) updates) =>
      super.copyWith((message) => updates(message as FeedCreateResponse))
          as FeedCreateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FeedCreateResponse create() => FeedCreateResponse._();
  @$core.override
  FeedCreateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FeedCreateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FeedCreateResponse>(create);
  static FeedCreateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Feed get feed => $_getN(0);
  @$pb.TagNumber(1)
  set feed(Feed value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasFeed() => $_has(0);
  @$pb.TagNumber(1)
  void clearFeed() => $_clearField(1);
  @$pb.TagNumber(1)
  Feed ensureFeed() => $_ensure(0);
}

class FeedUpdateRequest extends $pb.GeneratedMessage {
  factory FeedUpdateRequest({
    Feed? feed,
  }) {
    final result = create();
    if (feed != null) result.feed = feed;
    return result;
  }

  FeedUpdateRequest._();

  factory FeedUpdateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FeedUpdateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FeedUpdateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'rss'),
      createEmptyInstance: create)
    ..aOM<Feed>(1, _omitFieldNames ? '' : 'feed', subBuilder: Feed.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FeedUpdateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FeedUpdateRequest copyWith(void Function(FeedUpdateRequest) updates) =>
      super.copyWith((message) => updates(message as FeedUpdateRequest))
          as FeedUpdateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FeedUpdateRequest create() => FeedUpdateRequest._();
  @$core.override
  FeedUpdateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FeedUpdateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FeedUpdateRequest>(create);
  static FeedUpdateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Feed get feed => $_getN(0);
  @$pb.TagNumber(1)
  set feed(Feed value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasFeed() => $_has(0);
  @$pb.TagNumber(1)
  void clearFeed() => $_clearField(1);
  @$pb.TagNumber(1)
  Feed ensureFeed() => $_ensure(0);
}

class FeedUpdateResponse extends $pb.GeneratedMessage {
  factory FeedUpdateResponse({
    Feed? feed,
  }) {
    final result = create();
    if (feed != null) result.feed = feed;
    return result;
  }

  FeedUpdateResponse._();

  factory FeedUpdateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FeedUpdateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FeedUpdateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'rss'),
      createEmptyInstance: create)
    ..aOM<Feed>(1, _omitFieldNames ? '' : 'feed', subBuilder: Feed.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FeedUpdateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FeedUpdateResponse copyWith(void Function(FeedUpdateResponse) updates) =>
      super.copyWith((message) => updates(message as FeedUpdateResponse))
          as FeedUpdateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FeedUpdateResponse create() => FeedUpdateResponse._();
  @$core.override
  FeedUpdateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FeedUpdateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FeedUpdateResponse>(create);
  static FeedUpdateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Feed get feed => $_getN(0);
  @$pb.TagNumber(1)
  set feed(Feed value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasFeed() => $_has(0);
  @$pb.TagNumber(1)
  void clearFeed() => $_clearField(1);
  @$pb.TagNumber(1)
  Feed ensureFeed() => $_ensure(0);
}

class FeedDeleteRequest extends $pb.GeneratedMessage {
  factory FeedDeleteRequest() => create();

  FeedDeleteRequest._();

  factory FeedDeleteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FeedDeleteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FeedDeleteRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'rss'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FeedDeleteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FeedDeleteRequest copyWith(void Function(FeedDeleteRequest) updates) =>
      super.copyWith((message) => updates(message as FeedDeleteRequest))
          as FeedDeleteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FeedDeleteRequest create() => FeedDeleteRequest._();
  @$core.override
  FeedDeleteRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FeedDeleteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FeedDeleteRequest>(create);
  static FeedDeleteRequest? _defaultInstance;
}

class FeedDeleteResponse extends $pb.GeneratedMessage {
  factory FeedDeleteResponse({
    Feed? feed,
  }) {
    final result = create();
    if (feed != null) result.feed = feed;
    return result;
  }

  FeedDeleteResponse._();

  factory FeedDeleteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FeedDeleteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FeedDeleteResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'rss'),
      createEmptyInstance: create)
    ..aOM<Feed>(1, _omitFieldNames ? '' : 'feed', subBuilder: Feed.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FeedDeleteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FeedDeleteResponse copyWith(void Function(FeedDeleteResponse) updates) =>
      super.copyWith((message) => updates(message as FeedDeleteResponse))
          as FeedDeleteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FeedDeleteResponse create() => FeedDeleteResponse._();
  @$core.override
  FeedDeleteResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FeedDeleteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FeedDeleteResponse>(create);
  static FeedDeleteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Feed get feed => $_getN(0);
  @$pb.TagNumber(1)
  set feed(Feed value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasFeed() => $_has(0);
  @$pb.TagNumber(1)
  void clearFeed() => $_clearField(1);
  @$pb.TagNumber(1)
  Feed ensureFeed() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
