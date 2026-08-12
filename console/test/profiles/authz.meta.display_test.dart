import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/profiles.dart' as profiles;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('AuthzMetaDisplay', () {
    late meta.Token testToken;

    setUp(() {
      testToken = meta.Token()..libraryRead = true;
    });

    group('layout', () {
      testWidgets('renders without overflow in constrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 400, maxHeight: 500),
            child: profiles.AuthzMetaDisplay(testToken),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('User Management'), findsOneWidget);
        expect(find.text('Remote Control'), findsOneWidget);
        expect(find.text('Library Read'), findsOneWidget);
        expect(find.text('Library Modify'), findsOneWidget);
        expect(find.text('Community Modify'), findsOneWidget);
        expect(find.text('Billing Read'), findsOneWidget);
        expect(find.text('Billing Modify'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in Column', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Column(
            children: [
              profiles.AuthzMetaDisplay(testToken),
              const Spacer(),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('User Management'), findsOneWidget);
        expect(find.text('Library Read'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in constrained Column', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 400, maxHeight: 560),
            child: Column(
              children: [profiles.AuthzMetaDisplay(testToken)],
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('User Management'), findsOneWidget);
        expect(find.text('Library Read'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('resolutions', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          SingleChildScrollView(child: profiles.AuthzMetaDisplay(testToken)),
        );
        await tester.pumpAndSettle();

        expect(find.text('User Management'), findsOneWidget);
        expect(find.text('Remote Control'), findsOneWidget);
        expect(find.text('Library Read'), findsOneWidget);
        expect(find.text('Library Modify'), findsOneWidget);
        expect(find.text('Community Modify'), findsOneWidget);
        expect(find.text('Billing Read'), findsOneWidget);
        expect(find.text('Billing Modify'), findsOneWidget);
        expect(find.text('Can manage user access'), findsOneWidget);
        expect(find.text("Can use remote control to connect to and drive another device's playback"), findsOneWidget);
        expect(find.text('Can view library content'), findsOneWidget);
        expect(find.text('Can modify library content'), findsOneWidget);
        expect(find.text('Can modify community content'), findsOneWidget);
        expect(find.text('Can view billing information'), findsOneWidget);
        expect(find.text('Can modify billing settings'), findsOneWidget);
        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('display values', () {
      testWidgets('displays initial permission values correctly', (
        WidgetTester tester,
      ) async {
        final token =
            meta.Token()
              ..usermanagement = true
              ..remoteControl = false
              ..libraryRead = true
              ..libraryModify = false
              ..communityModify = true
              ..billingRead = false
              ..billingModify = false;

        await tester.pumpApp(profiles.AuthzMetaDisplay(token));
        await tester.pumpAndSettle();

        final checkboxes = tester.widgetList<Checkbox>(find.byType(Checkbox));
        final checkboxList = checkboxes.toList();

        expect(checkboxList.length, equals(7));
        expect(checkboxList[0].value, equals(true)); // usermanagement
        expect(checkboxList[1].value, equals(false)); // remoteControl
        expect(checkboxList[2].value, equals(true)); // libraryRead
        expect(checkboxList[3].value, equals(false)); // libraryModify
        expect(checkboxList[4].value, equals(true)); // communityModify
        expect(checkboxList[5].value, equals(false)); // billingRead
        expect(checkboxList[6].value, equals(false)); // billingModify
        expect(tester.takeException(), isNull);
      });

      testWidgets('all checkboxes are non-interactive', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(profiles.AuthzMetaDisplay(testToken));
        await tester.pumpAndSettle();

        final checkboxes = tester.widgetList<Checkbox>(find.byType(Checkbox));
        for (final checkbox in checkboxes) {
          expect(checkbox.onChanged, isNull);
        }
        expect(tester.takeException(), isNull);
      });

      testWidgets('handles all permissions disabled', (
        WidgetTester tester,
      ) async {
        final token =
            meta.Token()
              ..usermanagement = false
              ..remoteControl = false
              ..libraryRead = false
              ..libraryModify = false
              ..communityModify = false
              ..billingRead = false
              ..billingModify = false;

        await tester.pumpApp(profiles.AuthzMetaDisplay(token));
        await tester.pumpAndSettle();

        final checkboxes = tester.widgetList<Checkbox>(find.byType(Checkbox));
        for (final checkbox in checkboxes) {
          expect(checkbox.value, equals(false));
        }
        expect(tester.takeException(), isNull);
      });

      testWidgets('handles all permissions enabled', (
        WidgetTester tester,
      ) async {
        final token =
            meta.Token()
              ..usermanagement = true
              ..remoteControl = true
              ..libraryRead = true
              ..libraryModify = true
              ..communityModify = true
              ..billingRead = true
              ..billingModify = true;

        await tester.pumpApp(profiles.AuthzMetaDisplay(token));
        await tester.pumpAndSettle();

        final checkboxes = tester.widgetList<Checkbox>(find.byType(Checkbox));
        for (final checkbox in checkboxes) {
          expect(checkbox.value, equals(true));
        }
        expect(tester.takeException(), isNull);
      });
    });
  });
}
