// This is a generated file - do not edit.
//
// Generated from meta/meta.network.proto.

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

class WireguardDiagnostics extends $pb.GeneratedMessage {
  factory WireguardDiagnostics({
    $core.String? peerKey,
    $fixnum.Int64? keepaliveInterval,
    $fixnum.Int64? txBytes,
    $fixnum.Int64? rxBytes,
    $fixnum.Int64? lastHandshakeSec,
    $core.String? status,
  }) {
    final result = create();
    if (peerKey != null) result.peerKey = peerKey;
    if (keepaliveInterval != null) result.keepaliveInterval = keepaliveInterval;
    if (txBytes != null) result.txBytes = txBytes;
    if (rxBytes != null) result.rxBytes = rxBytes;
    if (lastHandshakeSec != null) result.lastHandshakeSec = lastHandshakeSec;
    if (status != null) result.status = status;
    return result;
  }

  WireguardDiagnostics._();

  factory WireguardDiagnostics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WireguardDiagnostics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardDiagnostics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'peer_key')
    ..a<$fixnum.Int64>(
        2, _omitFieldNames ? '' : 'keepalive_interval', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        3, _omitFieldNames ? '' : 'tx_bytes', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        4, _omitFieldNames ? '' : 'rx_bytes', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aInt64(5, _omitFieldNames ? '' : 'last_handshake_sec')
    ..aOS(6, _omitFieldNames ? '' : 'status')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardDiagnostics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardDiagnostics copyWith(void Function(WireguardDiagnostics) updates) =>
      super.copyWith((message) => updates(message as WireguardDiagnostics))
          as WireguardDiagnostics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WireguardDiagnostics create() => WireguardDiagnostics._();
  @$core.override
  WireguardDiagnostics createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WireguardDiagnostics getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardDiagnostics>(create);
  static WireguardDiagnostics? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get peerKey => $_getSZ(0);
  @$pb.TagNumber(1)
  set peerKey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasPeerKey() => $_has(0);
  @$pb.TagNumber(1)
  void clearPeerKey() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get keepaliveInterval => $_getI64(1);
  @$pb.TagNumber(2)
  set keepaliveInterval($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasKeepaliveInterval() => $_has(1);
  @$pb.TagNumber(2)
  void clearKeepaliveInterval() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get txBytes => $_getI64(2);
  @$pb.TagNumber(3)
  set txBytes($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasTxBytes() => $_has(2);
  @$pb.TagNumber(3)
  void clearTxBytes() => $_clearField(3);

  @$pb.TagNumber(4)
  $fixnum.Int64 get rxBytes => $_getI64(3);
  @$pb.TagNumber(4)
  set rxBytes($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(4)
  $core.bool hasRxBytes() => $_has(3);
  @$pb.TagNumber(4)
  void clearRxBytes() => $_clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get lastHandshakeSec => $_getI64(4);
  @$pb.TagNumber(5)
  set lastHandshakeSec($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(5)
  $core.bool hasLastHandshakeSec() => $_has(4);
  @$pb.TagNumber(5)
  void clearLastHandshakeSec() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get status => $_getSZ(5);
  @$pb.TagNumber(6)
  set status($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasStatus() => $_has(5);
  @$pb.TagNumber(6)
  void clearStatus() => $_clearField(6);
}

class NetworkInterface extends $pb.GeneratedMessage {
  factory NetworkInterface({
    $core.String? name,
    $core.String? ip,
    $core.bool? metered,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (ip != null) result.ip = ip;
    if (metered != null) result.metered = metered;
    return result;
  }

  NetworkInterface._();

  factory NetworkInterface.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NetworkInterface.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NetworkInterface',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'name')
    ..aOS(2, _omitFieldNames ? '' : 'ip')
    ..aOB(3, _omitFieldNames ? '' : 'metered')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NetworkInterface clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NetworkInterface copyWith(void Function(NetworkInterface) updates) =>
      super.copyWith((message) => updates(message as NetworkInterface))
          as NetworkInterface;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NetworkInterface create() => NetworkInterface._();
  @$core.override
  NetworkInterface createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NetworkInterface getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NetworkInterface>(create);
  static NetworkInterface? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(1)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearName() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get ip => $_getSZ(1);
  @$pb.TagNumber(2)
  set ip($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasIp() => $_has(1);
  @$pb.TagNumber(2)
  void clearIp() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.bool get metered => $_getBF(2);
  @$pb.TagNumber(3)
  set metered($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(3)
  $core.bool hasMetered() => $_has(2);
  @$pb.TagNumber(3)
  void clearMetered() => $_clearField(3);
}

class Network extends $pb.GeneratedMessage {
  factory Network({
    $core.Iterable<NetworkInterface>? interfaces,
    $core.bool? haveV4,
    $core.bool? haveV6,
    $core.String? defaultInterface,
  }) {
    final result = create();
    if (interfaces != null) result.interfaces.addAll(interfaces);
    if (haveV4 != null) result.haveV4 = haveV4;
    if (haveV6 != null) result.haveV6 = haveV6;
    if (defaultInterface != null) result.defaultInterface = defaultInterface;
    return result;
  }

  Network._();

  factory Network.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Network.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Network',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..pPM<NetworkInterface>(1, _omitFieldNames ? '' : 'interfaces',
        subBuilder: NetworkInterface.create)
    ..aOB(2, _omitFieldNames ? '' : 'have_v4')
    ..aOB(3, _omitFieldNames ? '' : 'have_v6')
    ..aOS(4, _omitFieldNames ? '' : 'default_interface')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Network clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Network copyWith(void Function(Network) updates) =>
      super.copyWith((message) => updates(message as Network)) as Network;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Network create() => Network._();
  @$core.override
  Network createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Network getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Network>(create);
  static Network? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<NetworkInterface> get interfaces => $_getList(0);

  @$pb.TagNumber(2)
  $core.bool get haveV4 => $_getBF(1);
  @$pb.TagNumber(2)
  set haveV4($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasHaveV4() => $_has(1);
  @$pb.TagNumber(2)
  void clearHaveV4() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.bool get haveV6 => $_getBF(2);
  @$pb.TagNumber(3)
  set haveV6($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(3)
  $core.bool hasHaveV6() => $_has(2);
  @$pb.TagNumber(3)
  void clearHaveV6() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get defaultInterface => $_getSZ(3);
  @$pb.TagNumber(4)
  set defaultInterface($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasDefaultInterface() => $_has(3);
  @$pb.TagNumber(4)
  void clearDefaultInterface() => $_clearField(4);
}

class NetworkMetricsResponse extends $pb.GeneratedMessage {
  factory NetworkMetricsResponse({
    WireguardDiagnostics? wireguard,
    Network? network,
  }) {
    final result = create();
    if (wireguard != null) result.wireguard = wireguard;
    if (network != null) result.network = network;
    return result;
  }

  NetworkMetricsResponse._();

  factory NetworkMetricsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NetworkMetricsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NetworkMetricsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<WireguardDiagnostics>(1, _omitFieldNames ? '' : 'wireguard',
        subBuilder: WireguardDiagnostics.create)
    ..aOM<Network>(2, _omitFieldNames ? '' : 'network',
        subBuilder: Network.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NetworkMetricsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NetworkMetricsResponse copyWith(
          void Function(NetworkMetricsResponse) updates) =>
      super.copyWith((message) => updates(message as NetworkMetricsResponse))
          as NetworkMetricsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NetworkMetricsResponse create() => NetworkMetricsResponse._();
  @$core.override
  NetworkMetricsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NetworkMetricsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NetworkMetricsResponse>(create);
  static NetworkMetricsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  WireguardDiagnostics get wireguard => $_getN(0);
  @$pb.TagNumber(1)
  set wireguard(WireguardDiagnostics value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasWireguard() => $_has(0);
  @$pb.TagNumber(1)
  void clearWireguard() => $_clearField(1);
  @$pb.TagNumber(1)
  WireguardDiagnostics ensureWireguard() => $_ensure(0);

  @$pb.TagNumber(2)
  Network get network => $_getN(1);
  @$pb.TagNumber(2)
  set network(Network value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasNetwork() => $_has(1);
  @$pb.TagNumber(2)
  void clearNetwork() => $_clearField(2);
  @$pb.TagNumber(2)
  Network ensureNetwork() => $_ensure(1);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
