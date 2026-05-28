// This is a generated file - do not edit.
//
// Generated from torrent.proto.

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

class Peers extends $pb.GeneratedMessage {
  factory Peers({
    $core.int? min,
    $core.int? max,
  }) {
    final result = create();
    if (min != null) result.min = min;
    if (max != null) result.max = max;
    return result;
  }

  Peers._();

  factory Peers.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Peers.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Peers',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'torrents'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'min', fieldType: $pb.PbFieldType.OU3)
    ..aI(2, _omitFieldNames ? '' : 'max', fieldType: $pb.PbFieldType.OU3)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Peers clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Peers copyWith(void Function(Peers) updates) =>
      super.copyWith((message) => updates(message as Peers)) as Peers;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Peers create() => Peers._();
  @$core.override
  Peers createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Peers getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Peers>(create);
  static Peers? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get min => $_getIZ(0);
  @$pb.TagNumber(1)
  set min($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasMin() => $_has(0);
  @$pb.TagNumber(1)
  void clearMin() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get max => $_getIZ(1);
  @$pb.TagNumber(2)
  set max($core.int value) => $_setUnsignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasMax() => $_has(1);
  @$pb.TagNumber(2)
  void clearMax() => $_clearField(2);
}

class Limit extends $pb.GeneratedMessage {
  factory Limit({
    $core.int? rate,
    $core.int? burst,
  }) {
    final result = create();
    if (rate != null) result.rate = rate;
    if (burst != null) result.burst = burst;
    return result;
  }

  Limit._();

  factory Limit.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Limit.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Limit',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'torrents'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'rate', fieldType: $pb.PbFieldType.OU3)
    ..aI(2, _omitFieldNames ? '' : 'burst', fieldType: $pb.PbFieldType.OU3)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Limit clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Limit copyWith(void Function(Limit) updates) =>
      super.copyWith((message) => updates(message as Limit)) as Limit;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Limit create() => Limit._();
  @$core.override
  Limit createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Limit getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Limit>(create);
  static Limit? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get rate => $_getIZ(0);
  @$pb.TagNumber(1)
  set rate($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasRate() => $_has(0);
  @$pb.TagNumber(1)
  void clearRate() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get burst => $_getIZ(1);
  @$pb.TagNumber(2)
  set burst($core.int value) => $_setUnsignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasBurst() => $_has(1);
  @$pb.TagNumber(2)
  void clearBurst() => $_clearField(2);
}

class TorrentSettings extends $pb.GeneratedMessage {
  factory TorrentSettings({
    $core.bool? seed,
    $core.bool? pex,
    $core.bool? autoBootstrap,
    $core.bool? autoLocateMedia,
    $core.bool? firewalled,
    $core.bool? resumable,
    $core.String? ip4,
    $core.String? ip6,
    $core.int? port,
    $core.bool? log,
    $core.bool? debug,
    Limit? download,
    Limit? upload,
    Limit? inbound,
    Limit? outbound,
    Peers? peers,
    $fixnum.Int64? maximumRequests,
    $fixnum.Int64? connections,
  }) {
    final result = create();
    if (seed != null) result.seed = seed;
    if (pex != null) result.pex = pex;
    if (autoBootstrap != null) result.autoBootstrap = autoBootstrap;
    if (autoLocateMedia != null) result.autoLocateMedia = autoLocateMedia;
    if (firewalled != null) result.firewalled = firewalled;
    if (resumable != null) result.resumable = resumable;
    if (ip4 != null) result.ip4 = ip4;
    if (ip6 != null) result.ip6 = ip6;
    if (port != null) result.port = port;
    if (log != null) result.log = log;
    if (debug != null) result.debug = debug;
    if (download != null) result.download = download;
    if (upload != null) result.upload = upload;
    if (inbound != null) result.inbound = inbound;
    if (outbound != null) result.outbound = outbound;
    if (peers != null) result.peers = peers;
    if (maximumRequests != null) result.maximumRequests = maximumRequests;
    if (connections != null) result.connections = connections;
    return result;
  }

  TorrentSettings._();

