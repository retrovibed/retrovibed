import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/community/qr.scanner.dart';
import 'package:retrovibed/community/community.pb.dart';
import 'package:retrovibed/community/qr.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('QRScannerModal', () {
    group('resolutions', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          QRScannerModal(
            onScanned: (Community community, String attribution) => Future.value(),
            onCancel: () {},
            camera: (onDetect, onError) => const SizedBox(),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('Scan QR Code'), findsOneWidget);
        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('confirmation', () {
      final community = Community(
        url: 'https://testdomain.community.retrovibe.space',
        description: 'A test community',
        createdAt: '2024-01-15T14:30:00Z',
      );

      testWidgets('shows confirmation after valid QR scan', (tester) async {
        late void Function(String) capturedOnDetect;

        await tester.pumpApp(
          QRScannerModal(
            onScanned: (_, __) => Future.value(),
            onCancel: () {},
            camera: (onDetect, onError) {
              capturedOnDetect = onDetect;
              return const SizedBox();
            },
          ),
        );
        await tester.pumpAndSettle();

        capturedOnDetect(encodeQRPayload(community));
        await tester.pumpAndSettle();

        expect(find.text('Yes'), findsOneWidget);
        expect(find.text('No'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('confirmation shows community detail', (tester) async {
        late void Function(String) capturedOnDetect;

        await tester.pumpApp(
          QRScannerModal(
            onScanned: (_, __) => Future.value(),
            onCancel: () {},
            camera: (onDetect, onError) {
              capturedOnDetect = onDetect;
              return const SizedBox();
            },
          ),
        );
        await tester.pumpAndSettle();

        capturedOnDetect(encodeQRPayload(community));
        await tester.pumpAndSettle();

        expect(find.text('testdomain'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('confirming calls onScanned with community', (tester) async {
        late void Function(String) capturedOnDetect;
        Community? scanned;

        await tester.pumpApp(
          QRScannerModal(
            onScanned: (c, a) {
              scanned = c;
              return Future.value();
            },
            onCancel: () {},
            camera: (onDetect, onError) {
              capturedOnDetect = onDetect;
              return const SizedBox();
            },
          ),
        );
        await tester.pumpAndSettle();

        capturedOnDetect(encodeQRPayload(community));
        await tester.pumpAndSettle();

        await tester.tap(find.text('Yes'));
        await tester.pumpAndSettle();

        expect(scanned, isNotNull);
        expect(scanned!.url, 'https://testdomain.community.retrovibe.space');
        expect(tester.takeException(), isNull);
      });

      testWidgets('cancelling dismisses confirmation', (tester) async {
        late void Function(String) capturedOnDetect;

        await tester.pumpApp(
          QRScannerModal(
            onScanned: (_, __) => Future.value(),
            onCancel: () {},
            camera: (onDetect, onError) {
              capturedOnDetect = onDetect;
              return const SizedBox();
            },
          ),
        );
        await tester.pumpAndSettle();

        capturedOnDetect(encodeQRPayload(community));
        await tester.pumpAndSettle();

        expect(find.text('Yes'), findsOneWidget);

        await tester.tap(find.text('No'));
        await tester.pumpAndSettle();

        expect(find.text('Yes'), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('invalid QR data shows error', (tester) async {
        late void Function(String) capturedOnDetect;

        await tester.pumpApp(
          QRScannerModal(
            onScanned: (_, __) => Future.value(),
            onCancel: () {},
            camera: (onDetect, onError) {
              capturedOnDetect = onDetect;
              return const SizedBox();
            },
          ),
        );
        await tester.pumpAndSettle();

        capturedOnDetect('not valid json at all');
        await tester.pumpAndSettle();

        expect(find.text('Yes'), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('ignores duplicate scans while processing', (tester) async {
        late void Function(String) capturedOnDetect;

        await tester.pumpApp(
          QRScannerModal(
            onScanned: (_, __) => Future.value(),
            onCancel: () {},
            camera: (onDetect, onError) {
              capturedOnDetect = onDetect;
              return const SizedBox();
            },
          ),
        );
        await tester.pumpAndSettle();

        final payload = encodeQRPayload(community);
        capturedOnDetect(payload);
        capturedOnDetect(payload);
        await tester.pumpAndSettle();

        expect(find.text('Yes'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('camera error resets processing allowing rescan', (tester) async {
        late void Function(String) capturedOnDetect;
        late void Function(Object, StackTrace) capturedOnError;

        await tester.pumpApp(
          QRScannerModal(
            onScanned: (_, __) => Future.value(),
            onCancel: () {},
            camera: (onDetect, onError) {
              capturedOnDetect = onDetect;
              capturedOnError = onError;
              return const SizedBox();
            },
          ),
        );
        await tester.pumpAndSettle();

        // Start a scan so _processing = true
        capturedOnDetect(encodeQRPayload(community));
        await tester.pumpAndSettle();

        // Dismiss via error reset
        capturedOnError(Exception('camera failed'), StackTrace.empty);
        await tester.pumpAndSettle();

        expect(find.text('Yes'), findsNothing);

        // Should be able to scan again
        capturedOnDetect(encodeQRPayload(community));
        await tester.pumpAndSettle();

        expect(find.text('Yes'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('onScanned error is handled locally', (tester) async {
        late void Function(String) capturedOnDetect;
        Community? scanned;

        await tester.pumpApp(
          QRScannerModal(
            onScanned: (_, __) => Future.error('scan failed').then((_) => scanned = community),
            onCancel: () {},
            camera: (onDetect, onError) {
              capturedOnDetect = onDetect;
              return const SizedBox();
            },
          ),
        );
        await tester.pumpAndSettle();

        capturedOnDetect(encodeQRPayload(community));
        await tester.pumpAndSettle();

        await tester.tap(find.text('Yes'));
        await tester.pumpAndSettle();

        // Confirmation should be dismissed and error shown
        expect(find.text('Yes'), findsNothing);
        expect(find.text('an unexpected problem has occurred'), findsOneWidget);
        expect(scanned, isNull); // parent side effect didn't execute
        expect(tester.takeException(), isNull);
      });
    });
  });
}
