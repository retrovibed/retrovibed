// This is a generated file - do not edit.
//
// Generated from meta.discovery.proto.

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

class DiscoveryDiagnostics extends $pb.GeneratedMessage {
  factory DiscoveryDiagnostics({
    $core.bool? enabled,
    $core.int? ratio,
    $core.int? partitions,
    $core.int? workloads,
    $core.String? localPartition,
    $fixnum.Int64? peers,
    $fixnum.Int64? peersDdisc,
    $fixnum.Int64? peersBep51,
    $fixnum.Int64? unknownHashes,
  }) {
    final result = create();
    if (enabled != null) result.enabled = enabled;
    if (ratio != null) result.ratio = ratio;
    if (partitions != null) result.partitions = partitions;
    if (workloads != null) result.workloads = workloads;
    if (localPartition != null) result.localPartition = localPartition;
    if (peers != null) result.peers = peers;
    if (peersDdisc != null) result.peersDdisc = peersDdisc;
    if (peersBep51 != null) result.peersBep51 = peersBep51;
    if (unknownHashes != null) result.unknownHashes = unknownHashes;
    return result;
  }

  DiscoveryDiagnostics._();

  factory DiscoveryDiagnostics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DiscoveryDiagnostics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DiscoveryDiagnostics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'enabled')
    ..aI(2, _omitFieldNames ? '' : 'ratio', fieldType: $pb.PbFieldType.OU3)
    ..aI(3, _omitFieldNames ? '' : 'partitions', fieldType: $pb.PbFieldType.OU3)
    ..aI(4, _omitFieldNames ? '' : 'workloads', fieldType: $pb.PbFieldType.OU3)
    ..aOS(5, _omitFieldNames ? '' : 'local_partition')
    ..a<$fixnum.Int64>(6, _omitFieldNames ? '' : 'peers', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        7, _omitFieldNames ? '' : 'peers_ddisc', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        8, _omitFieldNames ? '' : 'peers_bep51', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        9, _omitFieldNames ? '' : 'unknown_hashes', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoveryDiagnostics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoveryDiagnostics copyWith(void Function(DiscoveryDiagnostics) updates) =>
      super.copyWith((message) => updates(message as DiscoveryDiagnostics))
          as DiscoveryDiagnostics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DiscoveryDiagnostics create() => DiscoveryDiagnostics._();
  @$core.override
  DiscoveryDiagnostics createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DiscoveryDiagnostics getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DiscoveryDiagnostics>(create);
  static DiscoveryDiagnostics? _defaultInstance;

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
  $core.String get localPartition => $_getSZ(4);
  @$pb.TagNumber(5)
  set localPartition($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasLocalPartition() => $_has(4);
  @$pb.TagNumber(5)
  void clearLocalPartition() => $_clearField(5);

  @$pb.TagNumber(6)
  $fixnum.Int64 get peers => $_getI64(5);
  @$pb.TagNumber(6)
  set peers($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(6)
  $core.bool hasPeers() => $_has(5);
  @$pb.TagNumber(6)
  void clearPeers() => $_clearField(6);

  @$pb.TagNumber(7)
  $fixnum.Int64 get peersDdisc => $_getI64(6);
  @$pb.TagNumber(7)
  set peersDdisc($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(7)
  $core.bool hasPeersDdisc() => $_has(6);
  @$pb.TagNumber(7)
  void clearPeersDdisc() => $_clearField(7);

  @$pb.TagNumber(8)
  $fixnum.Int64 get peersBep51 => $_getI64(7);
  @$pb.TagNumber(8)
  set peersBep51($fixnum.Int64 value) => $_setInt64(7, value);
  @$pb.TagNumber(8)
  $core.bool hasPeersBep51() => $_has(7);
  @$pb.TagNumber(8)
  void clearPeersBep51() => $_clearField(8);

  @$pb.TagNumber(9)
  $fixnum.Int64 get unknownHashes => $_getI64(8);
  @$pb.TagNumber(9)
  set unknownHashes($fixnum.Int64 value) => $_setInt64(8, value);
  @$pb.TagNumber(9)
  $core.bool hasUnknownHashes() => $_has(8);
  @$pb.TagNumber(9)
  void clearUnknownHashes() => $_clearField(9);
}

class DiscoveryMetricsResponse extends $pb.GeneratedMessage {
  factory DiscoveryMetricsResponse({
    DiscoveryDiagnostics? discovery,
  }) {
    final result = create();
    if (discovery != null) result.discovery = discovery;
    return result;
  }

  DiscoveryMetricsResponse._();

  factory DiscoveryMetricsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DiscoveryMetricsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DiscoveryMetricsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<DiscoveryDiagnostics>(1, _omitFieldNames ? '' : 'discovery',
        subBuilder: DiscoveryDiagnostics.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoveryMetricsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiscoveryMetricsResponse copyWith(
          void Function(DiscoveryMetricsResponse) updates) =>
      super.copyWith((message) => updates(message as DiscoveryMetricsResponse))
          as DiscoveryMetricsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DiscoveryMetricsResponse create() => DiscoveryMetricsResponse._();
  @$core.override
  DiscoveryMetricsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DiscoveryMetricsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DiscoveryMetricsResponse>(create);
  static DiscoveryMetricsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  DiscoveryDiagnostics get discovery => $_getN(0);
  @$pb.TagNumber(1)
  set discovery(DiscoveryDiagnostics value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDiscovery() => $_has(0);
  @$pb.TagNumber(1)
  void clearDiscovery() => $_clearField(1);
  @$pb.TagNumber(1)
  DiscoveryDiagnostics ensureDiscovery() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
