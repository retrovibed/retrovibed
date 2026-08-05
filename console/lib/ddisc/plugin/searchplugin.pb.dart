// This is a generated file - do not edit.
//
// Generated from searchplugin.proto.

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

class Plugin extends $pb.GeneratedMessage {
  factory Plugin({
    $core.String? id,
    $core.String? name,
    $fixnum.Int64? size,
    $core.String? installedAt,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (name != null) result.name = name;
    if (size != null) result.size = size;
    if (installedAt != null) result.installedAt = installedAt;
    return result;
  }

  Plugin._();

  factory Plugin.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Plugin.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Plugin',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'searchplugin'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'name')
    ..a<$fixnum.Int64>(3, _omitFieldNames ? '' : 'size', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOS(4, _omitFieldNames ? '' : 'installed_at')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Plugin clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Plugin copyWith(void Function(Plugin) updates) =>
      super.copyWith((message) => updates(message as Plugin)) as Plugin;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Plugin create() => Plugin._();
  @$core.override
  Plugin createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Plugin getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Plugin>(create);
  static Plugin? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(2)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(2)
  void clearName() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get size => $_getI64(2);
  @$pb.TagNumber(3)
  set size($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasSize() => $_has(2);
  @$pb.TagNumber(3)
  void clearSize() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get installedAt => $_getSZ(3);
  @$pb.TagNumber(4)
  set installedAt($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasInstalledAt() => $_has(3);
  @$pb.TagNumber(4)
  void clearInstalledAt() => $_clearField(4);
}

class PluginSearchRequest extends $pb.GeneratedMessage {
  factory PluginSearchRequest({
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = create();
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    return result;
  }

  PluginSearchRequest._();

  factory PluginSearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PluginSearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PluginSearchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'searchplugin'),
      createEmptyInstance: create)
    ..a<$fixnum.Int64>(1, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginSearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginSearchRequest copyWith(void Function(PluginSearchRequest) updates) =>
      super.copyWith((message) => updates(message as PluginSearchRequest))
          as PluginSearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PluginSearchRequest create() => PluginSearchRequest._();
  @$core.override
  PluginSearchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PluginSearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PluginSearchRequest>(create);
  static PluginSearchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get offset => $_getI64(0);
  @$pb.TagNumber(1)
  set offset($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasOffset() => $_has(0);
  @$pb.TagNumber(1)
  void clearOffset() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get limit => $_getI64(1);
  @$pb.TagNumber(2)
  set limit($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(2)
  void clearLimit() => $_clearField(2);
}

class PluginSearchResponse extends $pb.GeneratedMessage {
  factory PluginSearchResponse({
    PluginSearchRequest? next,
    $core.Iterable<Plugin>? items,
  }) {
    final result = create();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  PluginSearchResponse._();

  factory PluginSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PluginSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PluginSearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'searchplugin'),
      createEmptyInstance: create)
    ..aOM<PluginSearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: PluginSearchRequest.create)
    ..pPM<Plugin>(2, _omitFieldNames ? '' : 'items', subBuilder: Plugin.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginSearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginSearchResponse copyWith(void Function(PluginSearchResponse) updates) =>
      super.copyWith((message) => updates(message as PluginSearchResponse))
          as PluginSearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PluginSearchResponse create() => PluginSearchResponse._();
  @$core.override
  PluginSearchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PluginSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PluginSearchResponse>(create);
  static PluginSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  PluginSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(PluginSearchRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  PluginSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<Plugin> get items => $_getList(1);
}

class PluginCreateRequest extends $pb.GeneratedMessage {
  factory PluginCreateRequest() => create();

  PluginCreateRequest._();

  factory PluginCreateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PluginCreateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PluginCreateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'searchplugin'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginCreateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginCreateRequest copyWith(void Function(PluginCreateRequest) updates) =>
      super.copyWith((message) => updates(message as PluginCreateRequest))
          as PluginCreateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PluginCreateRequest create() => PluginCreateRequest._();
  @$core.override
  PluginCreateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PluginCreateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PluginCreateRequest>(create);
  static PluginCreateRequest? _defaultInstance;
}

class PluginCreateResponse extends $pb.GeneratedMessage {
  factory PluginCreateResponse({
    Plugin? plugin,
  }) {
    final result = create();
    if (plugin != null) result.plugin = plugin;
    return result;
  }

  PluginCreateResponse._();

  factory PluginCreateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PluginCreateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PluginCreateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'searchplugin'),
      createEmptyInstance: create)
    ..aOM<Plugin>(1, _omitFieldNames ? '' : 'plugin', subBuilder: Plugin.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginCreateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginCreateResponse copyWith(void Function(PluginCreateResponse) updates) =>
      super.copyWith((message) => updates(message as PluginCreateResponse))
          as PluginCreateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PluginCreateResponse create() => PluginCreateResponse._();
  @$core.override
  PluginCreateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PluginCreateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PluginCreateResponse>(create);
  static PluginCreateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Plugin get plugin => $_getN(0);
  @$pb.TagNumber(1)
  set plugin(Plugin value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPlugin() => $_has(0);
  @$pb.TagNumber(1)
  void clearPlugin() => $_clearField(1);
  @$pb.TagNumber(1)
  Plugin ensurePlugin() => $_ensure(0);
}

class PluginFindRequest extends $pb.GeneratedMessage {
  factory PluginFindRequest() => create();

  PluginFindRequest._();

  factory PluginFindRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PluginFindRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PluginFindRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'searchplugin'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginFindRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginFindRequest copyWith(void Function(PluginFindRequest) updates) =>
      super.copyWith((message) => updates(message as PluginFindRequest))
          as PluginFindRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PluginFindRequest create() => PluginFindRequest._();
  @$core.override
  PluginFindRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PluginFindRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PluginFindRequest>(create);
  static PluginFindRequest? _defaultInstance;
}

class PluginFindResponse extends $pb.GeneratedMessage {
  factory PluginFindResponse({
    Plugin? plugin,
  }) {
    final result = create();
    if (plugin != null) result.plugin = plugin;
    return result;
  }

  PluginFindResponse._();

  factory PluginFindResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PluginFindResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PluginFindResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'searchplugin'),
      createEmptyInstance: create)
    ..aOM<Plugin>(1, _omitFieldNames ? '' : 'plugin', subBuilder: Plugin.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginFindResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginFindResponse copyWith(void Function(PluginFindResponse) updates) =>
      super.copyWith((message) => updates(message as PluginFindResponse))
          as PluginFindResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PluginFindResponse create() => PluginFindResponse._();
  @$core.override
  PluginFindResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PluginFindResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PluginFindResponse>(create);
  static PluginFindResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Plugin get plugin => $_getN(0);
  @$pb.TagNumber(1)
  set plugin(Plugin value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPlugin() => $_has(0);
  @$pb.TagNumber(1)
  void clearPlugin() => $_clearField(1);
  @$pb.TagNumber(1)
  Plugin ensurePlugin() => $_ensure(0);
}

class PluginDeleteRequest extends $pb.GeneratedMessage {
  factory PluginDeleteRequest() => create();

  PluginDeleteRequest._();

  factory PluginDeleteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PluginDeleteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PluginDeleteRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'searchplugin'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginDeleteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginDeleteRequest copyWith(void Function(PluginDeleteRequest) updates) =>
      super.copyWith((message) => updates(message as PluginDeleteRequest))
          as PluginDeleteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PluginDeleteRequest create() => PluginDeleteRequest._();
  @$core.override
  PluginDeleteRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PluginDeleteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PluginDeleteRequest>(create);
  static PluginDeleteRequest? _defaultInstance;
}

class PluginDeleteResponse extends $pb.GeneratedMessage {
  factory PluginDeleteResponse({
    Plugin? plugin,
  }) {
    final result = create();
    if (plugin != null) result.plugin = plugin;
    return result;
  }

  PluginDeleteResponse._();

  factory PluginDeleteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PluginDeleteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PluginDeleteResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'searchplugin'),
      createEmptyInstance: create)
    ..aOM<Plugin>(1, _omitFieldNames ? '' : 'plugin', subBuilder: Plugin.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginDeleteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PluginDeleteResponse copyWith(void Function(PluginDeleteResponse) updates) =>
      super.copyWith((message) => updates(message as PluginDeleteResponse))
          as PluginDeleteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PluginDeleteResponse create() => PluginDeleteResponse._();
  @$core.override
  PluginDeleteResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PluginDeleteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PluginDeleteResponse>(create);
  static PluginDeleteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Plugin get plugin => $_getN(0);
  @$pb.TagNumber(1)
  set plugin(Plugin value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPlugin() => $_has(0);
  @$pb.TagNumber(1)
  void clearPlugin() => $_clearField(1);
  @$pb.TagNumber(1)
  Plugin ensurePlugin() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
