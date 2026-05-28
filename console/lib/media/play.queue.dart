import 'dart:async';
import 'dart:collection';
import 'dart:math';

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
}

class PlayQueue {
  final ValueNotifier<PlayableMedia?> current = ValueNotifier(null);
  final RingBuffer<PlayableMedia> _upcoming = RingBuffer(128);
  final RingBuffer<PlayableMedia> _previous = RingBuffer(128);
  StreamIterator<PlayableMedia> _stream = StreamIterator(Stream.empty());

  int get upcoming => _upcoming.length;
  int get previous => _previous.length;

  Known get known => current.value.known;
  Duration get currentStart => current.value?.pos ?? Duration.zero;

  void reset(Stream<PlayableMedia> stream) {
    _stream = StreamIterator(stream);
    _upcoming.clear();
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
      _previous.insert(current.value);
      current.value = m;
      return m;
    } catch (cause) {
      print("unable to pull next media from playlist $cause");
      return null;
    }
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
  );

  mediakit.Media playable(String auth) => mediakit.Media(
    api.media.download_uri(current.id),
    extras: Map.of(<String, String>{
      "id": current.id,
      "title": current.description,
    }),
    start: pos,
    httpHeaders: <String, String>{"Authorization": httpx.auto_bearer_host()},
  );
}

extension PlayableMediaNullable on PlayableMedia? {
  Known get known =>
      this?.known ??
      Known(
        id: uuidx.min(),
        description: "",
      );
}

Stream<PlayableMedia> range(
  MediaSearchRequest req,
  Media current, {
  Duration pos = const Duration(milliseconds: 0),
  List<httpx.Option> Function() options = httpx.Request.empty,
  api.FnMediaSearch search = api.media.search,
  api.FnMediaRandom random = api.media.random,
}) async* {
  MediaSearchResponse i = await search(req, options: options());
  final initial = i.items.sublist(
    max(i.items.indexWhere((m) => m.id == current.id), 0),
  );

  for (var (idx, m) in initial.indexed) {
    print("DERP DERP 0");
    yield PlayableMedia(
      m,
      pos: idx == 0 ? pos : const Duration(milliseconds: 0),
    );
  }

  while (i.items.length == i.next.limit.toInt()) {
    i = await search(i.next, options: options());
    for (var m in i.items) {
      print("DERP DERP 1");
      yield PlayableMedia(m);
    }
    i.next..offset += 1;
  }

  // at this point we've run out of content from the provided search.
  // lets play random content. using things like the mimetypes from
  // from the initial request. we'll eventually add in more coherent
  // results to keep a trend going.
  while (true) {
    final v = await random(i.next, options: options());
    print("DERP DERP 2");
    yield PlayableMedia(v.media);
  }
}
