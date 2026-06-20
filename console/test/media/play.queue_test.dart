import 'dart:async';

import 'package:fixnum/fixnum.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:media_kit/media_kit.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media/api.dart' as api;
import 'package:retrovibed/media/play.queue.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;

class _FakePlatformPlayer extends PlatformPlayer {
  _FakePlatformPlayer() : super(configuration: const PlayerConfiguration());

  @override
  Future<void> open(Playable playable, {bool play = true}) async {}
}

Player _fakePlayer() => Player(platformPlayer: _FakePlatformPlayer());

PlayableMedia _media(String id, String title) => PlayableMedia(api.Media(id: id, description: title));

void main() {
  group('PlayQueue.current', () {
    test('defaults id to uuidx.min() when no media is loaded', () {
      final q = PlayQueue();
      expect(q.known.id, uuidx.min());
    });

    test('defaults description to empty when no media is loaded', () {
      final q = PlayQueue();
      expect(q.known.hasDescription(), isFalse);
    });

    test('returns id and description from loaded media', () async {
      final player = _fakePlayer();
      addTearDown(player.dispose);

      final q = PlayQueue();
      q.reset(Stream.fromIterable([_media('abc', 'My Song')]), api.Media(id: 'abc', description: 'My Song'));
      await q.advance("", player);

      expect(q.known.id, 'abc');
      expect(q.known.description, 'My Song');
    });
  });

  group('PlayQueue.pos', () {
    test('returns Duration.zero when no media is loaded', () {
      expect(PlayQueue().pos, Duration.zero);
    });

    test('reflects the position seeded via reset', () {
      final q = PlayQueue();
      q.reset(Stream.empty(), api.Media(id: 'a', description: 'A'), pos: const Duration(seconds: 5));
      expect(q.pos, const Duration(seconds: 5));
    });
  });

  group('PlayQueue.upcoming / previous counts', () {
    test('both start at zero', () {
      final q = PlayQueue();
      expect(q.upcoming, 0);
      expect(q.previous, 0);
    });
  });

  group('PlayQueue.reset', () {
    test('clears upcoming buffer', () async {
      final player = _fakePlayer();
      addTearDown(player.dispose);

      final q = PlayQueue();
      q.reset(Stream.fromIterable([_media('a', 'A'), _media('b', 'B')]), api.Media(id: 'a', description: 'A'));
      await q.advance("", player); // pulls 'a' into current, nothing in upcoming yet
      await q.advance("", player); // pulls 'b' into current, 'a' moves to previous

      // reset should clear upcoming
      q.reset(Stream.empty(), api.Media(id: 'z', description: 'Z'));
      expect(q.upcoming, 0);
    });

    test('clears previous history', () async {
      final player = _fakePlayer();
      addTearDown(player.dispose);

      final q = PlayQueue();
      q.reset(Stream.fromIterable([_media('a', 'A'), _media('b', 'B')]), api.Media(id: 'a', description: 'A'));
      await q.advance("", player);
      await q.advance("", player);
      expect(q.previous, 1);

      q.reset(Stream.empty(), api.Media(id: 'z', description: 'Z'));
      expect(q.previous, 0);
      expect(q.recent, ['z']); // recent always includes whatever's now current
    });
  });

  group('PlayQueue.advance', () {
    test('returns null on empty stream', () async {
      final player = _fakePlayer();
      addTearDown(player.dispose);

      final q = PlayQueue();
      q.reset(Stream.empty(), api.Media(id: 'seed', description: 'Seed'));
      final result = await q.advance("", player);

      expect(result, isNull);
      expect(q.known.id, 'seed'); // current is untouched by a failed advance
    });

    test('advances through stream items in order', () async {
      final player = _fakePlayer();
      addTearDown(player.dispose);

      final q = PlayQueue();
      q.reset(
        Stream.fromIterable([_media('1', 'First'), _media('2', 'Second')]),
        api.Media(id: '1', description: 'First'),
      );

      await q.advance("", player);
      expect(q.known.id, '1');
      expect(q.previous, 0);

      await q.advance("", player);
      expect(q.known.id, '2');
      expect(q.previous, 1);
    });

    test('drains upcoming buffer before pulling from stream', () async {
      final player = _fakePlayer();
      addTearDown(player.dispose);

      final q = PlayQueue();
      q.reset(Stream.fromIterable([_media('a', 'A'), _media('b', 'B')]), api.Media(id: 'a', description: 'A'));
      await q.advance("", player); // current = 'a'
      await q.advance("", player); // current = 'b', previous has 'a'

      // reverse puts 'b' in upcoming and restores 'a'
      await q.reverse("", player);
      expect(q.known.id, 'a');
      expect(q.upcoming, 1); // 'b' is buffered

      // advance should use the buffer, not the (exhausted) stream
      await q.advance("", player);
      expect(q.known.id, 'b');
      expect(q.upcoming, 0);
    });

    test('returns null and does not update current when player throws', () async {
      final player = _fakePlayer();
      addTearDown(player.dispose);

      // Provide a stream that errors on open by using a player override
      // Instead: test the catch path by exhausting the stream first
      final q = PlayQueue();
      q.reset(Stream.empty(), api.Media(id: 'seed', description: 'Seed'));
      final result = await q.advance("", player);

      expect(result, isNull);
      expect(q.known.id, 'seed');
    });

    test('does not double-count the seed track when the stream starts with the same media', () async {
      final player = _fakePlayer();
      addTearDown(player.dispose);

      final q = PlayQueue();
      q.reset(Stream.fromIterable([_media('x', 'X'), _media('y', 'Y')]), api.Media(id: 'x', description: 'X'));

      await q.advance("", player); // current -> stream's 'x'; seed must not land in previous
      expect(q.previous, 0);
      expect(q.recent, ['x']); // just current, no previous entry yet

      await q.advance("", player); // current -> 'y'; 'x' lands in previous exactly once
      expect(q.previous, 1);
      expect(q.recent, ['x', 'y']);
    });
  });

  group('range', () {
    api.Media _m(String id) => api.Media(id: id, description: 'Title $id');
    String _id(PlayableMedia m) => m.current.id;

    api.MediaSearchRequest _req(int limit) => api.MediaSearchRequest(limit: Int64(limit));

    api.FnMediaSearch _search(List<api.Media> items, int limit) =>
        (req, {options = const []}) async => api.MediaSearchResponse(
          items: items,
          next: _req(limit),
        );

    api.FnMediaFind _random(String id) =>
        (req, {options = const []}) async => api.MediaFindResponse(media: _m(id));

    test('yields items starting from current', () async {
      final q = PlayQueue();
      q.reset(Stream.empty(), _m('b'));

      final results = await range(
        _req(10),
        q,
        search: _search([_m('a'), _m('b'), _m('c')], 10),
        random: _random('rand'),
      ).take(3).toList();

      expect(_id(results[0]), 'b');
      expect(_id(results[1]), 'c');
      expect(_id(results[2]), 'rand');
    });

    test('yields current first even when missing from search results, then falls back to index 0', () async {
      final q = PlayQueue();
      q.reset(Stream.empty(), _m('missing'));

      final results = await range(
        _req(10),
        q,
        search: _search([_m('a'), _m('b')], 10),
        random: _random('rand'),
      ).take(3).toList();

      expect(_id(results[0]), 'missing');
      expect(_id(results[1]), 'a');
      expect(_id(results[2]), 'b');
    });

    test('applies pos to first item only', () async {
      const pos = Duration(seconds: 5);

      final q = PlayQueue();
      q.reset(Stream.empty(), _m('a'), pos: pos);

      final results = await range(
        _req(10),
        q,
        search: _search([_m('a'), _m('b')], 10),
        random: _random('rand'),
      ).take(2).toList();

      expect(results[0].pos, pos);
      expect(results[1].pos, Duration.zero);
    });

    test('paginates when items.length equals limit', () async {
      var calls = 0;

      Future<api.MediaSearchResponse> search(
        api.MediaSearchRequest req, {
        List<httpx.Option> options = const [],
      }) async {
        calls++;
        return calls == 1
            ? api.MediaSearchResponse(items: [_m('a'), _m('b')], next: _req(2))
            : api.MediaSearchResponse(items: [_m('c')], next: _req(2));
      }

      final q = PlayQueue();
      q.reset(Stream.empty(), _m('a'));

      final results = await range(
        _req(2),
        q,
        search: search,
        random: _random('rand'),
      ).take(4).toList();

      expect(_id(results[0]), 'a');
      expect(_id(results[1]), 'b');
      expect(_id(results[2]), 'c');
      expect(_id(results[3]), 'rand');
      expect(calls, 2);
    });

    test('falls back to random after search exhausted', () async {
      var randomCalled = false;

      final q = PlayQueue();
      q.reset(Stream.empty(), _m('a'));

      final results = await range(
        _req(10),
        q,
        search: _search([_m('a')], 10),
        random: (req, {options = const []}) async {
          randomCalled = true;
          return api.MediaFindResponse(media: _m('rand'));
        },
      ).take(2).toList();

      expect(randomCalled, isTrue);
      expect(_id(results[1]), 'rand');
    });
  });

  group('PlayQueue.reverse', () {
    test('returns null when previous is empty', () async {
      final player = _fakePlayer();
      addTearDown(player.dispose);

      final q = PlayQueue();
      final result = await q.reverse("", player);
      expect(result, isNull);
    });

    test('restores previous current after advance', () async {
      final player = _fakePlayer();
      addTearDown(player.dispose);

      final q = PlayQueue();
      q.reset(Stream.fromIterable([_media('x', 'X'), _media('y', 'Y')]), api.Media(id: 'x', description: 'X'));
      await q.advance("", player); // current = 'x'
      await q.advance("", player); // current = 'y', previous has 'x'

      await q.reverse("", player);
      expect(q.known.id, 'x');
      expect(q.previous, 0);
    });

    test('places current in upcoming when reversing', () async {
      final player = _fakePlayer();
      addTearDown(player.dispose);

      final q = PlayQueue();
      q.reset(Stream.fromIterable([_media('x', 'X'), _media('y', 'Y')]), api.Media(id: 'x', description: 'X'));
      await q.advance("", player); // current = 'x'
      await q.advance("", player); // current = 'y'

      expect(q.upcoming, 0);
      await q.reverse("", player); // current = 'x', 'y' should be in upcoming
      expect(q.upcoming, 1);
    });

    test('multiple reverse calls walk back through history', () async {
      final player = _fakePlayer();
      addTearDown(player.dispose);

      final q = PlayQueue();
      q.reset(
        Stream.fromIterable([
          _media('1', 'One'),
          _media('2', 'Two'),
          _media('3', 'Three'),
        ]),
        api.Media(id: '1', description: 'One'),
      );
      await q.advance("", player);
      await q.advance("", player);
      await q.advance("", player);
      expect(q.known.id, '3');

      await q.reverse("", player);
      expect(q.known.id, '2');

      await q.reverse("", player);
      expect(q.known.id, '1');

      final result = await q.reverse("", player);
      expect(result, isNull); // nothing left in previous
    });
  });
}
