import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/storage/minimal.settings.dart';
import 'package:retrovibed/storage/api.dart' as api;
import 'package:retrovibed/storage/local.storage.dart';
import 'package:retrovibed/storage/archive.storage.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/quotas.dart' as quotas;

Future<quotas.QuotaFindResponse> _mockArchiveQuota(
  String sku, {
  List<httpx.Option> options = const [],
}) async =>
    quotas.QuotaFindResponse(quota: quotas.Quota(consumed: fixnum.Int64(0)));

void main() {
  group('MinimalSettings', () {
    testWidgets('renders within 384x384', (WidgetTester tester) async {
      final initialSettings = api.StorageSettingsResponse()..local = api.Local();

      await tester.pumpApp(
        Center(
          child: ConstrainedBox(
            constraints: BoxConstraints(maxWidth: 384, maxHeight: 384),
            child: MinimalSettings(
              initialSettings,
              onChange: (settings) async {
                return settings;
              },
              archiveQuota: _mockArchiveQuota,
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('library'), findsOneWidget);
      expect(find.text('cache storage'), findsOneWidget);
      expect(find.text('archive'), findsOneWidget);
      expect(find.text('cloud storage'), findsOneWidget);

      final RenderBox localRenderBox = tester.renderObject(
        find.byType(LocalStorageSettings),
      );
      expect(localRenderBox.size.width, 384.0);
      expect(localRenderBox.size.height, 100.0);

      final RenderBox archiveRenderBox = tester.renderObject(
        find.byType(ArchiveStorage),
      );
      expect(archiveRenderBox.size.width, 384.0);
      expect(archiveRenderBox.size.height, 91.0);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders correctly with initial data', (
      WidgetTester tester,
    ) async {
      final initialSettings = api.StorageSettingsResponse()..local = api.Local();

      await tester.pumpApp(
        MinimalSettings(
          initialSettings,
          onChange: (settings) async {
            return settings;
          },
          archiveQuota: _mockArchiveQuota,
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('library'), findsOneWidget);
      expect(find.text('cache storage'), findsOneWidget);
      expect(find.text('archive'), findsOneWidget);
      expect(find.text('cloud storage'), findsOneWidget);
    });

    testWidgets('handles local storage changes correctly', (
      WidgetTester tester,
    ) async {
      final initialSettings = api.StorageSettingsResponse()..local = api.Local();

      bool onChangeCalled = false;

      await tester.pumpApp(
        MinimalSettings(
          initialSettings,
          onChange: (settings) async {
            onChangeCalled = true;
            return settings;
          },
          archiveQuota: _mockArchiveQuota,
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('library'), findsOneWidget);
      expect(find.text('cache storage'), findsOneWidget);

      expect(onChangeCalled, isFalse);
    });

    testWidgets('renders LocalStorageSettings and ArchiveStorage components', (
      WidgetTester tester,
    ) async {
      final initialSettings = api.StorageSettingsResponse()..local = api.Local();

      await tester.pumpApp(
        MinimalSettings(
          initialSettings,
          onChange: (settings) async {
            return settings;
          },
          archiveQuota: _mockArchiveQuota,
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(LocalStorageSettings), findsOneWidget);
      expect(find.byType(ArchiveStorage), findsOneWidget);
    });
  });
}
