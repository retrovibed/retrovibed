import 'dart:async';
import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/authn/cache.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/uuidx.dart' as uuidx;

meta.AuthzResponse _makeResponse() => meta.AuthzResponse(
  bearer: uuidx.min(),
  token: meta.Token()..expires = fixnum.Int64(DateTime.now().millisecondsSinceEpoch + 3600000),
);

meta.AuthzResponse _response(String bearer, fixnum.Int64 expires) =>
    meta.AuthzResponse(bearer: bearer, token: meta.Token()..expires = expires);

fixnum.Int64 _pastExpiry() =>
    fixnum.Int64(DateTime.now().subtract(const Duration(hours: 1)).millisecondsSinceEpoch ~/ 1000);

fixnum.Int64 _futureExpiry() =>
    fixnum.Int64(DateTime.now().add(const Duration(hours: 1)).millisecondsSinceEpoch ~/ 1000);

void main() {
  group('AuthzCache', () {
    testWidgets('does not render child until token resolves', (WidgetTester tester) async {
      final completer = Completer<meta.AuthzResponse>();

      await tester.pumpWidget(
        MaterialApp(
          home: Material(
            child: ds.LoadingGuard(
              AuthzCache(
                const Text('protected content'),
                current: ({String? host}) => completer.future,
              ),
            ),
          ),
        ),
      );
      await tester.pump();
      await tester.pump();

      expect(find.text('protected content'), findsNothing);

      completer.complete(_makeResponse());
      await tester.pumpAndSettle();

      expect(find.text('protected content'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    group('failed refresh', () {
      testWidgets('propagates the underlying error instead of caching an empty bearer', (WidgetTester tester) async {
        BuildContext? capturedCtx;
        final failure = Exception('authz unavailable');

        await tester.pumpWidget(
          MaterialApp(
            home: Material(
              child: ds.LoadingGuard(
                AuthzCache(
                  Builder(
                    builder: (ctx) {
                      capturedCtx = ctx;
                      return const Text('protected content');
                    },
                  ),
                  current: ({String? host}) => Future.error(failure),
                ),
              ),
            ),
          ),
        );
        await tester.pumpAndSettle();

        final cache = AuthzCache.of(capturedCtx!);
        await expectLater(cache.meta.auto(), throwsA(failure));
        expect(cache.meta.current.bearer, isEmpty);
        expect(tester.takeException(), isNull);
      });

      testWidgets('does not overwrite a previously cached bearer on a later failure', (WidgetTester tester) async {
        var callCount = 0;
        BuildContext? capturedCtx;

        await tester.pumpWidget(
          MaterialApp(
            home: Material(
              child: ds.LoadingGuard(
                AuthzCache(
                  Builder(
                    builder: (ctx) {
                      capturedCtx = ctx;
                      return const Text('protected content');
                    },
                  ),
                  current: ({String? host}) {
                    callCount++;
                    if (callCount == 1) {
                      return Future.value(_response('good-bearer', _pastExpiry()));
                    }
                    return Future.error(Exception('authz unavailable'));
                  },
                ),
              ),
            ),
          ),
        );
        await tester.pumpAndSettle();

        final cache = AuthzCache.of(capturedCtx!);
        expect(cache.meta.current.bearer, equals('good-bearer'));

        await expectLater(cache.meta.auto(), throwsException);
        await tester.pumpAndSettle();

        expect(cache.meta.current.bearer, equals('good-bearer'));
        expect(tester.takeException(), isNull);
      });
    });

    group('overlapping refresh', () {
      testWidgets('a refresh triggered mid-fetch is queued behind it instead of racing it', (
        WidgetTester tester,
      ) async {
        BuildContext? capturedCtx;
        final completers = <Completer<meta.AuthzResponse>>[];

        await tester.pumpWidget(
          MaterialApp(
            home: Material(
              child: ds.LoadingGuard(
                AuthzCache(
                  Builder(
                    builder: (ctx) {
                      capturedCtx = ctx;
                      return const Text('protected content');
                    },
                  ),
                  current: ({String? host}) {
                    final c = Completer<meta.AuthzResponse>();
                    completers.add(c);
                    return c.future;
                  },
                ),
              ),
            ),
          ),
        );
        await tester.pump();
        await tester.pump();

        final cache = AuthzCache.of(capturedCtx!);

        // trigger a second refresh before the first (session 1) fetch has
        // resolved.
        cache.refresh();
        await tester.pump(); // let the new Cached reach AuthzTokenData before forcing a fetch on it
        AuthzCache.meta(capturedCtx!).auto();

        // the second refresh must not start its own fetch concurrently with
        // the first — it queues behind the shared lock instead.
        expect(completers.length, equals(1));

        // session 1's fetch resolves, but with an already-expired token, so
        // the queued session-2 refresh performs its own fetch next rather
        // than treating this stale result as good enough.
        completers[0].complete(_response('stale-bearer', _pastExpiry()));
        await tester.pumpAndSettle();

        expect(completers.length, equals(2));
        expect(cache.meta.current.bearer, equals('stale-bearer'));

        // session 2's fetch resolves last, deterministically winning because
        // the two fetches are serialized rather than raced.
        completers[1].complete(_response('new-correct-bearer', _futureExpiry()));
        await tester.pumpAndSettle();

        expect(cache.meta.current.bearer, equals('new-correct-bearer'));
        expect(tester.takeException(), isNull);
      });
    });
  });
}
