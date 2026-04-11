import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/profiles.dart' as profiles;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('Profile Create Widget Tests', () {
    late meta.Profile testProfile;
    late meta.Token testToken;
    const testPublicKey = 'ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAITest test@test';

    setUp(() {
      testProfile = meta.Profile()..display = 'Test User';
      testToken = meta.Token()..libraryRead = true;
    });

    group('layout', () {
      testWidgets('renders without overflow in unconstrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Center(
            child: profiles.Create(testProfile, testPublicKey, testToken),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
        expect(find.byType(profiles.Create), findsOneWidget);
      });

      testWidgets('renders without overflow in constrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          ConstrainedBox(
            constraints: BoxConstraints(maxWidth: 400, maxHeight: 600),
            child: SingleChildScrollView(
              child: profiles.Create(
                testProfile,
                testPublicKey,
                testToken,
              ),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in SingleChildScrollView', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SingleChildScrollView(
            child: profiles.Create(testProfile, testPublicKey, testToken),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });
    });

    group('resolutions', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          SingleChildScrollView(
            child: profiles.Create(testProfile, testPublicKey, testToken),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('callbacks', () {
      testWidgets('calls onChange when profile is modified', (
        WidgetTester tester,
      ) async {
        meta.Profile? changedProfile;
        String? changedKey;
        meta.Token? changedToken;

        await tester.pumpApp(
          profiles.Create(
            testProfile,
            testPublicKey,
            testToken,
            onChange: (p, k, t) {
              changedProfile = p;
              changedKey = k;
              changedToken = t;
            },
          ),
        );
        await tester.pumpAndSettle();

        // Find and modify the name field
        final nameField = find.widgetWithText(TextFormField, 'Test User');
        await tester.enterText(nameField, 'Updated User');

        expect(changedProfile, isNotNull);
        expect(changedProfile?.display, 'Updated User');
        expect(changedKey, testPublicKey);
        expect(changedToken, testToken);
      });

      testWidgets('calls onChange when public key is modified', (
        WidgetTester tester,
      ) async {
        meta.Profile? changedProfile;
        String? changedKey;
        meta.Token? changedToken;

        await tester.pumpApp(
          profiles.Create(
            testProfile,
            testPublicKey,
            testToken,
            onChange: (p, k, t) {
              changedProfile = p;
              changedKey = k;
              changedToken = t;
            },
          ),
        );
        await tester.pumpAndSettle();

        // Find and modify the public key field
        final keyField = find.widgetWithText(TextFormField, testPublicKey);
        await tester.enterText(keyField, 'new-key');

        expect(changedProfile, testProfile);
        expect(changedKey, 'new-key');
        expect(changedToken, testToken);
      });

      testWidgets('calls onChange when permissions are modified', (
        WidgetTester tester,
      ) async {
        meta.Profile? changedProfile;
        String? changedKey;
        meta.Token? changedToken;

        await tester.pumpApp(
          profiles.Create(
            testProfile,
            testPublicKey,
            testToken,
            onChange: (p, k, t) {
              changedProfile = p;
              changedKey = k;
              changedToken = t;
            },
          ),
        );
        await tester.pumpAndSettle();

        // Toggle a permission checkbox
        final libraryModifyCheckbox = find.widgetWithText(
          forms.Checkbox,
          'Library Modify',
        );
        await tester.tap(libraryModifyCheckbox);

        expect(changedProfile, testProfile);
        expect(changedKey, testPublicKey);
        expect(changedToken, isNotNull);
        expect(changedToken?.libraryModify, true);
      });

      testWidgets('handles onChange being null', (WidgetTester tester) async {
        await tester.pumpApp(
          profiles.Create(testProfile, testPublicKey, testToken),
        );
        await tester.pumpAndSettle();

        // Modify a field - should not throw
        final nameField = find.widgetWithText(TextFormField, 'Test User');
        await tester.enterText(nameField, 'Updated User');

        expect(tester.takeException(), isNull);
      });
    });
  });
}
