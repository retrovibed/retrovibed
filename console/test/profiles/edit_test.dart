import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/profiles.dart' as profiles;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('Profile Edit Widget Tests', () {
    late meta.Profile testProfile;

    setUp(() {
      testProfile = meta.Profile()..display = 'Test User';
    });

    group('layout', () {
      testWidgets('renders without overflow in unconstrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Center(
            child: profiles.Edit(
              testProfile,
              pkey: 'ssh-ed25519 AAAA...',
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('name'), findsOneWidget);
        expect(find.text('public key'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in constrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          ConstrainedBox(
            constraints: BoxConstraints(maxWidth: 400, maxHeight: 300),
            child: profiles.Edit(
              testProfile,
              pkey: 'ssh-ed25519 AAAA...',
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('name'), findsOneWidget);
        expect(find.text('public key'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in Column', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Column(
            children: [
              profiles.Edit(
                testProfile,
                pkey: 'ssh-ed25519 AAAA...',
              ),
              Spacer(),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('name'), findsOneWidget);
        expect(find.text('public key'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in constrained Column', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          ConstrainedBox(
            constraints: BoxConstraints(maxWidth: 400, maxHeight: 300),
            child: Column(
              children: [
                profiles.Edit(
                  testProfile,
                  pkey: 'ssh-ed25519 AAAA...',
                ),
              ],
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('name'), findsOneWidget);
        expect(find.text('public key'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in SingleChildScrollView', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SingleChildScrollView(
            child: profiles.Edit(
              testProfile,
              pkey: 'ssh-ed25519 AAAA...',
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('name'), findsOneWidget);
        expect(find.text('public key'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('resolutions', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          SingleChildScrollView(
            child: profiles.Edit(
              testProfile,
              pkey: 'ssh-ed25519 AAAA...',
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('name'), findsOneWidget);
        expect(find.text('public key'), findsOneWidget);
        expect(tester.takeException(), isNull);
      }, variant: _resolutions);

      testWidgets('renders with long public key without overflow', (
        WidgetTester tester,
      ) async {
        final entry = _resolutions.currentValue!;
        final longPublicKey =
            'ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILong'
            'PublicKeyThatExtendsAcrossMultipleLinesAndMightCauseIssuesOnSmall'
            'ScreensIfNotHandledProperly user@example.com';

        await tester.pumpApp(
          physicalSize: entry.value,
          SingleChildScrollView(
            child: profiles.Edit(
              testProfile,
              pkey: longPublicKey,
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('name'), findsOneWidget);
        expect(find.text('public key'), findsOneWidget);
        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('form fields', () {
      testWidgets('displays initial values correctly', (
        WidgetTester tester,
      ) async {
        const testName = 'John Doe';
        const testKey = 'ssh-ed25519 AAAA... john@example.com';

        await tester.pumpApp(
          profiles.Edit(
            meta.Profile()..display = testName,
            pkey: testKey,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text(testName), findsOneWidget);
        expect(find.text(testKey), findsOneWidget);
      });

      testWidgets('handles empty initial values', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          profiles.Edit(
            meta.Profile(),
            pkey: '',
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('name'), findsOneWidget);
        expect(find.text('public key'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('calls onChange when name is modified', (
        WidgetTester tester,
      ) async {
        meta.Profile? changedProfile;
        String? changedKey;

        await tester.pumpApp(
          profiles.Edit(
            testProfile,
            pkey: 'ssh-ed25519 AAAA...',
            onChange: (profile, key) {
              changedProfile = profile;
              changedKey = key;
            },
          ),
        );
        await tester.pumpAndSettle();

        // Find the name field and enter text
        final nameField = find.byType(TextFormField).first;
        await tester.enterText(nameField, 'New Name');

        expect(changedProfile, isNotNull);
        expect(changedProfile?.display, equals('New Name'));
        expect(changedKey, equals('ssh-ed25519 AAAA...'));
      });

      testWidgets('calls onChange when public key is modified', (
        WidgetTester tester,
      ) async {
        meta.Profile? changedProfile;
        String? changedKey;

        await tester.pumpApp(
          profiles.Edit(
            testProfile,
            pkey: 'ssh-ed25519 AAAA...',
            onChange: (profile, key) {
              changedProfile = profile;
              changedKey = key;
            },
          ),
        );
        await tester.pumpAndSettle();

        // Find the public key field and enter text
        final keyField = find.byType(TextFormField).last;
        await tester.enterText(keyField, 'ssh-ed25519 NEWKEY...');

        expect(changedProfile, isNotNull);
        expect(changedKey, equals('ssh-ed25519 NEWKEY...'));
      });

      testWidgets('public key field allows multiline input', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          profiles.Edit(
            testProfile,
            pkey: 'Line 1\nLine 2\nLine 3',
          ),
        );
        await tester.pumpAndSettle();

        // Verify the public key field renders multiline content without overflow
        expect(find.text('Line 1\nLine 2\nLine 3'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('accessibility', () {
      testWidgets('form fields have proper helper text', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          profiles.Edit(
            testProfile,
            pkey: 'ssh-ed25519 AAAA...',
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('name'), findsOneWidget);
        expect(find.text('public key'), findsOneWidget);
      });
    });
  });
}
