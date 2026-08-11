// This is a generated file - do not edit.
//
// Generated from media.remote.control.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import 'media.pb.dart' as $0;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

/// enqueue a bit of media for playback.
class Queue extends $pb.GeneratedMessage {
  factory Queue({
    $0.Media? media,
  }) {
    final result = create();
    if (media != null) result.media = media;
    return result;
  }

  Queue._();

  factory Queue.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Queue.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Queue',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<$0.Media>(1, _omitFieldNames ? '' : 'media',
        subBuilder: $0.Media.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Queue clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Queue copyWith(void Function(Queue) updates) =>
      super.copyWith((message) => updates(message as Queue)) as Queue;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Queue create() => Queue._();
  @$core.override
  Queue createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Queue getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Queue>(create);
  static Queue? _defaultInstance;

  @$pb.TagNumber(1)
  $0.Media get media => $_getN(0);
  @$pb.TagNumber(1)
  set media($0.Media value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMedia() => $_has(0);
  @$pb.TagNumber(1)
  void clearMedia() => $_clearField(1);
  @$pb.TagNumber(1)
  $0.Media ensureMedia() => $_ensure(0);
}

/// dequeue the media with the given id.
class Dequeue extends $pb.GeneratedMessage {
  factory Dequeue({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  Dequeue._();

  factory Dequeue.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Dequeue.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Dequeue',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Dequeue clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Dequeue copyWith(void Function(Dequeue) updates) =>
      super.copyWith((message) => updates(message as Dequeue)) as Dequeue;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Dequeue create() => Dequeue._();
  @$core.override
  Dequeue createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Dequeue getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Dequeue>(create);
  static Dequeue? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);
}

/// playpause commands the device to toggle media playback.
/// the paused flag is to
class PlayPause extends $pb.GeneratedMessage {
  factory PlayPause({
    $core.bool? paused,
  }) {
    final result = create();
    if (paused != null) result.paused = paused;
    return result;
  }

  PlayPause._();

  factory PlayPause.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PlayPause.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PlayPause',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'paused')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PlayPause clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PlayPause copyWith(void Function(PlayPause) updates) =>
      super.copyWith((message) => updates(message as PlayPause)) as PlayPause;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PlayPause create() => PlayPause._();
  @$core.override
  PlayPause createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PlayPause getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<PlayPause>(create);
  static PlayPause? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get paused => $_getBF(0);
  @$pb.TagNumber(1)
  set paused($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasPaused() => $_has(0);
  @$pb.TagNumber(1)
  void clearPaused() => $_clearField(1);
}

/// Seek - move forward/back the given amount of time (milliseconds), relative
/// to the current position. offset == int32 max/min is a sentinel meaning
/// "skip to next/previous track" rather than a literal seek.
class Seek extends $pb.GeneratedMessage {
  factory Seek({
    $core.int? offset,
  }) {
    final result = create();
    if (offset != null) result.offset = offset;
    return result;
  }

  Seek._();

  factory Seek.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Seek.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Seek',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'offset')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Seek clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Seek copyWith(void Function(Seek) updates) =>
      super.copyWith((message) => updates(message as Seek)) as Seek;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Seek create() => Seek._();
  @$core.override
  Seek createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Seek getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Seek>(create);
  static Seek? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get offset => $_getIZ(0);
  @$pb.TagNumber(1)
  set offset($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasOffset() => $_has(0);
  @$pb.TagNumber(1)
  void clearOffset() => $_clearField(1);
}

enum Stream_Command { queue, dequeue, playpause, seek, notSet }

/// represents a stream of commands / responses for the remote control.
/// each command / response will contain a 'sid' representing the sequentialish
/// id.
class Stream extends $pb.GeneratedMessage {
  factory Stream({
    $core.String? sid,
    Queue? queue,
    Dequeue? dequeue,
    PlayPause? playpause,
    Seek? seek,
  }) {
    final result = create();
    if (sid != null) result.sid = sid;
    if (queue != null) result.queue = queue;
    if (dequeue != null) result.dequeue = dequeue;
    if (playpause != null) result.playpause = playpause;
    if (seek != null) result.seek = seek;
    return result;
  }

