import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/community.dart';
import 'package:retrovibed/community/qr.scanner.dart';
import 'package:retrovibed/design.kit/theme.defaults.dart';

void main() {
  group('Management', () {
    testWidgets('displays QRScannerModal when QR button is clicked', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Management(),
        theme: ThemeData().copyWith(
          extensions: [
            Defaults(desktop: false),
          ],
        ),
      );
      await tester.pumpAndSettle();

      // Verify the QR button exists
      expect(
        find.byIcon(Icons.qr_code_scanner),
        findsOneWidget,
        reason: 'QR scanner button should be visible',
      );

      // Click the QR button
      await tester.tap(find.byIcon(Icons.qr_code_scanner));
      await tester.pumpAndSettle();

      // Verify QRScannerModal is displayed
      expect(find.byType(QRScannerModal), findsOneWidget);
    });
  });
}
