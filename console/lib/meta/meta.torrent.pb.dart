// This is a generated file - do not edit.
//
// Generated from meta/meta.torrent.proto.

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

class TorrentDiagnostics extends $pb.GeneratedMessage {
  factory TorrentDiagnostics({
    $fixnum.Int64? total,
    $fixnum.Int64? seeding,
    $fixnum.Int64? bytes,
    $fixnum.Int64? downloaded,
    $fixnum.Int64? uploaded,
    $fixnum.Int64? peers,
  }) {
    final result = TorrentDiagnostics._();
    if (total != null) result.total = total;
    if (seeding != null) result.seeding = seeding;
    if (bytes != null) result.bytes = bytes;
    if (downloaded != null) result.downloaded = downloaded;
    if (uploaded != null) result.uploaded = uploaded;
    if (peers != null) result.peers = peers;
    return result;
  }

  TorrentDiagnostics._();

  factory TorrentDiagnostics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TorrentDiagnostics()..mergeFromBuffer(data, registry);
  factory TorrentDiagnostics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TorrentDiagnostics()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TorrentDiagnostics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: TorrentDiagnostics.$_createMessage)
    ..a<$fixnum.Int64>(1, _omitFieldNames ? '' : 'total', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'seeding', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(3, _omitFieldNames ? '' : 'bytes', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        4, _omitFieldNames ? '' : 'downloaded', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        5, _omitFieldNames ? '' : 'uploaded', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(6, _omitFieldNames ? '' : 'peers', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TorrentDiagnostics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TorrentDiagnostics copyWith(void Function(TorrentDiagnostics) updates) =>
      super.copyWith((message) => updates(message as TorrentDiagnostics))
          as TorrentDiagnostics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use TorrentDiagnostics() / TorrentDiagnostics.new instead')
  static TorrentDiagnostics create() => TorrentDiagnostics._();
  static $pb.GeneratedMessage $_createMessage() => TorrentDiagnostics._();
  @$core.override
  TorrentDiagnostics createEmptyInstance() => TorrentDiagnostics._();
  @$core.pragma('dart2js:noInline')
  static TorrentDiagnostics getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TorrentDiagnostics>(
          TorrentDiagnostics.$_createMessage);
  static TorrentDiagnostics? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get total => $_getI64(0);
  @$pb.TagNumber(1)
  set total($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasTotal() => $_has(0);
  @$pb.TagNumber(1)
  void clearTotal() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get seeding => $_getI64(1);
  @$pb.TagNumber(2)
  set seeding($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasSeeding() => $_has(1);
  @$pb.TagNumber(2)
  void clearSeeding() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get bytes => $_getI64(2);
  @$pb.TagNumber(3)
  set bytes($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasBytes() => $_has(2);
  @$pb.TagNumber(3)
  void clearBytes() => $_clearField(3);

  @$pb.TagNumber(4)
  $fixnum.Int64 get downloaded => $_getI64(3);
  @$pb.TagNumber(4)
  set downloaded($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(4)
  $core.bool hasDownloaded() => $_has(3);
  @$pb.TagNumber(4)
  void clearDownloaded() => $_clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get uploaded => $_getI64(4);
  @$pb.TagNumber(5)
  set uploaded($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(5)
  $core.bool hasUploaded() => $_has(4);
  @$pb.TagNumber(5)
  void clearUploaded() => $_clearField(5);

  @$pb.TagNumber(6)
  $fixnum.Int64 get peers => $_getI64(5);
  @$pb.TagNumber(6)
  set peers($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(6)
  $core.bool hasPeers() => $_has(5);
  @$pb.TagNumber(6)
  void clearPeers() => $_clearField(6);
}

class TorrentMetricsResponse extends $pb.GeneratedMessage {
  factory TorrentMetricsResponse({
    TorrentDiagnostics? torrent,
  }) {
    final result = TorrentMetricsResponse._();
    if (torrent != null) result.torrent = torrent;
    return result;
  }

  TorrentMetricsResponse._();

  factory TorrentMetricsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TorrentMetricsResponse()..mergeFromBuffer(data, registry);
  factory TorrentMetricsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TorrentMetricsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TorrentMetricsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: TorrentMetricsResponse.$_createMessage)
    ..aOM<TorrentDiagnostics>(1, _omitFieldNames ? '' : 'torrent',
        subBuilder: TorrentDiagnostics.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TorrentMetricsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TorrentMetricsResponse copyWith(
          void Function(TorrentMetricsResponse) updates) =>
      super.copyWith((message) => updates(message as TorrentMetricsResponse))
          as TorrentMetricsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use TorrentMetricsResponse() / TorrentMetricsResponse.new instead')
  static TorrentMetricsResponse create() => TorrentMetricsResponse._();
  static $pb.GeneratedMessage $_createMessage() => TorrentMetricsResponse._();
  @$core.override
  TorrentMetricsResponse createEmptyInstance() => TorrentMetricsResponse._();
  @$core.pragma('dart2js:noInline')
  static TorrentMetricsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TorrentMetricsResponse>(
          TorrentMetricsResponse.$_createMessage);
  static TorrentMetricsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  TorrentDiagnostics get torrent => $_getN(0);
  @$pb.TagNumber(1)
  set torrent(TorrentDiagnostics value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasTorrent() => $_has(0);
  @$pb.TagNumber(1)
  void clearTorrent() => $_clearField(1);
  @$pb.TagNumber(1)
  TorrentDiagnostics ensureTorrent() => $_ensure(0);
}

class TorrentInfoResponse extends $pb.GeneratedMessage {
  factory TorrentInfoResponse({
    TorrentMeta? meta,
    TorrentDetails? details,
    $core.Iterable<TorrentFile>? files,
  }) {
    final result = TorrentInfoResponse._();
    if (meta != null) result.meta = meta;
    if (details != null) result.details = details;
    if (files != null) result.files.addAll(files);
    return result;
  }

  TorrentInfoResponse._();

  factory TorrentInfoResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TorrentInfoResponse()..mergeFromBuffer(data, registry);
  factory TorrentInfoResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TorrentInfoResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TorrentInfoResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: TorrentInfoResponse.$_createMessage)
    ..aOM<TorrentMeta>(1, _omitFieldNames ? '' : 'meta',
        subBuilder: TorrentMeta.$_createMessage)
    ..aOM<TorrentDetails>(2, _omitFieldNames ? '' : 'details',
        subBuilder: TorrentDetails.$_createMessage)
    ..pPM<TorrentFile>(1000, _omitFieldNames ? '' : 'files',
        subBuilder: TorrentFile.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TorrentInfoResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TorrentInfoResponse copyWith(void Function(TorrentInfoResponse) updates) =>
      super.copyWith((message) => updates(message as TorrentInfoResponse))
          as TorrentInfoResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use TorrentInfoResponse() / TorrentInfoResponse.new instead')
  static TorrentInfoResponse create() => TorrentInfoResponse._();
  static $pb.GeneratedMessage $_createMessage() => TorrentInfoResponse._();
  @$core.override
  TorrentInfoResponse createEmptyInstance() => TorrentInfoResponse._();
  @$core.pragma('dart2js:noInline')
  static TorrentInfoResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TorrentInfoResponse>(
          TorrentInfoResponse.$_createMessage);
  static TorrentInfoResponse? _defaultInstance;

  @$pb.TagNumber(1)
  TorrentMeta get meta => $_getN(0);
  @$pb.TagNumber(1)
  set meta(TorrentMeta value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMeta() => $_has(0);
  @$pb.TagNumber(1)
  void clearMeta() => $_clearField(1);
  @$pb.TagNumber(1)
  TorrentMeta ensureMeta() => $_ensure(0);

  @$pb.TagNumber(2)
  TorrentDetails get details => $_getN(1);
  @$pb.TagNumber(2)
  set details(TorrentDetails value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasDetails() => $_has(1);
  @$pb.TagNumber(2)
  void clearDetails() => $_clearField(2);
  @$pb.TagNumber(2)
  TorrentDetails ensureDetails() => $_ensure(1);

  @$pb.TagNumber(1000)
  $pb.PbList<TorrentFile> get files => $_getList(2);
}

class TorrentFile extends $pb.GeneratedMessage {
  factory TorrentFile({
    $core.String? name,
    $fixnum.Int64? length,
    $core.String? path,
  }) {
    final result = TorrentFile._();
    if (name != null) result.name = name;
    if (length != null) result.length = length;
    if (path != null) result.path = path;
    return result;
  }

  TorrentFile._();

  factory TorrentFile.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TorrentFile()..mergeFromBuffer(data, registry);
  factory TorrentFile.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TorrentFile()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TorrentFile',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: TorrentFile.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'name')
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'length', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOS(3, _omitFieldNames ? '' : 'path')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TorrentFile clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TorrentFile copyWith(void Function(TorrentFile) updates) =>
      super.copyWith((message) => updates(message as TorrentFile))
          as TorrentFile;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use TorrentFile() / TorrentFile.new instead')
  static TorrentFile create() => TorrentFile._();
  static $pb.GeneratedMessage $_createMessage() => TorrentFile._();
  @$core.override
  TorrentFile createEmptyInstance() => TorrentFile._();
  @$core.pragma('dart2js:noInline')
  static TorrentFile getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<TorrentFile>(
          TorrentFile.$_createMessage);
  static TorrentFile? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(1)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearName() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get length => $_getI64(1);
  @$pb.TagNumber(2)
  set length($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasLength() => $_has(1);
  @$pb.TagNumber(2)
  void clearLength() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get path => $_getSZ(2);
  @$pb.TagNumber(3)
  set path($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasPath() => $_has(2);
  @$pb.TagNumber(3)
  void clearPath() => $_clearField(3);
}

class TorrentMeta extends $pb.GeneratedMessage {
  factory TorrentMeta({
    $core.String? comment,
    $core.String? encoding,
    $core.String? createdBy,
    $fixnum.Int64? creationDate,
    $core.Iterable<$core.String>? announceList,
    $core.Iterable<$core.String>? urlList,
  }) {
    final result = TorrentMeta._();
    if (comment != null) result.comment = comment;
    if (encoding != null) result.encoding = encoding;
    if (createdBy != null) result.createdBy = createdBy;
    if (creationDate != null) result.creationDate = creationDate;
    if (announceList != null) result.announceList.addAll(announceList);
    if (urlList != null) result.urlList.addAll(urlList);
    return result;
  }

  TorrentMeta._();

  factory TorrentMeta.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TorrentMeta()..mergeFromBuffer(data, registry);
  factory TorrentMeta.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TorrentMeta()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TorrentMeta',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: TorrentMeta.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'comment')
    ..aOS(2, _omitFieldNames ? '' : 'encoding')
    ..aOS(3, _omitFieldNames ? '' : 'created_by')
    ..a<$fixnum.Int64>(
        4, _omitFieldNames ? '' : 'creation_date', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..pPS(1000, _omitFieldNames ? '' : 'announce_list')
    ..pPS(1001, _omitFieldNames ? '' : 'url_list')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TorrentMeta clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TorrentMeta copyWith(void Function(TorrentMeta) updates) =>
      super.copyWith((message) => updates(message as TorrentMeta))
          as TorrentMeta;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use TorrentMeta() / TorrentMeta.new instead')
  static TorrentMeta create() => TorrentMeta._();
  static $pb.GeneratedMessage $_createMessage() => TorrentMeta._();
  @$core.override
  TorrentMeta createEmptyInstance() => TorrentMeta._();
  @$core.pragma('dart2js:noInline')
  static TorrentMeta getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<TorrentMeta>(
          TorrentMeta.$_createMessage);
  static TorrentMeta? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get comment => $_getSZ(0);
  @$pb.TagNumber(1)
  set comment($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasComment() => $_has(0);
  @$pb.TagNumber(1)
  void clearComment() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get encoding => $_getSZ(1);
  @$pb.TagNumber(2)
  set encoding($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasEncoding() => $_has(1);
  @$pb.TagNumber(2)
  void clearEncoding() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get createdBy => $_getSZ(2);
  @$pb.TagNumber(3)
  set createdBy($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasCreatedBy() => $_has(2);
  @$pb.TagNumber(3)
  void clearCreatedBy() => $_clearField(3);

  @$pb.TagNumber(4)
  $fixnum.Int64 get creationDate => $_getI64(3);
  @$pb.TagNumber(4)
  set creationDate($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(4)
  $core.bool hasCreationDate() => $_has(3);
  @$pb.TagNumber(4)
  void clearCreationDate() => $_clearField(4);

  @$pb.TagNumber(1000)
  $pb.PbList<$core.String> get announceList => $_getList(4);

  @$pb.TagNumber(1001)
  $pb.PbList<$core.String> get urlList => $_getList(5);
}

class TorrentDetails extends $pb.GeneratedMessage {
  factory TorrentDetails({
    $core.String? name,
    $fixnum.Int64? length,
    $core.String? source,
    $core.bool? private,
  }) {
    final result = TorrentDetails._();
    if (name != null) result.name = name;
    if (length != null) result.length = length;
    if (source != null) result.source = source;
    if (private != null) result.private = private;
    return result;
  }

  TorrentDetails._();

  factory TorrentDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TorrentDetails()..mergeFromBuffer(data, registry);
  factory TorrentDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TorrentDetails()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TorrentDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: TorrentDetails.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'name')
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'length', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOS(3, _omitFieldNames ? '' : 'source')
    ..aOB(4, _omitFieldNames ? '' : 'private')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TorrentDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TorrentDetails copyWith(void Function(TorrentDetails) updates) =>
      super.copyWith((message) => updates(message as TorrentDetails))
          as TorrentDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use TorrentDetails() / TorrentDetails.new instead')
  static TorrentDetails create() => TorrentDetails._();
  static $pb.GeneratedMessage $_createMessage() => TorrentDetails._();
  @$core.override
  TorrentDetails createEmptyInstance() => TorrentDetails._();
  @$core.pragma('dart2js:noInline')
  static TorrentDetails getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<TorrentDetails>(
          TorrentDetails.$_createMessage);
  static TorrentDetails? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(1)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearName() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get length => $_getI64(1);
  @$pb.TagNumber(2)
  set length($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasLength() => $_has(1);
  @$pb.TagNumber(2)
  void clearLength() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get source => $_getSZ(2);
  @$pb.TagNumber(3)
  set source($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasSource() => $_has(2);
  @$pb.TagNumber(3)
  void clearSource() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.bool get private => $_getBF(3);
  @$pb.TagNumber(4)
  set private($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(4)
  $core.bool hasPrivate() => $_has(3);
  @$pb.TagNumber(4)
  void clearPrivate() => $_clearField(4);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
