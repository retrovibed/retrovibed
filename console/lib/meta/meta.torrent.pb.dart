// This is a generated file - do not edit.
//
// Generated from meta.torrent.proto.

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
    final result = create();
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
      create()..mergeFromBuffer(data, registry);
  factory TorrentDiagnostics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TorrentDiagnostics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
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
  static TorrentDiagnostics create() => TorrentDiagnostics._();
  @$core.override
  TorrentDiagnostics createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TorrentDiagnostics getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TorrentDiagnostics>(create);
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
    final result = create();
    if (torrent != null) result.torrent = torrent;
    return result;
  }

  TorrentMetricsResponse._();

  factory TorrentMetricsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TorrentMetricsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TorrentMetricsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<TorrentDiagnostics>(1, _omitFieldNames ? '' : 'torrent',
        subBuilder: TorrentDiagnostics.create)
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
  static TorrentMetricsResponse create() => TorrentMetricsResponse._();
  @$core.override
  TorrentMetricsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TorrentMetricsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TorrentMetricsResponse>(create);
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

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
