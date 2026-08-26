// This is a generated file - do not edit.
//
// Generated from media/media.remote.control.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import '../meta/meta.daemon.pb.dart' as $1;
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

/// pause commands the device to toggle media playback.
class Pause extends $pb.GeneratedMessage {
  factory Pause() => create();

  Pause._();

  factory Pause.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Pause.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Pause',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Pause clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Pause copyWith(void Function(Pause) updates) =>
      super.copyWith((message) => updates(message as Pause)) as Pause;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Pause create() => Pause._();
  @$core.override
  Pause createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Pause getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Pause>(create);
  static Pause? _defaultInstance;
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

/// Sync requests (when sent with no fields set) or reports (when sent
/// populated) the listener's current library and playback queue. A connect
/// client sends an empty Sync to ask the listener to reply with its state;
/// the listener responds on the same tag with library + queue populated.
class Sync extends $pb.GeneratedMessage {
  factory Sync({
    $1.Daemon? library,
    $core.int? capacity,
    $0.Media? current,
    $core.String? token,
    $fixnum.Int64? expiration,
    $core.double? volume,
    $core.bool? muted,
    $core.bool? paused,
    $core.bool? fullscreen,
    $core.Iterable<$0.Media>? queue,
  }) {
    final result = create();
    if (library != null) result.library = library;
    if (capacity != null) result.capacity = capacity;
    if (current != null) result.current = current;
    if (token != null) result.token = token;
    if (expiration != null) result.expiration = expiration;
    if (volume != null) result.volume = volume;
    if (muted != null) result.muted = muted;
    if (paused != null) result.paused = paused;
    if (fullscreen != null) result.fullscreen = fullscreen;
    if (queue != null) result.queue.addAll(queue);
    return result;
  }

  Sync._();

