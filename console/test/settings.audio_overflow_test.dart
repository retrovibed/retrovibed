import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/discovery/settings.audio.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('SettingsAudio overflow', () {
    final resolutions = Resolutions.variant();
    testWidgets(
      'renders without overflow',
      (tester) async {
        await tester.pumpApp(
          const SettingsAudio(),
          physicalSize: resolutions.currentValue!.value,
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      },
      variant: resolutions,
    );

    testWidgets('dropdown items render without horizontal overflow', (
      tester,
    ) async {
      await tester.pumpApp(
        const SettingsAudio(),
        physicalSize: Size(360, 640),
        alignment: Alignment.topLeft,
      );
      await tester.pumpAndSettle();

      expect(find.text('Bitrate'), findsOneWidget);
      expect(find.text('Language'), findsOneWidget);
      expect(find.text('High (256 kbps)'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('bitrate dropdown can be opened without overflow', (
      tester,
    ) async {
      await tester.pumpApp(const SettingsAudio(), physicalSize: Size(360, 640));
      await tester.pumpAndSettle();

      final bitrateDropdown = find.text('High (256 kbps)');
      expect(bitrateDropdown, findsOneWidget);

      await tester.tap(bitrateDropdown);
      await tester.pumpAndSettle();

      // All bitrate options should be visible in dropdown
      expect(find.text('Auto'), findsOneWidget);
      expect(find.text('Low (64 kbps)'), findsOneWidget);
      expect(find.text('Medium (128 kbps)'), findsOneWidget);
      expect(find.text('High (256 kbps)'), findsExactly(2));
      expect(find.text('Lossless (FLAC/ALAC)'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('respects margin parameter', (tester) async {
      await tester.pumpApp(
        const SettingsAudio(margin: EdgeInsets.all(20)),
        physicalSize: Size(390, 844),
      );
      await tester.pumpAndSettle();

      expect(find.byType(SettingsAudio), findsOneWidget);
    });
  });
}
