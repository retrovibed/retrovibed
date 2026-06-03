// This is a generated file - do not edit.
//
// Generated from community.publish.proto.

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

class PublishedContent extends $pb.GeneratedMessage {
  factory PublishedContent({
    $core.String? id,
    $core.String? communityId,
    $core.String? knownMediaId,
    $core.String? magnetUri,
    $core.String? publishedAt,
    $core.String? createdAt,
    $core.String? updatedAt,
    $core.String? archivedId,
    $core.String? title,
    $core.String? description,
    $core.String? mimetype,
    $core.String? encryptionSeed,
    $fixnum.Int64? bytes,
    $core.String? libraryId,
    $core.String? oauthGoogleId,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (communityId != null) result.communityId = communityId;
    if (knownMediaId != null) result.knownMediaId = knownMediaId;
    if (magnetUri != null) result.magnetUri = magnetUri;
    if (publishedAt != null) result.publishedAt = publishedAt;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (archivedId != null) result.archivedId = archivedId;
    if (title != null) result.title = title;
    if (description != null) result.description = description;
    if (mimetype != null) result.mimetype = mimetype;
    if (encryptionSeed != null) result.encryptionSeed = encryptionSeed;
    if (bytes != null) result.bytes = bytes;
    if (libraryId != null) result.libraryId = libraryId;
    if (oauthGoogleId != null) result.oauthGoogleId = oauthGoogleId;
    return result;
  }

  PublishedContent._();

  factory PublishedContent.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishedContent.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishedContent',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'community_id')
    ..aOS(3, _omitFieldNames ? '' : 'known_media_id')
    ..aOS(4, _omitFieldNames ? '' : 'magnet_uri')
    ..aOS(5, _omitFieldNames ? '' : 'published_at')
    ..aOS(6, _omitFieldNames ? '' : 'created_at')
    ..aOS(7, _omitFieldNames ? '' : 'updated_at')
    ..aOS(8, _omitFieldNames ? '' : 'archived_id')
    ..aOS(9, _omitFieldNames ? '' : 'title')
    ..aOS(10, _omitFieldNames ? '' : 'description')
    ..aOS(11, _omitFieldNames ? '' : 'mimetype')
    ..aOS(12, _omitFieldNames ? '' : 'encryption_seed')
    ..a<$fixnum.Int64>(13, _omitFieldNames ? '' : 'bytes', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOS(1000, _omitFieldNames ? '' : 'library_id')
    ..aOS(1001, _omitFieldNames ? '' : 'oauth_google_id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedContent clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedContent copyWith(void Function(PublishedContent) updates) =>
      super.copyWith((message) => updates(message as PublishedContent))
          as PublishedContent;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishedContent create() => PublishedContent._();
  @$core.override
  PublishedContent createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishedContent getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishedContent>(create);
  static PublishedContent? _defaultInstance;

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
  $core.String get knownMediaId => $_getSZ(2);
  @$pb.TagNumber(3)
  set knownMediaId($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasKnownMediaId() => $_has(2);
  @$pb.TagNumber(3)
  void clearKnownMediaId() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get magnetUri => $_getSZ(3);
  @$pb.TagNumber(4)
  set magnetUri($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasMagnetUri() => $_has(3);
  @$pb.TagNumber(4)
  void clearMagnetUri() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get publishedAt => $_getSZ(4);
  @$pb.TagNumber(5)
  set publishedAt($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasPublishedAt() => $_has(4);
  @$pb.TagNumber(5)
  void clearPublishedAt() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get createdAt => $_getSZ(5);
  @$pb.TagNumber(6)
  set createdAt($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasCreatedAt() => $_has(5);
  @$pb.TagNumber(6)
  void clearCreatedAt() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.String get updatedAt => $_getSZ(6);
  @$pb.TagNumber(7)
  set updatedAt($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasUpdatedAt() => $_has(6);
  @$pb.TagNumber(7)
  void clearUpdatedAt() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.String get archivedId => $_getSZ(7);
  @$pb.TagNumber(8)
  set archivedId($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasArchivedId() => $_has(7);
  @$pb.TagNumber(8)
  void clearArchivedId() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.String get title => $_getSZ(8);
  @$pb.TagNumber(9)
  set title($core.String value) => $_setString(8, value);
  @$pb.TagNumber(9)
  $core.bool hasTitle() => $_has(8);
  @$pb.TagNumber(9)
  void clearTitle() => $_clearField(9);

  @$pb.TagNumber(10)
  $core.String get description => $_getSZ(9);
  @$pb.TagNumber(10)
  set description($core.String value) => $_setString(9, value);
  @$pb.TagNumber(10)
  $core.bool hasDescription() => $_has(9);
  @$pb.TagNumber(10)
  void clearDescription() => $_clearField(10);

  @$pb.TagNumber(11)
  $core.String get mimetype => $_getSZ(10);
  @$pb.TagNumber(11)
  set mimetype($core.String value) => $_setString(10, value);
  @$pb.TagNumber(11)
  $core.bool hasMimetype() => $_has(10);
  @$pb.TagNumber(11)
  void clearMimetype() => $_clearField(11);

  @$pb.TagNumber(12)
  $core.String get encryptionSeed => $_getSZ(11);
  @$pb.TagNumber(12)
  set encryptionSeed($core.String value) => $_setString(11, value);
  @$pb.TagNumber(12)
  $core.bool hasEncryptionSeed() => $_has(11);
  @$pb.TagNumber(12)
  void clearEncryptionSeed() => $_clearField(12);

  @$pb.TagNumber(13)
  $fixnum.Int64 get bytes => $_getI64(12);
  @$pb.TagNumber(13)
  set bytes($fixnum.Int64 value) => $_setInt64(12, value);
  @$pb.TagNumber(13)
  $core.bool hasBytes() => $_has(12);
  @$pb.TagNumber(13)
  void clearBytes() => $_clearField(13);

  /// private fields for retrovibed use only, not populated by clients.
  @$pb.TagNumber(1000)
  $core.String get libraryId => $_getSZ(13);
  @$pb.TagNumber(1000)
  set libraryId($core.String value) => $_setString(13, value);
  @$pb.TagNumber(1000)
  $core.bool hasLibraryId() => $_has(13);
  @$pb.TagNumber(1000)
  void clearLibraryId() => $_clearField(1000);

  @$pb.TagNumber(1001)
  $core.String get oauthGoogleId => $_getSZ(14);
  @$pb.TagNumber(1001)
  set oauthGoogleId($core.String value) => $_setString(14, value);
  @$pb.TagNumber(1001)
  $core.bool hasOauthGoogleId() => $_has(14);
  @$pb.TagNumber(1001)
  void clearOauthGoogleId() => $_clearField(1001);
}

class PublishContentRequest extends $pb.GeneratedMessage {
  factory PublishContentRequest({
    PublishedContent? publishedContent,
    $0.PublishMode? publishMode,
  }) {
    final result = create();
    if (publishedContent != null) result.publishedContent = publishedContent;
    if (publishMode != null) result.publishMode = publishMode;
    return result;
  }

  PublishContentRequest._();

  factory PublishContentRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishContentRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishContentRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<PublishedContent>(1, _omitFieldNames ? '' : 'published_content',
        subBuilder: PublishedContent.create)
    ..aE<$0.PublishMode>(2, _omitFieldNames ? '' : 'publish_mode',
        enumValues: $0.PublishMode.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishContentRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishContentRequest copyWith(
          void Function(PublishContentRequest) updates) =>
      super.copyWith((message) => updates(message as PublishContentRequest))
          as PublishContentRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishContentRequest create() => PublishContentRequest._();
  @$core.override
  PublishContentRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishContentRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishContentRequest>(create);
  static PublishContentRequest? _defaultInstance;

  @$pb.TagNumber(1)
  PublishedContent get publishedContent => $_getN(0);
  @$pb.TagNumber(1)
  set publishedContent(PublishedContent value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPublishedContent() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublishedContent() => $_clearField(1);
  @$pb.TagNumber(1)
  PublishedContent ensurePublishedContent() => $_ensure(0);

  @$pb.TagNumber(2)
  $0.PublishMode get publishMode => $_getN(1);
  @$pb.TagNumber(2)
  set publishMode($0.PublishMode value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasPublishMode() => $_has(1);
  @$pb.TagNumber(2)
  void clearPublishMode() => $_clearField(2);
}

class PublishContentResponse extends $pb.GeneratedMessage {
  factory PublishContentResponse({
    PublishedContent? publishedContent,
  }) {
    final result = create();
    if (publishedContent != null) result.publishedContent = publishedContent;
    return result;
  }

  PublishContentResponse._();

  factory PublishContentResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishContentResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishContentResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<PublishedContent>(1, _omitFieldNames ? '' : 'published_content',
        subBuilder: PublishedContent.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishContentResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishContentResponse copyWith(
          void Function(PublishContentResponse) updates) =>
      super.copyWith((message) => updates(message as PublishContentResponse))
          as PublishContentResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishContentResponse create() => PublishContentResponse._();
  @$core.override
  PublishContentResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishContentResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishContentResponse>(create);
  static PublishContentResponse? _defaultInstance;

  @$pb.TagNumber(1)
  PublishedContent get publishedContent => $_getN(0);
  @$pb.TagNumber(1)
  set publishedContent(PublishedContent value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPublishedContent() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublishedContent() => $_clearField(1);
  @$pb.TagNumber(1)
  PublishedContent ensurePublishedContent() => $_ensure(0);
}

class PublishContentDeleteRequest extends $pb.GeneratedMessage {
  factory PublishContentDeleteRequest({
    PublishedContent? publishedContent,
  }) {
    final result = create();
    if (publishedContent != null) result.publishedContent = publishedContent;
    return result;
  }

  PublishContentDeleteRequest._();

  factory PublishContentDeleteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishContentDeleteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishContentDeleteRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<PublishedContent>(1, _omitFieldNames ? '' : 'published_content',
        subBuilder: PublishedContent.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishContentDeleteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishContentDeleteRequest copyWith(
          void Function(PublishContentDeleteRequest) updates) =>
      super.copyWith(
              (message) => updates(message as PublishContentDeleteRequest))
          as PublishContentDeleteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishContentDeleteRequest create() =>
      PublishContentDeleteRequest._();
  @$core.override
  PublishContentDeleteRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishContentDeleteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishContentDeleteRequest>(create);
  static PublishContentDeleteRequest? _defaultInstance;

  @$pb.TagNumber(1)
  PublishedContent get publishedContent => $_getN(0);
  @$pb.TagNumber(1)
  set publishedContent(PublishedContent value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPublishedContent() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublishedContent() => $_clearField(1);
  @$pb.TagNumber(1)
  PublishedContent ensurePublishedContent() => $_ensure(0);
}

class PublishContentDeleteResponse extends $pb.GeneratedMessage {
  factory PublishContentDeleteResponse({
    PublishedContent? publishedContent,
  }) {
    final result = create();
    if (publishedContent != null) result.publishedContent = publishedContent;
    return result;
  }

  PublishContentDeleteResponse._();

  factory PublishContentDeleteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishContentDeleteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishContentDeleteResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<PublishedContent>(1, _omitFieldNames ? '' : 'published_content',
        subBuilder: PublishedContent.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishContentDeleteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishContentDeleteResponse copyWith(
          void Function(PublishContentDeleteResponse) updates) =>
      super.copyWith(
              (message) => updates(message as PublishContentDeleteResponse))
          as PublishContentDeleteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishContentDeleteResponse create() =>
      PublishContentDeleteResponse._();
  @$core.override
  PublishContentDeleteResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishContentDeleteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishContentDeleteResponse>(create);
  static PublishContentDeleteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  PublishedContent get publishedContent => $_getN(0);
  @$pb.TagNumber(1)
  set publishedContent(PublishedContent value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPublishedContent() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublishedContent() => $_clearField(1);
  @$pb.TagNumber(1)
  PublishedContent ensurePublishedContent() => $_ensure(0);
}

class PublishedContentSearchRequest extends $pb.GeneratedMessage {
  factory PublishedContentSearchRequest({
    $core.String? communityId,
    $core.String? query,
    $core.String? sync,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = create();
    if (communityId != null) result.communityId = communityId;
    if (query != null) result.query = query;
    if (sync != null) result.sync = sync;
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    return result;
  }

  PublishedContentSearchRequest._();

  factory PublishedContentSearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishedContentSearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishedContentSearchRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'community_id')
    ..aOS(2, _omitFieldNames ? '' : 'query')
    ..aOS(3, _omitFieldNames ? '' : 'sync')
    ..a<$fixnum.Int64>(
        900, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(901, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedContentSearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedContentSearchRequest copyWith(
          void Function(PublishedContentSearchRequest) updates) =>
      super.copyWith(
              (message) => updates(message as PublishedContentSearchRequest))
          as PublishedContentSearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishedContentSearchRequest create() =>
      PublishedContentSearchRequest._();
  @$core.override
  PublishedContentSearchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishedContentSearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishedContentSearchRequest>(create);
  static PublishedContentSearchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get communityId => $_getSZ(0);
  @$pb.TagNumber(1)
  set communityId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasCommunityId() => $_has(0);
  @$pb.TagNumber(1)
  void clearCommunityId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get query => $_getSZ(1);
  @$pb.TagNumber(2)
  set query($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasQuery() => $_has(1);
  @$pb.TagNumber(2)
  void clearQuery() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get sync => $_getSZ(2);
  @$pb.TagNumber(3)
  set sync($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasSync() => $_has(2);
  @$pb.TagNumber(3)
  void clearSync() => $_clearField(3);

  @$pb.TagNumber(900)
  $fixnum.Int64 get offset => $_getI64(3);
  @$pb.TagNumber(900)
  set offset($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(900)
  $core.bool hasOffset() => $_has(3);
  @$pb.TagNumber(900)
  void clearOffset() => $_clearField(900);

  @$pb.TagNumber(901)
  $fixnum.Int64 get limit => $_getI64(4);
  @$pb.TagNumber(901)
  set limit($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(901)
  $core.bool hasLimit() => $_has(4);
  @$pb.TagNumber(901)
  void clearLimit() => $_clearField(901);
}

class PublishedContentSearchResponse extends $pb.GeneratedMessage {
  factory PublishedContentSearchResponse({
    $0.Community? community,
    PublishedContentSearchRequest? next,
    $core.Iterable<PublishedContent>? items,
  }) {
    final result = create();
    if (community != null) result.community = community;
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  PublishedContentSearchResponse._();

  factory PublishedContentSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishedContentSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishedContentSearchResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<$0.Community>(1, _omitFieldNames ? '' : 'community',
        subBuilder: $0.Community.create)
    ..aOM<PublishedContentSearchRequest>(2, _omitFieldNames ? '' : 'next',
        subBuilder: PublishedContentSearchRequest.create)
    ..pPM<PublishedContent>(3, _omitFieldNames ? '' : 'items',
        subBuilder: PublishedContent.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedContentSearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedContentSearchResponse copyWith(
          void Function(PublishedContentSearchResponse) updates) =>
      super.copyWith(
              (message) => updates(message as PublishedContentSearchResponse))
          as PublishedContentSearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishedContentSearchResponse create() =>
      PublishedContentSearchResponse._();
  @$core.override
  PublishedContentSearchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishedContentSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishedContentSearchResponse>(create);
  static PublishedContentSearchResponse? _defaultInstance;

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

  @$pb.TagNumber(2)
  PublishedContentSearchRequest get next => $_getN(1);
  @$pb.TagNumber(2)
  set next(PublishedContentSearchRequest value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasNext() => $_has(1);
  @$pb.TagNumber(2)
  void clearNext() => $_clearField(2);
  @$pb.TagNumber(2)
  PublishedContentSearchRequest ensureNext() => $_ensure(1);

  @$pb.TagNumber(3)
  $pb.PbList<PublishedContent> get items => $_getList(2);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