  Stream._();

  factory Stream.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Stream.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static const $core.Map<$core.int, Stream_Command> _Stream_CommandByTag = {
    1000: Stream_Command.queue,
    1002: Stream_Command.dequeue,
    1003: Stream_Command.playpause,
    1004: Stream_Command.seek,
    0: Stream_Command.notSet
  };
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Stream',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..oo(0, [1000, 1002, 1003, 1004])
    ..aOS(1, _omitFieldNames ? '' : 'sid')
    ..aOM<Queue>(1000, _omitFieldNames ? '' : 'queue', subBuilder: Queue.create)
    ..aOM<Dequeue>(1002, _omitFieldNames ? '' : 'dequeue',
        subBuilder: Dequeue.create)
    ..aOM<PlayPause>(1003, _omitFieldNames ? '' : 'playpause',
        subBuilder: PlayPause.create)
    ..aOM<Seek>(1004, _omitFieldNames ? '' : 'seek', subBuilder: Seek.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Stream clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Stream copyWith(void Function(Stream) updates) =>
      super.copyWith((message) => updates(message as Stream)) as Stream;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Stream create() => Stream._();
  @$core.override
  Stream createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Stream getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Stream>(create);
  static Stream? _defaultInstance;

  @$pb.TagNumber(1000)
  @$pb.TagNumber(1002)
  @$pb.TagNumber(1003)
  @$pb.TagNumber(1004)
  Stream_Command whichCommand() => _Stream_CommandByTag[$_whichOneof(0)]!;
  @$pb.TagNumber(1000)
  @$pb.TagNumber(1002)
  @$pb.TagNumber(1003)
  @$pb.TagNumber(1004)
  void clearCommand() => $_clearField($_whichOneof(0));

  @$pb.TagNumber(1)
  $core.String get sid => $_getSZ(0);
  @$pb.TagNumber(1)
  set sid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSid() => $_has(0);
  @$pb.TagNumber(1)
  void clearSid() => $_clearField(1);

  @$pb.TagNumber(1000)
  Queue get queue => $_getN(1);
  @$pb.TagNumber(1000)
  set queue(Queue value) => $_setField(1000, value);
  @$pb.TagNumber(1000)
  $core.bool hasQueue() => $_has(1);
  @$pb.TagNumber(1000)
  void clearQueue() => $_clearField(1000);
  @$pb.TagNumber(1000)
  Queue ensureQueue() => $_ensure(1);

  @$pb.TagNumber(1002)
  Dequeue get dequeue => $_getN(2);
  @$pb.TagNumber(1002)
  set dequeue(Dequeue value) => $_setField(1002, value);
  @$pb.TagNumber(1002)
  $core.bool hasDequeue() => $_has(2);
  @$pb.TagNumber(1002)
  void clearDequeue() => $_clearField(1002);
  @$pb.TagNumber(1002)
  Dequeue ensureDequeue() => $_ensure(2);

  @$pb.TagNumber(1003)
  PlayPause get playpause => $_getN(3);
  @$pb.TagNumber(1003)
  set playpause(PlayPause value) => $_setField(1003, value);
  @$pb.TagNumber(1003)
  $core.bool hasPlaypause() => $_has(3);
  @$pb.TagNumber(1003)
  void clearPlaypause() => $_clearField(1003);
  @$pb.TagNumber(1003)
  PlayPause ensurePlaypause() => $_ensure(3);

  @$pb.TagNumber(1004)
  Seek get seek => $_getN(4);
  @$pb.TagNumber(1004)
  set seek(Seek value) => $_setField(1004, value);
  @$pb.TagNumber(1004)
  $core.bool hasSeek() => $_has(4);
  @$pb.TagNumber(1004)
  void clearSeek() => $_clearField(1004);
  @$pb.TagNumber(1004)
  Seek ensureSeek() => $_ensure(4);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
