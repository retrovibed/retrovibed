//
//  Generated code. Do not modify.
//  source: media.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

class Media extends $pb.GeneratedMessage {
  factory Media({
    $core.String? id,
    $core.String? description,
    $core.String? mimetype,
    $core.String? image,
  }) {
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    if (description != null) {
      $result.description = description;
    }
    if (mimetype != null) {
      $result.mimetype = mimetype;
    }
    if (image != null) {
      $result.image = image;
    }
    return $result;
  }
  Media._() : super();
  factory Media.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory Media.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'Media', package: const $pb.PackageName(_omitMessageNames ? '' : 'media'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'description')
    ..aOS(3, _omitFieldNames ? '' : 'mimetype')
    ..aOS(4, _omitFieldNames ? '' : 'image')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  Media clone() => Media()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  Media copyWith(void Function(Media) updates) => super.copyWith((message) => updates(message as Media)) as Media;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Media create() => Media._();
  Media createEmptyInstance() => create();
  static $pb.PbList<Media> createRepeated() => $pb.PbList<Media>();
  @$core.pragma('dart2js:noInline')
  static Media getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Media>(create);
  static Media? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(2)
  set description($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(2)
  void clearDescription() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get mimetype => $_getSZ(2);
  @$pb.TagNumber(3)
  set mimetype($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasMimetype() => $_has(2);
  @$pb.TagNumber(3)
  void clearMimetype() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get image => $_getSZ(3);
  @$pb.TagNumber(4)
  set image($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasImage() => $_has(3);
  @$pb.TagNumber(4)
  void clearImage() => clearField(4);
}

class MediaSearchRequest extends $pb.GeneratedMessage {
  factory MediaSearchRequest({
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
  MediaSearchRequest._() : super();
  factory MediaSearchRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory MediaSearchRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MediaSearchRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'media'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..a<$fixnum.Int64>(900, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(901, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  MediaSearchRequest clone() => MediaSearchRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  MediaSearchRequest copyWith(void Function(MediaSearchRequest) updates) => super.copyWith((message) => updates(message as MediaSearchRequest)) as MediaSearchRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MediaSearchRequest create() => MediaSearchRequest._();
  MediaSearchRequest createEmptyInstance() => create();
  static $pb.PbList<MediaSearchRequest> createRepeated() => $pb.PbList<MediaSearchRequest>();
  @$core.pragma('dart2js:noInline')
  static MediaSearchRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MediaSearchRequest>(create);
  static MediaSearchRequest? _defaultInstance;

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

class MediaSearchResponse extends $pb.GeneratedMessage {
  factory MediaSearchResponse({
    MediaSearchRequest? next,
    $core.Iterable<Media>? items,
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
  MediaSearchResponse._() : super();
  factory MediaSearchResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory MediaSearchResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MediaSearchResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'media'), createEmptyInstance: create)
    ..aOM<MediaSearchRequest>(1, _omitFieldNames ? '' : 'next', subBuilder: MediaSearchRequest.create)
    ..pc<Media>(2, _omitFieldNames ? '' : 'items', $pb.PbFieldType.PM, subBuilder: Media.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  MediaSearchResponse clone() => MediaSearchResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  MediaSearchResponse copyWith(void Function(MediaSearchResponse) updates) => super.copyWith((message) => updates(message as MediaSearchResponse)) as MediaSearchResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MediaSearchResponse create() => MediaSearchResponse._();
  MediaSearchResponse createEmptyInstance() => create();
  static $pb.PbList<MediaSearchResponse> createRepeated() => $pb.PbList<MediaSearchResponse>();
  @$core.pragma('dart2js:noInline')
  static MediaSearchResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MediaSearchResponse>(create);
  static MediaSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  MediaSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(MediaSearchRequest v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => clearField(1);
  @$pb.TagNumber(1)
  MediaSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.List<Media> get items => $_getList(1);
}

class Download extends $pb.GeneratedMessage {
  factory Download({
    Media? media,
    $fixnum.Int64? bytes,
    $fixnum.Int64? downloaded,
    $core.String? initiatedAt,
    $core.String? pausedAt,
    $core.int? peers,
  }) {
    final $result = create();
    if (media != null) {
      $result.media = media;
    }
    if (bytes != null) {
      $result.bytes = bytes;
    }
    if (downloaded != null) {
      $result.downloaded = downloaded;
    }
    if (initiatedAt != null) {
      $result.initiatedAt = initiatedAt;
    }
    if (pausedAt != null) {
      $result.pausedAt = pausedAt;
    }
    if (peers != null) {
      $result.peers = peers;
    }
    return $result;
  }
  Download._() : super();
  factory Download.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory Download.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'Download', package: const $pb.PackageName(_omitMessageNames ? '' : 'media'), createEmptyInstance: create)
    ..aOM<Media>(1, _omitFieldNames ? '' : 'media', subBuilder: Media.create)
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'bytes', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(3, _omitFieldNames ? '' : 'downloaded', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOS(4, _omitFieldNames ? '' : 'initiated_at')
    ..aOS(5, _omitFieldNames ? '' : 'paused_at')
    ..a<$core.int>(6, _omitFieldNames ? '' : 'peers', $pb.PbFieldType.OU3)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  Download clone() => Download()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  Download copyWith(void Function(Download) updates) => super.copyWith((message) => updates(message as Download)) as Download;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Download create() => Download._();
  Download createEmptyInstance() => create();
  static $pb.PbList<Download> createRepeated() => $pb.PbList<Download>();
  @$core.pragma('dart2js:noInline')
  static Download getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Download>(create);
  static Download? _defaultInstance;

  @$pb.TagNumber(1)
  Media get media => $_getN(0);
  @$pb.TagNumber(1)
  set media(Media v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasMedia() => $_has(0);
  @$pb.TagNumber(1)
  void clearMedia() => clearField(1);
  @$pb.TagNumber(1)
  Media ensureMedia() => $_ensure(0);

  @$pb.TagNumber(2)
  $fixnum.Int64 get bytes => $_getI64(1);
  @$pb.TagNumber(2)
  set bytes($fixnum.Int64 v) { $_setInt64(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasBytes() => $_has(1);
  @$pb.TagNumber(2)
  void clearBytes() => clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get downloaded => $_getI64(2);
  @$pb.TagNumber(3)
  set downloaded($fixnum.Int64 v) { $_setInt64(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasDownloaded() => $_has(2);
  @$pb.TagNumber(3)
  void clearDownloaded() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get initiatedAt => $_getSZ(3);
  @$pb.TagNumber(4)
  set initiatedAt($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasInitiatedAt() => $_has(3);
  @$pb.TagNumber(4)
  void clearInitiatedAt() => clearField(4);

  @$pb.TagNumber(5)
  $core.String get pausedAt => $_getSZ(4);
  @$pb.TagNumber(5)
  set pausedAt($core.String v) { $_setString(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasPausedAt() => $_has(4);
  @$pb.TagNumber(5)
  void clearPausedAt() => clearField(5);

  @$pb.TagNumber(6)
  $core.int get peers => $_getIZ(5);
  @$pb.TagNumber(6)
  set peers($core.int v) { $_setUnsignedInt32(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasPeers() => $_has(5);
  @$pb.TagNumber(6)
  void clearPeers() => clearField(6);
}

class DownloadSearchRequest extends $pb.GeneratedMessage {
  factory DownloadSearchRequest({
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
  DownloadSearchRequest._() : super();
  factory DownloadSearchRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DownloadSearchRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DownloadSearchRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'media'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..a<$fixnum.Int64>(900, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(901, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DownloadSearchRequest clone() => DownloadSearchRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DownloadSearchRequest copyWith(void Function(DownloadSearchRequest) updates) => super.copyWith((message) => updates(message as DownloadSearchRequest)) as DownloadSearchRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadSearchRequest create() => DownloadSearchRequest._();
  DownloadSearchRequest createEmptyInstance() => create();
  static $pb.PbList<DownloadSearchRequest> createRepeated() => $pb.PbList<DownloadSearchRequest>();
  @$core.pragma('dart2js:noInline')
  static DownloadSearchRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DownloadSearchRequest>(create);
  static DownloadSearchRequest? _defaultInstance;

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

class DownloadSearchResponse extends $pb.GeneratedMessage {
  factory DownloadSearchResponse({
    DownloadSearchRequest? next,
    $core.Iterable<Download>? items,
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
  DownloadSearchResponse._() : super();
  factory DownloadSearchResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DownloadSearchResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DownloadSearchResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'media'), createEmptyInstance: create)
    ..aOM<DownloadSearchRequest>(1, _omitFieldNames ? '' : 'next', subBuilder: DownloadSearchRequest.create)
    ..pc<Download>(2, _omitFieldNames ? '' : 'items', $pb.PbFieldType.PM, subBuilder: Download.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DownloadSearchResponse clone() => DownloadSearchResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DownloadSearchResponse copyWith(void Function(DownloadSearchResponse) updates) => super.copyWith((message) => updates(message as DownloadSearchResponse)) as DownloadSearchResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadSearchResponse create() => DownloadSearchResponse._();
  DownloadSearchResponse createEmptyInstance() => create();
  static $pb.PbList<DownloadSearchResponse> createRepeated() => $pb.PbList<DownloadSearchResponse>();
  @$core.pragma('dart2js:noInline')
  static DownloadSearchResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DownloadSearchResponse>(create);
  static DownloadSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  DownloadSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(DownloadSearchRequest v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => clearField(1);
  @$pb.TagNumber(1)
  DownloadSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.List<Download> get items => $_getList(1);
}

class DownloadBeginRequest extends $pb.GeneratedMessage {
  factory DownloadBeginRequest() => create();
  DownloadBeginRequest._() : super();
  factory DownloadBeginRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DownloadBeginRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DownloadBeginRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'media'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DownloadBeginRequest clone() => DownloadBeginRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DownloadBeginRequest copyWith(void Function(DownloadBeginRequest) updates) => super.copyWith((message) => updates(message as DownloadBeginRequest)) as DownloadBeginRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadBeginRequest create() => DownloadBeginRequest._();
  DownloadBeginRequest createEmptyInstance() => create();
  static $pb.PbList<DownloadBeginRequest> createRepeated() => $pb.PbList<DownloadBeginRequest>();
  @$core.pragma('dart2js:noInline')
  static DownloadBeginRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DownloadBeginRequest>(create);
  static DownloadBeginRequest? _defaultInstance;
}

class DownloadBeginResponse extends $pb.GeneratedMessage {
  factory DownloadBeginResponse({
    Download? download,
  }) {
    final $result = create();
    if (download != null) {
      $result.download = download;
    }
    return $result;
  }
  DownloadBeginResponse._() : super();
  factory DownloadBeginResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DownloadBeginResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DownloadBeginResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'media'), createEmptyInstance: create)
    ..aOM<Download>(1, _omitFieldNames ? '' : 'download', subBuilder: Download.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DownloadBeginResponse clone() => DownloadBeginResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DownloadBeginResponse copyWith(void Function(DownloadBeginResponse) updates) => super.copyWith((message) => updates(message as DownloadBeginResponse)) as DownloadBeginResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadBeginResponse create() => DownloadBeginResponse._();
  DownloadBeginResponse createEmptyInstance() => create();
  static $pb.PbList<DownloadBeginResponse> createRepeated() => $pb.PbList<DownloadBeginResponse>();
  @$core.pragma('dart2js:noInline')
  static DownloadBeginResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DownloadBeginResponse>(create);
  static DownloadBeginResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Download get download => $_getN(0);
  @$pb.TagNumber(1)
  set download(Download v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasDownload() => $_has(0);
  @$pb.TagNumber(1)
  void clearDownload() => clearField(1);
  @$pb.TagNumber(1)
  Download ensureDownload() => $_ensure(0);
}

class DownloadPauseRequest extends $pb.GeneratedMessage {
  factory DownloadPauseRequest() => create();
  DownloadPauseRequest._() : super();
  factory DownloadPauseRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DownloadPauseRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DownloadPauseRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'media'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DownloadPauseRequest clone() => DownloadPauseRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DownloadPauseRequest copyWith(void Function(DownloadPauseRequest) updates) => super.copyWith((message) => updates(message as DownloadPauseRequest)) as DownloadPauseRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadPauseRequest create() => DownloadPauseRequest._();
  DownloadPauseRequest createEmptyInstance() => create();
  static $pb.PbList<DownloadPauseRequest> createRepeated() => $pb.PbList<DownloadPauseRequest>();
  @$core.pragma('dart2js:noInline')
  static DownloadPauseRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DownloadPauseRequest>(create);
  static DownloadPauseRequest? _defaultInstance;
}

class DownloadPauseResponse extends $pb.GeneratedMessage {
  factory DownloadPauseResponse({
    Download? download,
  }) {
    final $result = create();
    if (download != null) {
      $result.download = download;
    }
    return $result;
  }
  DownloadPauseResponse._() : super();
  factory DownloadPauseResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DownloadPauseResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DownloadPauseResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'media'), createEmptyInstance: create)
    ..aOM<Download>(1, _omitFieldNames ? '' : 'download', subBuilder: Download.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DownloadPauseResponse clone() => DownloadPauseResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DownloadPauseResponse copyWith(void Function(DownloadPauseResponse) updates) => super.copyWith((message) => updates(message as DownloadPauseResponse)) as DownloadPauseResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DownloadPauseResponse create() => DownloadPauseResponse._();
  DownloadPauseResponse createEmptyInstance() => create();
  static $pb.PbList<DownloadPauseResponse> createRepeated() => $pb.PbList<DownloadPauseResponse>();
  @$core.pragma('dart2js:noInline')
  static DownloadPauseResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DownloadPauseResponse>(create);
  static DownloadPauseResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Download get download => $_getN(0);
  @$pb.TagNumber(1)
  set download(Download v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasDownload() => $_has(0);
  @$pb.TagNumber(1)
  void clearDownload() => clearField(1);
  @$pb.TagNumber(1)
  Download ensureDownload() => $_ensure(0);
}


const _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
