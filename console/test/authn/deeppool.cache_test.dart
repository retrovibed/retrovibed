import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/authn/cache.dart';
import 'package:retrovibed/authn/deeppool.cache.dart';
import 'package:retrovibed/authn/local.only.guard.dart';
import 'package:retrovibed/billing/meta.billing.pb.dart' as billing;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/meta/api.deeppool.dart' as deeppool;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

fixnum.Int64 _futureExpiry() =>
    fixnum.Int64(DateTime.now().add(const Duration(hours: 1)).millisecondsSinceEpoch ~/ 1000);

fixnum.Int64 _pastExpiry() =>
    fixnum.Int64(DateTime.now().subtract(const Duration(hours: 1)).millisecondsSinceEpoch ~/ 1000);

deeppool.AuthzResponse _response(String bearer, fixnum.Int64 expires) {
  return deeppool.AuthzResponse(
    bearer: bearer,
    token: deeppool.Token(expires: expires),
  );
}

void main() {
  group('DeeppoolAuthzCache', () {
    testWidgets('never returns an empty token', (WidgetTester tester) async {
      BuildContext? capturedCtx;

      await tester.pumpApp(
        DeeppoolAuthzCache(
          Builder(
            builder: (ctx) {
              capturedCtx = ctx;
              return const Text('child');
            },
          ),
          apideeppoolauthz: ({options = const []}) => Future.value(_response('bearer-token', _futureExpiry())),
        ),
      );
      await tester.pumpAndSettle();

      final cache = DeeppoolAuthzCache.of(capturedCtx!);
      final token = await cache.meta.auto();
      expect(token.bearer, isNotEmpty);
      expect(tester.takeException(), isNull);
    });

    group('of()', () {
      testWidgets('returns non-empty token on first call', (WidgetTester tester) async {
        BuildContext? capturedCtx;

        await tester.pumpApp(
          DeeppoolAuthzCache(
            Builder(
              builder: (ctx) {
                capturedCtx = ctx;
                return const Text('child');
              },
            ),
            apideeppoolauthz: ({options = const []}) =>
                Future.value(_response('bearer-from-deeppool', _futureExpiry())),
          ),
        );
        await tester.pumpAndSettle();

        final cache = DeeppoolAuthzCache.of(capturedCtx!);
        expect(cache, isNotNull);

        final token = await cache.meta.auto();
        expect(token.bearer, isNotEmpty);
        expect(token.bearer, equals('bearer-from-deeppool'));
        expect(tester.takeException(), isNull);
      });

      testWidgets('meta.current reflects api response after first token call', (WidgetTester tester) async {
        BuildContext? capturedCtx;

        await tester.pumpApp(
          DeeppoolAuthzCache(
            Builder(
              builder: (ctx) {
                capturedCtx = ctx;
                return const Text('child');
              },
            ),
            apideeppoolauthz: ({options = const []}) => Future.value(_response('cached-bearer', _futureExpiry())),
          ),
        );
        await tester.pumpAndSettle();

        final cache = DeeppoolAuthzCache.of(capturedCtx!);
        final token = await cache.meta.auto();
        await tester.pumpAndSettle();

        expect(token.bearer, equals('cached-bearer'));
        expect(cache.meta.current.bearer, equals('cached-bearer'));
        expect(tester.takeException(), isNull);
      });
    });

    group('bearer()', () {
      testWidgets('returns non-empty bearer string on first call', (WidgetTester tester) async {
        BuildContext? capturedCtx;

        await tester.pumpApp(
          DeeppoolAuthzCache(
            Builder(
              builder: (ctx) {
                capturedCtx = ctx;
                return const Text('child');
              },
            ),
            apideeppoolauthz: ({options = const []}) => Future.value(_response('bearer-option-token', _futureExpiry())),
          ),
        );
        await tester.pumpAndSettle();

        final option = DeeppoolAuthzCache.bearer(capturedCtx!);
        final request = httpx.Request();
        await option(request);

        expect(request.headers['Authorization'], isNotNull);
        expect(request.headers['Authorization'], contains('bearer-option-token'));
        expect(tester.takeException(), isNull);
      });

      testWidgets('bearer option header is non-empty', (WidgetTester tester) async {
        BuildContext? capturedCtx;

        await tester.pumpApp(
          DeeppoolAuthzCache(
            Builder(
              builder: (ctx) {
                capturedCtx = ctx;
                return const Text('child');
              },
            ),
            apideeppoolauthz: ({options = const []}) => Future.value(_response('my-token', _futureExpiry())),
          ),
        );
        await tester.pumpAndSettle();

        final option = DeeppoolAuthzCache.bearer(capturedCtx!);
        final request = httpx.Request();
        await option(request);

        expect(request.headers['Authorization'], isNotEmpty);
        expect(tester.takeException(), isNull);
      });
    });

    group('attributionToken()', () {
      testWidgets('returns token from api', (WidgetTester tester) async {
        BuildContext? capturedCtx;

        await tester.pumpApp(
          DeeppoolAuthzCache(
            Builder(
              builder: (ctx) {
                capturedCtx = ctx;
                return const Text('child');
              },
            ),
            apideeppoolauthz: ({options = const []}) => Future.value(_response('bearer', _futureExpiry())),
            apibillingattribution: ({options = const []}) =>
                Future.value(billing.AttributionTokenResponse(token: 'attr-jwt-token')),
          ),
        );
        await tester.pumpAndSettle();

        final token = await DeeppoolAuthzCache.attributionToken(capturedCtx!);
        expect(token, equals('attr-jwt-token'));
        expect(tester.takeException(), isNull);
      });

      testWidgets('caches token across calls', (WidgetTester tester) async {
        var callCount = 0;
        BuildContext? capturedCtx;

        await tester.pumpApp(
          DeeppoolAuthzCache(
            Builder(
              builder: (ctx) {
                capturedCtx = ctx;
                return const Text('child');
              },
            ),
            apideeppoolauthz: ({options = const []}) => Future.value(_response('bearer', _futureExpiry())),
            apibillingattribution: ({options = const []}) {
              callCount++;
              return Future.value(billing.AttributionTokenResponse(token: 'attr-v$callCount'));
            },
          ),
        );
        await tester.pumpAndSettle();

        final first = await DeeppoolAuthzCache.attributionToken(capturedCtx!);
        final second = await DeeppoolAuthzCache.attributionToken(capturedCtx!);

        expect(first, equals('attr-v1'));
        expect(second, equals('attr-v1'));
        expect(tester.takeException(), isNull);
      });
    });

    group('expiration', () {
      testWidgets('re-fetches token when current token is expired', (WidgetTester tester) async {
        var callCount = 0;
        BuildContext? capturedCtx;

        await tester.pumpApp(
          DeeppoolAuthzCache(
            Builder(
              builder: (ctx) {
                capturedCtx = ctx;
                return const Text('child');
              },
            ),
            apideeppoolauthz: ({options = const []}) {
              callCount++;
              return Future.value(_response('bearer-call-$callCount', _pastExpiry()));
            },
          ),
        );
        await tester.pumpAndSettle();

        final cache = DeeppoolAuthzCache.of(capturedCtx!);

        // First token() call fetches since token starts expired
        final firstToken = await cache.meta.auto();
        await tester.pumpAndSettle();
        expect(callCount, equals(1));
        expect(firstToken.bearer, equals('bearer-call-1'));

        // Second token() call re-fetches because the returned token is also expired
        final secondToken = await cache.meta.auto();
        await tester.pumpAndSettle();

        expect(callCount, equals(2));
        expect(secondToken.bearer, equals('bearer-call-2'));
        expect(tester.takeException(), isNull);
      });

      testWidgets('does not re-fetch when token is still valid', (WidgetTester tester) async {
        var callCount = 0;
        BuildContext? capturedCtx;

        await tester.pumpApp(
          DeeppoolAuthzCache(
            Builder(
              builder: (ctx) {
                capturedCtx = ctx;
                return const Text('child');
              },
            ),
            apideeppoolauthz: ({options = const []}) {
              callCount++;
              return Future.value(_response('valid-bearer', _futureExpiry()));
            },
          ),
        );
        await tester.pumpAndSettle();

        final cache = DeeppoolAuthzCache.of(capturedCtx!);

        // First token() call fetches the token
        final firstToken = await cache.meta.auto();
        await tester.pumpAndSettle();
        expect(callCount, equals(1));
        expect(firstToken.bearer, equals('valid-bearer'));

        // Second token() call uses the cached non-expired token, no re-fetch
        final token = await cache.meta.auto();
        await tester.pumpAndSettle();

        expect(callCount, equals(1));
        expect(token.bearer, equals('valid-bearer'));
        expect(tester.takeException(), isNull);
      });

      testWidgets('bearer string changes after expiry and refresh', (WidgetTester tester) async {
        var callCount = 0;
        BuildContext? capturedCtx;

        await tester.pumpApp(
          DeeppoolAuthzCache(
            Builder(
              builder: (ctx) {
                capturedCtx = ctx;
                return const Text('child');
              },
            ),
            apideeppoolauthz: ({options = const []}) {
              callCount++;
              return Future.value(_response('bearer-v$callCount', _pastExpiry()));
            },
          ),
        );
        await tester.pumpAndSettle();

        final cache = DeeppoolAuthzCache.of(capturedCtx!);

        // First token() fetch — token comes back expired
        final firstToken = await cache.meta.auto();
        await tester.pumpAndSettle();
        expect(callCount, equals(1));
        expect(firstToken.bearer, equals('bearer-v1'));

        // Second token() re-fetches because the token is expired
        final secondToken = await cache.meta.auto();
        await tester.pumpAndSettle();

        expect(callCount, equals(2));
        expect(secondToken.bearer, equals('bearer-v2'));
        expect(tester.takeException(), isNull);
      });
    });
  });

  group('LocalOnlyGuard', () {
    meta.AuthzResponse localMetaResponse(bool localOnly) => meta.AuthzResponse(
      bearer: 'meta-bearer',
      token: meta.Token()
        ..localOnly = localOnly
        ..expires = _futureExpiry(),
    );

    testWidgets('does not mount DeeppoolAuthzCache while the local identity is local_only', (
      WidgetTester tester,
    ) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Material(
            child: ds.LoadingGuard(
              AuthzCache(
                LocalOnlyGuard(const Text('protected'), (child) => DeeppoolAuthzCache(child)),
                current: ({String? host}) => Future.value(localMetaResponse(true)),
              ),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('protected'), findsOneWidget);
      expect(find.byType(DeeppoolAuthzCache), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('mounts DeeppoolAuthzCache once the local token is not local_only', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Material(
            child: ds.LoadingGuard(
              AuthzCache(
                LocalOnlyGuard(const Text('protected'), (child) => DeeppoolAuthzCache(child)),
                current: ({String? host}) => Future.value(localMetaResponse(false)),
              ),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('protected'), findsOneWidget);
      expect(find.byType(DeeppoolAuthzCache), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('reacts when the meta token refreshes from local_only to a registered account', (
      WidgetTester tester,
    ) async {
      var localOnly = true;

      await tester.pumpWidget(
        MaterialApp(
          home: Material(
            child: ds.LoadingGuard(
              AuthzCache(
                LocalOnlyGuard(const Text('protected'), (child) => DeeppoolAuthzCache(child)),
                current: ({String? host}) => Future.value(localMetaResponse(localOnly)),
              ),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(DeeppoolAuthzCache), findsNothing);

      localOnly = false;
      final ctx = tester.element(find.text('protected'));
      AuthzCache.of(ctx).refresh();
      await tester.pump(); // let the new Cached reach AuthzTokenData before forcing a fetch on it
      await AuthzCache.meta(ctx).auto();
      await tester.pumpAndSettle();

      expect(find.byType(DeeppoolAuthzCache), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('reacts when the meta token refreshes from a registered account to local_only', (
      WidgetTester tester,
    ) async {
      var localOnly = false;

      await tester.pumpWidget(
        MaterialApp(
          home: Material(
            child: ds.LoadingGuard(
              AuthzCache(
                LocalOnlyGuard(const Text('protected'), (child) => DeeppoolAuthzCache(child)),
                current: ({String? host}) => Future.value(localMetaResponse(localOnly)),
              ),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(DeeppoolAuthzCache), findsOneWidget);

      localOnly = true;
      final ctx = tester.element(find.text('protected'));
      AuthzCache.of(ctx).refresh();
      await tester.pump(); // let the new Cached reach AuthzTokenData before forcing a fetch on it
      await AuthzCache.meta(ctx).auto();
      await tester.pumpAndSettle();

      expect(find.byType(DeeppoolAuthzCache), findsNothing);
      expect(tester.takeException(), isNull);
    });
  });
}
