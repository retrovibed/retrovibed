// This is a generated file - do not edit.
//
// Generated from meta.dht.proto.

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

class DHTDiagnostics extends $pb.GeneratedMessage {
  factory DHTDiagnostics({
    $core.String? nodeId,
    $core.int? goodNodes,
    $core.int? nodes,
    $core.int? outstandingTransactions,
    $fixnum.Int64? successfulOutboundAnnouncePeerQueries,
    $core.int? badNodes,
    $fixnum.Int64? outboundQueriesAttempted,
  }) {
    final result = create();
    if (nodeId != null) result.nodeId = nodeId;
    if (goodNodes != null) result.goodNodes = goodNodes;
    if (nodes != null) result.nodes = nodes;
    if (outstandingTransactions != null)
      result.outstandingTransactions = outstandingTransactions;
    if (successfulOutboundAnnouncePeerQueries != null)
      result.successfulOutboundAnnouncePeerQueries =
          successfulOutboundAnnouncePeerQueries;
    if (badNodes != null) result.badNodes = badNodes;
    if (outboundQueriesAttempted != null)
      result.outboundQueriesAttempted = outboundQueriesAttempted;
    return result;
  }

  DHTDiagnostics._();

  factory DHTDiagnostics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DHTDiagnostics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DHTDiagnostics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'node_id')
    ..aI(2, _omitFieldNames ? '' : 'good_nodes')
    ..aI(3, _omitFieldNames ? '' : 'nodes')
    ..aI(4, _omitFieldNames ? '' : 'outstanding_transactions')
    ..aInt64(
        5, _omitFieldNames ? '' : 'successful_outbound_announce_peer_queries')
    ..aI(6, _omitFieldNames ? '' : 'bad_nodes', fieldType: $pb.PbFieldType.OU3)
    ..aInt64(7, _omitFieldNames ? '' : 'outbound_queries_attempted')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DHTDiagnostics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DHTDiagnostics copyWith(void Function(DHTDiagnostics) updates) =>
      super.copyWith((message) => updates(message as DHTDiagnostics))
          as DHTDiagnostics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DHTDiagnostics create() => DHTDiagnostics._();
  @$core.override
  DHTDiagnostics createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DHTDiagnostics getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DHTDiagnostics>(create);
  static DHTDiagnostics? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get nodeId => $_getSZ(0);
  @$pb.TagNumber(1)
  set nodeId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasNodeId() => $_has(0);
  @$pb.TagNumber(1)
  void clearNodeId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get goodNodes => $_getIZ(1);
  @$pb.TagNumber(2)
  set goodNodes($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasGoodNodes() => $_has(1);
  @$pb.TagNumber(2)
  void clearGoodNodes() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.int get nodes => $_getIZ(2);
  @$pb.TagNumber(3)
  set nodes($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasNodes() => $_has(2);
  @$pb.TagNumber(3)
  void clearNodes() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.int get outstandingTransactions => $_getIZ(3);
  @$pb.TagNumber(4)
  set outstandingTransactions($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(4)
  $core.bool hasOutstandingTransactions() => $_has(3);
  @$pb.TagNumber(4)
  void clearOutstandingTransactions() => $_clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get successfulOutboundAnnouncePeerQueries => $_getI64(4);
  @$pb.TagNumber(5)
  set successfulOutboundAnnouncePeerQueries($fixnum.Int64 value) =>
      $_setInt64(4, value);
  @$pb.TagNumber(5)
  $core.bool hasSuccessfulOutboundAnnouncePeerQueries() => $_has(4);
  @$pb.TagNumber(5)
  void clearSuccessfulOutboundAnnouncePeerQueries() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.int get badNodes => $_getIZ(5);
  @$pb.TagNumber(6)
  set badNodes($core.int value) => $_setUnsignedInt32(5, value);
  @$pb.TagNumber(6)
  $core.bool hasBadNodes() => $_has(5);
  @$pb.TagNumber(6)
  void clearBadNodes() => $_clearField(6);

  @$pb.TagNumber(7)
  $fixnum.Int64 get outboundQueriesAttempted => $_getI64(6);
  @$pb.TagNumber(7)
  set outboundQueriesAttempted($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(7)
  $core.bool hasOutboundQueriesAttempted() => $_has(6);
  @$pb.TagNumber(7)
  void clearOutboundQueriesAttempted() => $_clearField(7);
}

class DHTMetricsResponse extends $pb.GeneratedMessage {
  factory DHTMetricsResponse({
    DHTDiagnostics? dht,
  }) {
    final result = create();
    if (dht != null) result.dht = dht;
    return result;
  }

  DHTMetricsResponse._();

  factory DHTMetricsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DHTMetricsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DHTMetricsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<DHTDiagnostics>(1, _omitFieldNames ? '' : 'dht',
        subBuilder: DHTDiagnostics.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DHTMetricsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DHTMetricsResponse copyWith(void Function(DHTMetricsResponse) updates) =>
      super.copyWith((message) => updates(message as DHTMetricsResponse))
          as DHTMetricsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DHTMetricsResponse create() => DHTMetricsResponse._();
  @$core.override
  DHTMetricsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DHTMetricsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DHTMetricsResponse>(create);
  static DHTMetricsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  DHTDiagnostics get dht => $_getN(0);
  @$pb.TagNumber(1)
  set dht(DHTDiagnostics value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDht() => $_has(0);
  @$pb.TagNumber(1)
  void clearDht() => $_clearField(1);
  @$pb.TagNumber(1)
  DHTDiagnostics ensureDht() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
