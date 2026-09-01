// This is a generated file - do not edit.
//
// Generated from media/media.proto.

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
    $core.String? description,
    $core.String? mimetype,
    $core.String? image,
    $core.String? archiveId,
    $core.String? torrentId,
    $core.String? createdAt,
    $core.String? updatedAt,
    $core.String? knownMediaId,
    $core.String? encryptionSeed,
    $core.String? uri,
    $core.String? directoryId,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (description != null) result.description = description;
    if (mimetype != null) result.mimetype = mimetype;
    if (image != null) result.image = image;
    if (archiveId != null) result.archiveId = archiveId;
    if (torrentId != null) result.torrentId = torrentId;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (knownMediaId != null) result.knownMediaId = knownMediaId;
    if (encryptionSeed != null) result.encryptionSeed = encryptionSeed;
    if (uri != null) result.uri = uri;
    if (directoryId != null) result.directoryId = directoryId;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'description')
    ..aOS(3, _omitFieldNames ? '' : 'mimetype')
    ..aOS(4, _omitFieldNames ? '' : 'image')
    ..aOS(5, _omitFieldNames ? '' : 'archive_id')
    ..aOS(6, _omitFieldNames ? '' : 'torrent_id')
    ..aOS(7, _omitFieldNames ? '' : 'created_at')
    ..aOS(8, _omitFieldNames ? '' : 'updated_at')
    ..aOS(9, _omitFieldNames ? '' : 'known_media_id')
    ..aOS(10, _omitFieldNames ? '' : 'encryption_seed')
    ..aOS(11, _omitFieldNames ? '' : 'uri')
    ..aOS(12, _omitFieldNames ? '' : 'directory_id')
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
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(2)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(2)
  void clearDescription() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get mimetype => $_getSZ(2);
  @$pb.TagNumber(3)
  set mimetype($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasMimetype() => $_has(2);
  @$pb.TagNumber(3)
  void clearMimetype() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get image => $_getSZ(3);
  @$pb.TagNumber(4)
  set image($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasImage() => $_has(3);
  @$pb.TagNumber(4)
  void clearImage() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get archiveId => $_getSZ(4);
  @$pb.TagNumber(5)
  set archiveId($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasArchiveId() => $_has(4);
  @$pb.TagNumber(5)
  void clearArchiveId() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get torrentId => $_getSZ(5);
  @$pb.TagNumber(6)
  set torrentId($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasTorrentId() => $_has(5);
  @$pb.TagNumber(6)
  void clearTorrentId() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.String get createdAt => $_getSZ(6);
  @$pb.TagNumber(7)
  set createdAt($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasCreatedAt() => $_has(6);
  @$pb.TagNumber(7)
  void clearCreatedAt() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.String get updatedAt => $_getSZ(7);
  @$pb.TagNumber(8)
  set updatedAt($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasUpdatedAt() => $_has(7);
  @$pb.TagNumber(8)
  void clearUpdatedAt() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.String get knownMediaId => $_getSZ(8);
  @$pb.TagNumber(9)
  set knownMediaId($core.String value) => $_setString(8, value);
  @$pb.TagNumber(9)
  $core.bool hasKnownMediaId() => $_has(8);
  @$pb.TagNumber(9)
  void clearKnownMediaId() => $_clearField(9);

  @$pb.TagNumber(10)
  $core.String get encryptionSeed => $_getSZ(9);
  @$pb.TagNumber(10)
  set encryptionSeed($core.String value) => $_setString(9, value);
  @$pb.TagNumber(10)
  $core.bool hasEncryptionSeed() => $_has(9);
  @$pb.TagNumber(10)
  void clearEncryptionSeed() => $_clearField(10);

  @$pb.TagNumber(11)
  $core.String get uri => $_getSZ(10);
  @$pb.TagNumber(11)
  set uri($core.String value) => $_setString(10, value);
  @$pb.TagNumber(11)
  $core.bool hasUri() => $_has(10);
  @$pb.TagNumber(11)
  void clearUri() => $_clearField(11);

  @$pb.TagNumber(12)
  $core.String get directoryId => $_getSZ(11);
  @$pb.TagNumber(12)
  set directoryId($core.String value) => $_setString(11, value);
  @$pb.TagNumber(12)
  $core.bool hasDirectoryId() => $_has(11);
  @$pb.TagNumber(12)
  void clearDirectoryId() => $_clearField(12);
}

class MediaSearchRequest extends $pb.GeneratedMessage {
  factory MediaSearchRequest({
    $core.String? query,
    $core.Iterable<$core.String>? mimetypes,
    $core.bool? adult,
    $core.bool? hidden,
    $core.Iterable<$core.String>? excluded,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = create();
    if (query != null) result.query = query;
    if (mimetypes != null) result.mimetypes.addAll(mimetypes);
    if (adult != null) result.adult = adult;
    if (hidden != null) result.hidden = hidden;
    if (excluded != null) result.excluded.addAll(excluded);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..pPS(2, _omitFieldNames ? '' : 'mimetypes')
    ..aOB(3, _omitFieldNames ? '' : 'adult')
    ..aOB(4, _omitFieldNames ? '' : 'hidden')
    ..pPS(5, _omitFieldNames ? '' : 'excluded')
    ..a<$fixnum.Int64>(
        900, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(901, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
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
  $pb.PbList<$core.String> get mimetypes => $_getList(1);

  @$pb.TagNumber(3)
  $core.bool get adult => $_getBF(2);
  @$pb.TagNumber(3)
  set adult($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(3)
  $core.bool hasAdult() => $_has(2);
  @$pb.TagNumber(3)
  void clearAdult() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.bool get hidden => $_getBF(3);
  @$pb.TagNumber(4)
  set hidden($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(4)
  $core.bool hasHidden() => $_has(3);
  @$pb.TagNumber(4)
  void clearHidden() => $_clearField(4);

  @$pb.TagNumber(5)
  $pb.PbList<$core.String> get excluded => $_getList(4);

  @$pb.TagNumber(900)
  $fixnum.Int64 get offset => $_getI64(5);
  @$pb.TagNumber(900)
  set offset($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(900)
  $core.bool hasOffset() => $_has(5);
  @$pb.TagNumber(900)
  void clearOffset() => $_clearField(900);

  @$pb.TagNumber(901)
  $fixnum.Int64 get limit => $_getI64(6);
  @$pb.TagNumber(901)
  set limit($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(901)
  $core.bool hasLimit() => $_has(6);
  @$pb.TagNumber(901)
  void clearLimit() => $_clearField(901);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
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

class MediaUpdateRequest extends $pb.GeneratedMessage {
  factory MediaUpdateRequest({
    Media? media,
  }) {
    final result = create();
    if (media != null) result.media = media;
    return result;
  }

  MediaUpdateRequest._();

  factory MediaUpdateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MediaUpdateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MediaUpdateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Media>(1, _omitFieldNames ? '' : 'media', subBuilder: Media.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaUpdateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaUpdateRequest copyWith(void Function(MediaUpdateRequest) updates) =>
      super.copyWith((message) => updates(message as MediaUpdateRequest))
          as MediaUpdateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MediaUpdateRequest create() => MediaUpdateRequest._();
  @$core.override
  MediaUpdateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MediaUpdateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MediaUpdateRequest>(create);
  static MediaUpdateRequest? _defaultInstance;

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

class MediaUpdateResponse extends $pb.GeneratedMessage {
  factory MediaUpdateResponse({
    Media? media,
  }) {
    final result = create();
    if (media != null) result.media = media;
    return result;
  }

  MediaUpdateResponse._();

  factory MediaUpdateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MediaUpdateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MediaUpdateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Media>(1, _omitFieldNames ? '' : 'media', subBuilder: Media.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaUpdateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MediaUpdateResponse copyWith(void Function(MediaUpdateResponse) updates) =>
      super.copyWith((message) => updates(message as MediaUpdateResponse))
          as MediaUpdateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MediaUpdateResponse create() => MediaUpdateResponse._();
  @$core.override
  MediaUpdateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MediaUpdateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MediaUpdateResponse>(create);
  static MediaUpdateResponse? _defaultInstance;

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
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

class Download extends $pb.GeneratedMessage {
  factory Download({
    Media? media,
    $fixnum.Int64? bytes,
    $fixnum.Int64? downloaded,
    $core.String? initiatedAt,
    $core.String? pausedAt,
    $core.int? peers,
    $core.bool? distributing,
    $core.String? path,
    $core.int? peersSeeders,
    $core.int? peersHalf,
    $core.int? peersAvailable,
    $core.String? completedAt,
    $core.String? verifyAt,
  }) {
    final result = create();
    if (media != null) result.media = media;
    if (bytes != null) result.bytes = bytes;
    if (downloaded != null) result.downloaded = downloaded;
    if (initiatedAt != null) result.initiatedAt = initiatedAt;
    if (pausedAt != null) result.pausedAt = pausedAt;
    if (peers != null) result.peers = peers;
    if (distributing != null) result.distributing = distributing;
    if (path != null) result.path = path;
    if (peersSeeders != null) result.peersSeeders = peersSeeders;
    if (peersHalf != null) result.peersHalf = peersHalf;
    if (peersAvailable != null) result.peersAvailable = peersAvailable;
    if (completedAt != null) result.completedAt = completedAt;
    if (verifyAt != null) result.verifyAt = verifyAt;
    return result;
  }

  Download._();

  factory Download.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Download.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Download',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Media>(1, _omitFieldNames ? '' : 'media', subBuilder: Media.create)
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'bytes', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        3, _omitFieldNames ? '' : 'downloaded', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOS(4, _omitFieldNames ? '' : 'initiated_at')
    ..aOS(5, _omitFieldNames ? '' : 'paused_at')
    ..aI(6, _omitFieldNames ? '' : 'peers', fieldType: $pb.PbFieldType.OU3)
    ..aOB(7, _omitFieldNames ? '' : 'distributing')
    ..aOS(8, _omitFieldNames ? '' : 'path')
    ..aI(9, _omitFieldNames ? '' : 'peers_seeders',
        fieldType: $pb.PbFieldType.OU3)
    ..aI(10, _omitFieldNames ? '' : 'peers_half',
        fieldType: $pb.PbFieldType.OU3)
    ..aI(11, _omitFieldNames ? '' : 'peers_available',
        fieldType: $pb.PbFieldType.OU3)
    ..aOS(12, _omitFieldNames ? '' : 'completed_at')
    ..aOS(13, _omitFieldNames ? '' : 'verify_at')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Download clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Download copyWith(void Function(Download) updates) =>
      super.copyWith((message) => updates(message as Download)) as Download;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Download create() => Download._();
  @$core.override
  Download createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Download getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Download>(create);
  static Download? _defaultInstance;

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

  @$pb.TagNumber(2)
  $fixnum.Int64 get bytes => $_getI64(1);
  @$pb.TagNumber(2)
  set bytes($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasBytes() => $_has(1);
  @$pb.TagNumber(2)
  void clearBytes() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get downloaded => $_getI64(2);
  @$pb.TagNumber(3)
  set downloaded($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasDownloaded() => $_has(2);
  @$pb.TagNumber(3)
  void clearDownloaded() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get initiatedAt => $_getSZ(3);
  @$pb.TagNumber(4)
  set initiatedAt($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasInitiatedAt() => $_has(3);
  @$pb.TagNumber(4)
  void clearInitiatedAt() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get pausedAt => $_getSZ(4);
  @$pb.TagNumber(5)
  set pausedAt($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasPausedAt() => $_has(4);
  @$pb.TagNumber(5)
  void clearPausedAt() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.int get peers => $_getIZ(5);
  @$pb.TagNumber(6)
  set peers($core.int value) => $_setUnsignedInt32(5, value);
  @$pb.TagNumber(6)
  $core.bool hasPeers() => $_has(5);
  @$pb.TagNumber(6)
  void clearPeers() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.bool get distributing => $_getBF(6);
  @$pb.TagNumber(7)
  set distributing($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(7)
  $core.bool hasDistributing() => $_has(6);
  @$pb.TagNumber(7)
  void clearDistributing() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.String get path => $_getSZ(7);
  @$pb.TagNumber(8)
  set path($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasPath() => $_has(7);
  @$pb.TagNumber(8)
  void clearPath() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.int get peersSeeders => $_getIZ(8);
  @$pb.TagNumber(9)
  set peersSeeders($core.int value) => $_setUnsignedInt32(8, value);
  @$pb.TagNumber(9)
  $core.bool hasPeersSeeders() => $_has(8);
  @$pb.TagNumber(9)
  void clearPeersSeeders() => $_clearField(9);

  @$pb.TagNumber(10)
  $core.int get peersHalf => $_getIZ(9);
  @$pb.TagNumber(10)
  set peersHalf($core.int value) => $_setUnsignedInt32(9, value);
  @$pb.TagNumber(10)
  $core.bool hasPeersHalf() => $_has(9);
  @$pb.TagNumber(10)
  void clearPeersHalf() => $_clearField(10);

  @$pb.TagNumber(11)
  $core.int get peersAvailable => $_getIZ(10);
  @$pb.TagNumber(11)
  set peersAvailable($core.int value) => $_setUnsignedInt32(10, value);
  @$pb.TagNumber(11)
  $core.bool hasPeersAvailable() => $_has(10);
  @$pb.TagNumber(11)
  void clearPeersAvailable() => $_clearField(11);

  @$pb.TagNumber(12)
  $core.String get completedAt => $_getSZ(11);
  @$pb.TagNumber(12)
  set completedAt($core.String value) => $_setString(11, value);
  @$pb.TagNumber(12)
  $core.bool hasCompletedAt() => $_has(11);
  @$pb.TagNumber(12)
  void clearCompletedAt() => $_clearField(12);

  @$pb.TagNumber(13)
  $core.String get verifyAt => $_getSZ(12);
  @$pb.TagNumber(13)
  set verifyAt($core.String value) => $_setString(12, value);
  @$pb.TagNumber(13)
  $core.bool hasVerifyAt() => $_has(12);
  @$pb.TagNumber(13)
  void clearVerifyAt() => $_clearField(13);
}

class MagnetCreateRequest extends $pb.GeneratedMessage {
  factory MagnetCreateRequest({
    $core.String? uri,
  }) {
    final result = create();
    if (uri != null) result.uri = uri;
    return result;
  }

  MagnetCreateRequest._();

  factory MagnetCreateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MagnetCreateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MagnetCreateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'uri')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MagnetCreateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MagnetCreateRequest copyWith(void Function(MagnetCreateRequest) updates) =>
      super.copyWith((message) => updates(message as MagnetCreateRequest))
          as MagnetCreateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MagnetCreateRequest create() => MagnetCreateRequest._();
  @$core.override
  MagnetCreateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MagnetCreateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MagnetCreateRequest>(create);
  static MagnetCreateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get uri => $_getSZ(0);
  @$pb.TagNumber(1)
  set uri($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasUri() => $_has(0);
  @$pb.TagNumber(1)
  void clearUri() => $_clearField(1);
}

class MagnetCreateResponse extends $pb.GeneratedMessage {
  factory MagnetCreateResponse({
    Download? download,
  }) {
    final result = create();
    if (download != null) result.download = download;
    return result;
  }

  MagnetCreateResponse._();

  factory MagnetCreateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MagnetCreateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MagnetCreateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Download>(1, _omitFieldNames ? '' : 'download',
        subBuilder: Download.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MagnetCreateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MagnetCreateResponse copyWith(void Function(MagnetCreateResponse) updates) =>
      super.copyWith((message) => updates(message as MagnetCreateResponse))
          as MagnetCreateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MagnetCreateResponse create() => MagnetCreateResponse._();
  @$core.override
  MagnetCreateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MagnetCreateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MagnetCreateResponse>(create);
  static MagnetCreateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Download get download => $_getN(0);
  @$pb.TagNumber(1)
  set download(Download value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDownload() => $_has(0);
  @$pb.TagNumber(1)
  void clearDownload() => $_clearField(1);
  @$pb.TagNumber(1)
  Download ensureDownload() => $_ensure(0);
}

class DownloadSearchRequest extends $pb.GeneratedMessage {
  factory DownloadSearchRequest({
    $core.String? query,
    $core.Iterable<$core.String>? mimetypes,
    $core.bool? completed,
    $core.bool? hidden,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = create();
    if (query != null) result.query = query;
    if (mimetypes != null) result.mimetypes.addAll(mimetypes);
    if (completed != null) result.completed = completed;
    if (hidden != null) result.hidden = hidden;
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    return result;
  }

  DownloadSearchRequest._();

  factory DownloadSearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DownloadSearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DownloadSearchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..pPS(2, _omitFieldNames ? '' : 'mimetypes')
    ..aOB(3, _omitFieldNames ? '' : 'completed')
    ..aOB(4, _omitFieldNames ? '' : 'hidden')
    ..a<$fixnum.Int64>(
        900, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(901, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadSearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadSearchRequest copyWith(
          void Function(DownloadSearchRequest) updates) =>
      super.copyWith((message) => updates(message as DownloadSearchRequest))
          as DownloadSearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadSearchRequest create() => DownloadSearchRequest._();
  @$core.override
  DownloadSearchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DownloadSearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DownloadSearchRequest>(create);
  static DownloadSearchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get query => $_getSZ(0);
  @$pb.TagNumber(1)
  set query($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasQuery() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuery() => $_clearField(1);

  @$pb.TagNumber(2)
  $pb.PbList<$core.String> get mimetypes => $_getList(1);

  @$pb.TagNumber(3)
  $core.bool get completed => $_getBF(2);
  @$pb.TagNumber(3)
  set completed($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(3)
  $core.bool hasCompleted() => $_has(2);
  @$pb.TagNumber(3)
  void clearCompleted() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.bool get hidden => $_getBF(3);
  @$pb.TagNumber(4)
  set hidden($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(4)
  $core.bool hasHidden() => $_has(3);
  @$pb.TagNumber(4)
  void clearHidden() => $_clearField(4);

  @$pb.TagNumber(900)
  $fixnum.Int64 get offset => $_getI64(4);
  @$pb.TagNumber(900)
  set offset($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(900)
  $core.bool hasOffset() => $_has(4);
  @$pb.TagNumber(900)
  void clearOffset() => $_clearField(900);

  @$pb.TagNumber(901)
  $fixnum.Int64 get limit => $_getI64(5);
  @$pb.TagNumber(901)
  set limit($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(901)
  $core.bool hasLimit() => $_has(5);
  @$pb.TagNumber(901)
  void clearLimit() => $_clearField(901);
}

class DownloadSearchResponse extends $pb.GeneratedMessage {
  factory DownloadSearchResponse({
    DownloadSearchRequest? next,
    $core.Iterable<Download>? items,
  }) {
    final result = create();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  DownloadSearchResponse._();

  factory DownloadSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DownloadSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DownloadSearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<DownloadSearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: DownloadSearchRequest.create)
    ..pPM<Download>(2, _omitFieldNames ? '' : 'items',
        subBuilder: Download.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadSearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadSearchResponse copyWith(
          void Function(DownloadSearchResponse) updates) =>
      super.copyWith((message) => updates(message as DownloadSearchResponse))
          as DownloadSearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadSearchResponse create() => DownloadSearchResponse._();
  @$core.override
  DownloadSearchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DownloadSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DownloadSearchResponse>(create);
  static DownloadSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  DownloadSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(DownloadSearchRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  DownloadSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<Download> get items => $_getList(1);
}

class DownloadUpdateRequest extends $pb.GeneratedMessage {
  factory DownloadUpdateRequest({
    Download? download,
  }) {
    final result = create();
    if (download != null) result.download = download;
    return result;
  }

  DownloadUpdateRequest._();

  factory DownloadUpdateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DownloadUpdateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DownloadUpdateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Download>(1, _omitFieldNames ? '' : 'download',
        subBuilder: Download.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadUpdateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadUpdateRequest copyWith(
          void Function(DownloadUpdateRequest) updates) =>
      super.copyWith((message) => updates(message as DownloadUpdateRequest))
          as DownloadUpdateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadUpdateRequest create() => DownloadUpdateRequest._();
  @$core.override
  DownloadUpdateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DownloadUpdateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DownloadUpdateRequest>(create);
  static DownloadUpdateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Download get download => $_getN(0);
  @$pb.TagNumber(1)
  set download(Download value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDownload() => $_has(0);
  @$pb.TagNumber(1)
  void clearDownload() => $_clearField(1);
  @$pb.TagNumber(1)
  Download ensureDownload() => $_ensure(0);
}

class DownloadUpdateResponse extends $pb.GeneratedMessage {
  factory DownloadUpdateResponse({
    Download? download,
  }) {
    final result = create();
    if (download != null) result.download = download;
    return result;
  }

  DownloadUpdateResponse._();

  factory DownloadUpdateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DownloadUpdateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DownloadUpdateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Download>(1, _omitFieldNames ? '' : 'download',
        subBuilder: Download.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadUpdateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadUpdateResponse copyWith(
          void Function(DownloadUpdateResponse) updates) =>
      super.copyWith((message) => updates(message as DownloadUpdateResponse))
          as DownloadUpdateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadUpdateResponse create() => DownloadUpdateResponse._();
  @$core.override
  DownloadUpdateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DownloadUpdateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DownloadUpdateResponse>(create);
  static DownloadUpdateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Download get download => $_getN(0);
  @$pb.TagNumber(1)
  set download(Download value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDownload() => $_has(0);
  @$pb.TagNumber(1)
  void clearDownload() => $_clearField(1);
  @$pb.TagNumber(1)
  Download ensureDownload() => $_ensure(0);
}

class DownloadMetadataRequest extends $pb.GeneratedMessage {
  factory DownloadMetadataRequest() => create();

  DownloadMetadataRequest._();

  factory DownloadMetadataRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DownloadMetadataRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DownloadMetadataRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadMetadataRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadMetadataRequest copyWith(
          void Function(DownloadMetadataRequest) updates) =>
      super.copyWith((message) => updates(message as DownloadMetadataRequest))
          as DownloadMetadataRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadMetadataRequest create() => DownloadMetadataRequest._();
  @$core.override
  DownloadMetadataRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DownloadMetadataRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DownloadMetadataRequest>(create);
  static DownloadMetadataRequest? _defaultInstance;
}

class DownloadMetadataResponse extends $pb.GeneratedMessage {
  factory DownloadMetadataResponse({
    Download? download,
  }) {
    final result = create();
    if (download != null) result.download = download;
    return result;
  }

  DownloadMetadataResponse._();

  factory DownloadMetadataResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DownloadMetadataResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DownloadMetadataResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Download>(1, _omitFieldNames ? '' : 'download',
        subBuilder: Download.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadMetadataResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadMetadataResponse copyWith(
          void Function(DownloadMetadataResponse) updates) =>
      super.copyWith((message) => updates(message as DownloadMetadataResponse))
          as DownloadMetadataResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadMetadataResponse create() => DownloadMetadataResponse._();
  @$core.override
  DownloadMetadataResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DownloadMetadataResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DownloadMetadataResponse>(create);
  static DownloadMetadataResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Download get download => $_getN(0);
  @$pb.TagNumber(1)
  set download(Download value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDownload() => $_has(0);
  @$pb.TagNumber(1)
  void clearDownload() => $_clearField(1);
  @$pb.TagNumber(1)
  Download ensureDownload() => $_ensure(0);
}

class DownloadBeginRequest extends $pb.GeneratedMessage {
  factory DownloadBeginRequest() => create();

  DownloadBeginRequest._();

  factory DownloadBeginRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DownloadBeginRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DownloadBeginRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadBeginRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadBeginRequest copyWith(void Function(DownloadBeginRequest) updates) =>
      super.copyWith((message) => updates(message as DownloadBeginRequest))
          as DownloadBeginRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadBeginRequest create() => DownloadBeginRequest._();
  @$core.override
  DownloadBeginRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DownloadBeginRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DownloadBeginRequest>(create);
  static DownloadBeginRequest? _defaultInstance;
}

class DownloadBeginResponse extends $pb.GeneratedMessage {
  factory DownloadBeginResponse({
    Download? download,
  }) {
    final result = create();
    if (download != null) result.download = download;
    return result;
  }

  DownloadBeginResponse._();

  factory DownloadBeginResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DownloadBeginResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DownloadBeginResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Download>(1, _omitFieldNames ? '' : 'download',
        subBuilder: Download.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadBeginResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadBeginResponse copyWith(
          void Function(DownloadBeginResponse) updates) =>
      super.copyWith((message) => updates(message as DownloadBeginResponse))
          as DownloadBeginResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadBeginResponse create() => DownloadBeginResponse._();
  @$core.override
  DownloadBeginResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DownloadBeginResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DownloadBeginResponse>(create);
  static DownloadBeginResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Download get download => $_getN(0);
  @$pb.TagNumber(1)
  set download(Download value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDownload() => $_has(0);
  @$pb.TagNumber(1)
  void clearDownload() => $_clearField(1);
  @$pb.TagNumber(1)
  Download ensureDownload() => $_ensure(0);
}

class DownloadPauseRequest extends $pb.GeneratedMessage {
  factory DownloadPauseRequest() => create();

  DownloadPauseRequest._();

  factory DownloadPauseRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DownloadPauseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DownloadPauseRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadPauseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadPauseRequest copyWith(void Function(DownloadPauseRequest) updates) =>
      super.copyWith((message) => updates(message as DownloadPauseRequest))
          as DownloadPauseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadPauseRequest create() => DownloadPauseRequest._();
  @$core.override
  DownloadPauseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DownloadPauseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DownloadPauseRequest>(create);
  static DownloadPauseRequest? _defaultInstance;
}

class DownloadPauseResponse extends $pb.GeneratedMessage {
  factory DownloadPauseResponse({
    Download? download,
  }) {
    final result = create();
    if (download != null) result.download = download;
    return result;
  }

  DownloadPauseResponse._();

  factory DownloadPauseResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DownloadPauseResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DownloadPauseResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Download>(1, _omitFieldNames ? '' : 'download',
        subBuilder: Download.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadPauseResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadPauseResponse copyWith(
          void Function(DownloadPauseResponse) updates) =>
      super.copyWith((message) => updates(message as DownloadPauseResponse))
          as DownloadPauseResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadPauseResponse create() => DownloadPauseResponse._();
  @$core.override
  DownloadPauseResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DownloadPauseResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DownloadPauseResponse>(create);
  static DownloadPauseResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Download get download => $_getN(0);
  @$pb.TagNumber(1)
  set download(Download value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDownload() => $_has(0);
  @$pb.TagNumber(1)
  void clearDownload() => $_clearField(1);
  @$pb.TagNumber(1)
  Download ensureDownload() => $_ensure(0);
}

class DownloadTuneRequest extends $pb.GeneratedMessage {
  factory DownloadTuneRequest({
    $core.Iterable<$core.String>? peers,
  }) {
    final result = create();
    if (peers != null) result.peers.addAll(peers);
    return result;
  }

  DownloadTuneRequest._();

  factory DownloadTuneRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DownloadTuneRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DownloadTuneRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..pPS(1, _omitFieldNames ? '' : 'peers')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadTuneRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadTuneRequest copyWith(void Function(DownloadTuneRequest) updates) =>
      super.copyWith((message) => updates(message as DownloadTuneRequest))
          as DownloadTuneRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadTuneRequest create() => DownloadTuneRequest._();
  @$core.override
  DownloadTuneRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DownloadTuneRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DownloadTuneRequest>(create);
  static DownloadTuneRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<$core.String> get peers => $_getList(0);
}

class DownloadTuneResponse extends $pb.GeneratedMessage {
  factory DownloadTuneResponse() => create();

  DownloadTuneResponse._();

  factory DownloadTuneResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DownloadTuneResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DownloadTuneResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadTuneResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadTuneResponse copyWith(void Function(DownloadTuneResponse) updates) =>
      super.copyWith((message) => updates(message as DownloadTuneResponse))
          as DownloadTuneResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadTuneResponse create() => DownloadTuneResponse._();
  @$core.override
  DownloadTuneResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DownloadTuneResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DownloadTuneResponse>(create);
  static DownloadTuneResponse? _defaultInstance;
}

class DownloadDeleteRequest extends $pb.GeneratedMessage {
  factory DownloadDeleteRequest() => create();

  DownloadDeleteRequest._();

  factory DownloadDeleteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DownloadDeleteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DownloadDeleteRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadDeleteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadDeleteRequest copyWith(
          void Function(DownloadDeleteRequest) updates) =>
      super.copyWith((message) => updates(message as DownloadDeleteRequest))
          as DownloadDeleteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadDeleteRequest create() => DownloadDeleteRequest._();
  @$core.override
  DownloadDeleteRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DownloadDeleteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DownloadDeleteRequest>(create);
  static DownloadDeleteRequest? _defaultInstance;
}

class DownloadDeleteResponse extends $pb.GeneratedMessage {
  factory DownloadDeleteResponse({
    Download? download,
  }) {
    final result = create();
    if (download != null) result.download = download;
    return result;
  }

  DownloadDeleteResponse._();

  factory DownloadDeleteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DownloadDeleteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DownloadDeleteResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Download>(1, _omitFieldNames ? '' : 'download',
        subBuilder: Download.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadDeleteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DownloadDeleteResponse copyWith(
          void Function(DownloadDeleteResponse) updates) =>
      super.copyWith((message) => updates(message as DownloadDeleteResponse))
          as DownloadDeleteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadDeleteResponse create() => DownloadDeleteResponse._();
  @$core.override
  DownloadDeleteResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DownloadDeleteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DownloadDeleteResponse>(create);
  static DownloadDeleteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Download get download => $_getN(0);
  @$pb.TagNumber(1)
  set download(Download value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDownload() => $_has(0);
  @$pb.TagNumber(1)
  void clearDownload() => $_clearField(1);
  @$pb.TagNumber(1)
  Download ensureDownload() => $_ensure(0);
}

class Published extends $pb.GeneratedMessage {
  factory Published({
    $core.String? id,
    $core.String? mimetype,
    $core.String? description,
    $fixnum.Int64? bytes,
    $core.String? entropy,
    $core.String? expiresAt,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (mimetype != null) result.mimetype = mimetype;
    if (description != null) result.description = description;
    if (bytes != null) result.bytes = bytes;
    if (entropy != null) result.entropy = entropy;
    if (expiresAt != null) result.expiresAt = expiresAt;
    return result;
  }

  Published._();

  factory Published.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Published.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Published',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'mimetype')
    ..aOS(3, _omitFieldNames ? '' : 'description')
    ..a<$fixnum.Int64>(4, _omitFieldNames ? '' : 'bytes', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOS(5, _omitFieldNames ? '' : 'entropy')
    ..aOS(6, _omitFieldNames ? '' : 'expires_at')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Published clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Published copyWith(void Function(Published) updates) =>
      super.copyWith((message) => updates(message as Published)) as Published;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Published create() => Published._();
  @$core.override
  Published createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Published getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Published>(create);
  static Published? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get mimetype => $_getSZ(1);
  @$pb.TagNumber(2)
  set mimetype($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasMimetype() => $_has(1);
  @$pb.TagNumber(2)
  void clearMimetype() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(3)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(3)
  void clearDescription() => $_clearField(3);

  @$pb.TagNumber(4)
  $fixnum.Int64 get bytes => $_getI64(3);
  @$pb.TagNumber(4)
  set bytes($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(4)
  $core.bool hasBytes() => $_has(3);
  @$pb.TagNumber(4)
  void clearBytes() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get entropy => $_getSZ(4);
  @$pb.TagNumber(5)
  set entropy($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasEntropy() => $_has(4);
  @$pb.TagNumber(5)
  void clearEntropy() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get expiresAt => $_getSZ(5);
  @$pb.TagNumber(6)
  set expiresAt($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasExpiresAt() => $_has(5);
  @$pb.TagNumber(6)
  void clearExpiresAt() => $_clearField(6);
}

class PublishedUploadRequest extends $pb.GeneratedMessage {
  factory PublishedUploadRequest({
    $core.String? entropy,
    $core.String? mimetype,
    $fixnum.Int64? ttl,
  }) {
    final result = create();
    if (entropy != null) result.entropy = entropy;
    if (mimetype != null) result.mimetype = mimetype;
    if (ttl != null) result.ttl = ttl;
    return result;
  }

  PublishedUploadRequest._();

  factory PublishedUploadRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishedUploadRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishedUploadRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'entropy')
    ..aOS(2, _omitFieldNames ? '' : 'mimetype')
    ..a<$fixnum.Int64>(3, _omitFieldNames ? '' : 'ttl', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedUploadRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedUploadRequest copyWith(
          void Function(PublishedUploadRequest) updates) =>
      super.copyWith((message) => updates(message as PublishedUploadRequest))
          as PublishedUploadRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishedUploadRequest create() => PublishedUploadRequest._();
  @$core.override
  PublishedUploadRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishedUploadRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishedUploadRequest>(create);
  static PublishedUploadRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get entropy => $_getSZ(0);
  @$pb.TagNumber(1)
  set entropy($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasEntropy() => $_has(0);
  @$pb.TagNumber(1)
  void clearEntropy() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get mimetype => $_getSZ(1);
  @$pb.TagNumber(2)
  set mimetype($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasMimetype() => $_has(1);
  @$pb.TagNumber(2)
  void clearMimetype() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get ttl => $_getI64(2);
  @$pb.TagNumber(3)
  set ttl($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasTtl() => $_has(2);
  @$pb.TagNumber(3)
  void clearTtl() => $_clearField(3);
}

class PublishedUploadResponse extends $pb.GeneratedMessage {
  factory PublishedUploadResponse({
    Published? published,
  }) {
    final result = create();
    if (published != null) result.published = published;
    return result;
  }

  PublishedUploadResponse._();

  factory PublishedUploadResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishedUploadResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishedUploadResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Published>(1, _omitFieldNames ? '' : 'published',
        subBuilder: Published.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedUploadResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedUploadResponse copyWith(
          void Function(PublishedUploadResponse) updates) =>
      super.copyWith((message) => updates(message as PublishedUploadResponse))
          as PublishedUploadResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishedUploadResponse create() => PublishedUploadResponse._();
  @$core.override
  PublishedUploadResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishedUploadResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishedUploadResponse>(create);
  static PublishedUploadResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Published get published => $_getN(0);
  @$pb.TagNumber(1)
  set published(Published value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPublished() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublished() => $_clearField(1);
  @$pb.TagNumber(1)
  Published ensurePublished() => $_ensure(0);
}

class MetadataSyncRequest extends $pb.GeneratedMessage {
  factory MetadataSyncRequest({
    Media? media,
  }) {
    final result = create();
    if (media != null) result.media = media;
    return result;
  }

  MetadataSyncRequest._();

  factory MetadataSyncRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetadataSyncRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetadataSyncRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Media>(1, _omitFieldNames ? '' : 'media', subBuilder: Media.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetadataSyncRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetadataSyncRequest copyWith(void Function(MetadataSyncRequest) updates) =>
      super.copyWith((message) => updates(message as MetadataSyncRequest))
          as MetadataSyncRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetadataSyncRequest create() => MetadataSyncRequest._();
  @$core.override
  MetadataSyncRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetadataSyncRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetadataSyncRequest>(create);
  static MetadataSyncRequest? _defaultInstance;

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

class MetadataSyncResponse extends $pb.GeneratedMessage {
  factory MetadataSyncResponse({
    Media? media,
  }) {
    final result = create();
    if (media != null) result.media = media;
    return result;
  }

  MetadataSyncResponse._();

  factory MetadataSyncResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetadataSyncResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetadataSyncResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Media>(1, _omitFieldNames ? '' : 'media', subBuilder: Media.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetadataSyncResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetadataSyncResponse copyWith(void Function(MetadataSyncResponse) updates) =>
      super.copyWith((message) => updates(message as MetadataSyncResponse))
          as MetadataSyncResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetadataSyncResponse create() => MetadataSyncResponse._();
  @$core.override
  MetadataSyncResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetadataSyncResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetadataSyncResponse>(create);
  static MetadataSyncResponse? _defaultInstance;

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

class PublishedRequest extends $pb.GeneratedMessage {
  factory PublishedRequest({
    $core.String? knownMediaId,
  }) {
    final result = create();
    if (knownMediaId != null) result.knownMediaId = knownMediaId;
    return result;
  }

  PublishedRequest._();

  factory PublishedRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishedRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishedRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'known_media_id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedRequest copyWith(void Function(PublishedRequest) updates) =>
      super.copyWith((message) => updates(message as PublishedRequest))
          as PublishedRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishedRequest create() => PublishedRequest._();
  @$core.override
  PublishedRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishedRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishedRequest>(create);
  static PublishedRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get knownMediaId => $_getSZ(0);
  @$pb.TagNumber(1)
  set knownMediaId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasKnownMediaId() => $_has(0);
  @$pb.TagNumber(1)
  void clearKnownMediaId() => $_clearField(1);
}

class PublishedResponse extends $pb.GeneratedMessage {
  factory PublishedResponse({
    Media? media,
    $core.String? magnetUri,
  }) {
    final result = create();
    if (media != null) result.media = media;
    if (magnetUri != null) result.magnetUri = magnetUri;
    return result;
  }

  PublishedResponse._();

  factory PublishedResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishedResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishedResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Media>(1, _omitFieldNames ? '' : 'media', subBuilder: Media.create)
    ..aOS(2, _omitFieldNames ? '' : 'magnet_uri')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedResponse copyWith(void Function(PublishedResponse) updates) =>
      super.copyWith((message) => updates(message as PublishedResponse))
          as PublishedResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishedResponse create() => PublishedResponse._();
  @$core.override
  PublishedResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishedResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishedResponse>(create);
  static PublishedResponse? _defaultInstance;

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

  @$pb.TagNumber(2)
  $core.String get magnetUri => $_getSZ(1);
  @$pb.TagNumber(2)
  set magnetUri($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasMagnetUri() => $_has(1);
  @$pb.TagNumber(2)
  void clearMagnetUri() => $_clearField(2);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
