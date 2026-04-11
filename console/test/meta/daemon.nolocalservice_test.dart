import 'dart:io';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/meta/api.dart' as api;
import 'package:retrovibed/meta/daemon.auto.dart';
import 'package:retrovibed/meta/daemon.mdns.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

Future<api.DaemonLookupResponse> _offlineLatest() async {
  throw SocketException('', osError: OSError('', 111));
}

void main() {
  group('NoLocalService loose constraints resolutions', () {
    for (final entry in Resolutions.all.entries) {
      testWidgets('renders without overflow at ${entry.key}', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: entry.value,
          fit: FlexFit.loose,
          EndpointAuto(
            latest: _offlineLatest,
            backoff: httpx.Backoff.constant(Duration.zero),
            const SizedBox(),
          ),
        );

        await tester.pumpAndSettle();
        expect(tester.takeException(), isNull);
        final size = tester.getSize(find.byType(NoLocalService));
        expect(size.width, lessThanOrEqualTo(entry.value.width));
        expect(size.height, lessThanOrEqualTo(entry.value.height));
      });
    }
  });

  group('NoLocalService tight constraints resolutions', () {
    for (final entry in Resolutions.all.entries) {
      testWidgets('renders without overflow at ${entry.key}', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: entry.value,
          fit: FlexFit.tight,
          EndpointAuto(
            latest: _offlineLatest,
            backoff: httpx.Backoff.constant(Duration.zero),
            const SizedBox(),
          ),
        );

        await tester.pumpAndSettle();
        expect(tester.takeException(), isNull);
        final size = tester.getSize(find.byType(NoLocalService));
        expect(size.width, lessThanOrEqualTo(entry.value.width));
        expect(size.height, lessThanOrEqualTo(entry.value.height));
      });
    }
  });
}
