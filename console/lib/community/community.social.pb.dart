// This is a generated file - do not edit.
//
// Generated from community/community.social.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'community.pb.dart' as $0;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class PluginPublisher extends $pb.GeneratedMessage {
  factory PluginPublisher({
    $core.String? id,
    $core.String? path,
    $core.String? description,
    $core.String? mimetype,
    $core.String? createdAt,
    $core.String? updatedAt,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (path != null) result.path = path;
    if (description != null) result.description = description;
    if (mimetype != null) result.mimetype = mimetype;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    return result;
  }

  PluginPublisher._();

  factory PluginPublisher.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PluginPublisher.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PluginPublisher',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'path')
    ..aOS(3, _omitFieldNames ? '' : 'description')
    ..aOS(4, _omitFieldNames ? '' : 'mimetype')
    ..aOS(5, _omitFieldNames ? '' : 'created_at')
    ..aOS(6, _omitFieldNames ? '' : 'updated_at')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginPublisher clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginPublisher copyWith(void Function(PluginPublisher) updates) =>
      super.copyWith((message) => updates(message as PluginPublisher))
          as PluginPublisher;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PluginPublisher create() => PluginPublisher._();
  @$core.override
  PluginPublisher createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PluginPublisher getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PluginPublisher>(create);
  static PluginPublisher? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get path => $_getSZ(1);
  @$pb.TagNumber(2)
  set path($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPath() => $_has(1);
  @$pb.TagNumber(2)
  void clearPath() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(3)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(3)
  void clearDescription() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get mimetype => $_getSZ(3);
  @$pb.TagNumber(4)
  set mimetype($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasMimetype() => $_has(3);
  @$pb.TagNumber(4)
  void clearMimetype() => $_clearField(4);

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
}

class CommunityPublisher extends $pb.GeneratedMessage {
  factory CommunityPublisher({
    $core.String? id,
    $core.String? communityId,
    $core.String? publisherId,
    $core.String? createdAt,
    $core.String? updatedAt,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (communityId != null) result.communityId = communityId;
    if (publisherId != null) result.publisherId = publisherId;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    return result;
  }

  CommunityPublisher._();

  factory CommunityPublisher.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityPublisher.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityPublisher',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'community_id')
    ..aOS(3, _omitFieldNames ? '' : 'publisher_id')
    ..aOS(4, _omitFieldNames ? '' : 'created_at')
    ..aOS(5, _omitFieldNames ? '' : 'updated_at')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityPublisher clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityPublisher copyWith(void Function(CommunityPublisher) updates) =>
      super.copyWith((message) => updates(message as CommunityPublisher))
          as CommunityPublisher;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityPublisher create() => CommunityPublisher._();
  @$core.override
  CommunityPublisher createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityPublisher getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityPublisher>(create);
  static CommunityPublisher? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get communityId => $_getSZ(1);
  @$pb.TagNumber(2)
  set communityId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasCommunityId() => $_has(1);
  @$pb.TagNumber(2)
  void clearCommunityId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get publisherId => $_getSZ(2);
  @$pb.TagNumber(3)
  set publisherId($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasPublisherId() => $_has(2);
  @$pb.TagNumber(3)
  void clearPublisherId() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get createdAt => $_getSZ(3);
  @$pb.TagNumber(4)
  set createdAt($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasCreatedAt() => $_has(3);
  @$pb.TagNumber(4)
  void clearCreatedAt() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get updatedAt => $_getSZ(4);
  @$pb.TagNumber(5)
  set updatedAt($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasUpdatedAt() => $_has(4);
  @$pb.TagNumber(5)
  void clearUpdatedAt() => $_clearField(5);
}

class CommunitySocial extends $pb.GeneratedMessage {
  factory CommunitySocial({
    $0.Community? community,
    $core.Iterable<CommunityPublisher>? publishers,
  }) {
    final result = create();
    if (community != null) result.community = community;
    if (publishers != null) result.publishers.addAll(publishers);
    return result;
  }

  CommunitySocial._();

  factory CommunitySocial.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunitySocial.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunitySocial',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<$0.Community>(1, _omitFieldNames ? '' : 'community',
        subBuilder: $0.Community.create)
    ..pPM<CommunityPublisher>(1000, _omitFieldNames ? '' : 'publishers',
        subBuilder: CommunityPublisher.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunitySocial clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunitySocial copyWith(void Function(CommunitySocial) updates) =>
      super.copyWith((message) => updates(message as CommunitySocial))
          as CommunitySocial;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunitySocial create() => CommunitySocial._();
  @$core.override
  CommunitySocial createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunitySocial getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunitySocial>(create);
  static CommunitySocial? _defaultInstance;

  @$pb.TagNumber(1)
  $0.Community get community => $_getN(0);
  @$pb.TagNumber(1)
  set community($0.Community value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasCommunity() => $_has(0);
  @$pb.TagNumber(1)
  void clearCommunity() => $_clearField(1);
  @$pb.TagNumber(1)
  $0.Community ensureCommunity() => $_ensure(0);

  @$pb.TagNumber(1000)
  $pb.PbList<CommunityPublisher> get publishers => $_getList(1);
}

class SocialsSearchRequest extends $pb.GeneratedMessage {
  factory SocialsSearchRequest({
    $core.String? query,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
    $core.Iterable<$core.String>? communities,
  }) {
    final result = create();
    if (query != null) result.query = query;
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    if (communities != null) result.communities.addAll(communities);
    return result;
  }

  SocialsSearchRequest._();

  factory SocialsSearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SocialsSearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SocialsSearchRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..a<$fixnum.Int64>(
        900, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(901, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..pPS(1000, _omitFieldNames ? '' : 'communities')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SocialsSearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SocialsSearchRequest copyWith(void Function(SocialsSearchRequest) updates) =>
      super.copyWith((message) => updates(message as SocialsSearchRequest))
          as SocialsSearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SocialsSearchRequest create() => SocialsSearchRequest._();
  @$core.override
  SocialsSearchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SocialsSearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SocialsSearchRequest>(create);
  static SocialsSearchRequest? _defaultInstance;

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

  @$pb.TagNumber(1000)
  $pb.PbList<$core.String> get communities => $_getList(3);
}

class SocialsSearchResponse extends $pb.GeneratedMessage {
  factory SocialsSearchResponse({
    SocialsSearchRequest? next,
    $core.Iterable<CommunitySocial>? items,
  }) {
    final result = create();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  SocialsSearchResponse._();

  factory SocialsSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SocialsSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SocialsSearchResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<SocialsSearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: SocialsSearchRequest.create)
    ..pPM<CommunitySocial>(2, _omitFieldNames ? '' : 'items',
        subBuilder: CommunitySocial.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SocialsSearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SocialsSearchResponse copyWith(
          void Function(SocialsSearchResponse) updates) =>
      super.copyWith((message) => updates(message as SocialsSearchResponse))
          as SocialsSearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SocialsSearchResponse create() => SocialsSearchResponse._();
  @$core.override
  SocialsSearchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SocialsSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SocialsSearchResponse>(create);
  static SocialsSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  SocialsSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(SocialsSearchRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  SocialsSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<CommunitySocial> get items => $_getList(1);
}

class PluginPublisherSearchRequest extends $pb.GeneratedMessage {
  factory PluginPublisherSearchRequest({
    $core.String? query,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
    $core.Iterable<$core.String>? excluded,
  }) {
    final result = create();
    if (query != null) result.query = query;
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    if (excluded != null) result.excluded.addAll(excluded);
    return result;
  }

  PluginPublisherSearchRequest._();

  factory PluginPublisherSearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PluginPublisherSearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PluginPublisherSearchRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..a<$fixnum.Int64>(
        900, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(901, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..pPS(1000, _omitFieldNames ? '' : 'excluded')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginPublisherSearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginPublisherSearchRequest copyWith(
          void Function(PluginPublisherSearchRequest) updates) =>
      super.copyWith(
              (message) => updates(message as PluginPublisherSearchRequest))
          as PluginPublisherSearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PluginPublisherSearchRequest create() =>
      PluginPublisherSearchRequest._();
  @$core.override
  PluginPublisherSearchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PluginPublisherSearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PluginPublisherSearchRequest>(create);
  static PluginPublisherSearchRequest? _defaultInstance;

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

  /// publishers to leave out - what a community has already attached.
  @$pb.TagNumber(1000)
  $pb.PbList<$core.String> get excluded => $_getList(3);
}

class PluginPublisherSearchResponse extends $pb.GeneratedMessage {
  factory PluginPublisherSearchResponse({
    PluginPublisherSearchRequest? next,
    $core.Iterable<PluginPublisher>? items,
  }) {
    final result = create();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  PluginPublisherSearchResponse._();

  factory PluginPublisherSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PluginPublisherSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PluginPublisherSearchResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<PluginPublisherSearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: PluginPublisherSearchRequest.create)
    ..pPM<PluginPublisher>(2, _omitFieldNames ? '' : 'items',
        subBuilder: PluginPublisher.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginPublisherSearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginPublisherSearchResponse copyWith(
          void Function(PluginPublisherSearchResponse) updates) =>
      super.copyWith(
              (message) => updates(message as PluginPublisherSearchResponse))
          as PluginPublisherSearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PluginPublisherSearchResponse create() =>
      PluginPublisherSearchResponse._();
  @$core.override
  PluginPublisherSearchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PluginPublisherSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PluginPublisherSearchResponse>(create);
  static PluginPublisherSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  PluginPublisherSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(PluginPublisherSearchRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  PluginPublisherSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<PluginPublisher> get items => $_getList(1);
}

class PluginPublisherCreateResponse extends $pb.GeneratedMessage {
  factory PluginPublisherCreateResponse({
    PluginPublisher? publisher,
  }) {
    final result = create();
    if (publisher != null) result.publisher = publisher;
    return result;
  }

  PluginPublisherCreateResponse._();

  factory PluginPublisherCreateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PluginPublisherCreateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PluginPublisherCreateResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<PluginPublisher>(1, _omitFieldNames ? '' : 'publisher',
        subBuilder: PluginPublisher.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginPublisherCreateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginPublisherCreateResponse copyWith(
          void Function(PluginPublisherCreateResponse) updates) =>
      super.copyWith(
              (message) => updates(message as PluginPublisherCreateResponse))
          as PluginPublisherCreateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PluginPublisherCreateResponse create() =>
      PluginPublisherCreateResponse._();
  @$core.override
  PluginPublisherCreateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PluginPublisherCreateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PluginPublisherCreateResponse>(create);
  static PluginPublisherCreateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  PluginPublisher get publisher => $_getN(0);
  @$pb.TagNumber(1)
  set publisher(PluginPublisher value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPublisher() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublisher() => $_clearField(1);
  @$pb.TagNumber(1)
  PluginPublisher ensurePublisher() => $_ensure(0);
}

class PluginPublisherFindResponse extends $pb.GeneratedMessage {
  factory PluginPublisherFindResponse({
    PluginPublisher? publisher,
  }) {
    final result = create();
    if (publisher != null) result.publisher = publisher;
    return result;
  }

  PluginPublisherFindResponse._();

  factory PluginPublisherFindResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PluginPublisherFindResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PluginPublisherFindResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<PluginPublisher>(1, _omitFieldNames ? '' : 'publisher',
        subBuilder: PluginPublisher.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginPublisherFindResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginPublisherFindResponse copyWith(
          void Function(PluginPublisherFindResponse) updates) =>
      super.copyWith(
              (message) => updates(message as PluginPublisherFindResponse))
          as PluginPublisherFindResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PluginPublisherFindResponse create() =>
      PluginPublisherFindResponse._();
  @$core.override
  PluginPublisherFindResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PluginPublisherFindResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PluginPublisherFindResponse>(create);
  static PluginPublisherFindResponse? _defaultInstance;

  @$pb.TagNumber(1)
  PluginPublisher get publisher => $_getN(0);
  @$pb.TagNumber(1)
  set publisher(PluginPublisher value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPublisher() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublisher() => $_clearField(1);
  @$pb.TagNumber(1)
  PluginPublisher ensurePublisher() => $_ensure(0);
}

class PluginPublisherUpdateRequest extends $pb.GeneratedMessage {
  factory PluginPublisherUpdateRequest({
    PluginPublisher? publisher,
  }) {
    final result = create();
    if (publisher != null) result.publisher = publisher;
    return result;
  }

  PluginPublisherUpdateRequest._();

  factory PluginPublisherUpdateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PluginPublisherUpdateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PluginPublisherUpdateRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<PluginPublisher>(1, _omitFieldNames ? '' : 'publisher',
        subBuilder: PluginPublisher.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginPublisherUpdateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginPublisherUpdateRequest copyWith(
          void Function(PluginPublisherUpdateRequest) updates) =>
      super.copyWith(
              (message) => updates(message as PluginPublisherUpdateRequest))
          as PluginPublisherUpdateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PluginPublisherUpdateRequest create() =>
      PluginPublisherUpdateRequest._();
  @$core.override
  PluginPublisherUpdateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PluginPublisherUpdateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PluginPublisherUpdateRequest>(create);
  static PluginPublisherUpdateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  PluginPublisher get publisher => $_getN(0);
  @$pb.TagNumber(1)
  set publisher(PluginPublisher value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPublisher() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublisher() => $_clearField(1);
  @$pb.TagNumber(1)
  PluginPublisher ensurePublisher() => $_ensure(0);
}

class PluginPublisherUpdateResponse extends $pb.GeneratedMessage {
  factory PluginPublisherUpdateResponse({
    PluginPublisher? publisher,
  }) {
    final result = create();
    if (publisher != null) result.publisher = publisher;
    return result;
  }

  PluginPublisherUpdateResponse._();

  factory PluginPublisherUpdateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PluginPublisherUpdateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PluginPublisherUpdateResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<PluginPublisher>(1, _omitFieldNames ? '' : 'publisher',
        subBuilder: PluginPublisher.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginPublisherUpdateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginPublisherUpdateResponse copyWith(
          void Function(PluginPublisherUpdateResponse) updates) =>
      super.copyWith(
              (message) => updates(message as PluginPublisherUpdateResponse))
          as PluginPublisherUpdateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PluginPublisherUpdateResponse create() =>
      PluginPublisherUpdateResponse._();
  @$core.override
  PluginPublisherUpdateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PluginPublisherUpdateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PluginPublisherUpdateResponse>(create);
  static PluginPublisherUpdateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  PluginPublisher get publisher => $_getN(0);
  @$pb.TagNumber(1)
  set publisher(PluginPublisher value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPublisher() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublisher() => $_clearField(1);
  @$pb.TagNumber(1)
  PluginPublisher ensurePublisher() => $_ensure(0);
}

/// the clone of a publisher is the same module under a second identity, so the
/// request carries nothing beyond the id in the path.
class PluginPublisherCloneResponse extends $pb.GeneratedMessage {
  factory PluginPublisherCloneResponse({
    PluginPublisher? publisher,
  }) {
    final result = create();
    if (publisher != null) result.publisher = publisher;
    return result;
  }

  PluginPublisherCloneResponse._();

  factory PluginPublisherCloneResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PluginPublisherCloneResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PluginPublisherCloneResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<PluginPublisher>(1, _omitFieldNames ? '' : 'publisher',
        subBuilder: PluginPublisher.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginPublisherCloneResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginPublisherCloneResponse copyWith(
          void Function(PluginPublisherCloneResponse) updates) =>
      super.copyWith(
              (message) => updates(message as PluginPublisherCloneResponse))
          as PluginPublisherCloneResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PluginPublisherCloneResponse create() =>
      PluginPublisherCloneResponse._();
  @$core.override
  PluginPublisherCloneResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PluginPublisherCloneResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PluginPublisherCloneResponse>(create);
  static PluginPublisherCloneResponse? _defaultInstance;

  @$pb.TagNumber(1)
  PluginPublisher get publisher => $_getN(0);
  @$pb.TagNumber(1)
  set publisher(PluginPublisher value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPublisher() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublisher() => $_clearField(1);
  @$pb.TagNumber(1)
  PluginPublisher ensurePublisher() => $_ensure(0);
}

class PluginPublisherDeleteResponse extends $pb.GeneratedMessage {
  factory PluginPublisherDeleteResponse({
    PluginPublisher? publisher,
  }) {
    final result = create();
    if (publisher != null) result.publisher = publisher;
    return result;
  }

  PluginPublisherDeleteResponse._();

  factory PluginPublisherDeleteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PluginPublisherDeleteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PluginPublisherDeleteResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<PluginPublisher>(1, _omitFieldNames ? '' : 'publisher',
        subBuilder: PluginPublisher.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginPublisherDeleteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginPublisherDeleteResponse copyWith(
          void Function(PluginPublisherDeleteResponse) updates) =>
      super.copyWith(
              (message) => updates(message as PluginPublisherDeleteResponse))
          as PluginPublisherDeleteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PluginPublisherDeleteResponse create() =>
      PluginPublisherDeleteResponse._();
  @$core.override
  PluginPublisherDeleteResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PluginPublisherDeleteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PluginPublisherDeleteResponse>(create);
  static PluginPublisherDeleteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  PluginPublisher get publisher => $_getN(0);
  @$pb.TagNumber(1)
  set publisher(PluginPublisher value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPublisher() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublisher() => $_clearField(1);
  @$pb.TagNumber(1)
  PluginPublisher ensurePublisher() => $_ensure(0);
}

class CommunityPublisherEnableResponse extends $pb.GeneratedMessage {
  factory CommunityPublisherEnableResponse({
    CommunityPublisher? enabled,
  }) {
    final result = create();
    if (enabled != null) result.enabled = enabled;
    return result;
  }

  CommunityPublisherEnableResponse._();

  factory CommunityPublisherEnableResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityPublisherEnableResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityPublisherEnableResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<CommunityPublisher>(1, _omitFieldNames ? '' : 'enabled',
        subBuilder: CommunityPublisher.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityPublisherEnableResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityPublisherEnableResponse copyWith(
          void Function(CommunityPublisherEnableResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CommunityPublisherEnableResponse))
          as CommunityPublisherEnableResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityPublisherEnableResponse create() =>
      CommunityPublisherEnableResponse._();
  @$core.override
  CommunityPublisherEnableResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityPublisherEnableResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityPublisherEnableResponse>(
          create);
  static CommunityPublisherEnableResponse? _defaultInstance;

  @$pb.TagNumber(1)
  CommunityPublisher get enabled => $_getN(0);
  @$pb.TagNumber(1)
  set enabled(CommunityPublisher value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasEnabled() => $_has(0);
  @$pb.TagNumber(1)
  void clearEnabled() => $_clearField(1);
  @$pb.TagNumber(1)
  CommunityPublisher ensureEnabled() => $_ensure(0);
}

class CommunityPublisherDisableResponse extends $pb.GeneratedMessage {
  factory CommunityPublisherDisableResponse({
    CommunityPublisher? disabled,
  }) {
    final result = create();
    if (disabled != null) result.disabled = disabled;
    return result;
  }

  CommunityPublisherDisableResponse._();

  factory CommunityPublisherDisableResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityPublisherDisableResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityPublisherDisableResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<CommunityPublisher>(1, _omitFieldNames ? '' : 'disabled',
        subBuilder: CommunityPublisher.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityPublisherDisableResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityPublisherDisableResponse copyWith(
          void Function(CommunityPublisherDisableResponse) updates) =>
      super.copyWith((message) =>
              updates(message as CommunityPublisherDisableResponse))
          as CommunityPublisherDisableResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityPublisherDisableResponse create() =>
      CommunityPublisherDisableResponse._();
  @$core.override
  CommunityPublisherDisableResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityPublisherDisableResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityPublisherDisableResponse>(
          create);
  static CommunityPublisherDisableResponse? _defaultInstance;

  @$pb.TagNumber(1)
  CommunityPublisher get disabled => $_getN(0);
  @$pb.TagNumber(1)
  set disabled(CommunityPublisher value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDisabled() => $_has(0);
  @$pb.TagNumber(1)
  void clearDisabled() => $_clearField(1);
  @$pb.TagNumber(1)
  CommunityPublisher ensureDisabled() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
