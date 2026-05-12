import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/profiles.dart' as profiles;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('AuthzMetaEdit Widget Tests', () {
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
            constraints: BoxConstraints(maxWidth: 400, maxHeight: 500),
            child: profiles.AuthzMetaEdit(testToken),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('User Management'), findsOneWidget);
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
              profiles.AuthzMetaEdit(testToken),
              Spacer(),
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
            constraints: BoxConstraints(maxWidth: 400, maxHeight: 500),
            child: Column(
              children: [
                profiles.AuthzMetaEdit(testToken),
              ],
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
          SingleChildScrollView(child: profiles.AuthzMetaEdit(testToken)),
        );
        await tester.pumpAndSettle();

        expect(find.text('User Management'), findsOneWidget);
        expect(find.text('Library Read'), findsOneWidget);
        expect(find.text('Library Modify'), findsOneWidget);
        expect(find.text('Community Modify'), findsOneWidget);
        expect(find.text('Billing Read'), findsOneWidget);
        expect(find.text('Billing Modify'), findsOneWidget);
        expect(find.text('Can manage user access'), findsOneWidget);
        expect(find.text('Can view library content'), findsOneWidget);
        expect(find.text('Can modify library content'), findsOneWidget);
        expect(find.text('Can modify community content'), findsOneWidget);
        expect(find.text('Can view billing information'), findsOneWidget);
        expect(find.text('Can modify billing settings'), findsOneWidget);
        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('permission checkboxes', () {
      testWidgets('displays initial permission values correctly', (
        WidgetTester tester,
      ) async {
        final token =
            meta.Token()
              ..usermanagement = true
              ..libraryRead = true
              ..libraryModify = false
              ..communityModify = true
              ..billingRead = false
              ..billingModify = false;

        await tester.pumpApp(profiles.AuthzMetaEdit(token));
        await tester.pumpAndSettle();

        final checkboxes = tester.widgetList<Checkbox>(find.byType(Checkbox));
        final checkboxList = checkboxes.toList();

        expect(checkboxList.length, equals(6));
        expect(checkboxList[0].value, equals(true)); // usermanagement
        expect(checkboxList[1].value, equals(true)); // libraryRead
        expect(checkboxList[2].value, equals(false)); // libraryModify
        expect(checkboxList[3].value, equals(true)); // communityModify
        expect(checkboxList[4].value, equals(false)); // billingRead
        expect(checkboxList[5].value, equals(false)); // billingModify
      });

      testWidgets('handles all permissions disabled', (
        WidgetTester tester,
      ) async {
        final token =
            meta.Token()
              ..usermanagement = false
              ..libraryRead = false
              ..libraryModify = false
              ..communityModify = false
              ..billingRead = false
              ..billingModify = false;

        await tester.pumpApp(profiles.AuthzMetaEdit(token));
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
              ..libraryRead = true
              ..libraryModify = true
              ..communityModify = true
              ..billingRead = true
              ..billingModify = true;

        await tester.pumpApp(profiles.AuthzMetaEdit(token));
        await tester.pumpAndSettle();

        final checkboxes = tester.widgetList<Checkbox>(find.byType(Checkbox));
        for (final checkbox in checkboxes) {
          expect(checkbox.value, equals(true));
        }
        expect(tester.takeException(), isNull);
      });

      testWidgets('calls onChange when usermanagement is toggled', (
        WidgetTester tester,
      ) async {
        meta.Token? changedToken;

        await tester.pumpApp(
          profiles.AuthzMetaEdit(
            testToken,
            onChange: (token) {
              changedToken = token;
            },
          ),
          authzCurrent: authn.AuthzCache.fakeWith(meta.Token()..usermanagement = true),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.byType(Checkbox).first);

        expect(changedToken, isNotNull);
        expect(changedToken?.usermanagement, isTrue);
      });

      testWidgets('calls onChange when libraryRead is toggled', (
        WidgetTester tester,
      ) async {
        meta.Token? changedToken;
        final token = meta.Token()..libraryRead = false;

        await tester.pumpApp(
          profiles.AuthzMetaEdit(
            token,
            onChange: (t) {
              changedToken = t;
            },
          ),
          authzCurrent: authn.AuthzCache.fakeWith(meta.Token()..usermanagement = true),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.byType(Checkbox).at(1));

        expect(changedToken, isNotNull);
        expect(changedToken?.libraryRead, isTrue);
      });

      testWidgets('calls onChange when billingModify is toggled', (
        WidgetTester tester,
      ) async {
        meta.Token? changedToken;

        await tester.pumpApp(
          profiles.AuthzMetaEdit(
            testToken,
            onChange: (token) {
              changedToken = token;
            },
          ),
          authzCurrent: authn.AuthzCache.fakeWith(meta.Token()..usermanagement = true),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.byType(Checkbox).last);

        expect(changedToken, isNotNull);
        expect(changedToken?.billingModify, isTrue);
      });

      testWidgets('multiple permission changes work correctly', (
        WidgetTester tester,
      ) async {
        final List<meta.Token> changedTokens = [];

        await tester.pumpApp(
          profiles.AuthzMetaEdit(
            testToken,
            onChange: (token) {
              changedTokens.add(token);
            },
          ),
          authzCurrent: authn.AuthzCache.fakeWith(meta.Token()..usermanagement = true),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.byType(Checkbox).at(0)); // User Management
        await tester.tap(find.byType(Checkbox).at(2)); // Library Modify
        await tester.tap(find.byType(Checkbox).at(4)); // Billing Read

        expect(changedTokens.length, equals(3));
        expect(changedTokens[0].usermanagement, isTrue);
        expect(changedTokens[1].libraryModify, isTrue);
        expect(changedTokens[2].billingRead, isTrue);
      });
    });
  });
}
