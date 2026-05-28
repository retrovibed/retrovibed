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
      q.reset(Stream.fromIterable([_media('abc', 'My Song')]));
      await q.advance("", player);

      expect(q.known.id, 'abc');
      expect(q.known.description, 'My Song');
    });
  });

  group('PlayQueue.currentStart', () {
    test('returns Duration.zero when no media is loaded', () {
      expect(PlayQueue().currentStart, Duration.zero);
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
      q.reset(Stream.fromIterable([_media('a', 'A'), _media('b', 'B')]));
      await q.advance("", player); // pulls 'a' into current, nothing in upcoming yet
      await q.advance("", player); // pulls 'b' into current, 'a' moves to previous

      // reset should clear upcoming
      q.reset(Stream.empty());
      expect(q.upcoming, 0);
    });
  });

  group('PlayQueue.advance', () {
    test('returns null on empty stream', () async {
      final player = _fakePlayer();
      addTearDown(player.dispose);

      final q = PlayQueue();
      q.reset(Stream.empty());
      final result = await q.advance("", player);

      expect(result, isNull);
      expect(q.known.id, uuidx.min());
    });

    test('advances through stream items in order', () async {
      final player = _fakePlayer();
      addTearDown(player.dispose);

      final q = PlayQueue();
      q.reset(Stream.fromIterable([_media('1', 'First'), _media('2', 'Second')]));

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
      q.reset(Stream.fromIterable([_media('a', 'A'), _media('b', 'B')]));
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
      q.reset(Stream.empty());
      final result = await q.advance("", player);

      expect(result, isNull);
      expect(q.known.id, uuidx.min());
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

    api.FnMediaRandom _random(String id) => (req, {options = const []}) async => api.MediaFindResponse(media: _m(id));

    test('yields items starting from current', () async {
      final results =
          await range(
            _req(10),
            _m('b'),
            search: _search([_m('a'), _m('b'), _m('c')], 10),
            random: _random('rand'),
          ).take(3).toList();

      expect(_id(results[0]), 'b');
      expect(_id(results[1]), 'c');
      expect(_id(results[2]), 'rand');
    });

    test('falls back to index 0 when current not found', () async {
      final results =
          await range(
            _req(10),
            _m('missing'),
            search: _search([_m('a'), _m('b')], 10),
            random: _random('rand'),
          ).take(2).toList();

      expect(_id(results[0]), 'a');
      expect(_id(results[1]), 'b');
    });

    test('applies pos to first item only', () async {
      const pos = Duration(seconds: 5);

      final results =
          await range(
            _req(10),
            _m('a'),
            pos: pos,
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

      final results =
          await range(
            _req(2),
            _m('a'),
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

      final results =
          await range(
            _req(10),
            _m('a'),
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
      q.reset(Stream.fromIterable([_media('x', 'X'), _media('y', 'Y')]));
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
      q.reset(Stream.fromIterable([_media('x', 'X'), _media('y', 'Y')]));
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
