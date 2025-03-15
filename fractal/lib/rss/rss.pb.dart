//
//  Generated code. Do not modify.
//  source: rss.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

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
  }) {
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    if (createdAt != null) {
      $result.createdAt = createdAt;
    }
    if (updatedAt != null) {
      $result.updatedAt = updatedAt;
    }
    if (nextCheck != null) {
      $result.nextCheck = nextCheck;
    }
    if (description != null) {
      $result.description = description;
    }
    if (url != null) {
      $result.url = url;
    }
    if (autodownload != null) {
      $result.autodownload = autodownload;
    }
    if (autoarchive != null) {
      $result.autoarchive = autoarchive;
    }
    return $result;
  }
  Feed._() : super();
  factory Feed.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory Feed.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'Feed', package: const $pb.PackageName(_omitMessageNames ? '' : 'rss'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'created_at')
    ..aOS(3, _omitFieldNames ? '' : 'updated_at')
    ..aOS(4, _omitFieldNames ? '' : 'next_check')
    ..aOS(5, _omitFieldNames ? '' : 'description')
    ..aOS(6, _omitFieldNames ? '' : 'url')
    ..aOB(7, _omitFieldNames ? '' : 'autodownload')
    ..aOB(8, _omitFieldNames ? '' : 'autoarchive')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  Feed clone() => Feed()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  Feed copyWith(void Function(Feed) updates) => super.copyWith((message) => updates(message as Feed)) as Feed;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Feed create() => Feed._();
  Feed createEmptyInstance() => create();
  static $pb.PbList<Feed> createRepeated() => $pb.PbList<Feed>();
  @$core.pragma('dart2js:noInline')
  static Feed getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Feed>(create);
  static Feed? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get createdAt => $_getSZ(1);
  @$pb.TagNumber(2)
  set createdAt($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasCreatedAt() => $_has(1);
  @$pb.TagNumber(2)
  void clearCreatedAt() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get updatedAt => $_getSZ(2);
  @$pb.TagNumber(3)
  set updatedAt($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasUpdatedAt() => $_has(2);
  @$pb.TagNumber(3)
  void clearUpdatedAt() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get nextCheck => $_getSZ(3);
  @$pb.TagNumber(4)
  set nextCheck($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasNextCheck() => $_has(3);
  @$pb.TagNumber(4)
  void clearNextCheck() => clearField(4);

  @$pb.TagNumber(5)
  $core.String get description => $_getSZ(4);
  @$pb.TagNumber(5)
  set description($core.String v) { $_setString(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasDescription() => $_has(4);
  @$pb.TagNumber(5)
  void clearDescription() => clearField(5);

  @$pb.TagNumber(6)
  $core.String get url => $_getSZ(5);
  @$pb.TagNumber(6)
  set url($core.String v) { $_setString(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasUrl() => $_has(5);
  @$pb.TagNumber(6)
  void clearUrl() => clearField(6);

  @$pb.TagNumber(7)
  $core.bool get autodownload => $_getBF(6);
  @$pb.TagNumber(7)
  set autodownload($core.bool v) { $_setBool(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasAutodownload() => $_has(6);
  @$pb.TagNumber(7)
  void clearAutodownload() => clearField(7);

  @$pb.TagNumber(8)
  $core.bool get autoarchive => $_getBF(7);
  @$pb.TagNumber(8)
  set autoarchive($core.bool v) { $_setBool(7, v); }
  @$pb.TagNumber(8)
  $core.bool hasAutoarchive() => $_has(7);
  @$pb.TagNumber(8)
  void clearAutoarchive() => clearField(8);
}

class FeedSearchRequest extends $pb.GeneratedMessage {
  factory FeedSearchRequest({
    $core.String? query,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final $result = create();
    if (query != null) {
      $result.query = query;
    }
    if (offset != null) {
      $result.offset = offset;
    }
    if (limit != null) {
      $result.limit = limit;
    }
    return $result;
  }
  FeedSearchRequest._() : super();
  factory FeedSearchRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory FeedSearchRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'FeedSearchRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'rss'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..a<$fixnum.Int64>(900, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(901, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  FeedSearchRequest clone() => FeedSearchRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  FeedSearchRequest copyWith(void Function(FeedSearchRequest) updates) => super.copyWith((message) => updates(message as FeedSearchRequest)) as FeedSearchRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FeedSearchRequest create() => FeedSearchRequest._();
  FeedSearchRequest createEmptyInstance() => create();
  static $pb.PbList<FeedSearchRequest> createRepeated() => $pb.PbList<FeedSearchRequest>();
  @$core.pragma('dart2js:noInline')
  static FeedSearchRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<FeedSearchRequest>(create);
  static FeedSearchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get query => $_getSZ(0);
  @$pb.TagNumber(1)
  set query($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasQuery() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuery() => clearField(1);

  @$pb.TagNumber(900)
  $fixnum.Int64 get offset => $_getI64(1);
  @$pb.TagNumber(900)
  set offset($fixnum.Int64 v) { $_setInt64(1, v); }
  @$pb.TagNumber(900)
  $core.bool hasOffset() => $_has(1);
  @$pb.TagNumber(900)
  void clearOffset() => clearField(900);

  @$pb.TagNumber(901)
  $fixnum.Int64 get limit => $_getI64(2);
  @$pb.TagNumber(901)
  set limit($fixnum.Int64 v) { $_setInt64(2, v); }
  @$pb.TagNumber(901)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(901)
  void clearLimit() => clearField(901);
}

class FeedSearchResponse extends $pb.GeneratedMessage {
  factory FeedSearchResponse({
    FeedSearchRequest? next,
    $core.Iterable<Feed>? items,
  }) {
    final $result = create();
    if (next != null) {
      $result.next = next;
    }
    if (items != null) {
      $result.items.addAll(items);
    }
    return $result;
  }
  FeedSearchResponse._() : super();
  factory FeedSearchResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory FeedSearchResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'FeedSearchResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'rss'), createEmptyInstance: create)
    ..aOM<FeedSearchRequest>(1, _omitFieldNames ? '' : 'next', subBuilder: FeedSearchRequest.create)
    ..pc<Feed>(2, _omitFieldNames ? '' : 'items', $pb.PbFieldType.PM, subBuilder: Feed.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  FeedSearchResponse clone() => FeedSearchResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  FeedSearchResponse copyWith(void Function(FeedSearchResponse) updates) => super.copyWith((message) => updates(message as FeedSearchResponse)) as FeedSearchResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FeedSearchResponse create() => FeedSearchResponse._();
  FeedSearchResponse createEmptyInstance() => create();
  static $pb.PbList<FeedSearchResponse> createRepeated() => $pb.PbList<FeedSearchResponse>();
  @$core.pragma('dart2js:noInline')
  static FeedSearchResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<FeedSearchResponse>(create);
  static FeedSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  FeedSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(FeedSearchRequest v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => clearField(1);
  @$pb.TagNumber(1)
  FeedSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.List<Feed> get items => $_getList(1);
}

class FeedCreateRequest extends $pb.GeneratedMessage {
  factory FeedCreateRequest({
    Feed? feed,
  }) {
    final $result = create();
    if (feed != null) {
      $result.feed = feed;
    }
    return $result;
  }
  FeedCreateRequest._() : super();
  factory FeedCreateRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory FeedCreateRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'FeedCreateRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'rss'), createEmptyInstance: create)
    ..aOM<Feed>(1, _omitFieldNames ? '' : 'feed', subBuilder: Feed.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  FeedCreateRequest clone() => FeedCreateRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  FeedCreateRequest copyWith(void Function(FeedCreateRequest) updates) => super.copyWith((message) => updates(message as FeedCreateRequest)) as FeedCreateRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FeedCreateRequest create() => FeedCreateRequest._();
  FeedCreateRequest createEmptyInstance() => create();
  static $pb.PbList<FeedCreateRequest> createRepeated() => $pb.PbList<FeedCreateRequest>();
  @$core.pragma('dart2js:noInline')
  static FeedCreateRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<FeedCreateRequest>(create);
  static FeedCreateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Feed get feed => $_getN(0);
  @$pb.TagNumber(1)
  set feed(Feed v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasFeed() => $_has(0);
  @$pb.TagNumber(1)
  void clearFeed() => clearField(1);
  @$pb.TagNumber(1)
  Feed ensureFeed() => $_ensure(0);
}

class FeedCreateResponse extends $pb.GeneratedMessage {
  factory FeedCreateResponse({
    Feed? feed,
  }) {
    final $result = create();
    if (feed != null) {
      $result.feed = feed;
    }
    return $result;
  }
  FeedCreateResponse._() : super();
  factory FeedCreateResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory FeedCreateResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'FeedCreateResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'rss'), createEmptyInstance: create)
    ..aOM<Feed>(1, _omitFieldNames ? '' : 'feed', subBuilder: Feed.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  FeedCreateResponse clone() => FeedCreateResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  FeedCreateResponse copyWith(void Function(FeedCreateResponse) updates) => super.copyWith((message) => updates(message as FeedCreateResponse)) as FeedCreateResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FeedCreateResponse create() => FeedCreateResponse._();
  FeedCreateResponse createEmptyInstance() => create();
  static $pb.PbList<FeedCreateResponse> createRepeated() => $pb.PbList<FeedCreateResponse>();
  @$core.pragma('dart2js:noInline')
  static FeedCreateResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<FeedCreateResponse>(create);
  static FeedCreateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Feed get feed => $_getN(0);
  @$pb.TagNumber(1)
  set feed(Feed v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasFeed() => $_has(0);
  @$pb.TagNumber(1)
  void clearFeed() => clearField(1);
  @$pb.TagNumber(1)
  Feed ensureFeed() => $_ensure(0);
}

class FeedUpdateRequest extends $pb.GeneratedMessage {
  factory FeedUpdateRequest({
    Feed? feed,
  }) {
    final $result = create();
    if (feed != null) {
      $result.feed = feed;
    }
    return $result;
  }
  FeedUpdateRequest._() : super();
  factory FeedUpdateRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory FeedUpdateRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'FeedUpdateRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'rss'), createEmptyInstance: create)
    ..aOM<Feed>(1, _omitFieldNames ? '' : 'feed', subBuilder: Feed.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  FeedUpdateRequest clone() => FeedUpdateRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  FeedUpdateRequest copyWith(void Function(FeedUpdateRequest) updates) => super.copyWith((message) => updates(message as FeedUpdateRequest)) as FeedUpdateRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FeedUpdateRequest create() => FeedUpdateRequest._();
  FeedUpdateRequest createEmptyInstance() => create();
  static $pb.PbList<FeedUpdateRequest> createRepeated() => $pb.PbList<FeedUpdateRequest>();
  @$core.pragma('dart2js:noInline')
  static FeedUpdateRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<FeedUpdateRequest>(create);
  static FeedUpdateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Feed get feed => $_getN(0);
  @$pb.TagNumber(1)
  set feed(Feed v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasFeed() => $_has(0);
  @$pb.TagNumber(1)
  void clearFeed() => clearField(1);
  @$pb.TagNumber(1)
  Feed ensureFeed() => $_ensure(0);
}

class FeedUpdateResponse extends $pb.GeneratedMessage {
  factory FeedUpdateResponse({
    Feed? feed,
  }) {
    final $result = create();
    if (feed != null) {
      $result.feed = feed;
    }
    return $result;
  }
  FeedUpdateResponse._() : super();
  factory FeedUpdateResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory FeedUpdateResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'FeedUpdateResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'rss'), createEmptyInstance: create)
    ..aOM<Feed>(1, _omitFieldNames ? '' : 'feed', subBuilder: Feed.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  FeedUpdateResponse clone() => FeedUpdateResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  FeedUpdateResponse copyWith(void Function(FeedUpdateResponse) updates) => super.copyWith((message) => updates(message as FeedUpdateResponse)) as FeedUpdateResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FeedUpdateResponse create() => FeedUpdateResponse._();
  FeedUpdateResponse createEmptyInstance() => create();
  static $pb.PbList<FeedUpdateResponse> createRepeated() => $pb.PbList<FeedUpdateResponse>();
  @$core.pragma('dart2js:noInline')
  static FeedUpdateResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<FeedUpdateResponse>(create);
  static FeedUpdateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Feed get feed => $_getN(0);
  @$pb.TagNumber(1)
  set feed(Feed v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasFeed() => $_has(0);
  @$pb.TagNumber(1)
  void clearFeed() => clearField(1);
  @$pb.TagNumber(1)
  Feed ensureFeed() => $_ensure(0);
}

class FeedDeleteRequest extends $pb.GeneratedMessage {
  factory FeedDeleteRequest() => create();
  FeedDeleteRequest._() : super();
  factory FeedDeleteRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory FeedDeleteRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'FeedDeleteRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'rss'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  FeedDeleteRequest clone() => FeedDeleteRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  FeedDeleteRequest copyWith(void Function(FeedDeleteRequest) updates) => super.copyWith((message) => updates(message as FeedDeleteRequest)) as FeedDeleteRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FeedDeleteRequest create() => FeedDeleteRequest._();
  FeedDeleteRequest createEmptyInstance() => create();
  static $pb.PbList<FeedDeleteRequest> createRepeated() => $pb.PbList<FeedDeleteRequest>();
  @$core.pragma('dart2js:noInline')
  static FeedDeleteRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<FeedDeleteRequest>(create);
  static FeedDeleteRequest? _defaultInstance;
}

class FeedDeleteResponse extends $pb.GeneratedMessage {
  factory FeedDeleteResponse({
    Feed? feed,
  }) {
    final $result = create();
    if (feed != null) {
      $result.feed = feed;
    }
    return $result;
  }
  FeedDeleteResponse._() : super();
  factory FeedDeleteResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory FeedDeleteResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'FeedDeleteResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'rss'), createEmptyInstance: create)
    ..aOM<Feed>(1, _omitFieldNames ? '' : 'feed', subBuilder: Feed.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  FeedDeleteResponse clone() => FeedDeleteResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  FeedDeleteResponse copyWith(void Function(FeedDeleteResponse) updates) => super.copyWith((message) => updates(message as FeedDeleteResponse)) as FeedDeleteResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FeedDeleteResponse create() => FeedDeleteResponse._();
  FeedDeleteResponse createEmptyInstance() => create();
  static $pb.PbList<FeedDeleteResponse> createRepeated() => $pb.PbList<FeedDeleteResponse>();
  @$core.pragma('dart2js:noInline')
  static FeedDeleteResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<FeedDeleteResponse>(create);
  static FeedDeleteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Feed get feed => $_getN(0);
  @$pb.TagNumber(1)
  set feed(Feed v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasFeed() => $_has(0);
  @$pb.TagNumber(1)
  void clearFeed() => clearField(1);
  @$pb.TagNumber(1)
  Feed ensureFeed() => $_ensure(0);
}


const _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
