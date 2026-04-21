import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/authz.dart';

void main() {
  group('Bearer', () {
    test('holds metadata and bearer string', () {
      final b = Bearer('meta', 'tok');
      expect(b.metadata, equals('meta'));
      expect(b.bearer, equals('tok'));
    });
  });

  group('Cached', () {
    test('pending throws when token is called', () {
      final c = Cached(Bearer('m', 'tok'), Cached.pending);
      expect(() => c.token(), throwsArgumentError);
    });

    test('noprefresh returns current value immediately', () async {
      final b = Bearer('m', 'tok');
      final c = Cached(b, Cached.noprefresh);
      final result = await c.token();
      expect(result.bearer, equals('tok'));
      expect(result.metadata, equals('m'));
    });

    test('token() updates current via refresh function', () async {
      final updated = Bearer('m2', 'new-tok');
      final c = Cached(Bearer('m', 'old-tok'), (_) => Future.value(updated));
      final result = await c.token();
      expect(result.bearer, equals('new-tok'));
      expect(c.current.bearer, equals('new-tok'));
    });

    test('token() reflects last refresh result on repeated calls', () async {
      var callCount = 0;
      final c = Cached<String>(
        Bearer('m', 'tok-0'),
        (_) {
          callCount++;
          return Future.value(Bearer('m', 'tok-$callCount'));
        },
      );

      final first = await c.token();
      final second = await c.token();

      expect(first.bearer, equals('tok-1'));
      expect(second.bearer, equals('tok-2'));
      expect(callCount, equals(2));
    });
  });

  group('refresh', () {
    bool _expired(String meta, DateTime ts) => meta == 'expired';
    bool _notExpired(String meta, DateTime ts) => false;

    test('returns current value when not expired', () async {
      final current = Bearer('live', 'cached-tok');
      final fn = refresh<String>(
        (_) => Future.value(Bearer('new', 'new-tok')),
        _notExpired,
      );
      final c = Cached(current, fn);
      final result = await c.token();
      expect(result.bearer, equals('cached-tok'));
    });

    test('calls fn and returns new value when expired', () async {
      final fn = refresh<String>(
        (_) => Future.value(Bearer('new', 'refreshed-tok')),
        _expired,
      );
      final c = Cached(Bearer('expired', 'old-tok'), fn);
      final result = await c.token();
      expect(result.bearer, equals('refreshed-tok'));
    });

    test('passes current metadata to fn on expiry', () async {
      String? received;
      final fn = refresh<String>(
        (meta) {
          received = meta;
          return Future.value(Bearer(meta, 'tok'));
        },
        _expired,
      );
      final c = Cached(Bearer('expired', 'tok'), fn);
      await c.token();
      expect(received, equals('expired'));
    });

    test('checks expiry against current DateTime', () async {
      DateTime? capturedTs;
      final fn = refresh<String>(
        (_) => Future.value(Bearer('m', 'tok')),
        (meta, ts) {
          capturedTs = ts;
          return false;
        },
      );
      final before = DateTime.now();
      final c = Cached(Bearer('m', 'tok'), fn);
      await c.token();
      final after = DateTime.now();

      expect(capturedTs, isNotNull);
      expect(capturedTs!.isAfter(before) || capturedTs!.isAtSameMomentAs(before), isTrue);
      expect(capturedTs!.isBefore(after) || capturedTs!.isAtSameMomentAs(after), isTrue);
    });

    test('expired token with future expiry timestamp is not refreshed', () async {
      var refreshed = false;
      final futureTs = DateTime.now().add(const Duration(hours: 1));
      final fn = refresh<DateTime>(
        (_) {
          refreshed = true;
          return Future.value(Bearer(futureTs, 'new-tok'));
        },
        (meta, ts) => meta.isBefore(ts),
      );

      // meta is a future timestamp — not before now, so not expired
      final c = Cached(Bearer(futureTs, 'cached-tok'), fn);
      final result = await c.token();

      expect(refreshed, isFalse);
      expect(result.bearer, equals('cached-tok'));
    });

    test('expired token with past expiry timestamp is refreshed', () async {
      var refreshed = false;
      final pastTs = DateTime.now().subtract(const Duration(hours: 1));
      final fn = refresh<DateTime>(
        (_) {
          refreshed = true;
          return Future.value(Bearer(DateTime.now().add(const Duration(hours: 1)), 'new-tok'));
        },
        (meta, ts) => meta.isBefore(ts),
      );

      final c = Cached(Bearer(pastTs, 'old-tok'), fn);
      final result = await c.token();

      expect(refreshed, isTrue);
      expect(result.bearer, equals('new-tok'));
    });
  });
}