  factory Sync.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Sync.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Sync',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<$1.Daemon>(1, _omitFieldNames ? '' : 'library',
        subBuilder: $1.Daemon.create)
    ..aI(2, _omitFieldNames ? '' : 'capacity', fieldType: $pb.PbFieldType.OU3)
    ..aOM<$0.Media>(3, _omitFieldNames ? '' : 'current',
        subBuilder: $0.Media.create)
    ..aOS(4, _omitFieldNames ? '' : 'token')
    ..aInt64(5, _omitFieldNames ? '' : 'expiration')
    ..aD(6, _omitFieldNames ? '' : 'volume', fieldType: $pb.PbFieldType.OF)
    ..aOB(7, _omitFieldNames ? '' : 'muted')
    ..aOB(8, _omitFieldNames ? '' : 'paused')
    ..aOB(9, _omitFieldNames ? '' : 'fullscreen')
    ..pPM<$0.Media>(1000, _omitFieldNames ? '' : 'queue',
        subBuilder: $0.Media.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Sync clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Sync copyWith(void Function(Sync) updates) =>
      super.copyWith((message) => updates(message as Sync)) as Sync;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Sync create() => Sync._();
  @$core.override
  Sync createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Sync getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Sync>(create);
  static Sync? _defaultInstance;

  @$pb.TagNumber(1)
  $1.Daemon get library => $_getN(0);
  @$pb.TagNumber(1)
  set library($1.Daemon value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasLibrary() => $_has(0);
  @$pb.TagNumber(1)
  void clearLibrary() => $_clearField(1);
  @$pb.TagNumber(1)
  $1.Daemon ensureLibrary() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.int get capacity => $_getIZ(1);
  @$pb.TagNumber(2)
  set capacity($core.int value) => $_setUnsignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasCapacity() => $_has(1);
  @$pb.TagNumber(2)
  void clearCapacity() => $_clearField(2);

  @$pb.TagNumber(3)
  $0.Media get current => $_getN(2);
  @$pb.TagNumber(3)
  set current($0.Media value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasCurrent() => $_has(2);
  @$pb.TagNumber(3)
  void clearCurrent() => $_clearField(3);
  @$pb.TagNumber(3)
  $0.Media ensureCurrent() => $_ensure(2);

  @$pb.TagNumber(4)
  $core.String get token => $_getSZ(3);
  @$pb.TagNumber(4)
  set token($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasToken() => $_has(3);
  @$pb.TagNumber(4)
  void clearToken() => $_clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get expiration => $_getI64(4);
  @$pb.TagNumber(5)
  set expiration($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(5)
  $core.bool hasExpiration() => $_has(4);
  @$pb.TagNumber(5)
  void clearExpiration() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.double get volume => $_getN(5);
  @$pb.TagNumber(6)
  set volume($core.double value) => $_setFloat(5, value);
  @$pb.TagNumber(6)
  $core.bool hasVolume() => $_has(5);
  @$pb.TagNumber(6)
  void clearVolume() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.bool get muted => $_getBF(6);
  @$pb.TagNumber(7)
  set muted($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(7)
  $core.bool hasMuted() => $_has(6);
  @$pb.TagNumber(7)
  void clearMuted() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.bool get paused => $_getBF(7);
  @$pb.TagNumber(8)
  set paused($core.bool value) => $_setBool(7, value);
  @$pb.TagNumber(8)
  $core.bool hasPaused() => $_has(7);
  @$pb.TagNumber(8)
  void clearPaused() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.bool get fullscreen => $_getBF(8);
  @$pb.TagNumber(9)
  set fullscreen($core.bool value) => $_setBool(8, value);
  @$pb.TagNumber(9)
  $core.bool hasFullscreen() => $_has(8);
  @$pb.TagNumber(9)
  void clearFullscreen() => $_clearField(9);

  @$pb.TagNumber(1000)
  $pb.PbList<$0.Media> get queue => $_getList(9);
}

/// Fullscreen toggles fullscreen on the receiving device. No payload -
/// each Fullscreen command flips the device's current state; ordering
/// against concurrent/stale commands is resolved by the sender using
/// Stream.sid (a uuidv7) as a vector clock, same as Volume.
class Fullscreen extends $pb.GeneratedMessage {
  factory Fullscreen() => create();

  Fullscreen._();

  factory Fullscreen.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Fullscreen.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Fullscreen',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Fullscreen clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Fullscreen copyWith(void Function(Fullscreen) updates) =>
      super.copyWith((message) => updates(message as Fullscreen)) as Fullscreen;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Fullscreen create() => Fullscreen._();
  @$core.override
  Fullscreen createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Fullscreen getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Fullscreen>(create);
  static Fullscreen? _defaultInstance;
}

/// Mute toggles the receiving device's audio between silent and its prior
/// level. No payload - same shape/ordering semantics as Fullscreen.
class Mute extends $pb.GeneratedMessage {
  factory Mute() => create();

  Mute._();

  factory Mute.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Mute.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Mute',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Mute clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Mute copyWith(void Function(Mute) updates) =>
      super.copyWith((message) => updates(message as Mute)) as Mute;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Mute create() => Mute._();
  @$core.override
  Mute createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Mute getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Mute>(create);
  static Mute? _defaultInstance;
}

enum Stream_Command {
  queue,
  dequeue,
  pause,
  seek,
  sync,
  volume,
  fullscreen,
  mute,
  notSet
}

/// represents a stream of commands / responses for the remote control.
/// each command / response will contain a 'sid' representing the sequentialish
/// id.
class Stream extends $pb.GeneratedMessage {
  factory Stream({
    $core.String? sid,
    $fixnum.Int64? vid,
    Queue? queue,
    Dequeue? dequeue,
    Pause? pause,
    Seek? seek,
    Sync? sync,
    Seek? volume,
    Fullscreen? fullscreen,
    Mute? mute,
  }) {
    final result = create();
    if (sid != null) result.sid = sid;
    if (vid != null) result.vid = vid;
    if (queue != null) result.queue = queue;
    if (dequeue != null) result.dequeue = dequeue;
    if (pause != null) result.pause = pause;
    if (seek != null) result.seek = seek;
    if (sync != null) result.sync = sync;
    if (volume != null) result.volume = volume;
    if (fullscreen != null) result.fullscreen = fullscreen;
    if (mute != null) result.mute = mute;
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
    1003: Stream_Command.pause,
    1004: Stream_Command.seek,
    1005: Stream_Command.sync,
    1006: Stream_Command.volume,
    1007: Stream_Command.fullscreen,
    1008: Stream_Command.mute,
    0: Stream_Command.notSet
  };
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Stream',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..oo(0, [1000, 1002, 1003, 1004, 1005, 1006, 1007, 1008])
    ..aOS(1, _omitFieldNames ? '' : 'sid')
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'vid', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOM<Queue>(1000, _omitFieldNames ? '' : 'queue', subBuilder: Queue.create)
    ..aOM<Dequeue>(1002, _omitFieldNames ? '' : 'dequeue',
        subBuilder: Dequeue.create)
    ..aOM<Pause>(1003, _omitFieldNames ? '' : 'pause', subBuilder: Pause.create)
    ..aOM<Seek>(1004, _omitFieldNames ? '' : 'seek', subBuilder: Seek.create)
    ..aOM<Sync>(1005, _omitFieldNames ? '' : 'sync', subBuilder: Sync.create)
    ..aOM<Seek>(1006, _omitFieldNames ? '' : 'volume', subBuilder: Seek.create)
    ..aOM<Fullscreen>(1007, _omitFieldNames ? '' : 'fullscreen',
        subBuilder: Fullscreen.create)
    ..aOM<Mute>(1008, _omitFieldNames ? '' : 'mute', subBuilder: Mute.create)
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
  @$pb.TagNumber(1005)
  @$pb.TagNumber(1006)
  @$pb.TagNumber(1007)
  @$pb.TagNumber(1008)
  Stream_Command whichCommand() => _Stream_CommandByTag[$_whichOneof(0)]!;
  @$pb.TagNumber(1000)
  @$pb.TagNumber(1002)
  @$pb.TagNumber(1003)
  @$pb.TagNumber(1004)
  @$pb.TagNumber(1005)
  @$pb.TagNumber(1006)
  @$pb.TagNumber(1007)
  @$pb.TagNumber(1008)
  void clearCommand() => $_clearField($_whichOneof(0));

  @$pb.TagNumber(1)
  $core.String get sid => $_getSZ(0);
  @$pb.TagNumber(1)
  set sid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSid() => $_has(0);
  @$pb.TagNumber(1)
  void clearSid() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get vid => $_getI64(1);
  @$pb.TagNumber(2)
  set vid($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasVid() => $_has(1);
  @$pb.TagNumber(2)
  void clearVid() => $_clearField(2);

  @$pb.TagNumber(1000)
  Queue get queue => $_getN(2);
  @$pb.TagNumber(1000)
  set queue(Queue value) => $_setField(1000, value);
  @$pb.TagNumber(1000)
  $core.bool hasQueue() => $_has(2);
  @$pb.TagNumber(1000)
  void clearQueue() => $_clearField(1000);
  @$pb.TagNumber(1000)
  Queue ensureQueue() => $_ensure(2);

  @$pb.TagNumber(1002)
  Dequeue get dequeue => $_getN(3);
  @$pb.TagNumber(1002)
  set dequeue(Dequeue value) => $_setField(1002, value);
  @$pb.TagNumber(1002)
  $core.bool hasDequeue() => $_has(3);
  @$pb.TagNumber(1002)
  void clearDequeue() => $_clearField(1002);
  @$pb.TagNumber(1002)
  Dequeue ensureDequeue() => $_ensure(3);

  @$pb.TagNumber(1003)
  Pause get pause => $_getN(4);
  @$pb.TagNumber(1003)
  set pause(Pause value) => $_setField(1003, value);
  @$pb.TagNumber(1003)
  $core.bool hasPause() => $_has(4);
  @$pb.TagNumber(1003)
  void clearPause() => $_clearField(1003);
  @$pb.TagNumber(1003)
  Pause ensurePause() => $_ensure(4);

  @$pb.TagNumber(1004)
  Seek get seek => $_getN(5);
  @$pb.TagNumber(1004)
  set seek(Seek value) => $_setField(1004, value);
  @$pb.TagNumber(1004)
  $core.bool hasSeek() => $_has(5);
  @$pb.TagNumber(1004)
  void clearSeek() => $_clearField(1004);
  @$pb.TagNumber(1004)
  Seek ensureSeek() => $_ensure(5);

  @$pb.TagNumber(1005)
  Sync get sync => $_getN(6);
  @$pb.TagNumber(1005)
  set sync(Sync value) => $_setField(1005, value);
  @$pb.TagNumber(1005)
  $core.bool hasSync() => $_has(6);
  @$pb.TagNumber(1005)
  void clearSync() => $_clearField(1005);
  @$pb.TagNumber(1005)
  Sync ensureSync() => $_ensure(6);

  @$pb.TagNumber(1006)
  Seek get volume => $_getN(7);
  @$pb.TagNumber(1006)
  set volume(Seek value) => $_setField(1006, value);
  @$pb.TagNumber(1006)
  $core.bool hasVolume() => $_has(7);
  @$pb.TagNumber(1006)
  void clearVolume() => $_clearField(1006);
  @$pb.TagNumber(1006)
  Seek ensureVolume() => $_ensure(7);

  @$pb.TagNumber(1007)
  Fullscreen get fullscreen => $_getN(8);
  @$pb.TagNumber(1007)
  set fullscreen(Fullscreen value) => $_setField(1007, value);
  @$pb.TagNumber(1007)
  $core.bool hasFullscreen() => $_has(8);
  @$pb.TagNumber(1007)
  void clearFullscreen() => $_clearField(1007);
  @$pb.TagNumber(1007)
  Fullscreen ensureFullscreen() => $_ensure(8);

  @$pb.TagNumber(1008)
  Mute get mute => $_getN(9);
  @$pb.TagNumber(1008)
  set mute(Mute value) => $_setField(1008, value);
  @$pb.TagNumber(1008)
  $core.bool hasMute() => $_has(9);
  @$pb.TagNumber(1008)
  void clearMute() => $_clearField(1008);
  @$pb.TagNumber(1008)
  Mute ensureMute() => $_ensure(9);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