  factory TorrentSettings.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TorrentSettings.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TorrentSettings',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'torrents'),
      createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'seed')
    ..aOB(2, _omitFieldNames ? '' : 'pex')
    ..aOB(3, _omitFieldNames ? '' : 'auto_bootstrap')
    ..aOB(4, _omitFieldNames ? '' : 'auto_locate_media')
    ..aOB(5, _omitFieldNames ? '' : 'firewalled')
    ..aOB(6, _omitFieldNames ? '' : 'resumable')
    ..aOS(7, _omitFieldNames ? '' : 'ip4')
    ..aOS(8, _omitFieldNames ? '' : 'ip6')
    ..aI(9, _omitFieldNames ? '' : 'port', fieldType: $pb.PbFieldType.OU3)
    ..aOB(998, _omitFieldNames ? '' : 'log')
    ..aOB(999, _omitFieldNames ? '' : 'debug')
    ..aOM<Limit>(1000, _omitFieldNames ? '' : 'download',
        subBuilder: Limit.create)
    ..aOM<Limit>(1001, _omitFieldNames ? '' : 'upload',
        subBuilder: Limit.create)
    ..aOM<Limit>(1002, _omitFieldNames ? '' : 'inbound',
        subBuilder: Limit.create)
    ..aOM<Limit>(1003, _omitFieldNames ? '' : 'outbound',
        subBuilder: Limit.create)
    ..aOM<Peers>(1004, _omitFieldNames ? '' : 'peers', subBuilder: Peers.create)
    ..a<$fixnum.Int64>(
        1005, _omitFieldNames ? '' : 'maximum_requests', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        1006, _omitFieldNames ? '' : 'connections', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TorrentSettings clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TorrentSettings copyWith(void Function(TorrentSettings) updates) =>
      super.copyWith((message) => updates(message as TorrentSettings))
          as TorrentSettings;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TorrentSettings create() => TorrentSettings._();
  @$core.override
  TorrentSettings createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TorrentSettings getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TorrentSettings>(create);
  static TorrentSettings? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get seed => $_getBF(0);
  @$pb.TagNumber(1)
  set seed($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSeed() => $_has(0);
  @$pb.TagNumber(1)
  void clearSeed() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.bool get pex => $_getBF(1);
  @$pb.TagNumber(2)
  set pex($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPex() => $_has(1);
  @$pb.TagNumber(2)
  void clearPex() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.bool get autoBootstrap => $_getBF(2);
  @$pb.TagNumber(3)
  set autoBootstrap($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(3)
  $core.bool hasAutoBootstrap() => $_has(2);
  @$pb.TagNumber(3)
  void clearAutoBootstrap() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.bool get autoLocateMedia => $_getBF(3);
  @$pb.TagNumber(4)
  set autoLocateMedia($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(4)
  $core.bool hasAutoLocateMedia() => $_has(3);
  @$pb.TagNumber(4)
  void clearAutoLocateMedia() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.bool get firewalled => $_getBF(4);
  @$pb.TagNumber(5)
  set firewalled($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(5)
  $core.bool hasFirewalled() => $_has(4);
  @$pb.TagNumber(5)
  void clearFirewalled() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.bool get resumable => $_getBF(5);
  @$pb.TagNumber(6)
  set resumable($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(6)
  $core.bool hasResumable() => $_has(5);
  @$pb.TagNumber(6)
  void clearResumable() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.String get ip4 => $_getSZ(6);
  @$pb.TagNumber(7)
  set ip4($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasIp4() => $_has(6);
  @$pb.TagNumber(7)
  void clearIp4() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.String get ip6 => $_getSZ(7);
  @$pb.TagNumber(8)
  set ip6($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasIp6() => $_has(7);
  @$pb.TagNumber(8)
  void clearIp6() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.int get port => $_getIZ(8);
  @$pb.TagNumber(9)
  set port($core.int value) => $_setUnsignedInt32(8, value);
  @$pb.TagNumber(9)
  $core.bool hasPort() => $_has(8);
  @$pb.TagNumber(9)
  void clearPort() => $_clearField(9);

  @$pb.TagNumber(998)
  $core.bool get log => $_getBF(9);
  @$pb.TagNumber(998)
  set log($core.bool value) => $_setBool(9, value);
  @$pb.TagNumber(998)
  $core.bool hasLog() => $_has(9);
  @$pb.TagNumber(998)
  void clearLog() => $_clearField(998);

  @$pb.TagNumber(999)
  $core.bool get debug => $_getBF(10);
  @$pb.TagNumber(999)
  set debug($core.bool value) => $_setBool(10, value);
  @$pb.TagNumber(999)
  $core.bool hasDebug() => $_has(10);
  @$pb.TagNumber(999)
  void clearDebug() => $_clearField(999);

  @$pb.TagNumber(1000)
  Limit get download => $_getN(11);
  @$pb.TagNumber(1000)
  set download(Limit value) => $_setField(1000, value);
  @$pb.TagNumber(1000)
  $core.bool hasDownload() => $_has(11);
  @$pb.TagNumber(1000)
  void clearDownload() => $_clearField(1000);
  @$pb.TagNumber(1000)
  Limit ensureDownload() => $_ensure(11);

  @$pb.TagNumber(1001)
  Limit get upload => $_getN(12);
  @$pb.TagNumber(1001)
  set upload(Limit value) => $_setField(1001, value);
  @$pb.TagNumber(1001)
  $core.bool hasUpload() => $_has(12);
  @$pb.TagNumber(1001)
  void clearUpload() => $_clearField(1001);
  @$pb.TagNumber(1001)
  Limit ensureUpload() => $_ensure(12);

  @$pb.TagNumber(1002)
  Limit get inbound => $_getN(13);
  @$pb.TagNumber(1002)
  set inbound(Limit value) => $_setField(1002, value);
  @$pb.TagNumber(1002)
  $core.bool hasInbound() => $_has(13);
  @$pb.TagNumber(1002)
  void clearInbound() => $_clearField(1002);
  @$pb.TagNumber(1002)
  Limit ensureInbound() => $_ensure(13);

  @$pb.TagNumber(1003)
  Limit get outbound => $_getN(14);
  @$pb.TagNumber(1003)
  set outbound(Limit value) => $_setField(1003, value);
  @$pb.TagNumber(1003)
  $core.bool hasOutbound() => $_has(14);
  @$pb.TagNumber(1003)
  void clearOutbound() => $_clearField(1003);
  @$pb.TagNumber(1003)
  Limit ensureOutbound() => $_ensure(14);

  @$pb.TagNumber(1004)
  Peers get peers => $_getN(15);
  @$pb.TagNumber(1004)
  set peers(Peers value) => $_setField(1004, value);
  @$pb.TagNumber(1004)
  $core.bool hasPeers() => $_has(15);
  @$pb.TagNumber(1004)
  void clearPeers() => $_clearField(1004);
  @$pb.TagNumber(1004)
  Peers ensurePeers() => $_ensure(15);

  @$pb.TagNumber(1005)
  $fixnum.Int64 get maximumRequests => $_getI64(16);
  @$pb.TagNumber(1005)
  set maximumRequests($fixnum.Int64 value) => $_setInt64(16, value);
  @$pb.TagNumber(1005)
  $core.bool hasMaximumRequests() => $_has(16);
  @$pb.TagNumber(1005)
  void clearMaximumRequests() => $_clearField(1005);

  @$pb.TagNumber(1006)
  $fixnum.Int64 get connections => $_getI64(17);
  @$pb.TagNumber(1006)
  set connections($fixnum.Int64 value) => $_setInt64(17, value);
  @$pb.TagNumber(1006)
  $core.bool hasConnections() => $_has(17);
  @$pb.TagNumber(1006)
  void clearConnections() => $_clearField(1006);
}

class DiscoverySettings extends $pb.GeneratedMessage {
  factory DiscoverySettings({
    $core.bool? enabled,
    $core.int? ratio,
    $core.int? partitions,
    $core.int? workloads,
    $core.String? seed,
  }) {
    final result = create();
    if (enabled != null) result.enabled = enabled;
    if (ratio != null) result.ratio = ratio;
    if (partitions != null) result.partitions = partitions;
    if (workloads != null) result.workloads = workloads;
    if (seed != null) result.seed = seed;
    return result;
  }

  DiscoverySettings._();

  factory DiscoverySettings.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DiscoverySettings.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DiscoverySettings',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'torrents'),
      createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'enabled')
    ..aI(2, _omitFieldNames ? '' : 'ratio', fieldType: $pb.PbFieldType.OU3)
    ..aI(3, _omitFieldNames ? '' : 'partitions', fieldType: $pb.PbFieldType.OU3)
    ..aI(4, _omitFieldNames ? '' : 'workloads', fieldType: $pb.PbFieldType.OU3)
    ..aOS(5, _omitFieldNames ? '' : 'seed')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoverySettings clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoverySettings copyWith(void Function(DiscoverySettings) updates) =>
      super.copyWith((message) => updates(message as DiscoverySettings))
          as DiscoverySettings;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DiscoverySettings create() => DiscoverySettings._();
  @$core.override
  DiscoverySettings createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DiscoverySettings getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DiscoverySettings>(create);
  static DiscoverySettings? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get enabled => $_getBF(0);
  @$pb.TagNumber(1)
  set enabled($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasEnabled() => $_has(0);
  @$pb.TagNumber(1)
  void clearEnabled() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get ratio => $_getIZ(1);
  @$pb.TagNumber(2)
  set ratio($core.int value) => $_setUnsignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasRatio() => $_has(1);
  @$pb.TagNumber(2)
  void clearRatio() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.int get partitions => $_getIZ(2);
  @$pb.TagNumber(3)
  set partitions($core.int value) => $_setUnsignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasPartitions() => $_has(2);
  @$pb.TagNumber(3)
  void clearPartitions() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.int get workloads => $_getIZ(3);
  @$pb.TagNumber(4)
  set workloads($core.int value) => $_setUnsignedInt32(3, value);
  @$pb.TagNumber(4)
  $core.bool hasWorkloads() => $_has(3);
  @$pb.TagNumber(4)
  void clearWorkloads() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get seed => $_getSZ(4);
  @$pb.TagNumber(5)
  set seed($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasSeed() => $_has(4);
  @$pb.TagNumber(5)
  void clearSeed() => $_clearField(5);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
