// This is a generated file - do not edit.
//
// Generated from audio/meta.audio.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class AudioSink extends $pb.GeneratedMessage {
  factory AudioSink({
    $core.String? id,
    $core.String? name,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (name != null) result.name = name;
    return result;
  }

  AudioSink._();

  factory AudioSink.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AudioSink.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AudioSink',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AudioSink clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AudioSink copyWith(void Function(AudioSink) updates) =>
      super.copyWith((message) => updates(message as AudioSink)) as AudioSink;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AudioSink create() => AudioSink._();
  @$core.override
  AudioSink createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AudioSink getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AudioSink>(create);
  static AudioSink? _defaultInstance;

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
}

class AudioSinkSearchResponse extends $pb.GeneratedMessage {
  factory AudioSinkSearchResponse({
    $core.Iterable<AudioSink>? items,
  }) {
    final result = create();
    if (items != null) result.items.addAll(items);
    return result;
  }

  AudioSinkSearchResponse._();

  factory AudioSinkSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AudioSinkSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AudioSinkSearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..pPM<AudioSink>(1, _omitFieldNames ? '' : 'items',
        subBuilder: AudioSink.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AudioSinkSearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AudioSinkSearchResponse copyWith(
          void Function(AudioSinkSearchResponse) updates) =>
      super.copyWith((message) => updates(message as AudioSinkSearchResponse))
          as AudioSinkSearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AudioSinkSearchResponse create() => AudioSinkSearchResponse._();
  @$core.override
  AudioSinkSearchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AudioSinkSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AudioSinkSearchResponse>(create);
  static AudioSinkSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<AudioSink> get items => $_getList(0);
}

class AudioSinkCurrentResponse extends $pb.GeneratedMessage {
  factory AudioSinkCurrentResponse({
    AudioSink? sink,
  }) {
    final result = create();
    if (sink != null) result.sink = sink;
    return result;
  }

  AudioSinkCurrentResponse._();

  factory AudioSinkCurrentResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AudioSinkCurrentResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AudioSinkCurrentResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<AudioSink>(1, _omitFieldNames ? '' : 'sink',
        subBuilder: AudioSink.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AudioSinkCurrentResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AudioSinkCurrentResponse copyWith(
          void Function(AudioSinkCurrentResponse) updates) =>
      super.copyWith((message) => updates(message as AudioSinkCurrentResponse))
          as AudioSinkCurrentResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AudioSinkCurrentResponse create() => AudioSinkCurrentResponse._();
  @$core.override
  AudioSinkCurrentResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AudioSinkCurrentResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AudioSinkCurrentResponse>(create);
  static AudioSinkCurrentResponse? _defaultInstance;

  @$pb.TagNumber(1)
  AudioSink get sink => $_getN(0);
  @$pb.TagNumber(1)
  set sink(AudioSink value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasSink() => $_has(0);
  @$pb.TagNumber(1)
  void clearSink() => $_clearField(1);
  @$pb.TagNumber(1)
  AudioSink ensureSink() => $_ensure(0);
}

class AudioSinkTouchRequest extends $pb.GeneratedMessage {
  factory AudioSinkTouchRequest({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  AudioSinkTouchRequest._();

  factory AudioSinkTouchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AudioSinkTouchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AudioSinkTouchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AudioSinkTouchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AudioSinkTouchRequest copyWith(
          void Function(AudioSinkTouchRequest) updates) =>
      super.copyWith((message) => updates(message as AudioSinkTouchRequest))
          as AudioSinkTouchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AudioSinkTouchRequest create() => AudioSinkTouchRequest._();
  @$core.override
  AudioSinkTouchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AudioSinkTouchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AudioSinkTouchRequest>(create);
  static AudioSinkTouchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);
}

class AudioSinkTouchResponse extends $pb.GeneratedMessage {
  factory AudioSinkTouchResponse({
    AudioSink? sink,
  }) {
    final result = create();
    if (sink != null) result.sink = sink;
    return result;
  }

  AudioSinkTouchResponse._();

  factory AudioSinkTouchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AudioSinkTouchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AudioSinkTouchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<AudioSink>(1, _omitFieldNames ? '' : 'sink',
        subBuilder: AudioSink.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AudioSinkTouchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AudioSinkTouchResponse copyWith(
          void Function(AudioSinkTouchResponse) updates) =>
      super.copyWith((message) => updates(message as AudioSinkTouchResponse))
          as AudioSinkTouchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AudioSinkTouchResponse create() => AudioSinkTouchResponse._();
  @$core.override
  AudioSinkTouchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AudioSinkTouchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AudioSinkTouchResponse>(create);
  static AudioSinkTouchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  AudioSink get sink => $_getN(0);
  @$pb.TagNumber(1)
  set sink(AudioSink value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasSink() => $_has(0);
  @$pb.TagNumber(1)
  void clearSink() => $_clearField(1);
  @$pb.TagNumber(1)
  AudioSink ensureSink() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
