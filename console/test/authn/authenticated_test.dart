import 'dart:async';
import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/authn/authenticated.dart';
import 'package:retrovibed/authn/api.dart' as api;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

api.Session _session({int expiresSeconds = 0}) {
  return api.Session(token: 'session-token', expires: fixnum.Int64(expiresSeconds));
}

void main() {
  group('Authenticated', () {
    testWidgets('LoadingBoundary covers the initial session fetch, then clears', (
      WidgetTester tester,
    ) async {
      final completer = Completer<api.Authed>();

      await tester.pumpApp(
        ds.LoadingGuard(
          Authenticated(
            const Text('protected content'),
            apissh: () => completer.future,
            apisignup: () async => _session(),
            apicurrent: (token) async => _session(),
          ),
        ),
      );
      await tester.pump();
      await tester.pump();

      expect(find.byType(CircularProgressIndicator), findsOneWidget);
      expect(find.text('protected content'), findsNothing);

      completer.complete(api.Authed(profiles: const []));
      await tester.pumpAndSettle();

      expect(find.byType(CircularProgressIndicator), findsNothing);
      expect(find.text('protected content'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('a failed initial fetch still clears the loading boundary', (
      WidgetTester tester,
    ) async {
      final completer = Completer<api.Authed>();

      await tester.pumpApp(
        ds.LoadingGuard(
          Authenticated(
            const Text('protected content'),
            apissh: () => completer.future,
            apisignup: () async => _session(),
            apicurrent: (token) async => _session(),
          ),
        ),
      );
      await tester.pump();
      await tester.pump();

      expect(find.byType(CircularProgressIndicator), findsOneWidget);

      completer.completeError(Exception('boom'));
      await tester.pumpAndSettle();

      expect(find.byType(CircularProgressIndicator), findsNothing);
      expect(find.text('protected content'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets(
      'a later current() call after the initial fetch resolves does not re-arm the loading boundary',
      (WidgetTester tester) async {
        var sshCalls = 0;
        final second = Completer<api.Authed>();

        Future<api.Authed> apissh() {
          sshCalls++;
          // expiresSeconds: 0 below means _expires is already in the past, so the
          // first session immediately counts as expired and a later current()
          // call below triggers a genuine second fetch via this same function.
          if (sshCalls == 1) {
            return Future.value(api.Authed(profiles: [api.Authn(token: 'tok')]));
          }
          return second.future;
        }

        late BuildContext capturedContext;

        await tester.pumpApp(
          ds.LoadingGuard(
            Authenticated(
              Builder(
                builder: (context) {
                  capturedContext = context;
                  return const Text('protected content');
                },
              ),
              apissh: apissh,
              apisignup: () async => _session(),
              apicurrent: (token) async => _session(),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(CircularProgressIndicator), findsNothing);
        expect(find.text('protected content'), findsOneWidget);

        Authenticated.session(capturedContext);
        await tester.pump();
        await tester.pump();

        expect(sshCalls, equals(2));
        expect(find.byType(CircularProgressIndicator), findsNothing);
        expect(find.text('protected content'), findsOneWidget);

        second.complete(api.Authed(profiles: const []));
        await tester.pumpAndSettle();
        expect(tester.takeException(), isNull);
      },
    );
  });
}
