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
  });
}
