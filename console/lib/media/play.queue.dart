import 'dart:async';
import 'dart:collection';

import 'package:flutter/foundation.dart';
import 'package:media_kit/media_kit.dart' as mediakit;
import 'package:retrovibed/media.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/httpx.dart' as httpx;
import 'api.dart' as api;

class RingBuffer<T> {
  final ListQueue<T> _queue;
  final int capacity;

  RingBuffer(this.capacity) : assert(capacity > 0), _queue = ListQueue<T>(capacity);

  int get length => _queue.length;
  bool get isEmpty => _queue.isEmpty;
  bool get isFull => _queue.length == capacity;

  List<T> toList() {
    return _queue.toList();
  }

  void insert(T? element) {
    if (element == null) return;
    if (isFull) {
      _queue.removeFirst();
    }
    _queue.addLast(element);
  }

  T? remove() {
    if (isEmpty) {
      return null;
    }
    return _queue.removeLast();
  }

  T? peek() {
    return _queue.lastOrNull;
  }

  void clear() {
    _queue.clear();
  }

  void removeWhere(bool Function(T) test) {
    _queue.removeWhere(test);
  }
}

class PlayQueue {
  // current/pos/recent are read lazily from inside range()'s async generator
  // body - never synchronously at the point search()/random()/acoustic() are
  // called, since reset() (which seeds current) doesn't run until after
  // autoqueue(...) has already returned its Stream object.
  final ValueNotifier<PlayableMedia?> current = ValueNotifier(null);
  final RingBuffer<PlayableMedia> _upcoming = RingBuffer(128);
  final RingBuffer<PlayableMedia> _previous = RingBuffer(128);
  StreamIterator<PlayableMedia> _stream = StreamIterator(Stream.empty());
  final ValueNotifier<int> revision = ValueNotifier(0);

  int get upcoming => _upcoming.length;
  int get previous => _previous.length;
  int get capacity => _upcoming.capacity;

  List<PlayableMedia> get queued => _upcoming.toList();

  Known get known => current.value.known;
  Duration get pos => current.value?.pos ?? Duration.zero;

  // oldest -> newest, ending with the currently playing track (if any) - so
  // recent.last is always "what's playing right now", with no separate
  // current/seed parameter needed wherever this feeds an exclusion list.
  List<PlayableMedia> get recent => [
    ..._previous.toList(),
    if (current.value != null) current.value!,
  ];

  void reset(Stream<PlayableMedia> stream, Media media, {Duration pos = const Duration(milliseconds: 0)}) {
    _stream = StreamIterator(stream);
    _upcoming.clear();
    _previous.clear();
    current.value = PlayableMedia(media, pos: pos);
    revision.value++;
  }

  Future<PlayableMedia?> advance(String auth, mediakit.Player player) async {
    PlayableMedia? m = _upcoming.remove();
    if (m == null) {
      final hasNext = await _stream.moveNext();
      m = hasNext ? _stream.current : null;
    }

    if (m == null) {
      print("end of media reached");
      return null;
    }

    try {
      await player.open(m.playable(auth));
      // skips the push when the outgoing and incoming track are the same -
      // notably true for the very first advance after reset(), since the
      // seed media and the stream's first (re-fetched) item share an id.
      if (current.value?.current.id != m.current.id) {
        _previous.insert(current.value);
      }
      current.value = m;
      return m;
    } catch (cause) {
      print("unable to pull next media from playlist $cause");
      return null;
    }
  }

  void push(PlayableMedia media) {
    _upcoming.insert(media);
    revision.value++;
  }

  void remove(String id) {
    _upcoming.removeWhere((m) => m.current.id == id);
    revision.value++;
  }

  Future<PlayableMedia?> reverse(String auth, mediakit.Player player) async {
    final prev = _previous.remove();
    if (prev == null) return null;
    await player.open(prev.playable(auth));
    _upcoming.insert(current.value);
    current.value = prev;
    return current.value;
  }
}

class PlayableMedia {
  final Media current;
  final Duration pos;

  const PlayableMedia(this.current, {this.pos = const Duration(milliseconds: 0)});

  Known get known => Known(
    id: current.id,
    description: current.description,
    image: current.image,
  );

  mediakit.Media playable(String auth) {
    print("DERP DERP playable... ${current} ${auth}");
    return mediakit.Media(
      api.media.download_uri(current.id),
      extras: Map.of(<String, String>{
        "id": current.id,
        "title": current.description,
      }),
      start: pos,
      httpHeaders: <String, String>{"Authorization": auth},
    );
  }
}

extension PlayableMediaNullable on PlayableMedia? {
  Known get known =>
      this?.known ??
      Known(
        id: uuidx.min(),
      );
}

Stream<PlayableMedia> range(
  MediaSearchRequest req,
  PlayQueue queue, {
  List<httpx.Option> Function() options = httpx.Request.empty,
  api.FnMediaSearch search = api.media.search,
  api.FnMediaFind random = api.media.random,
}) async* {
  final playable = queue.current.value;

  MediaSearchResponse i = await search(req, options: options());

  // the anchor (queue.current) is always yielded directly below, using the
  // Media + position already known to the queue - it doesn't need to be
  // present in the search results. when it is present, skip it here so it
  // isn't yielded twice.
  final anchor = playable?.current.id ?? uuidx.min();
  final idx = i.items.indexWhere((m) => m.id == anchor);
  final initial = idx == -1 ? i.items : i.items.sublist(idx + 1);

  if (playable != null) {
    yield PlayableMedia(
      playable.current,
      pos: queue.pos,
    );
  }

  for (var m in initial) {
    yield PlayableMedia(m);
  }

  while (i.items.length == i.next.limit.toInt()) {
    i = await search(i.next, options: options());
    for (var m in i.items) {
      yield PlayableMedia(m);
    }
    i.next..offset += 1;
  }

  // at this point we've run out of content from the provided search.
  // lets play random content. using things like the mimetypes from
  // from the initial request. we'll eventually add in more coherent
  // results to keep a trend going.
  while (true) {
    i.next.excluded
      ..clear()
      ..addAll(queue.recent.map((m) => m.current.id));
    final v = await random(i.next, options: options());
    yield PlayableMedia(v.media);
  }
}

typedef RangeFn =
    Stream<PlayableMedia> Function(
      MediaSearchRequest req,
      PlayQueue queue, {
      List<httpx.Option> Function() options,
    });

Stream<PlayableMedia> search(
  MediaSearchRequest req,
  PlayQueue queue, {
  List<httpx.Option> Function() options = httpx.Request.empty,
}) {
  return range(req, queue, options: options, random: api.media.random);
}

Stream<PlayableMedia> random(
  MediaSearchRequest req,
  PlayQueue queue, {
  List<httpx.Option> Function() options = httpx.Request.empty,
}) {
  return range(
    req,
    queue,
    options: options,
    search: api.media.emptysearch,
    random: api.media.random,
  );
}

Stream<PlayableMedia> acoustic(
  MediaSearchRequest req,
  PlayQueue queue, {
  List<httpx.Option> Function() options = httpx.Request.empty,
}) {
  return range(
    req,
    queue,
    options: options,
    search: api.media.emptysearch,
    random: api.media.acoustic,
  );
}
