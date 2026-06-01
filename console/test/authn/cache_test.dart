import 'dart:async';
import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/authn/cache.dart';
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/uuidx.dart' as uuidx;

meta.AuthzResponse _makeResponse() => meta.AuthzResponse(
  bearer: uuidx.min(),
  token: meta.Token()..expires = fixnum.Int64(DateTime.now().millisecondsSinceEpoch + 3600000),
);

void main() {
  group('AuthzCache', () {
    testWidgets('does not render child until token resolves', (WidgetTester tester) async {
      final completer = Completer<meta.AuthzResponse>();

      await tester.pumpWidget(
        MaterialApp(
          home: Material(
            child: AuthzCache(
              const Text('protected content'),
              current: ({String? host}) => completer.future,
            ),
          ),
        ),
      );
      await tester.pump();

      expect(find.text('protected content'), findsNothing);

      completer.complete(_makeResponse());
      await tester.pumpAndSettle();

      expect(find.text('protected content'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
