// This is a generated file - do not edit.
//
// Generated from media/content.addressable.storage.proto.

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

class Media extends $pb.GeneratedMessage {
  factory Media({
    $core.String? id,
    $core.String? accountId,
    $core.String? uploadedBy,
    $core.String? createdAt,
    $core.String? updatedAt,
    $core.String? mimetype,
    $fixnum.Int64? bytes,
    $fixnum.Int64? usage,
    $core.String? tombstonedAt,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (accountId != null) result.accountId = accountId;
    if (uploadedBy != null) result.uploadedBy = uploadedBy;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (mimetype != null) result.mimetype = mimetype;
    if (bytes != null) result.bytes = bytes;
    if (usage != null) result.usage = usage;
    if (tombstonedAt != null) result.tombstonedAt = tombstonedAt;
    return result;
  }

  Media._();

  factory Media.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Media.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Media',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'retrovibed.cas'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'account_id')
    ..aOS(3, _omitFieldNames ? '' : 'uploaded_by')
    ..aOS(4, _omitFieldNames ? '' : 'created_at')
    ..aOS(5, _omitFieldNames ? '' : 'updated_at')
    ..aOS(6, _omitFieldNames ? '' : 'mimetype')
    ..a<$fixnum.Int64>(7, _omitFieldNames ? '' : 'bytes', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(8, _omitFieldNames ? '' : 'usage', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOS(9, _omitFieldNames ? '' : 'tombstoned_at')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Media clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Media copyWith(void Function(Media) updates) =>
      super.copyWith((message) => updates(message as Media)) as Media;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Media create() => Media._();
  @$core.override
  Media createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Media getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Media>(create);
  static Media? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get accountId => $_getSZ(1);
  @$pb.TagNumber(2)
  set accountId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasAccountId() => $_has(1);
  @$pb.TagNumber(2)
  void clearAccountId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get uploadedBy => $_getSZ(2);
  @$pb.TagNumber(3)
  set uploadedBy($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasUploadedBy() => $_has(2);
  @$pb.TagNumber(3)
  void clearUploadedBy() => $_clearField(3);

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

  @$pb.TagNumber(6)
  $core.String get mimetype => $_getSZ(5);
  @$pb.TagNumber(6)
  set mimetype($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasMimetype() => $_has(5);
  @$pb.TagNumber(6)
  void clearMimetype() => $_clearField(6);

  @$pb.TagNumber(7)
  $fixnum.Int64 get bytes => $_getI64(6);
  @$pb.TagNumber(7)
  set bytes($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(7)
  $core.bool hasBytes() => $_has(6);
  @$pb.TagNumber(7)
  void clearBytes() => $_clearField(7);

  @$pb.TagNumber(8)
  $fixnum.Int64 get usage => $_getI64(7);
  @$pb.TagNumber(8)
  set usage($fixnum.Int64 value) => $_setInt64(7, value);
  @$pb.TagNumber(8)
  $core.bool hasUsage() => $_has(7);
  @$pb.TagNumber(8)
  void clearUsage() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.String get tombstonedAt => $_getSZ(8);
  @$pb.TagNumber(9)
  set tombstonedAt($core.String value) => $_setString(8, value);
  @$pb.TagNumber(9)
  $core.bool hasTombstonedAt() => $_has(8);
  @$pb.TagNumber(9)
  void clearTombstonedAt() => $_clearField(9);
}

class MediaSearchRequest extends $pb.GeneratedMessage {
  factory MediaSearchRequest({
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

  MediaSearchRequest._();

  factory MediaSearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MediaSearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MediaSearchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'retrovibed.cas'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(3, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaSearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaSearchRequest copyWith(void Function(MediaSearchRequest) updates) =>
      super.copyWith((message) => updates(message as MediaSearchRequest))
          as MediaSearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MediaSearchRequest create() => MediaSearchRequest._();
  @$core.override
  MediaSearchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MediaSearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MediaSearchRequest>(create);
  static MediaSearchRequest? _defaultInstance;

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

class MediaSearchResponse extends $pb.GeneratedMessage {
  factory MediaSearchResponse({
    MediaSearchRequest? next,
    $core.Iterable<Media>? items,
  }) {
    final result = create();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  MediaSearchResponse._();

  factory MediaSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MediaSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MediaSearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'retrovibed.cas'),
      createEmptyInstance: create)
    ..aOM<MediaSearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: MediaSearchRequest.create)
    ..pPM<Media>(2, _omitFieldNames ? '' : 'items', subBuilder: Media.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaSearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaSearchResponse copyWith(void Function(MediaSearchResponse) updates) =>
      super.copyWith((message) => updates(message as MediaSearchResponse))
          as MediaSearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MediaSearchResponse create() => MediaSearchResponse._();
  @$core.override
  MediaSearchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MediaSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MediaSearchResponse>(create);
  static MediaSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  MediaSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(MediaSearchRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  MediaSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<Media> get items => $_getList(1);
}

class MediaCreateRequest extends $pb.GeneratedMessage {
  factory MediaCreateRequest({
    Media? media,
  }) {
    final result = create();
    if (media != null) result.media = media;
    return result;
  }

  MediaCreateRequest._();

  factory MediaCreateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MediaCreateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MediaCreateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'retrovibed.cas'),
      createEmptyInstance: create)
    ..aOM<Media>(1, _omitFieldNames ? '' : 'media', subBuilder: Media.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaCreateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaCreateRequest copyWith(void Function(MediaCreateRequest) updates) =>
      super.copyWith((message) => updates(message as MediaCreateRequest))
          as MediaCreateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MediaCreateRequest create() => MediaCreateRequest._();
  @$core.override
  MediaCreateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MediaCreateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MediaCreateRequest>(create);
  static MediaCreateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Media get media => $_getN(0);
  @$pb.TagNumber(1)
  set media(Media value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMedia() => $_has(0);
  @$pb.TagNumber(1)
  void clearMedia() => $_clearField(1);
  @$pb.TagNumber(1)
  Media ensureMedia() => $_ensure(0);
}

class MediaCreateResponse extends $pb.GeneratedMessage {
  factory MediaCreateResponse({
    Media? media,
  }) {
    final result = create();
    if (media != null) result.media = media;
    return result;
  }

  MediaCreateResponse._();

  factory MediaCreateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MediaCreateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MediaCreateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'retrovibed.cas'),
      createEmptyInstance: create)
    ..aOM<Media>(1, _omitFieldNames ? '' : 'media', subBuilder: Media.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaCreateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaCreateResponse copyWith(void Function(MediaCreateResponse) updates) =>
      super.copyWith((message) => updates(message as MediaCreateResponse))
          as MediaCreateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MediaCreateResponse create() => MediaCreateResponse._();
  @$core.override
  MediaCreateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MediaCreateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MediaCreateResponse>(create);
  static MediaCreateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Media get media => $_getN(0);
  @$pb.TagNumber(1)
  set media(Media value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMedia() => $_has(0);
  @$pb.TagNumber(1)
  void clearMedia() => $_clearField(1);
  @$pb.TagNumber(1)
  Media ensureMedia() => $_ensure(0);
}

class MediaUploadRequest extends $pb.GeneratedMessage {
  factory MediaUploadRequest({
    Media? media,
  }) {
    final result = create();
    if (media != null) result.media = media;
    return result;
  }

  MediaUploadRequest._();

  factory MediaUploadRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MediaUploadRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MediaUploadRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'retrovibed.cas'),
      createEmptyInstance: create)
    ..aOM<Media>(1, _omitFieldNames ? '' : 'media', subBuilder: Media.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaUploadRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaUploadRequest copyWith(void Function(MediaUploadRequest) updates) =>
      super.copyWith((message) => updates(message as MediaUploadRequest))
          as MediaUploadRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MediaUploadRequest create() => MediaUploadRequest._();
  @$core.override
  MediaUploadRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MediaUploadRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MediaUploadRequest>(create);
  static MediaUploadRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Media get media => $_getN(0);
  @$pb.TagNumber(1)
  set media(Media value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMedia() => $_has(0);
  @$pb.TagNumber(1)
  void clearMedia() => $_clearField(1);
  @$pb.TagNumber(1)
  Media ensureMedia() => $_ensure(0);
}

class MediaUploadResponse extends $pb.GeneratedMessage {
  factory MediaUploadResponse({
    Media? media,
  }) {
    final result = create();
    if (media != null) result.media = media;
    return result;
  }

  MediaUploadResponse._();

  factory MediaUploadResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MediaUploadResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MediaUploadResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'retrovibed.cas'),
      createEmptyInstance: create)
    ..aOM<Media>(1, _omitFieldNames ? '' : 'media', subBuilder: Media.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaUploadResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaUploadResponse copyWith(void Function(MediaUploadResponse) updates) =>
      super.copyWith((message) => updates(message as MediaUploadResponse))
          as MediaUploadResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MediaUploadResponse create() => MediaUploadResponse._();
  @$core.override
  MediaUploadResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MediaUploadResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MediaUploadResponse>(create);
  static MediaUploadResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Media get media => $_getN(0);
  @$pb.TagNumber(1)
  set media(Media value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMedia() => $_has(0);
  @$pb.TagNumber(1)
  void clearMedia() => $_clearField(1);
  @$pb.TagNumber(1)
  Media ensureMedia() => $_ensure(0);
}

class MediaDownloadRequest extends $pb.GeneratedMessage {
  factory MediaDownloadRequest() => create();

  MediaDownloadRequest._();

  factory MediaDownloadRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MediaDownloadRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MediaDownloadRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'retrovibed.cas'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaDownloadRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaDownloadRequest copyWith(void Function(MediaDownloadRequest) updates) =>
      super.copyWith((message) => updates(message as MediaDownloadRequest))
          as MediaDownloadRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MediaDownloadRequest create() => MediaDownloadRequest._();
  @$core.override
  MediaDownloadRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MediaDownloadRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MediaDownloadRequest>(create);
  static MediaDownloadRequest? _defaultInstance;
}

class MediaCompletedRequest extends $pb.GeneratedMessage {
  factory MediaCompletedRequest() => create();

  MediaCompletedRequest._();

  factory MediaCompletedRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MediaCompletedRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MediaCompletedRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'retrovibed.cas'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaCompletedRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaCompletedRequest copyWith(
          void Function(MediaCompletedRequest) updates) =>
      super.copyWith((message) => updates(message as MediaCompletedRequest))
          as MediaCompletedRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MediaCompletedRequest create() => MediaCompletedRequest._();
  @$core.override
  MediaCompletedRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MediaCompletedRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MediaCompletedRequest>(create);
  static MediaCompletedRequest? _defaultInstance;
}

class MediaCompletedResponse extends $pb.GeneratedMessage {
  factory MediaCompletedResponse({
    Media? media,
  }) {
    final result = create();
    if (media != null) result.media = media;
    return result;
  }

  MediaCompletedResponse._();

  factory MediaCompletedResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MediaCompletedResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MediaCompletedResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'retrovibed.cas'),
      createEmptyInstance: create)
    ..aOM<Media>(1, _omitFieldNames ? '' : 'media', subBuilder: Media.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaCompletedResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaCompletedResponse copyWith(
          void Function(MediaCompletedResponse) updates) =>
      super.copyWith((message) => updates(message as MediaCompletedResponse))
          as MediaCompletedResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MediaCompletedResponse create() => MediaCompletedResponse._();
  @$core.override
  MediaCompletedResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MediaCompletedResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MediaCompletedResponse>(create);
  static MediaCompletedResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Media get media => $_getN(0);
  @$pb.TagNumber(1)
  set media(Media value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMedia() => $_has(0);
  @$pb.TagNumber(1)
  void clearMedia() => $_clearField(1);
  @$pb.TagNumber(1)
  Media ensureMedia() => $_ensure(0);
}

class MediaFindRequest extends $pb.GeneratedMessage {
  factory MediaFindRequest() => create();

  MediaFindRequest._();

  factory MediaFindRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MediaFindRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MediaFindRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'retrovibed.cas'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaFindRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaFindRequest copyWith(void Function(MediaFindRequest) updates) =>
      super.copyWith((message) => updates(message as MediaFindRequest))
          as MediaFindRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MediaFindRequest create() => MediaFindRequest._();
  @$core.override
  MediaFindRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MediaFindRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MediaFindRequest>(create);
  static MediaFindRequest? _defaultInstance;
}

class MediaFindResponse extends $pb.GeneratedMessage {
  factory MediaFindResponse({
    Media? media,
  }) {
    final result = create();
    if (media != null) result.media = media;
    return result;
  }

  MediaFindResponse._();

  factory MediaFindResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MediaFindResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MediaFindResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'retrovibed.cas'),
      createEmptyInstance: create)
    ..aOM<Media>(1, _omitFieldNames ? '' : 'media', subBuilder: Media.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaFindResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaFindResponse copyWith(void Function(MediaFindResponse) updates) =>
      super.copyWith((message) => updates(message as MediaFindResponse))
          as MediaFindResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MediaFindResponse create() => MediaFindResponse._();
  @$core.override
  MediaFindResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MediaFindResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MediaFindResponse>(create);
  static MediaFindResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Media get media => $_getN(0);
  @$pb.TagNumber(1)
  set media(Media value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMedia() => $_has(0);
  @$pb.TagNumber(1)
  void clearMedia() => $_clearField(1);
  @$pb.TagNumber(1)
  Media ensureMedia() => $_ensure(0);
}

class MediaDeleteRequest extends $pb.GeneratedMessage {
  factory MediaDeleteRequest() => create();

  MediaDeleteRequest._();

  factory MediaDeleteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MediaDeleteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MediaDeleteRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'retrovibed.cas'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaDeleteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaDeleteRequest copyWith(void Function(MediaDeleteRequest) updates) =>
      super.copyWith((message) => updates(message as MediaDeleteRequest))
          as MediaDeleteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MediaDeleteRequest create() => MediaDeleteRequest._();
  @$core.override
  MediaDeleteRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MediaDeleteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MediaDeleteRequest>(create);
  static MediaDeleteRequest? _defaultInstance;
}

class MediaDeleteResponse extends $pb.GeneratedMessage {
  factory MediaDeleteResponse({
    Media? media,
  }) {
    final result = create();
    if (media != null) result.media = media;
    return result;
  }

  MediaDeleteResponse._();

  factory MediaDeleteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MediaDeleteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MediaDeleteResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'retrovibed.cas'),
      createEmptyInstance: create)
    ..aOM<Media>(1, _omitFieldNames ? '' : 'media', subBuilder: Media.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaDeleteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaDeleteResponse copyWith(void Function(MediaDeleteResponse) updates) =>
      super.copyWith((message) => updates(message as MediaDeleteResponse))
          as MediaDeleteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MediaDeleteResponse create() => MediaDeleteResponse._();
  @$core.override
  MediaDeleteResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MediaDeleteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MediaDeleteResponse>(create);
  static MediaDeleteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Media get media => $_getN(0);
  @$pb.TagNumber(1)
  set media(Media value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMedia() => $_has(0);
  @$pb.TagNumber(1)
  void clearMedia() => $_clearField(1);
  @$pb.TagNumber(1)
  Media ensureMedia() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
